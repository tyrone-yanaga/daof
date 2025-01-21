// internal/services/cart.go
package services

import (
	"context"
	"errors"
	"time"

	"your-project/internal/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type CartService struct {
	db             *gorm.DB
	redis          *redis.Client
	productService *ProductService
}

func NewCartService(db *gorm.DB, redis *redis.Client, productService *ProductService) *CartService {
	return &CartService{
		db:             db,
		redis:          redis,
		productService: productService,
	}
}

func (s *CartService) Get(ctx context.Context, cartID uint) (*models.Cart, error) {
	var cart models.Cart
	err := s.db.Preload("Items").First(&cart, cartID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCartNotFound
		}
		return nil, err
	}

	// Update cart items with current prices and stock levels
	for i, item := range cart.Items {
		product, err := s.productService.Get(ctx, item.ProductID)
		if err != nil {
			continue
		}
		cart.Items[i].Price = product.Price
	}

	return &cart, nil
}

func (s *CartService) Create(ctx context.Context) (*models.Cart, error) {
	cart := &models.Cart{
		Status:    "active",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.db.Create(cart).Error; err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *CartService) AddItem(ctx context.Context, cartID uint, item *models.CartItem) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check if product exists and has sufficient stock
		product, err := s.productService.Get(ctx, item.ProductID)
		if err != nil {
			return err
		}

		if product.StockLevel < item.Quantity {
			return ErrInsufficientStock
		}

		// Check if item already exists in cart
		var existingItem models.CartItem
		err = tx.Where("cart_id = ? AND product_id = ? AND variation_id = ?",
			cartID, item.ProductID, item.VariationID).First(&existingItem).Error

		if err == nil {
			// Update existing item
			existingItem.Quantity += item.Quantity
			existingItem.Price = product.Price
			return tx.Save(&existingItem).Error
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new item
			item.CartID = cartID
			item.Price = product.Price
			return tx.Create(item).Error
		}

		return err
	})
}

func (s *CartService) UpdateItem(ctx context.Context, cartID uint, itemID uint, quantity int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var item models.CartItem
		if err := tx.Where("cart_id = ? AND id = ?", cartID, itemID).First(&item).Error; err != nil {
			return err
		}

		// Check stock availability
		product, err := s.productService.Get(ctx, item.ProductID)
		if err != nil {
			return err
		}

		if product.StockLevel < quantity {
			return ErrInsufficientStock
		}

		// Update quantity and price
		item.Quantity = quantity
		item.Price = product.Price

		return tx.Save(&item).Error
	})
}

func (s *CartService) RemoveItem(ctx context.Context, cartID uint, itemID uint) error {
	return s.db.Where("cart_id = ? AND id = ?", cartID, itemID).Delete(&models.CartItem{}).Error
}

func (s *CartService) Clear(ctx context.Context, cartID uint) error {
	return s.db.Where("cart_id = ?", cartID).Delete(&models.CartItem{}).Error
}

func (s *CartService) ApplyCoupon(ctx context.Context, cartID uint, code string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var coupon models.Coupon
		if err := tx.Where("code = ? AND is_active = ? AND usage_count < usage_limit",
			code, true).First(&coupon).Error; err != nil {
			return err
		}

		// Validate coupon
		if time.Now().Before(coupon.StartDate) || time.Now().After(coupon.EndDate) {
			return ErrCouponExpired
		}

		// Get cart total
		cart, err := s.Get(ctx, cartID)
		if err != nil {
			return err
		}

		var total float64
		for _, item := range cart.Items {
			total += item.Price * float64(item.Quantity)
		}

		if total < coupon.MinimumSpend {
			return ErrMinimumSpendNotMet
		}

		// Apply discount
		coupon.UsageCount++
		if err := tx.Save(&coupon).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *CartService) GetTotals(ctx context.Context, cartID uint) (map[string]float64, error) {
	cart, err := s.Get(ctx, cartID)
	if err != nil {
		return nil, err
	}

	var subtotal float64
	for _, item := range cart.Items {
		subtotal += item.Price * float64(item.Quantity)
	}

	// Calculate tax (example rate)
	tax := subtotal * 0.1

	return map[string]float64{
		"subtotal": subtotal,
		"tax":      tax,
		"total":    subtotal + tax,
	}, nil
}
