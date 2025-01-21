package models

import (
	"errors"
	"time"
)

type Cart struct {
	ID        string     `json:"id"`
	UserID    *uint      `json:"user_id,omitempty"` // Optional, for guest checkouts
	Items     []CartItem `json:"items"`
	Subtotal  float64    `json:"subtotal"`
	Total     float64    `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ExpiresAt time.Time  `json:"expires_at"`
}

type CartItem struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Subtotal  float64 `json:"subtotal"`
	Name      string  `json:"name"` // Denormalized from product
	SKU       string  `json:"sku"`  // Denormalized from product
}

// Calculate updates the cart totals
func (c *Cart) Calculate() {
	var subtotal float64
	for i := range c.Items {
		c.Items[i].Subtotal = c.Items[i].Price * float64(c.Items[i].Quantity)
		subtotal += c.Items[i].Subtotal
	}
	c.Subtotal = subtotal
	c.Total = subtotal // Add tax, shipping, etc. here
}

// AddItem adds or updates an item in the cart
func (c *Cart) AddItem(item CartItem) error {
	if item.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Check if item already exists
	for i := range c.Items {
		if c.Items[i].ProductID == item.ProductID {
			c.Items[i].Quantity += item.Quantity
			c.Calculate()
			return nil
		}
	}

	// Add new item
	c.Items = append(c.Items, item)
	c.Calculate()
	return nil
}

// UpdateItem updates the quantity of an item
func (c *Cart) UpdateItem(productID uint, quantity int) error {
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	for i := range c.Items {
		if c.Items[i].ProductID == productID {
			if quantity == 0 {
				// Remove item
				c.Items = append(c.Items[:i], c.Items[i+1:]...)
			} else {
				c.Items[i].Quantity = quantity
			}
			c.Calculate()
			return nil
		}
	}

	return errors.New("item not found in cart")
}
