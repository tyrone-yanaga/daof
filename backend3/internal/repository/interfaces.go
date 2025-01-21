// internal/repository/interfaces.go
package repository

import (
    "context"
    "your-project/internal/models"
)

type ProductRepository interface {
    Create(ctx context.Context, product *models.Product) error
    GetByID(ctx context.Context, id uint) (*models.Product, error)
    List(ctx context.Context, page, limit int, filters ProductFilter) ([]models.Product, int64, error)
    Update(ctx context.Context, product *models.Product) error
    Delete(ctx context.Context, id uint) error
    UpdateStock(ctx context.Context, id uint, quantity int) error
}

type CartRepository interface {
    Create(ctx context.Context, cart *models.Cart) error
    GetByID(ctx context.Context, id uint) (*models.Cart, error)
    AddItem(ctx context.Context, item *models.CartItem) error
    UpdateItem(ctx context.Context, item *models.CartItem) error
    RemoveItem(ctx context.Context, cartID, itemID uint) error
    Clear(ctx context.Context, cartID uint) error
}

type OrderRepository interface {
    Create(ctx context.Context, order *models.Order) error
    GetByID(ctx context.Context, id uint) (*models.Order, error)
    List(ctx context.Context, userID uint, page, limit int) ([]models.Order, int64, error)
    UpdateStatus(ctx context.Context, id uint, status string) error
    GetByPaymentID(ctx context.Context, paymentID string) (*models.Order, error)
}

type UserRepository interface {
    Create(ctx context.Context, user *models.User) error
    GetByID(ctx context.Context, id uint) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    Update(ctx context.Context, user *models.User) error
    AddAddress(ctx context.Context, address *models.UserAddress) error
    ListAddresses(ctx context.Context, userID uint) ([]models.UserAddress, error)
}

type ProductFilter struct {
    CategoryID *uint
    MinPrice   *float64
    MaxPrice   *float64
    Search     string
    InStock    *bool
}