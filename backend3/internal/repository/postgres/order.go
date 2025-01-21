// internal/repository/postgres/order.go
package postgres

import (
    "context"
    "errors"
    "gorm.io/gorm"
    "your-project/internal/models"
    "your-project/internal/repository"
)

type orderRepository struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
    return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *models.Order) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(order).Error; err != nil {
            return err
        }

        // Create order items
        for i := range order.Items {
            order.Items[i].OrderID = order.ID
            if err := tx.Create(&order.Items[i]).Error; err != nil {
                return err
            }
        }

        return nil
    })
}

func (r *orderRepository) GetByID(ctx context.Context, id uint) (*models.Order, error) {
    var order models.Order
    err := r.db.WithContext(ctx).
        Preload("Items").
        Preload("Transactions").
        First(&order, id).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrNotFound
        }
        return nil, err
    }

    return &order, nil
}

func (r *orderRepository) List(ctx context.Context, userID uint, page, limit int) ([]models.Order, int64, error) {
    var orders []models.Order
    var total int64

    query := r.db.WithContext(ctx).Model(&models.Order{}).Where("user_id = ?", userID)

    // Get total count
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get orders with pagination
    err := query.
        Preload("Items").
        Preload("Transactions").
        Order("created_at DESC").
        Offset((page - 1) * limit).
        Limit(limit).
        Find(&orders).Error

    if err != nil {
        return nil, 0, err
    }

    return orders, total, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id uint, status string) error {
    return r.db.WithContext(ctx).
        Model(&models.Order{}).
        Where("id = ?", id).
        Update("status", status).Error
}

func (r *orderRepository) GetByPaymentID(ctx context.Context, paymentID string) (*models.Order, error) {
    var order models.Order
    err := r.db.WithContext(ctx).
        Where("payment_id = ?", paymentID).
        Preload("Items").
        Preload("Transactions").
        First(&order).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, repository.ErrNotFound
        }
        return nil, err
    }

    return &order, nil
}
