// internal/models/user.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID            uint           `json:"id" gorm:"primaryKey"`
    Email         string         `json:"email" gorm:"uniqueIndex;size:255"`
    Password      string         `json:"-" gorm:"size:255"`
    FirstName     string         `json:"firstName" gorm:"size:100"`
    LastName      string         `json:"lastName" gorm:"size:100"`
    Role          string         `json:"role" gorm:"size:50;default:'customer'"`
    Active        bool           `json:"active" gorm:"default:true"`
    EmailVerified bool           `json:"emailVerified" gorm:"default:false"`
    
    // Relations
    Addresses     []Address      `json:"addresses" gorm:"foreignKey:UserID"`
    Orders        []Order        `json:"orders" gorm:"foreignKey:UserID"`
    
    // Authentication
    LastLogin     *time.Time     `json:"lastLogin"`
    ResetToken    string         `json:"-" gorm:"size:255"`
    ResetExpires  *time.Time     `json:"-"`
    
    // Timestamps
    CreatedAt     time.Time      `json:"createdAt"`
    UpdatedAt     time.Time      `json:"updatedAt"`
    DeletedAt     gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type UserAddress struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    UserID    uint      `json:"userId"`
    Address   Address   `json:"address" gorm:"embedded"`
    IsDefault bool      `json:"isDefault" gorm:"default:false"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

type Session struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    UserID       uint      `json:"userId"`
    Token        string    `json:"token" gorm:"uniqueIndex;size:255"`
    RefreshToken string    `json:"refreshToken" gorm:"uniqueIndex;size:255"`
    UserAgent    string    `json:"userAgent" gorm:"size:255"`
    IP           string    `json:"ip" gorm:"size:45"`
    ExpiresAt    time.Time `json:"expiresAt"`
    CreatedAt    time.Time `json:"createdAt"`
    UpdatedAt    time.Time `json:"updatedAt"`
}