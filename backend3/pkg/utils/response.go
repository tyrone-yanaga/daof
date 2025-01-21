// internal/utils/response.go
package utils

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

type Response struct {
    Status  int         `json:"status"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Total      int64      `json:"total"`
    Page       int        `json:"page"`
    PageSize   int        `json:"pageSize"`
    TotalPages int        `json:"totalPages"`
}

func SuccessResponse(c *gin.Context, status int, message string, data interface{}) {
    c.JSON(status, Response{
        Status:  status,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c *gin.Context, status int, err error) {
    c.JSON(status, Response{
        Status: status,
        Error:  err.Error(),
    })
}

func PaginatedSuccessResponse(c *gin.Context, data interface{}, total int64, page, pageSize int) {
    totalPages := (int(total) + pageSize - 1) / pageSize
    
    c.JSON(http.StatusOK, PaginatedResponse{
        Data:       data,
        Total:      total,
        Page:       page,
        PageSize:   pageSize,
        TotalPages: totalPages,
    })
}

func ValidationErrorResponse(c *gin.Context, err *ValidationError) {
    c.JSON(http.StatusBadRequest, Response{
        Status: http.StatusBadRequest,
        Error:  fmt.Sprintf("%s: %s", err.Field, err.Message),
    })
}