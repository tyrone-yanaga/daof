// internal/router/router.go
package router

import (
    "github.com/gin-gonic/gin"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    
    "your-project/internal/config"
    "your-project/internal/handlers"
    "your-project/internal/middleware"
    "your-project/internal/services"
)

type Router struct {
    config  *config.Config
    db      *gorm.DB
    redis   *redis.Client
}

func NewRouter(config *config.Config, db *gorm.DB, redis *redis.Client) *Router {
    return &Router{
        config: config,
        db:     db,
        redis:  redis,
    }
}

func (r *Router) Setup() *gin.Engine {
    router := gin.Default()

    // Middleware
    router.Use(middleware.CORS(r.config.CORS))
    router.Use(middleware.Logger())
    router.Use(middleware.RateLimit(r.redis))

    // Initialize services
    productService := services.NewProductService(r.db)
    cartService := services.NewCartService(r.db, r.redis)
    checkoutService := services.NewCheckoutService(r.db, r.redis, r.config)
    odooService := services.NewOdooService(r.config.Odoo)
    
    // Initialize handlers
    productHandler := handlers.NewProductHandler(productService)
    cartHandler := handlers.NewCartHandler(cartService)
    checkoutHandler := handlers.NewCheckoutHandler(checkoutService, cartService, r.config)

    // API routes
    api := router.Group("/api")
    {
        // Product routes
        products := api.Group("/products")
        {
            products.GET("", productHandler.List)
            products.GET("/:id", productHandler.Get)
            products.GET("/category/:category", productHandler.ListByCategory)
            
            // Protected routes
            protected := products.Use(middleware.Auth())
            {
                protected.POST("", productHandler.Create)
                protected.PUT("/:id", productHandler.Update)
                protected.DELETE("/:id", productHandler.Delete)
            }
        }

        // Cart routes
        cart := api.Group("/cart")
        {
            cart.GET("", cartHandler.Get)
            cart.POST("/items", cartHandler.AddItem)
            cart.PUT("/items/:id", cartHandler.UpdateItem)
            cart.DELETE("/items/:id", cartHandler.RemoveItem)
            cart.DELETE("", cartHandler.Clear)
        }

        // Checkout routes
        checkout := api.Group("/checkout")
        {
            checkout.POST("/init", checkoutHandler.InitCheckout)
            checkout.POST("/complete", checkoutHandler.CompleteCheckout)
            checkout.POST("/webhook", checkoutHandler.WebhookHandler)
        }

        // Order routes
        orders := api.Group("/orders").Use(middleware.Auth())
        {
            orders.GET("", orderHandler.List)
            orders.GET("/:id", orderHandler.Get)
            orders.PUT("/:id/cancel", orderHandler.Cancel)
        }

        // Inventory sync routes (admin only)
        inventory := api.Group("/inventory").Use(middleware.AdminAuth())
        {
            inventory.POST("/sync", inventoryHandler.SyncWithOdoo)
            inventory.GET("/status", inventoryHandler.GetSyncStatus)
        }
    }

    // Health check
    router.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    return router
}