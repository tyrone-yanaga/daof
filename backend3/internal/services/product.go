// internal/services/product.go
package services

import (
	"context"
	"errors"
	"your-project/internal/models"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type ProductService struct {
	db         *gorm.DB
	odooClient OdooClient
}

type ProductFilters struct {
	Category string
	MinPrice float64
	MaxPrice float64
	Search   string
	InStock  *bool
}

func NewProductService(db *gorm.DB, odooClient OdooClient) *ProductService {
	return &ProductService{
		db:         db,
		odooClient: odooClient,
	}
}

func (s *ProductService) List(ctx context.Context, page, limit int, filters ProductFilters) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := s.db.Model(&models.Product{}).Where("status = ?", "active")

	// Apply filters
	if filters.Category != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.slug = ?", filters.Category)
	}

	if filters.MinPrice > 0 {
		query = query.Where("price >= ?", filters.MinPrice)
	}

	if filters.MaxPrice > 0 {
		query = query.Where("price <= ?", filters.MaxPrice)
	}

	if filters.Search != "" {
		search := "%" + filters.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
	}

	if filters.InStock != nil {
		if *filters.InStock {
			query = query.Where("stock_level > 0")
		}
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (page - 1) * limit
	err := query.Preload("Images").
		Preload("Category").
		Preload("Variations").
		Offset(offset).
		Limit(limit).
		Find(&products).Error

	if err != nil {
		return nil, 0, err
	}

	// Get stock levels from Odoo
	for i := range products {
		stockLevel, err := s.odooClient.GetStockLevel(ctx, products[i].SKU)
		if err == nil {
			products[i].StockLevel = stockLevel
		}
	}

	return products, total, nil
}

func (s *ProductService) Get(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product

	err := s.db.Preload("Images").
		Preload("Category").
		Preload("Variations").
		First(&product, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	// Get stock level from Odoo
	stockLevel, err := s.odooClient.GetStockLevel(ctx, product.SKU)
	if err == nil {
		product.StockLevel = stockLevel
	}

	return &product, nil
}

func (s *ProductService) Create(ctx context.Context, product *models.Product) error {
	// Generate slug
	product.Slug = slug.Make(product.Name)

	// Start transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create product
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		// Sync with Odoo
		if err := s.odooClient.CreateProduct(ctx, product); err != nil {
			return err
		}

		return nil
	})
}

func (s *ProductService) Update(ctx context.Context, product *models.Product) error {
	// Start transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update product
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		// Sync with Odoo
		if err := s.odooClient.UpdateProduct(ctx, product); err != nil {
			return err
		}

		return nil
	})
}

func (s *ProductService) Delete(ctx context.Context, id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		product := &models.Product{ID: id}

		// Soft delete the product
		if err := tx.Delete(product).Error; err != nil {
			return err
		}

		// Sync deletion with Odoo
		if err := s.odooClient.DeleteProduct(ctx, product.SKU); err != nil {
			return err
		}

		return nil
	})
}

func (s *ProductService) UpdateStock(ctx context.Context, id uint, quantity int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var product models.Product
		if err := tx.First(&product, id).Error; err != nil {
			return err
		}

		product.StockLevel += quantity
		if product.StockLevel < 0 {
			return ErrInsufficientStock
		}

		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		// Sync with Odoo
		if err := s.odooClient.UpdateStock(ctx, product.SKU, quantity); err != nil {
			return err
		}

		return nil
	})
}
