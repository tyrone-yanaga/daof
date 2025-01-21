// internal/repository/postgres/cart.go
package postgres

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "your-project/internal/models"
    "your-project/internal/repository"
)

type cartRepository struct {
    db *gorm.DB
}

func NewCartRepository(db *gorm.DB) repository.CartRepository {
    return &cartRepository{db: db}
}

func (r *cartRepository) Create(ctx context.Context, cart *models.Cart) error {
    return r.db.WithContext(ctx).Create(cart).Error
}

func (r *cartRepository) GetByID(ctx context.Context, id uint) (*models.Cart, error) {
    var cart models.Cart
    err := r.db.WithContext(ctx).
        Preload("Items").
        First(&cart, id).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrNotFound
        }
        return nil, err
    }

    return &cart, nil
}

func (r *cartRepository) AddItem(ctx context.Context, item *models.CartItem) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Check if item already exists in cart
        var existingItem models.CartItem
        err := tx.Where("cart_id = ? AND product_id = ?", 
            item.CartID, item.ProductID).First(&existingItem).Error

        if err == nil {
            // Update existing item quantity
            existingItem.Quantity += item.Quantity
            return tx.Save(&existingItem).Error
        } else if errors.Is(err, gorm.ErrRecordNotFound) {
            // Create new item
            return tx.Create(item).Error
        }

        return err
    })
}

func (r *cartRepository) UpdateItem(ctx context.Context, item *models.CartItem) error {
    return r.db.WithContext(ctx).Save(item).Error
}

func (r *cartRepository) RemoveItem(ctx context.Context, cartID, itemID uint) error {
    return r.db.WithContext(ctx).
        Where("cart_id = ? AND id = ?", cartID, itemID).
        Delete(&models.CartItem{}).Error
}

func (r *cartRepository) Clear(ctx context.Context, cartID uint) error {
    return r.db.WithContext(ctx).
        Where("cart_id = ?", cartID).
        Delete(&models.CartItem{}).Error
}
