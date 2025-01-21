package models

import (
	"time"
)

// Product represents the base product model
type Product struct {
	ID          uint               `json:"id" gorm:"primaryKey"`
	OdooID      int64              `json:"odoo_id" gorm:"unique"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	BasePrice   float64            `json:"base_price" gorm:"column:list_price"`
	SKU         string             `json:"sku" gorm:"column:default_code"`
	Active      bool               `json:"active" gorm:"default:true"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Variants    []ProductVariant   `json:"variants" gorm:"foreignKey:ProductID"`
	Attributes  []ProductAttribute `json:"attributes" gorm:"many2many:product_attributes"`
}

// ProductVariant represents a specific variant of a product
type ProductVariant struct {
	ID              uint                    `json:"id" gorm:"primaryKey"`
	OdooID          int64                   `json:"odoo_id" gorm:"unique"`
	ProductID       uint                    `json:"product_id"`
	Name            string                  `json:"name"`
	Price           float64                 `json:"price" gorm:"column:list_price"`
	Stock           float64                 `json:"stock" gorm:"column:qty_available"`
	SKU             string                  `json:"sku" gorm:"column:default_code"`
	Active          bool                    `json:"active" gorm:"default:true"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
	AttributeValues []ProductAttributeValue `json:"attribute_values" gorm:"many2many:variant_attribute_values"`
}

// ProductAttribute represents a type of variant attribute (e.g., "Size", "Color")
type ProductAttribute struct {
	ID        uint                    `json:"id" gorm:"primaryKey"`
	OdooID    int64                   `json:"odoo_id" gorm:"unique"`
	Name      string                  `json:"name"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	Values    []ProductAttributeValue `json:"values" gorm:"foreignKey:AttributeID"`
}

// ProductAttributeValue represents possible values for an attribute (e.g., "Small", "Red")
type ProductAttributeValue struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	OdooID      int64     `json:"odoo_id" gorm:"unique"`
	AttributeID uint      `json:"attribute_id"`
	Value       string    `json:"value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName sets the table name for GORM
func (Product) TableName() string {
	return "products"
}

func (ProductVariant) TableName() string {
	return "product_variants"
}

func (ProductAttribute) TableName() string {
	return "product_attributes"
}

func (ProductAttributeValue) TableName() string {
	return "product_attribute_values"
}

// FromOdooProduct converts Odoo product data to our Product model
func FromOdooProduct(odooProduct map[string]interface{}) Product {
	return Product{
		OdooID:      odooProduct["id"].(int64),
		Name:        odooProduct["name"].(string),
		Description: getStringOrEmpty(odooProduct["description"]),
		BasePrice:   odooProduct["list_price"].(float64),
		SKU:         getStringOrEmpty(odooProduct["default_code"]),
		Active:      true,
	}
}

// FromOdooProductVariant converts Odoo product variant data to our ProductVariant model
func FromOdooProductVariant(odooVariant map[string]interface{}, productID uint) ProductVariant {
	return ProductVariant{
		OdooID:    odooVariant["id"].(int64),
		ProductID: productID,
		Name:      odooVariant["name"].(string),
		Price:     odooVariant["list_price"].(float64),
		Stock:     odooVariant["qty_available"].(float64),
		SKU:       getStringOrEmpty(odooVariant["default_code"]),
		Active:    true,
	}
}

// FromOdooAttribute converts Odoo attribute data to our ProductAttribute model
func FromOdooAttribute(odooAttribute map[string]interface{}) ProductAttribute {
	return ProductAttribute{
		OdooID: odooAttribute["id"].(int64),
		Name:   odooAttribute["name"].(string),
	}
}

// FromOdooAttributeValue converts Odoo attribute value data to our ProductAttributeValue model
func FromOdooAttributeValue(odooValue map[string]interface{}, attributeID uint) ProductAttributeValue {
	return ProductAttributeValue{
		OdooID:      odooValue["id"].(int64),
		AttributeID: attributeID,
		Value:       odooValue["name"].(string),
	}
}

// Helper function to handle potentially nil string values from Odoo
func getStringOrEmpty(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}
