package services

import (
	"context"
	"ecommerce/internal/models"
	"ecommerce/pkg/redis"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CartService struct {
	redisClient    *redis.Client
	productService *ProductService
}

func NewCartService(redisClient *redis.Client, productService *ProductService) *CartService {
	return &CartService{
		redisClient:    redisClient,
		productService: productService,
	}
}

func (s *CartService) CreateCart(ctx context.Context, userID *uint) (*models.Cart, error) {
	cart := &models.Cart{
		ID:        uuid.New().String(),
		UserID:    userID,
		Items:     []models.CartItem{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // Cart expires in 24 hours
	}

	if err := s.saveCart(ctx, cart); err != nil {
		return nil, err
	}

	return cart, nil
}

func (s *CartService) GetCart(ctx context.Context, cartID string) (*models.Cart, error) {
	var cart models.Cart
	err := s.redisClient.Get(ctx, fmt.Sprintf("cart:%s", cartID), &cart)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if time.Now().After(cart.ExpiresAt) {
		return nil, fmt.Errorf("cart has expired")
	}

	return &cart, nil
}

func (s *CartService) AddToCart(ctx context.Context, cartID string, productID uint, quantity int) error {
	// Get cart
	cart, err := s.GetCart(ctx, cartID)
	if err != nil {
		return err
	}

	// Get product details
	product, err := s.productService.GetProduct(fmt.Sprint(productID))
	if err != nil {
		return fmt.Errorf("failed to get product: %w", err)
	}

	// Check stock availability
	if product.Stock < quantity {
		return fmt.Errorf("insufficient stock available")
	}

	// Create cart item
	item := models.CartItem{
		ProductID: productID,
		Quantity:  quantity,
		Price:     product.Price,
		Name:      product.Name,
		SKU:       product.SKU,
	}

	// Add to cart
	if err := cart.AddItem(item); err != nil {
		return err
	}

	// Save updated cart
	return s.saveCart(ctx, cart)
}

func (s *CartService) UpdateCartItem(ctx context.Context, cartID string, productID uint, quantity int) error {
	cart, err := s.GetCart(ctx, cartID)
	if err != nil {
		return err
	}

	if quantity > 0 {
		// Check stock availability
		product, err := s.productService.GetProduct(fmt.Sprint(productID))
		if err != nil {
			return fmt.Errorf("failed to get product: %w", err)
		}
		if product.Stock < quantity {
			return fmt.Errorf("insufficient stock available")
		}
	}

	if err := cart.UpdateItem(productID, quantity); err != nil {
		return err
	}

	return s.saveCart(ctx, cart)
}

func (s *CartService) saveCart(ctx context.Context, cart *models.Cart) error {
	cart.UpdatedAt = time.Now()
	return s.redisClient.Set(ctx, fmt.Sprintf("cart:%s", cart.ID), cart, 24*time.Hour)
}
