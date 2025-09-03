package catalog

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"not null"`
	Image       string  `json:"images" gorm:"column:images"`
	Description string  `json:"description"`
	SKU         string  `json:"sku" gorm:"uniqueIndex"`
	Price       float64 `json:"price" gorm:"not null"`
	Stock       int     `json:"stock" gorm:"default:0"`

	CategoryID uint `json:"category"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
                                                                                    
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
