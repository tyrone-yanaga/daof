// internal/middleware/ratelimit.go
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"your-project/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func RateLimit(redis *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf("ratelimit:%s", c.ClientIP())

		// Get current count
		ctx := context.Background()
		count, err := redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit check failed"})
			c.Abort()
			return
		}

		config := config.New()
		limit := config.Server.RateLimit.Requests

		if count >= limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": time.Until(time.Now().Add(time.Minute)),
			})
			c.Abort()
			return
		}

		// Increment counter
		pipe := redis.Pipeline()
		pipe.Incr(ctx, key)
		if count == 0 {
			pipe.Expire(ctx, key, time.Minute)
		}
		_, err = pipe.Exec(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit update failed"})
			c.Abort()
			return
		}

		c.Next()
	}
}
