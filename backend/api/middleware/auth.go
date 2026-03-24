package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateAccessToken creates a short-lived JWT (15 minutes).
func GenerateAccessToken(userID, email string, isAdmin bool) (string, error) {
	claims := Claims{
		UserID:  userID,
		Email:   email,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// RefreshClaims extends RegisteredClaims with token version for revocation.
type RefreshClaims struct {
	TokenVersion int `json:"token_version"`
	jwt.RegisteredClaims
}

// GenerateRefreshToken creates a long-lived JWT (7 days) with token version for revocation.
func GenerateRefreshToken(userID string, tokenVersion int) (string, error) {
	claims := RefreshClaims{
		TokenVersion: tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken validates and parses a JWT token.
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Explicit algorithm check — prevent alg-switch attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// ParseRefreshToken validates a refresh token and returns the user ID and token version.
func ParseRefreshToken(tokenString string) (string, int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return "", 0, err
	}
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims.Subject, claims.TokenVersion, nil
	}
	return "", 0, jwt.ErrTokenInvalidClaims
}

// JWTAuth is Gin middleware that validates JWT bearer tokens.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization_required"})
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			log.Printf("[security] JWT validation failed: ip=%s path=%s error=%v", c.ClientIP(), c.Request.URL.Path, err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("is_admin", claims.IsAdmin)
		c.Next()
	}
}

// GetUserID extracts user ID from gin context (set by JWTAuth middleware).
func GetUserID(c *gin.Context) string {
	if v, exists := c.Get("user_id"); exists {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetUserEmail extracts user email from gin context (set by JWTAuth middleware).
func GetUserEmail(c *gin.Context) string {
	if v, exists := c.Get("user_email"); exists {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
