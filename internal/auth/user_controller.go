package auth

import (
	"crypto/rand"
	"ecommerce/config"
	"ecommerce/database"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	//"log"
	//"strings"
	//"io"
	//"bytes"
)

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (c *UserController) GoogleLogin(ctx *gin.Context) {
	state := c.generateRandomState()

	url := config.GoogleOAuthConfig.AuthCodeURL(state)

	ctx.JSON(http.StatusOK, gin.H{
		"auth_url": url,
		"message":  "Redirect to this url to login with Google",
		"state":    state,
	})
}

func (c *UserController) GoogleCallBack(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization code not provided by google",
		})
		return
	}
	//exchange code for token
	token, err := config.GoogleOAuthConfig.Exchange(ctx, code)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to exchange code for token: " + err.Error(),
		})
		return
	}
	// get user info by using token

	userInfo, err := c.getUserInfoFromGoogle(token.AccessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user info from Google: " + err.Error(),
		})
		return
	}
	user, err := c.saveOrUpdateUser(userInfo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save user to database: " + err.Error(),
		})
		return
	}
	jwtToken, err :=c.userService.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate authentication token: " + err.Error(),
		})
		return
	}
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "https://alrizvan.com"
	}

	redirectURL := fmt.Sprintf("%s/auth/success?token=%s", frontendURL, jwtToken)
	ctx.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (c *UserController) generateRandomState() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func (c *UserController) getUserInfoFromGoogle(accessToken string) (*GoogleUserInfo, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", accessToken)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userInfo GoogleUserInfo

	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (c *UserController) saveOrUpdateUser(googleUser *GoogleUserInfo) (*User, error) {
	var user User

	// Try to find existing user by Google ID
	result := database.DB.Where("google_id = ?", googleUser.ID).First(&user)

	if result.Error != nil {
		// User doesn't exist, create new one
		user = User{
			GoogleID: googleUser.ID,
			Email:    googleUser.Email,
			Name:     googleUser.Name,
		}

		if err := database.DB.Create(&user).Error; err != nil {
			return nil, fmt.Errorf("failed to create user: %v", err)
		}

		fmt.Printf("✅ New user created: %s (%s)\n", user.Name, user.Email)
	} else {
		// User exists, check if any data has changed before updating
		needsUpdate := user.Email != googleUser.Email ||
			user.Name != googleUser.Name

		if needsUpdate {
			user.Email = googleUser.Email
			user.Name = googleUser.Name

			if err := database.DB.Save(&user).Error; err != nil {
				return nil, fmt.Errorf("failed to update user: %v", err)
			}

			fmt.Printf("✅ User updated: %s (%s)\n", user.Name, user.Email)
		} else {
			fmt.Printf("✅ User login: %s (%s) - no updates needed\n", user.Name, user.Email)
		}
	}

	return &user, nil
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)

	user, err := c.userService.GetProfile(userIDUint)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"profile": user,
	})
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)
	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := c.userService.UpdateProfile(userIDUint, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"profile": user,
	})
}

// ✅ Add simple address endpoints
func (c *UserController) GetAddresses(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)

	addresses, err := c.userService.GetAddresses(userIDUint)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get addresses",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"addresses": addresses,
	})
}

func (c *UserController) CreateAddress(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)

	var req CreateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	address, err := c.userService.CreateAddress(userIDUint, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Address created successfully",
		"address": address,
	})
}

func (c *UserController) UpdateAddress(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)

	addressID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid address ID",
		})
		return
	}
	var req CreateAddressRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	address, err := c.userService.UpdateAddress(userIDUint, uint(addressID), req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Address updated successfully",
		"address": address,
	})
}

func (c *UserController) DeleteAddress(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)
	addressID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid address ID",
		})
		return
	}

	err = c.userService.DeleteAddress(userIDUint, uint(addressID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Address deleted successfully",
	})
}

// Add this method to your user_controller.go

func (c *UserController) Logout(ctx *gin.Context) {
	// Get token ID from middleware
	tokenID, exists := ctx.Get("tokenID")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Token ID not found",
		})
		return
	}

	// Delete the token
	err := c.userService.Logout(tokenID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out (token deleted)",
	})
}

func GetVisitorCountByCity(ctx *gin.Context) {
	db := database.DB
	city := ctx.Query("city")
	date := ctx.Query("date")
	from := ctx.Query("from")
	to := ctx.Query("to")

	var count int64
	if date != "" {
		db.Raw("SELECT COUNT(*) FROM visitor_logs WHERE city = ? AND DATE(visited_at) = ?", city, date).Scan(&count)
	} else if from != "" && to != "" {
		db.Raw("SELECT COUNT(*) FROM visitor_logs WHERE city = ? AND visited_at BETWEEN ? AND ?", city, from, to).Scan(&count)
	} else {
		db.Raw("SELECT COUNT(*) FROM visitor_logs WHERE city = ?", city).Scan(&count)
	}

	ctx.JSON(200, gin.H{
		"city":  city,
		"count": count,
	})
}
