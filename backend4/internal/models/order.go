package models

import "time"

type Order struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	UserID       uint         `json:"user_id"`
	Status       string       `json:"status"`
	Total        float64      `json:"total"`
	PaymentID    string       `json:"payment_id"`
	Items        []OrderItem  `json:"items"`
	ShippingInfo ShippingInfo `json:"shipping_info"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type ShippingInfo struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	OrderID  uint   `json:"order_id"`
	Address  string `json:"address"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
	PostCode string `json:"post_code"`
}
