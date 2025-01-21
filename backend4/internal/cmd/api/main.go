package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecommerce/internal/api/handlers"
	"ecommerce/internal/api/routes"
	"ecommerce/internal/config"
	"ecommerce/internal/database"
	"ecommerce/internal/scheduler"
	"ecommerce/internal/services"
	"ecommerce/internal/sync"
	"ecommerce/pkg/adyen"
	"ecommerce/pkg/odoo"
	"ecommerce/pkg/queue"
	"ecommerce/pkg/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize Redis client
	redisClient, err := redis.NewClient(redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize RabbitMQ client
	queueClient, err := queue.NewClient(queue.Config{
		URL: cfg.RabbitMQ.URL,
	})
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer queueClient.Close()

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Odoo client
	odooClient, err := odoo.NewClient(odoo.Config{
		URL:      cfg.Odoo.URL,
		Database: cfg.Odoo.Database,
		Username: cfg.Odoo.Username,
		Password: cfg.Odoo.Password,
	})
	if err != nil {
		log.Fatalf("Failed to connect to Odoo: %v", err)
	}

	// Initialize Adyen client
	adyenClient, err := adyen.NewClient(&adyen.Config{
		ApiKey:      cfg.Adyen.ApiKey,
		Environment: cfg.Adyen.Environment,
	})

	// Initialize services
	productService := services.NewProductService(odooClient)
	cartService := services.NewCartService(redisClient, productService)
	orderService := services.NewOrderService(queueClient, odooClient, db)
	checkoutService := services.NewCheckoutService(
		cartService,
		orderService,
		redisClient,
		odooClient,
		queueClient,
		adyenClient,
		cfg.Server.BaseURL,
	)

	// Initialize handlers
	handlers := &handlers.Handlers{
		Product:  *handlers.NewProductHandler(productService),
		Cart:     *handlers.NewCartHandler(cartService),
		Checkout: *handlers.NewCheckoutHandler(checkoutService),
		Order:    *handlers.NewOrderHandler(orderService),
	}

	// Initialize sync service
	odooSync := sync.NewOdooSync(db, odooClient)

	// Initialize and start scheduler
	syncScheduler := scheduler.NewSyncScheduler(odooSync)
	syncScheduler.Start()
	defer syncScheduler.Stop()

	// Initialize Gin router
	r := gin.Default()

	// Configure CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"} // Add your frontend URL
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	r.Use(cors.New(corsConfig))

	// Setup routes
	routes.SetupRoutes(r, handlers)

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
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
