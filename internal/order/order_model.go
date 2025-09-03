package order

import (
	"ecommerce/internal/auth"
	"ecommerce/internal/catalog"
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"not null;index"`

	// Order Information
	OrderNumber   string  `json:"order_number" gorm:"not null;uniqueIndex"`
	Status        string  `json:"status" gorm:"not null;default:'pending'"`
	PaymentStatus string  `json:"payment_status" gorm:"default:'pending'"`
	Total         float64 `json:"total" gorm:"not null"`

	// Customer Information
	ShippingAddress string `json:"shipping_address" gorm:"not null;default:''"`
	CustomerName    string `json:"customer_name" gorm:"not null;default:''"`  // ✅ Add default
	CustomerPhone   string `json:"customer_phone" gorm:"not null;default:''"` // ✅ Add default
	PaymentMethod   string `json:"payment_method" gorm:"not null;default:''"` // ✅ Add default

	// Additional Information
	Notes string `json:"notes" gorm:"default:''"`

	User          auth.User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items         []OrderItem    `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	PaymentProofs []PaymentProof `json:"payment_proofs,omitempty" gorm:"foreignKey:OrderID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
type CreateOrderRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required" validate:"max=500"`
	CustomerName    string `json:"customer_name" binding:"required" validate:"max=100"`
	CustomerPhone   string `json:"customer_phone" binding:"required" validate:"max=20"`
	PaymentMethod   string `json:"payment_method" binding:"required" validate:"oneof=bkash nagad rocket cod"`
	Notes           string `json:"notes" validate:"max=1000"`
}
type OrderItem struct {
	ID uint `json:"id" gorm:"primaryKey"`

	// Foreign Key Fields
	OrderID   uint `json:"order_id" gorm:"not null;index"`
	ProductID uint `json:"product_id" gorm:"not null;index"`

	// Product Information Snapshot (stored at time of order)
	ProductName     string  `json:"name" gorm:"not null"`
	Price    float64 `json:"price" gorm:"not null;default:0"`    // ✅ Add default
	Quantity int     `json:"quantity" gorm:"not null;default:1"` // ✅ Add default
	Subtotal float64 `json:"subtotal" gorm:"not null;default:0"` // ✅ Add default

	// Product Details Snapshot
	ProductImage string `json:"product_image" gorm:"default:''"` // ✅ Add default
	ProductSKU   string `json:"product_sku" gorm:"default:''"`

	// Relationships
	Product catalog.Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`

	// Timestamps
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type PaymentProof struct {
	ID      uint `json:"id" gorm:"primaryKey"`
	OrderID uint `json:"order_id" gorm:"not null;index"`

	// Payment Details
	TransactionID string  `json:"transaction_id" gorm:"not null"`
	PaymentMethod string  `json:"payment_method" gorm:"not null"`
	Amount        float64 `json:"amount" gorm:"not null"`
	Screenshot    string  `json:"screenshot" gorm:"not null"` // Image URL

	// Customer Payment Information
	SenderNumber string `json:"sender_number" gorm:"not null"`
	SenderName   string `json:"sender_name" gorm:"not null"`
	PaymentDate  string `json:"payment_date" gorm:"not null"`

	// Admin Review
	Status     string     `json:"status" gorm:"default:'pending'"` // pending, approved, rejected
	AdminNotes string     `json:"admin_notes"`
	ReviewedBy uint       `json:"reviewed_by"` // Admin user ID
	ReviewedAt *time.Time `json:"reviewed_at"`

	// Relationships
	Order     Order     `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubmitPaymentProofRequest struct {
	TransactionID string  `json:"transaction_id" binding:"required"`
	PaymentMethod string  `json:"payment_method" binding:"required"`
	Amount        float64 `json:"amount" binding:"required,min=0"`
	Screenshot    string  `json:"screenshot" binding:"required"` // Image URL
	SenderNumber  string  `json:"sender_number" binding:"required"`
	SenderName    string  `json:"sender_name" binding:"required"`
	PaymentDate   string  `json:"payment_date" binding:"required"`
}
