package payment

import (
	"gorm.io/gorm"
	"time"
)

type PaymentRepository interface {
	Create(payment *Payment) error
	FindById(id uint) (*Payment, error)
	FindByOrderId(id uint) (*Payment, error)
	Update(payment *Payment) error

	FindByExternalPaymentID(externalID string) (*Payment, error)
	UpdateStatus(paymentID uint, status string) error
	UpdateExternalPaymentID(paymentID uint, externalID string) error
	UpdateTransactionID(paymentID uint, transactionID string) error

	FindPendingPayments() ([]Payment, error)
	FindByUserID(userID uint) ([]Payment, error) // ✅ Fixed: []Payment not []*Payment

	GetDB() *gorm.DB
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) Create(payment *Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindById(id uint) (*Payment, error) {
	var payment Payment
	err := r.db.First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByOrderId(id uint) (*Payment, error) {
	var payment Payment
	err := r.db.Where("order_id = ?", id).First(&payment).Error // ✅ Fixed: added .Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(payment *Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) FindByExternalPaymentID(externalID string) (*Payment, error) {
	var payment Payment
	err := r.db.Where("external_payment_id = ?", externalID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdateStatus(paymentID uint, status string) error {
	return r.db.Model(&Payment{}).Where("id = ?", paymentID).Update("status", status).Error // ✅ Simplified
}

// ✅ Fixed: Multiple field update using map
func (r *paymentRepository) UpdateExternalPaymentID(paymentID uint, externalID string) error {
	updates := map[string]interface{}{
		"external_payment_id": externalID,
		"status":              "initiated",
		"initiated_at":        time.Now(),
	}
	return r.db.Model(&Payment{}).Where("id = ?", paymentID).Updates(updates).Error
}

// ✅ Fixed: Multiple field update using map
func (r *paymentRepository) UpdateTransactionID(paymentID uint, transactionID string) error {
	updates := map[string]interface{}{
		"transaction_id": transactionID,
		"status":         "completed",
		"completed_at":   time.Now(),
	}
	return r.db.Model(&Payment{}).Where("id = ?", paymentID).Updates(updates).Error
}

// ✅ Fixed: Complete implementation
func (r *paymentRepository) FindPendingPayments() ([]Payment, error) {
	var payments []Payment
	err := r.db.Where("status IN ?", []string{"pending", "initiated"}).Find(&payments).Error
	return payments, err
}

// ✅ Fixed: Correct implementation with JOIN
func (r *paymentRepository) FindByUserID(userID uint) ([]Payment, error) {
	var payments []Payment // ✅ Fixed: plural name and correct type
	err := r.db.Table("payments").
		Joins("JOIN orders ON payments.order_id = orders.id"). // ✅ Using foreign key relationship
		Where("orders.user_id = ?", userID).
		Order("payments.created_at DESC").
		Find(&payments).Error // ✅ Fixed: added .Error
	if err != nil {
		return nil, err
	}
	return payments, nil // ✅ Fixed: return slice directly
}

func (r *paymentRepository) GetDB() *gorm.DB {
	return r.db
}
