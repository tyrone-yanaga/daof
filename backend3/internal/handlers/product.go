// internal/handlers/product.go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "your-project/internal/models"
    "your-project/internal/services"
)

type ProductHandler struct {
    productService services.ProductService
}

func NewProductHandler(ps services.ProductService) *ProductHandler {
    return &ProductHandler{
        productService: ps,
    }
}

type ProductResponse struct {
    ID           uint              `json:"id"`
    Name         string           `json:"name"`
    Description  string           `json:"description"`
    Price        float64          `json:"price"`
    Images       []models.Image   `json:"images"`
    Variations   []models.Variation `json:"variations"`
    StockLevel   int              `json:"stockLevel"`
}

// List handles GET /products
func (h *ProductHandler) List(c *gin.Context) {
    // Parse pagination parameters
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    // Get filters from query parameters
    filters := models.ProductFilters{
        Category: c.Query("category"),
        MinPrice: c.Query("minPrice"),
        MaxPrice: c.Query("maxPrice"),
        Search:   c.Query("search"),
    }

    products, total, err := h.productService.List(c.Request.Context(), page, limit, filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": products,
        "meta": gin.H{
            "total":       total,
            "page":       page,
            "limit":      limit,
            "totalPages": (total + limit - 1) / limit,
        },
    })
}

// Get handles GET /products/:id
func (h *ProductHandler) Get(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    product, err := h.productService.Get(c.Request.Context(), uint(id))
    if err != nil {
        if err == services.ErrProductNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Get stock level from Odoo
    stockLevel, err := h.productService.GetStockLevel(c.Request.Context(), uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    response := ProductResponse{
        ID:          product.ID,
        Name:        product.Name,
        Description: product.Description,
        Price:       product.Price,
        Images:      product.Images,
        Variations:  product.Variations,
        StockLevel:  stockLevel,
    }

    c.JSON(http.StatusOK, response)
}

// Create handles POST /products
func (h *ProductHandler) Create(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    createdProduct, err := h.productService.Create(c.Request.Context(), &product)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, createdProduct)
}

// Update handles PUT /products/:id
func (h *ProductHandler) Update(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product.ID = uint(id)
    updatedProduct, err := h.productService.Update(c.Request.Context(), &product)
    if err != nil {
        if err == services.ErrProductNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updatedProduct)
}

// Delete handles DELETE /products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    err = h.productService.Delete(c.Request.Context(), uint(id))
    if err != nil {
        if err == services.ErrProductNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}