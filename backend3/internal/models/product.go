// internal/models/product.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID          uint           `json:"id" gorm:"primaryKey"`
    SKU         string         `json:"sku" gorm:"uniqueIndex;size:50"`
    Name        string         `json:"name" gorm:"size:255;not null"`
    Description string         `json:"description" gorm:"type:text"`
    Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
    Cost        float64        `json:"cost" gorm:"type:decimal(10,2)"`
    Weight      float64        `json:"weight" gorm:"type:decimal(8,2)"`
    
    // Stock Information
    StockLevel  int            `json:"stockLevel" gorm:"default:0"`
    MinStock    int            `json:"minStock" gorm:"default:5"`
    MaxStock    int            `json:"maxStock" gorm:"default:100"`
    
    // Relationships
    CategoryID  uint           `json:"categoryId"`
    Category    Category       `json:"category" gorm:"foreignKey:CategoryID"`
    Images      []ProductImage `json:"images" gorm:"foreignKey:ProductID"`
    Variations  []Variation    `json:"variations" gorm:"foreignKey:ProductID"`
    
    // SEO
    MetaTitle       string `json:"metaTitle" gorm:"size:255"`
    MetaDescription string `json:"metaDescription" gorm:"type:text"`
    Slug            string `json:"slug" gorm:"uniqueIndex;size:255"`
    
    // Status
    Status     string         `json:"status" gorm:"default:'active'"`
    Featured   bool           `json:"featured" gorm:"default:false"`
    
    // Timestamps
    CreatedAt  time.Time      `json:"createdAt"`
    UpdatedAt  time.Time      `json:"updatedAt"`
    DeletedAt  gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

type ProductImage struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    ProductID uint      `json:"productId"`
    URL       string    `json:"url" gorm:"size:255;not null"`
    Alt       string    `json:"alt" gorm:"size:255"`
    Position  int       `json:"position" gorm:"default:0"`
    IsPrimary bool      `json:"isPrimary" gorm:"default:false"`
    CreatedAt time.Time `json:"createdAt"`
}

type Variation struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    ProductID uint      `json:"productId"`
    SKU       string    `json:"sku" gorm:"uniqueIndex;size:50"`
    Type      string    `json:"type" gorm:"size:50;not null"` // e.g., 'size', 'color'
    Value     string    `json:"value" gorm:"size:50;not null"` // e.g., 'XL', 'Red'
    Price     float64   `json:"price" gorm:"type:decimal(10,2)"`
    StockLevel int      `json:"stockLevel" gorm:"default:0"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

type Category struct {
    ID          uint       `json:"id" gorm:"primaryKey"`
    Name        string     `json:"name" gorm:"size:255;not null"`
    Description string     `json:"description" gorm:"type:text"`
    Slug        string     `json:"slug" gorm:"uniqueIndex;size:255"`
    ParentID    *uint      `json:"parentId"`
    Level       int        `json:"level" gorm:"default:0"`
    Order       int        `json:"order" gorm:"default:0"`
    IsActive    bool       `json:"isActive" gorm:"default:true"`
    CreatedAt   time.Time  `json:"createdAt"`
    UpdatedAt   time.Time  `json:"updatedAt"`
}