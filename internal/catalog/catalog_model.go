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

	CategoryID uint `json:"category_id"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Category struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
                                                                                    
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
	SubSubCategories []SubSubCategory `json:"sub_subcategories,omitempty" gorm:"foreignKey:SubCategoryID"`



	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type SubCategory struct {
	 ID  uint `json:"id" gorm:"primaryKey"`
	 Name string `json:"name" gorm:"not null"`
	 CategoryID  uint `json:"category_id" gorm:"not null"`
     SubSubCategories []SubSubCategory `json:"sub_subcategories,omitempty" gorm:"foreignKey:SubCategoryID"`
	 CreatedAt time.Time `json:"created_at"`
	 UpdatedAt time.Time `json:"updated_at"`
	 DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type SubSubCategory struct {
	 ID  uint `json:"id" gorm:"primaryKey"`
	 Name string `json:"name" gorm:"not null"`
	 SubCategoryID uint `json:"sub_category_id" gorm:"not null"`
	 ProductCount int `json:"product_count" gorm:"default:0"`

	 CreatedAt time.Time      `json:"created_at"`
     UpdatedAt time.Time      `json:"updated_at"`
     DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
