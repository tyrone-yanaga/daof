package models

import (
	"time"
)

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OdooID      int64     `json:"odoo_id" gorm:"unique"` // Added to track Odoo reference
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price" gorm:"column:list_price"`    // Matches Odoo's list_price
	Stock       float64   `json:"stock" gorm:"column:qty_available"` // Changed to float64 to match Odoo
	SKU         string    `json:"sku" gorm:"column:default_code"`    // Matches Odoo's default_code
	Active      bool      `json:"active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName sets the table name for GORM
func (Product) TableName() string {
	return "products"
}

// FromOdooProduct converts Odoo product data to our Product model
func FromOdooProduct(odooProduct map[string]interface{}) Product {
	return Product{
		OdooID:      odooProduct["id"].(int64),
		Name:        odooProduct["name"].(string),
		Description: getStringOrEmpty(odooProduct["description"]),
		Price:       odooProduct["list_price"].(float64),
		Stock:       odooProduct["qty_available"].(float64),
		SKU:         getStringOrEmpty(odooProduct["default_code"]),
		Active:      true,
	}
}

// Helper function to handle potentially nil string values from Odoo
func getStringOrEmpty(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}
