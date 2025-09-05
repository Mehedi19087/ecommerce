package auth

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)

	// ✅ Add simple profile methods
	FindByID(id uint) (*User, error)
	UpdateProfile(user *User) error

	// ✅ Add simple address methods
	CreateAddress(address *Address) error
	GetUserAddresses(userID uint) ([]Address, error)
	GetAddressByID(id uint, userID uint) (*Address, error)
	UpdateAddress(address *Address) error
	DeleteAddress(id uint, userID uint) error
}

type userRepository struct {
	db *gorm.DB
}

//constructor for assigning value

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email= ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uint) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateProfile(user *User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) CreateAddress(address *Address) error {
	return r.db.Create(address).Error
}

func (r *userRepository) GetUserAddresses(userID uint) ([]Address, error) {
	var addresses []Address
	err := r.db.Where("user_id = ?", userID).Find(&addresses).Error
	return addresses, err
}

func (r *userRepository) GetAddressByID(id uint, userID uint) (*Address, error) {
	var address Address
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&address).Error
	return &address, err
}

func (r *userRepository) UpdateAddress(address *Address) error {
	return r.db.Save(address).Error
}

func (r *userRepository) DeleteAddress(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&Address{}).Error
}

