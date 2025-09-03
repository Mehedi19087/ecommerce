package auth

import (
	"errors"
	"fmt" // ✅ Add for token ID generation
	"gorm.io/gorm"
	"log" // ✅ Add for logging
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	//"golang.org/x/crypto/bcrypt"
)

// ✅ Global token management (secure token deletion system)
var (
	// Blacklisted tokens (use Redis in production)
	blacklistedTokens = make(map[string]bool)

	// Active tokens per user (for single device login)
	activeUserTokens = make(map[uint]string) // userID -> tokenID
)

// ✅ Updated interface with logout method
type UserService interface {
	//Register(name, email, password string) (*User, string, error)
	//Login(email, password string) (string, error)
	Logout(tokenID string) error // ✅ NEW: Proper logout

	GenerateToken(userID uint) (string, error)


	// Profile methods
	GetProfile(userID uint) (*User, error)
	UpdateProfile(userID uint, req UpdateProfileRequest) (*User, error)

	// Address methods
	GetAddresses(userID uint) ([]Address, error)
	CreateAddress(userID uint, req CreateAddressRequest) (*Address, error)
	UpdateAddress(userID uint, addressID uint, req CreateAddressRequest) (*Address, error)
	DeleteAddress(userID uint, addressID uint) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

// ✅ SECURE: Complete token generation with all security fixes
func (s *userService) GenerateToken(userID uint) (string, error) {
	// ✅ Generate unique token ID (prevents token reuse attacks)
	tokenID := fmt.Sprintf("token_%d_%d", userID, time.Now().UnixNano())

	// ✅ Delete old tokens for single device login
	s.deleteOldUserTokens(userID)

	// ✅ Create secure token with proper claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"jti":     tokenID,                              // ✅ Unique token ID
		"iss":     "ecommerce-api",                      // ✅ Issuer verification
		"aud":     "ecommerce-app",                      // ✅ Audience restriction
		"exp":     time.Now().Add(time.Hour * 2).Unix(), // ✅ Short expiry (2 hours)
		"iat":     time.Now().Unix(),                    // ✅ Issued at time
		"nbf":     time.Now().Unix(),                    // ✅ Not valid before
	})

	// ✅ SECURE: Strong secret requirement (no fallback)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET environment variable is required")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	// ✅ Store active token for this user
	activeUserTokens[userID] = tokenID

	// ✅ Log token creation for security monitoring
	log.Printf("🔑 Token created: UserID=%d, TokenID=%s", userID, tokenID)

	return tokenString, nil
}

// ✅ SECURE: Delete old tokens when user logs in from new device
func (s *userService) deleteOldUserTokens(userID uint) {
	if oldTokenID, exists := activeUserTokens[userID]; exists {
		// Add old token to blacklist (this "deletes" it)
		blacklistedTokens[oldTokenID] = true
		log.Printf("🗑️ Old token deleted: UserID=%d, TokenID=%s", userID, oldTokenID)
	}
}

// ✅ SECURE: Proper logout that actually works
func (s *userService) Logout(tokenID string) error {
	if tokenID == "" {
		return errors.New("token ID is required")
	}

	// ✅ Delete: Add token to blacklist
	blacklistedTokens[tokenID] = true

	// ✅ Delete: Remove from active tokens using Go's built-in delete()
	for userID, activeTokenID := range activeUserTokens {
		if activeTokenID == tokenID {
			delete(activeUserTokens, userID) // ✅ This removes the userID key from map
			break
		}
	}
	log.Printf("🗑️ Token deleted (logout): TokenID=%s", tokenID)
	return nil
}

// ✅ SECURE: Check if token is deleted/blacklisted (for middleware)
func IsTokenDeleted(tokenID string) bool {
	return blacklistedTokens[tokenID]
}

// ✅ SECURE: Force delete all tokens for a user (for password change)
func (s *userService) DeleteAllUserTokens(userID uint) {
	if tokenID, exists := activeUserTokens[userID]; exists {
		blacklistedTokens[tokenID] = true
		delete(activeUserTokens, userID)
		log.Printf("🗑️ All tokens deleted for UserID=%d", userID)
	}
}

func (s *userService) GetProfile(userID uint) (*User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdateProfile(userID uint, req UpdateProfileRequest) (*User, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields
	user.Name = req.Name
	user.Phone = req.Phone
	user.Birthday = req.Birthday
	user.Gender = req.Gender
	err = s.repo.UpdateProfile(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetAddresses(userID uint) ([]Address, error) {
	return s.repo.GetUserAddresses(userID)
}

func (s *userService) CreateAddress(userID uint, req CreateAddressRequest) (*Address, error) {
	address := &Address{
		UserID:  userID,
		Name:    req.Name,
		Phone:   req.Phone,
		Address: req.Address,
		City:    req.City,
		Zone:    req.Zone,
		Label:   req.Label,
	}
	err := s.repo.CreateAddress(address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *userService) UpdateAddress(userID uint, addressID uint, req CreateAddressRequest) (*Address, error) {
	address, err := s.repo.GetAddressByID(addressID, userID)
	if err != nil {
		return nil, errors.New("address not found")
	}

	// Update fields
	address.Name = req.Name
	address.Phone = req.Phone
	address.Address = req.Address
	address.City = req.City
	address.Zone = req.Zone
	address.Label = req.Label
	err = s.repo.UpdateAddress(address)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *userService) DeleteAddress(userID uint, addressID uint) error {
	return s.repo.DeleteAddress(addressID, userID)
}
