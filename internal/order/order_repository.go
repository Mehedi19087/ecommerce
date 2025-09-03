package order

import (
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *Order) error
	CreateOrderItem(item *OrderItem) error
	GetByUserID(userID uint) ([]Order, error)
	GetByID(orderID uint, userID uint) (*Order, error)
	UpdateStatus(orderID uint, userID uint, status string) error
	CreatePaymentProof(proof *PaymentProof) error
	GetPaymentProofByOrderID(orderID uint, userID uint) (*PaymentProof, error)
	UpdatePaymentProof(orderID uint, userID uint, proofData SubmitPaymentProofRequest) error

	GetAllOrders() ([]Order, error)
	UpdateOrderStatusAdmin(orderID uint, status string) error
	GetPaymentProofByID(proofID uint) (*PaymentProof, error)
	ReviewPaymentProof(proofID uint, status string, adminNotes string, reviewerID uint) error
	UpdateOrderPaymentStatus(orderID uint, paymentStatus string) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) CreateOrderItem(item *OrderItem) error {
	return r.db.Create(item).Error
}

func (r *orderRepository) GetByUserID(userID uint) ([]Order, error) {
	var orders []Order
	err := r.db.Where("user_id = ?", userID).
		Preload("Items").
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetByID(orderID uint, userID uint) (*Order, error) {
	var order Order
	err := r.db.Where("id = ? AND user_id = ?", orderID, userID).
		Preload("Items").
		First(&order).Error
	return &order, err
}
func (r *orderRepository) UpdateStatus(orderID uint, userID uint, status string) error {
	return r.db.Model(&Order{}).
		Where("id = ? AND user_id = ?", orderID, userID).
		Update("status", status).Error
}

func (r *orderRepository) CreatePaymentProof(proof *PaymentProof) error {
	return r.db.Create(proof).Error
}

func (r *orderRepository) GetPaymentProofByOrderID(orderID uint, userID uint) (*PaymentProof, error) {
	// First verify order belongs to user
	var order Order
	err := r.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if err != nil {
		return nil, err
	}

	// Get payment proof for this order
	var proof PaymentProof
	err = r.db.Where("order_id = ?", orderID).First(&proof).Error
	return &proof, err
}

func (r *orderRepository) UpdatePaymentProof(orderID uint, userID uint, proofData SubmitPaymentProofRequest) error {
	// First verify order belongs to user
	var order Order
	err := r.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error
	if err != nil {
		return err
	}

	// Update payment proof
	updates := map[string]interface{}{
		"transaction_id": proofData.TransactionID,
		"payment_method": proofData.PaymentMethod,
		"amount":         proofData.Amount,
		"screenshot":     proofData.Screenshot,
		"sender_number":  proofData.SenderNumber,
		"sender_name":    proofData.SenderName,
		"payment_date":   proofData.PaymentDate,
	}

	return r.db.Model(&PaymentProof{}).Where("order_id = ?", orderID).Updates(updates).Error
}

func (r *orderRepository) GetAllOrders() ([]Order, error) {
	var orders []Order
	err := r.db.Preload("Items").
		Preload("PaymentProofs").
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *orderRepository) UpdateOrderStatusAdmin(orderID uint, status string) error {
	return r.db.Model(&Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

func (r *orderRepository) GetPaymentProofByID(proofID uint) (*PaymentProof, error) {
	var proof PaymentProof
	err := r.db.Preload("Order").First(&proof, proofID).Error
	return &proof, err
}

func (r *orderRepository) ReviewPaymentProof(proofID uint, status string, adminNotes string, reviewerID uint) error {
	updates := map[string]interface{}{
		"status":      status,
		"admin_notes": adminNotes,
		"reviewed_by": reviewerID,
		"reviewed_at": r.db.NowFunc(),
	}

	return r.db.Model(&PaymentProof{}).
		Where("id = ?", proofID).
		Updates(updates).Error
}

func (r *orderRepository) UpdateOrderPaymentStatus(orderID uint, paymentStatus string) error {
	return r.db.Model(&Order{}).
		Where("id = ?", orderID).
		Update("payment_status", paymentStatus).Error
}
