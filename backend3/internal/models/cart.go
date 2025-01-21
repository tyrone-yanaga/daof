// internal/models/cart.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    *uint          `json:"userId"`
	SessionID string         `json:"sessionId" gorm:"size:100"`
	Status    string         `json:"status" gorm:"size:50;default:'active'"`
	Items     []CartItem     `json:"items" gorm:"foreignKey:CartID"`
	ExpiresAt time.Time      `json:"expiresAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type CartItem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CartID      uint      `json:"cartId"`
	ProductID   uint      `json:"productId"`
	VariationID *uint     `json:"variationId"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Coupon struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Code         string    `json:"code" gorm:"uniqueIndex;size:50"`
	Type         string    `json:"type" gorm:"size:50"` // percentage, fixed
	Value        float64   `json:"value" gorm:"type:decimal(10,2)"`
	MinimumSpend float64   `json:"minimumSpend" gorm:"type:decimal(10,2)"`
	MaxDiscount  float64   `json:"maxDiscount" gorm:"type:decimal(10,2)"`
	StartDate    time.Time `json:"startDate"`
	EndDate      time.Time `json:"endDate"`
	UsageLimit   int       `json:"usageLimit"`
	UsageCount   int       `json:"usageCount"`
	IsActive     bool      `json:"isActive" gorm:"default:true"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
