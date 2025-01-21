// internal/handlers/cart.go
package handlers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "your-project/internal/models"
    "your-project/internal/services"
)

type CartHandler struct {
    cartService services.CartService
}

func NewCartHandler(cs services.CartService) *CartHandler {
    return &CartHandler{
        cartService: cs,
    }
}

// Get handles GET /cart
func (h *CartHandler) Get(c *gin.Context) {
    // Get cart ID from session/cookie
    cartID := getCartIDFromSession(c)
    if cartID == 0 {
        c.JSON(http.StatusOK, gin.H{"items": []models.CartItem{}})
        return
    }

    cart, err := h.cartService.Get(c.Request.Context(), cartID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, cart)
}

// AddItem handles POST /cart/items
func (h *CartHandler) AddItem(c *gin.Context) {
    var item models.CartItem
    if err := c.ShouldBindJSON(&item); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get or create cart ID
    cartID := getCartIDFromSession(c)
    if cartID == 0 {
        newCart, err := h.cartService.Create(c.Request.Context())
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        cartID = newCart.ID
        setCartIDInSession(c, cartID)
    }

    // Check stock availability
    available, err := h.cartService.CheckStockAvailability(c.Request.Context(), item.ProductID, item.Quantity)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if !available {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
        return
    }

    item.CartID = cartID
    addedItem, err := h.cartService.AddItem(c.Request.Context(), &item)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, addedItem)
}

// UpdateItem handles PUT /cart/items/:id
func (h *CartHandler) UpdateItem(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
        return
    }

    var update struct {
        Quantity int `json:"quantity" binding:"required,min=1"`
    }
    if err := c.ShouldBindJSON(&update); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    cartID := getCartIDFromSession(c)
    if cartID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
        return
    }

    updatedItem, err := h.cartService.UpdateItem(c.Request.Context(), cartID, uint(id), update.Quantity)
    if err != nil {
        if err == services.ErrItemNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, updatedItem)
}

// RemoveItem handles DELETE /cart/items/:id
func (h *CartHandler) RemoveItem(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
        return
    }

    cartID := getCartIDFromSession(c)
    if cartID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
        return
    }

    err = h.cartService.RemoveItem(c.Request.Context(), cartID, uint(id))
    if err != nil {
        if err == services.ErrItemNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Item removed successfully"})
}

// Helper functions for session management
func getCartIDFromSession(c *gin.Context) uint {
    cartID, _ := c.Cookie("cart_id")
    id, _ := strconv.ParseUint(cartID, 10, 64)
    return uint(id)
}

func setCartIDInSession(c *gin.Context, cartID uint) {
    c.SetCookie("cart_id", strconv.FormatUint(uint64(cartID), 10), 86400*30, "/", "", false, true)
}