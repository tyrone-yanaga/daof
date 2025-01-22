package handlers

import (
	"ecommerce/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
	}
}

type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	VariantID uint `json:"variant_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,gt=0"`
}

func (h *CartHandler) CreateCart(c *gin.Context) {
	// Get user ID from context if authenticated
	var userID *uint
	if user, exists := c.Get("user"); exists {
		id := user.(uint)
		userID = &id
	}

	cart, err := h.cartService.CreateCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cart)
}

func (h *CartHandler) GetCart(c *gin.Context) {
	cartID := c.Param("id")

	cart, err := h.cartService.GetCart(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	cartID := c.Param("id")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.AddToCart(c.Request.Context(), cartID, req.ProductID, req.VariantID, req.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get updated cart
	cart, err := h.cartService.GetCart(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}

func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	cartID := c.Param("id")

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.UpdateCartItem(c.Request.Context(), cartID, req.ProductID, req.VariantID, req.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get updated cart
	cart, err := h.cartService.GetCart(c.Request.Context(), cartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cart)
}
