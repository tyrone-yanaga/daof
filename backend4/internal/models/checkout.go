package models

import (
	"time"
)

type CheckoutSession struct {
	ID           string       `json:"id"`
	CartID       string       `json:"cart_id"`
	UserID       *uint        `json:"user_id,omitempty"`
	Status       string       `json:"status"` // "pending", "processing", "completed", "failed"
	PaymentID    string       `json:"payment_id,omitempty"`
	Total        float64      `json:"total"`
	Currency     string       `json:"currency"`
	PaymentData  PaymentData  `json:"paymentData"`
	ShippingInfo ShippingInfo `json:"shipping_info"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	ExpiresAt    time.Time    `json:"expires_at"`
}

type CheckoutRequest struct {
	CartID        string       `json:"cart_id" binding:"required"`
	ShippingInfo  ShippingInfo `json:"shipping_info" binding:"required"`
	PaymentMethod string       `json:"payment_method" binding:"required"`
	Currency      string       `json:"currency" binding:"required"`
}
