// internal/repository/postgres/product.go
package postgres

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "your-project/internal/models"
    "your-project/internal/repository"
)

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
    return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(product).Error; err != nil {
            return err
        }

        // Create associated images
        for i := range product.Images {
            product.Images[i].ProductID = product.ID
            if err := tx.Create(&product.Images[i]).Error; err != nil {
                return err
            }
        }

        // Create variations
        for i := range product.Variations {
            product.Variations[i].ProductID = product.ID
            if err := tx.Create(&product.Variations[i]).Error; err != nil {
                return err
            }
        }

        return nil
    })
}

func (r *productRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
    var product models.Product
    err := r.db.WithContext(ctx).
        Preload("Images").
        Preload("Variations").
        Preload("Category").
        First(&product, id).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrNotFound
        }
        return nil, err
    }

    return &product, nil
}

func (r *productRepository) List(ctx context.Context, page, limit int, filters repository.ProductFilter) ([]models.Product, int64, error) {
    var products []models.Product
    var total int64
    
    query := r.db.WithContext(ctx).Model(&models.Product{})

    // Apply filters
    if filters.CategoryID != nil {
        query = query.Where("category_id = ?", *filters.CategoryID)
    }

    if filters.MinPrice != nil {
        query = query.Where("price >= ?", *filters.MinPrice)
    }

    if filters.MaxPrice != nil {
        query = query.Where("price <= ?", *filters.MaxPrice)
    }

    if filters.Search != "" {
        search := "%" + filters.Search + "%"
        query = query.Where("name ILIKE ? OR description ILIKE ?", search, search)
    }

    if filters.InStock != nil && *filters.InStock {
        query = query.Where("stock_level > 0")
    }

    // Count total before pagination
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply pagination and get records
    err := query.
        Preload("Images").
        Preload("Variations").
        Preload("Category").
        Offset((page - 1) * limit).
        Limit(limit).
        Find(&products).Error

    if err != nil {
        return nil, 0, err
    }

    return products, total, nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Update main product
        if err := tx.Save(product).Error; err != nil {
            return err
        }

        // Update images
        if err := tx.Where("product_id = ?", product.ID).Delete(&models.ProductImage{}).Error; err != nil {
            return err
        }
        for i := range product.Images {
            product.Images[i].ProductID = product.ID
            if err := tx.Create(&product.Images[i]).Error; err != nil {
                return err
            }
        }

        // Update variations
        if err := tx.Where("product_id = ?", product.ID).Delete(&models.Variation{}).Error; err != nil {
            return err
        }
        for i := range product.Variations {
            product.Variations[i].ProductID = product.ID
            if err := tx.Create(&product.Variations[i]).Error; err != nil {
                return err
            }
        }

        return nil
    })
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Soft delete the product (will cascade to related entities)
        if err := tx.Delete(&models.Product{}, id).Error; err != nil {
            return err
        }
        return nil
    })
}

func (r *productRepository) UpdateStock(ctx context.Context, id uint, quantity int) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var product models.Product
        if err := tx.First(&product, id).Error; err != nil {
            return err
        }

        newStock := product.StockLevel + quantity
        if newStock < 0 {
            return repository.ErrInsufficientStock
        }

        if err := tx.Model(&product).Update("stock_level", newStock).Error; err != nil {
            return err
        }

        return nil
    })
}
