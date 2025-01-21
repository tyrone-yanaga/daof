// internal/models/order.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Order struct {
    ID              uint           `json:"id" gorm:"primaryKey"`
    OrderNumber     string         `json:"orderNumber" gorm:"uniqueIndex;size:50"`
    UserID          uint           `json:"userId"`
    User            User           `json:"user" gorm:"foreignKey:UserID"`
    
    // Amounts
    Subtotal        float64        `json:"subtotal" gorm:"type:decimal(10,2);not null"`
    Tax            float64        `json:"tax" gorm:"type:decimal(10,2);not null"`
    ShippingCost   float64        `json:"shippingCost" gorm:"type:decimal(10,2);not null"`
    Discount       float64        `json:"discount" gorm:"type:decimal(10,2);default:0"`
    Total          float64        `json:"total" gorm:"type:decimal(10,2);not null"`
    
    // Relations
    Items          []OrderItem    `json:"items" gorm:"foreignKey:OrderID"`
    Transactions   []Transaction  `json:"transactions" gorm:"foreignKey:OrderID"`
    
    // Addresses
    ShippingAddress Address        `json:"shippingAddress" gorm:"embedded;embeddedPrefix:shipping_"`
    BillingAddress  Address        `json:"billingAddress" gorm:"embedded;embeddedPrefix:billing_"`
    
    // Status
    Status          string         `json:"status" gorm:"size:50;default:'pending'"`
    PaymentStatus   string         `json:"paymentStatus" gorm:"size:50;default:'pending'"`
    FulfillmentStatus string       `json:"fulfillmentStatus" gorm:"size:50;default:'pending'"`
    
    // Tracking
    TrackingNumber string         `json:"trackingNumber" gorm:"size:100"`
    TrackingURL    string         `json:"trackingUrl" gorm:"size:255"`
    
    // Customer Info
    CustomerEmail  string         `json:"customerEmail" gorm:"size:255"`
    CustomerNotes  string         `json:"customerNotes" gorm:"type:text"`
    
    // Timestamps
    CreatedAt      time.Time      `json:"createdAt"`
    UpdatedAt      time.Time      `json:"updatedAt"`
    DeletedAt      gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type OrderItem struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    OrderID     uint      `json:"orderId"`
    ProductID   uint      `json:"productId"`
    VariationID *uint     `json:"variationId"`
    SKU         string    `json:"sku" gorm:"size:50"`
    Name        string    `json:"name" gorm:"size:255"`
    Quantity    int       `json:"quantity" gorm:"not null"`
    Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
    Cost        float64   `json:"cost" gorm:"type:decimal(10,2)"`
    Tax         float64   `json:"tax" gorm:"type:decimal(10,2)"`
    Discount    float64   `json:"discount" gorm:"type:decimal(10,2);default:0"`
    Total       float64   `json:"total" gorm:"type:decimal(10,2);not null"`
    CreatedAt   time.Time `json:"createdAt"`
}

type Transaction struct {
    ID            uint      `json:"id" gorm:"primaryKey"`
    OrderID       uint      `json:"orderId"`
    PaymentMethod string    `json:"paymentMethod" gorm:"size:50"`
    Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
    Currency      string    `json:"currency" gorm:"size:3;default:'USD'"`
    Status        string    `json:"status" gorm:"size:50"`
    PaymentID     string    `json:"paymentId" gorm:"size:100"` // External payment reference
    Error         string    `json:"error" gorm:"type:text"`
    CreatedAt     time.Time `json:"createdAt"`
}

type Address struct {
    FirstName   string `json:"firstName" gorm:"size:100"`
    LastName    string `json:"lastName" gorm:"size:100"`
    Company     string `json:"company" gorm:"size:100"`
    Address1    string `json:"address1" gorm:"size:255"`
    Address2    string `json:"address2" gorm:"size:255"`
    City        string `json:"city" gorm:"size:100"`
    State       string `json:"state" gorm:"size:100"`
    PostalCode  string `json:"postalCode" gorm:"size:20"`
    Country     string `json:"country" gorm:"size:2"`
    Phone       string `json:"phone" gorm:"size:50"`
}