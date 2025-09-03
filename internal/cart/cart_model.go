package cart

import (
	"ecommerce/internal/catalog"
	"time"
)

type Cart struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `json:"user_id"`
	
	Items     []CartItem `gorm:"foreignKey:CartID" json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (c *Cart) CalculateTotal() float64 {
	var total float64

	for _, item := range c.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

type CartItem struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CartID    uint            `json:"cart_id"`
	ProductID uint            `json:"product_id"`
	Product   catalog.Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int             `json:"quantity"`
	Price     float64         `json:"price"` // Price at time of adding
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
