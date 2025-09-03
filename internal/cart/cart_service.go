package cart

import (
	"ecommerce/internal/catalog"
	"errors"
)

type CartService interface {
	GetCartByUserID(userID uint) (*Cart, error)
	AddItemToCart(userID uint, productID uint, quantity int) (*Cart, error)
	UpdateCartItem(userID uint, itemID uint, quantity int) (*Cart, error)
	RemoveCartItem(userID uint, itemID uint) (*Cart, error)
	ClearCart(userID uint) error
}

type cartService struct {
	repo        CartRepository
	productRepo catalog.ProductRepository
}

func NewCartService(repo CartRepository, productRepo catalog.ProductRepository) CartService {
	return &cartService{repo: repo, productRepo: productRepo}
}

func (s *cartService) GetCartByUserID(userID uint) (*Cart, error) {
	cart, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return &Cart{UserID: userID, Items: []CartItem{}}, nil
	}

	return cart, nil
}

func (s *cartService) AddItemToCart(userID uint, productID uint, quantity int) (*Cart, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	// Check if product exists
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Get or create cart
	cart, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart = &Cart{UserID: userID}
		if err := s.repo.Create(cart); err != nil {
			return nil, err
		}
	}

	// Check if item already exists in cart
	for i, item := range cart.Items {
		if item.ProductID == productID {
			// Update quantity
			cart.Items[i].Quantity += quantity
			err := s.repo.UpdateItem(item.ID, cart.Items[i].Quantity)
			if err != nil {
				return nil, err
			}
			return s.GetCartByUserID(userID)
		}
	}

	// Add new item
	cartItem := &CartItem{
		CartID:    cart.ID,
		ProductID: productID,
		Product:   *product,
		Quantity:  quantity,
		Price:     product.Price,
	}

	if err := s.repo.AddItem(cartItem); err != nil {
		return nil, err
	}

	return s.GetCartByUserID(userID)
}

func (s *cartService) UpdateCartItem(userID uint, itemID uint, quantity int) (*Cart, error) {
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}

	// Check if cart exists
	cart, err := s.repo.FindByUserID(userID)
	if err != nil || cart == nil {
		return nil, errors.New("cart not found")
	}

	// Check if item belongs to user's cart
	item, err := s.repo.FindCartItemByID(itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if item.CartID != cart.ID {
		return nil, errors.New("item does not belong to user's cart")
	}

	// Update item quantity
	if err := s.repo.UpdateItem(itemID, quantity); err != nil {
		return nil, err
	}

	return s.GetCartByUserID(userID)
}

func (s *cartService) RemoveCartItem(userID uint, itemID uint) (*Cart, error) {
	// Check if cart exists
	cart, err := s.repo.FindByUserID(userID)
	if err != nil || cart == nil {
		return nil, errors.New("cart not found")
	}

	// Check if item belongs to user's cart
	item, err := s.repo.FindCartItemByID(itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if item.CartID != cart.ID {
		return nil, errors.New("item does not belong to user's cart")
	}

	// Remove item
	if err := s.repo.RemoveItem(itemID); err != nil {
		return nil, err
	}

	return s.GetCartByUserID(userID)
}

func (s *cartService) ClearCart(userID uint) error {
	cart, err := s.repo.FindByUserID(userID)
	if err != nil {
		return err
	}

	if cart == nil {
		return nil // No cart to clear
	}

	return s.repo.ClearCart(cart.ID)
}
