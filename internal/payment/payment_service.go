package payment

import (
	"errors"
	"fmt"
	"time"
)

type PaymentService interface {
	// Core payment operations
	InitiatePayment(orderID uint, paymentMethod string) (*Payment, error)
	CompletePayment(paymentID uint, transactionID string) error
	CancelPayment(paymentID uint) error

	// Payment queries
	GetPaymentByID(id uint) (*Payment, error)
	GetPaymentByOrderID(orderID uint) (*Payment, error)
	GetUserPayments(userID uint) ([]Payment, error)
	GetPendingPayments() ([]Payment, error)
}

type paymentService struct {
	paymentRepo PaymentRepository
}

func NewPaymentService(paymentRepo PaymentRepository) PaymentService {
	return &paymentService{
		paymentRepo: paymentRepo,
	}
}

func (s *paymentService) InitiatePayment(orderID uint, paymentMethod string) (*Payment, error) {
	// Validate payment method
	validMethods := map[string]bool{
		"bkash": true,
		"nagad": true,
		"cash":  true,
	}
	if !validMethods[paymentMethod] {
		return nil, errors.New("invalid payment method")
	}

	// Check if payment already exists
	existingPayment, err := s.paymentRepo.FindByOrderId(orderID)
	if err == nil && existingPayment != nil {
		return nil, errors.New("payment already exists for this order")
	}

	// Get order details (you might want to inject order service here)
	// For now, we'll assume the payment amount is already set in the order

	payment := &Payment{
		OrderID:       orderID,
		PaymentMethod: paymentMethod,
		Status:        "pending",
		Currency:      "BDT",
	}

	// Set expiry for digital payments
	if paymentMethod != "cash" {
		expiryTime := time.Now().Add(15 * time.Minute)
		payment.ExpiresAt = &expiryTime
	}

	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, fmt.Errorf("failed to create payment: %v", err)
	}

	return payment, nil
}

func (s *paymentService) CompletePayment(paymentID uint, transactionID string) error {
	payment, err := s.paymentRepo.FindById(paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	if payment.Status == "completed" {
		return errors.New("payment already completed")
	}

	if payment.Status == "expired" {
		return errors.New("payment has expired")
	}

	if payment.Status == "failed" {
		return errors.New("payment has failed")
	}

	// Update payment status and transaction ID
	if err := s.paymentRepo.UpdateTransactionID(paymentID, transactionID); err != nil {
		return fmt.Errorf("failed to complete payment: %v", err)
	}

	return nil
}

func (s *paymentService) CancelPayment(paymentID uint) error {
	payment, err := s.paymentRepo.FindById(paymentID)
	if err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	if payment.Status == "completed" {
		return errors.New("cannot cancel completed payment")
	}

	if err := s.paymentRepo.UpdateStatus(paymentID, "failed"); err != nil {
		return fmt.Errorf("failed to cancel payment: %v", err)
	}

	return nil
}

func (s *paymentService) GetPaymentByID(id uint) (*Payment, error) {
	if id == 0 {
		return nil, errors.New("payment ID is required")
	}
	return s.paymentRepo.FindById(id)
}

func (s *paymentService) GetPaymentByOrderID(orderID uint) (*Payment, error) {
	if orderID == 0 {
		return nil, errors.New("order ID is required")
	}
	return s.paymentRepo.FindByOrderId(orderID)
}

func (s *paymentService) GetUserPayments(userID uint) ([]Payment, error) {
	if userID == 0 {
		return nil, errors.New("user ID is required")
	}
	return s.paymentRepo.FindByUserID(userID)
}

func (s *paymentService) GetPendingPayments() ([]Payment, error) {
	return s.paymentRepo.FindPendingPayments()
}
