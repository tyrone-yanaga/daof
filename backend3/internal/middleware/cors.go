// internal/middleware/cors.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "your-project/internal/config"
)

func CORS(corsConfig config.CORSConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", strings.Join(corsConfig.AllowedOrigins, ","))
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowedHeaders, ","))
        c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowedMethods, ","))

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
