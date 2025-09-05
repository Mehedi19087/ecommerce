package order

import (
	"ecommerce/internal/cart"
	"errors"
	"fmt"
	"time"

	gomail "gopkg.in/gomail.v2"

	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrderFromCart(userID uint, orderData CreateOrderRequest) (*Order, error)
	GetUserOrders(userID uint) ([]Order, error)
	GetOrderByID(orderID uint, userID uint) (*Order, error)
	CancelOrder(orderID uint, userID uint) error
	SubmitPaymentProof(orderID uint, userID uint, proofData SubmitPaymentProofRequest) (*PaymentProof, error)
	GetPaymentProof(orderID uint, userID uint) (*PaymentProof, error)
	UpdatePaymentProof(orderID uint, userID uint, proofData SubmitPaymentProofRequest) (*PaymentProof, error)
	GetAllOrdersAdmin() ([]Order, error)
	UpdateOrderStatusAdmin(orderID uint, status string) error
	ReviewPaymentProofAdmin(proofID uint, status string, adminNotes string, reviewerID uint) error
}

type orderService struct {
	repo        OrderRepository
	cartService cart.CartService
}

func NewOrderService(repo OrderRepository, cartService cart.CartService) OrderService {
	return &orderService{
		repo:        repo,
		cartService: cartService,
	}
}

func (s *orderService) CreateOrderFromCart(userID uint, orderData CreateOrderRequest) (*Order, error) {
	// Get cart
	userCart, err := s.cartService.GetCartByUserID(userID)
	if err != nil {
		return nil, errors.New("failed to get cart")
	}

	if len(userCart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Create order
	order := &Order{
		UserID:          userID,
		OrderNumber:     s.generateOrderNumber(),
		Status:          "pending",
		PaymentStatus:   "pending",
		Total:           userCart.CalculateTotal(),
		ShippingAddress: orderData.ShippingAddress,
		CustomerName:    orderData.CustomerName,
		CustomerPhone:   orderData.CustomerPhone,
		PaymentMethod:   orderData.PaymentMethod,
		Notes:           orderData.Notes,
	}

	// Save order
	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	// Create order items
	for _, cartItem := range userCart.Items {
		var productImage string 
		if len(cartItem.Product.Image) > 0 {
        productImage = cartItem.Product.Image[0] // Take first image
    }
		orderItem := &OrderItem{
			OrderID:      order.ID,
			ProductID:    cartItem.ProductID,
			ProductName:  cartItem.Product.Name,
			Price:        cartItem.Product.Price,
			Quantity:     cartItem.Quantity,
			Subtotal:     float64(cartItem.Quantity) * cartItem.Product.Price,
			ProductImage: productImage,
			ProductSKU:   cartItem.Product.SKU,
		}

		if err := s.repo.CreateOrderItem(orderItem); err != nil {
			return nil, err
		}

		order.Items = append(order.Items, *orderItem)
	}
    
   go func() {
    // Send email to user
    m := gomail.NewMessage()
    m.SetHeader("From", "mehedimahbub706@gmail.com")
    m.SetHeader("To", "alrizvanthreads@yahoo.com") // Client email
    m.SetHeader("Subject", "Order Created - Pending Payment")
    m.SetBody("text/plain", "Your order has been created and is currently pending. Please complete your payment and submit the payment proof.")

    d := gomail.NewDialer("smtp.gmail.com", 587, "mehedimahbub706@gmail.com", "lvsadqnenykybjcu")
    if err := d.DialAndSend(m); err != nil {
        fmt.Println("Failed to send email:", err)
    }}()
	// Clear cart
	s.cartService.ClearCart(userID)

	return order, nil
}

func (s *orderService) GetUserOrders(userID uint) ([]Order, error) {
	return s.repo.GetByUserID(userID)
}

func (s *orderService) GetOrderByID(orderID uint, userID uint) (*Order, error) {
	order, err := s.repo.GetByID(orderID, userID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}
func (s *orderService) CancelOrder(orderID uint, userID uint) error {
	// Check if order exists and belongs to user
	order, err := s.repo.GetByID(orderID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return err
	}

	// Check if order can be cancelled
	if order.Status != "pending" {
		return errors.New("only pending orders can be cancelled")
	}

	// Update order status to cancelled
	return s.repo.UpdateStatus(orderID, userID, "cancelled")
}

func (s *orderService) SubmitPaymentProof(orderID uint, userID uint, proofData SubmitPaymentProofRequest) (*PaymentProof, error) {
	// Check if order exists and belongs to user
	order, err := s.repo.GetByID(orderID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Check if order status allows payment proof submission
	if order.Status == "cancelled" {
		return nil, errors.New("cannot submit payment proof for cancelled order")
	}

	if order.PaymentStatus == "paid" {
		return nil, errors.New("payment already confirmed for this order")
	}
	// Create payment proof
	proof := &PaymentProof{
		OrderID:       orderID,
		TransactionID: proofData.TransactionID,
		PaymentMethod: proofData.PaymentMethod,
		Amount:        proofData.Amount,
		Screenshot:    proofData.Screenshot,
		SenderNumber:  proofData.SenderNumber,
		SenderName:    proofData.SenderName,
		PaymentDate:   proofData.PaymentDate,
		Status:        "pending",
	}

	if err := s.repo.CreatePaymentProof(proof); err != nil {
		return nil, err
	}

	return proof, nil
}

func (s *orderService) GetPaymentProof(orderID uint, userID uint) (*PaymentProof, error) {
	proof, err := s.repo.GetPaymentProofByOrderID(orderID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment proof not found")
		}
		return nil, err
	}
	return proof, nil
}
func (s *orderService) UpdatePaymentProof(orderID uint, userID uint, proofData SubmitPaymentProofRequest) (*PaymentProof, error) {
	// Check if payment proof exists
	existingProof, err := s.repo.GetPaymentProofByOrderID(orderID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment proof not found")
		}
		return nil, err
	}

	// Check if payment proof can be updated
	if existingProof.Status != "pending" {
		return nil, errors.New("cannot update payment proof that has been reviewed")
	}

	// Update payment proof
	err = s.repo.UpdatePaymentProof(orderID, userID, proofData)
	if err != nil {
		return nil, err
	}

	// Return updated payment proof
	return s.repo.GetPaymentProofByOrderID(orderID, userID)
}

func (s *orderService) GetAllOrdersAdmin() ([]Order, error) {
	return s.repo.GetAllOrders()
}
func (s *orderService) UpdateOrderStatusAdmin(orderID uint, status string) error {
	// Check if order exists using repository method
	_, err := s.repo.GetByID(orderID, 0) // Use 0 for userID since admin can access any order
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found")
		}
		return err
	}

	return s.repo.UpdateOrderStatusAdmin(orderID, status)
}

func (s *orderService) ReviewPaymentProofAdmin(proofID uint, status string, adminNotes string, reviewerID uint) error {
	// Check if payment proof exists
	proof, err := s.repo.GetPaymentProofByID(proofID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("payment proof not found")
		}
		return err
	}

	if proof.Status != "pending" {
		return errors.New("payment proof has already been reviewed")
	}

	// Update payment proof status
	err = s.repo.ReviewPaymentProof(proofID, status, adminNotes, reviewerID)
	if err != nil {
		return err
	}

	// If approved, update order payment status using repository method
	if status == "approved" {
		return s.repo.UpdateOrderPaymentStatus(proof.OrderID, "paid") // ✅ Use repository method
	}

	return nil
}

func (s *orderService) generateOrderNumber() string {
	return fmt.Sprintf("ORD%d", time.Now().Unix())
}
