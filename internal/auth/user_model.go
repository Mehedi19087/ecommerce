package auth

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	GoogleID string `gorm:"not null" json:"google_id"`
	Email    string `gorm:"not null;uniqueIndex" json:"email"`

	Phone    string `json:"phone"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Address struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	City    string `json:"city"`
	Zone    string `json:"zone"`
	Label   string `json:"label"` // "Home", "Office"

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Birthday string `json:"birthday"` // Format: "1990-01-15"
	Gender   string `json:"gender"`
}

type CreateAddressRequest struct {
	Name    string `json:"name" binding:"required"`
	Phone   string `json:"phone" binding:"required"`
	Address string `json:"address" binding:"required"`
	City    string `json:"city" binding:"required"`
	Zone    string `json:"zone" binding:"required"`
	Label   string `json:"label" binding:"required"`
}
type VisitorLogs struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	IP        string    `json:"ip"`
	Country   string    `json:"country"`
	Region    string    `json:"region"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

type GoogleUserInfo struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
