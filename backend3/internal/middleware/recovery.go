// internal/middleware/recovery.go
package middleware

import (
    "fmt"
    "net/http"
    "runtime/debug"

    "github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                stack := debug.Stack()
                
                // Log the error and stack trace
                fmt.Printf("Recovery from panic: %v\nStack: %s\n", err, stack)

                c.JSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal server error",
                })
                c.Abort()
            }
        }()
        
        c.Next()
    }
}
