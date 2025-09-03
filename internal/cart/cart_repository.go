package cart

import (
	"gorm.io/gorm"
)

type CartRepository interface {
	FindByUserID(userID uint) (*Cart, error)
	Create(cart *Cart) error
	AddItem(item *CartItem) error
	UpdateItem(itemID uint, quantity int) error
	RemoveItem(itemID uint) error
	ClearCart(cartID uint) error
	FindCartItemByID(itemID uint) (*CartItem, error)
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) FindByUserID(userID uint) (*Cart, error) {
	var cart Cart

	err := r.db.Where("user_id = ?", userID).Preload("Items.Product").First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No cart found, but not an error
		}
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) Create(cart *Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) AddItem(item *CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) UpdateItem(itemID uint, quantity int) error {
	return r.db.Model(&CartItem{}).Where("id = ?", itemID).Update("quantity", quantity).Error
}

func (r *cartRepository) RemoveItem(itemID uint) error {
	return r.db.Delete(&CartItem{}, itemID).Error
}

func (r *cartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&CartItem{}).Error
}

func (r *cartRepository) FindCartItemByID(itemID uint) (*CartItem, error) {
	var item CartItem
	err := r.db.Preload("Product").First(&item, itemID).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
