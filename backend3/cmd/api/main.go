// cmd/api/main.go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"

    "your-project/internal/config"
    "your-project/internal/handlers"
    "your-project/internal/middleware"
    "your-project/internal/repository/postgres"
    "your-project/internal/services/odoo"
)

func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: .env file not found")
    }

    // Initialize configuration
    cfg := config.New()

    // Initialize database
    db, err := initDB(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Initialize repositories
    productRepo := postgres.NewProductRepository(db)
    cartRepo := postgres.NewCartRepository(db)

    // Initialize Odoo client
    odooClient := odoo.NewClient(cfg.OdooConfig)

    // Initialize services
    productService := services.NewProductService(productRepo, odooClient)

    // Initialize handlers
    productHandler := handlers.NewProductHandler(productService)
    cartHandler := handlers.NewCartHandler(cartRepo)

    // Setup router
    router := setupRouter(productHandler, cartHandler)

    // Create server
    srv := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: router,
    }

    // Graceful shutdown
    go func() {
        // Service connections
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}

func setupRouter(productHandler *handlers.ProductHandler, cartHandler *handlers.CartHandler) *gin.Engine {
    router := gin.Default()

    // Middleware
    router.Use(middleware.CORS())
    router.Use(middleware.Logger())

    // Routes
    api := router.Group("/api")
    {
        products := api.Group("/products")
        {
            products.GET("", productHandler.List)
            products.GET("/:id", productHandler.Get)
            products.POST("", middleware.Auth(), productHandler.Create)
            products.PUT("/:id", middleware.Auth(), productHandler.Update)
            products.DELETE("/:id", middleware.Auth(), productHandler.Delete)
        }

        cart := api.Group("/cart")
        {
            cart.GET("", cartHandler.Get)
            cart.POST("/items", cartHandler.AddItem)
            cart.PUT("/items/:id", cartHandler.UpdateItem)
            cart.DELETE("/items/:id", cartHandler.RemoveItem)
        }
    }

    return router
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Auto migrate models
    if err := db.AutoMigrate(&models.Product{}, &models.Cart{}, &models.CartItem{}); err != nil {
        return nil, err
    }

    return db, nil
}