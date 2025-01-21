package routes

import (
	"ecommerce/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, handlers *handlers.Handlers) {
	api := r.Group("/api")
	{
		// Product routes
		api.GET("/products", handlers.Product.GetProducts)
		api.GET("/products/:id", handlers.Product.GetProduct)

		// Order routes
		api.POST("/orders", handlers.Order.CreateOrder)
		api.GET("/orders/:id", handlers.Order.GetOrder)

		// Checkout routes
		api.POST("/checkout", handlers.Checkout.InitiateCheckout)
		api.POST("/checkout/:id/complete", handlers.Checkout.CompleteCheckout)

		//Cart routes
		api.POST("/carts", handlers.Cart.CreateCart)
		api.GET("/carts/:id", handlers.Cart.GetCart)
		api.POST("/carts/:id/items", handlers.Cart.AddToCart)
		api.PUT("/carts/:id/items", handlers.Cart.UpdateCartItem)
	}
}
