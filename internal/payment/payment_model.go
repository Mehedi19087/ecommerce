package payment

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	OrderID       uint    `json:"order_id" gorm:"not null;uniqueIndex"` // One payment per order
	PaymentMethod string  `json:"payment_method" gorm:"not null"`       // "cash", "bkash", "nagad"
	Status        string  `json:"status" gorm:"not null"`               // "pending", "initiated", "completed", "failed", "expired"
	Amount        float64 `json:"amount" gorm:"not null"`
	Currency      string  `json:"currency" gorm:"default:'BDT'"`

	// External payment provider data (bKash, Nagad, etc.)
	ExternalPaymentID string `json:"external_payment_id,omitempty"` // bKash payment ID
	TransactionID     string `json:"transaction_id,omitempty"`      // bKash transaction ID (after completion)
	PaymentURL        string `json:"payment_url,omitempty"`    
	Screenshot string `json:"screenshot,omitempty"`    

	// Timestamps for payment lifecycle
	InitiatedAt *time.Time `json:"initiated_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Error handling & debugging
	ProviderResponse string `json:"provider_response,omitempty" gorm:"type:text"` // Full bKash response
	FailureReason    string `json:"failure_reason,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
