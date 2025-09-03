package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// ✅ SECURE: Complete token validation with all security checks
func validateToken(tokenString string) (jwt.MapClaims, error) {
	// ✅ Require strong secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable required")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ✅ Strict signing method validation
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// ✅ Check if token is deleted/blacklisted
		if jti, exists := claims["jti"]; exists {
			tokenID := jti.(string)
			if IsTokenDeleted(tokenID) {
				log.Printf("🚫 Deleted token attempted: %s", tokenID)
				return nil, errors.New("token has been deleted")
			}

			// ✅ Log token usage for security monitoring
			userID := claims["user_id"]
			log.Printf("🔍 Token used: UserID=%v, TokenID=%s", userID, tokenID)
		}

		// ✅ Validate issuer (prevent cross-app token reuse)
		if iss, exists := claims["iss"]; exists {
			if iss != "ecommerce-api" {
				return nil, errors.New("invalid token issuer")
			}
		}

		// ✅ Validate audience (additional security layer)
		if aud, exists := claims["aud"]; exists {
			if aud != "ecommerce-app" {
				return nil, errors.New("invalid token audience")
			}
		}

		// ✅ Check token age for suspicious activity
		if iat, exists := claims["iat"]; exists {
			issuedAt := int64(iat.(float64))
			tokenAge := time.Now().Unix() - issuedAt

			// Alert on very old tokens (potential replay attack)
			if tokenAge > 7200 { // 2 hours
				log.Printf("⚠️ Old token used (age: %d seconds): %v", tokenAge, claims["jti"])
			}
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ✅ SECURE: Enhanced JWT middleware
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Validate Bearer format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Use: Bearer <token>",
			})
			c.Abort()
			return
		}

		// Validate token with all security checks
		claims, err := validateToken(tokenParts[1])
		if err != nil {
			log.Printf("🚨 Token validation failed: %s", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user context for controllers
		c.Set("userID", uint(claims["user_id"].(float64)))
		if jti, exists := claims["jti"]; exists {
			c.Set("tokenID", jti)
		}

		c.Next()
	}
}

func LocationTrackingMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if ip == "" || ip == "127.0.0.1" || ip == "::1" {
			ip = "8.8.8.8" // fallback for local testing
		}

		resp, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%s", ip))
		if err != nil {
			ctx.Next()
			return
		}
		defer resp.Body.Close()

		var data struct {
			Country    string `json:"country"`
			RegionName string `json:"regionName"`
			City       string `json:"city"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			ctx.Next()
			return
		}

		// Store in DB
		db.Exec(
			"INSERT INTO visitor_logs (ip, country, region, city, visited_at) VALUES (?, ?, ?, ?, ?)",
			ip, data.Country, data.RegionName, data.City, time.Now(),
		)

		ctx.Next()
	}
}
