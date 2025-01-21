package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint               `json:"id" gorm:"primaryKey"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Price       float64            `json:"price"`
	Images      []ProductImage     `json:"images" gorm:"foreignKey:ProductID"`
	Variations  []ProductVariation `json:"variations" gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ProductID uint   `json:"product_id"`
	URL       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
}

type ProductVariation struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	ProductID uint   `json:"product_id"`
	Type      string `json:"type"` // e.g., "size", "color"
	Value     string `json:"value"`
	SKU       string `json:"sku"`
}

type OdooClient struct {
	client *resty.Client
	url    string
	db     string
	apiKey string
}

func NewOdooClient(url, db, apiKey string) *OdooClient {
	return &OdooClient{
		client: resty.New(),
		url:    url,
		db:     db,
		apiKey: apiKey,
	}
}

func (oc *OdooClient) GetInventory(productID string) (int, error) {
	// Implementation for Odoo inventory check
	return 0, nil
}

func setupRouter(db *gorm.DB, odooClient *OdooClient) *gin.Engine {
	r := gin.Default()

	// Product routes
	r.GET("/api/products", getProducts(db))
	r.GET("/api/products/:catergory", getProduct(db))
	r.GET("/api/products/:id", getProduct(db))

	// Cart routes
	r.POST("/api/cart", updateCart(db))
	r.GET("/api/cart", getCart(db))

	// Checkout routes
	r.POST("/api/checkout/init", initCheckout())

	return r
}

func main() {
	// Initialize DB connection
	// Initialize Odoo client
	// Setup router
	// Start server
}
