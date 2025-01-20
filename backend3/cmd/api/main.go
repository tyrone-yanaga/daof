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

    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"

    "your-project/internal/config"
    "your-project/internal/router"
    "your-project/internal/database"
)

func main() {
    // Load configuration
    cfg := config.New()

    // Initialize database
    db, err := initDatabase(cfg)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Initialize Redis
    rdb := initRedis(cfg)

    // Setup router
    r := router.NewRouter(cfg, db, rdb)
    handler := r.Setup()

    // Configure server
    srv := &http.Server{
        Addr:    ":" + cfg.Server.Port,
        Handler: handler,
    }

    // Start server in goroutine
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")

    // Graceful shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }

    log.Println("Server exiting")
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Run migrations
    if err := database.RunMigrations(db); err != nil {
        return nil, err
    }

    return db, nil
}

func initRedis(cfg *config.Config) *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })

    return rdb
}