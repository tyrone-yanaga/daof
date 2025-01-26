package handlers

import (
	"ecommerce/internal/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.productService.GetProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	product, err := h.productService.GetProduct(id)
	if err != nil {
		log.Printf("Error fetching product ID %s: %v", id, err) //debug
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// In your handler file
func (h *ProductHandler) GetProductImage(c *gin.Context) {
	productID := c.Param("id")
	//TODO add productID validation
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	// Get image bytes from service
	imageBytes, err := h.productService.GetProductImage(productID)
	if err != nil {
		fmt.Printf("Image fetch error: %v\n", err) // Add debug logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product image"})
		return
	}

	fmt.Printf("Image bytes length: %d\n", len(imageBytes)) // Add debug logging
	// Set content type header
	c.Header("Content-Type", "image/jpeg")
	c.Header("Cache-Control", "public, max-age=3600")

	// Write image bytes to response
	c.Data(http.StatusOK, "image/jpeg", imageBytes)
}
