package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/m13ha/asiko/errors"
)

var (
	jwtKey          = []byte(os.Getenv("JWT_SECRET_KEY"))
	tokenExpiration = time.Hour * 24
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GetUserIDFromContext(c *gin.Context) string {
	if userID, exists := c.Get("userID"); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUUIDFromContext returns the authenticated user's UUID from context.
// It supports both the new strongly-typed key ("userUUID") and falls back to parsing
// the legacy string key ("userID") for backward compatibility.
func GetUUIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	if v, ok := c.Get("userUUID"); ok {
		if id, ok2 := v.(uuid.UUID); ok2 {
			return id, true
		}
	}
	// Fallback: parse legacy string if present
	if s := GetUserIDFromContext(c); s != "" {
		if id, err := uuid.Parse(s); err == nil {
			// Cache it for subsequent reads in this request
			c.Set("userUUID", id)
			return id, true
		}
	}
	return uuid.UUID{}, false
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c.Request)
		if tokenString == "" {
			errors.Unauthorized(c.Writer, "")
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			errors.Unauthorized(c.Writer, "Invalid token")
			c.Abort()
			return
		}

		// Parse and store strongly-typed UUID; if invalid, treat as unauthorized
		uid, err := uuid.Parse(claims.UserID)
		if err != nil {
			errors.Unauthorized(c.Writer, "Invalid token")
			c.Abort()
			return
		}

		// Set both for compatibility during transition
		c.Set("userUUID", uid)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) != 2 || strArr[0] != "Bearer" {
		return ""
	}
	return strArr[1]
}

func GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(tokenExpiration)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "appointment_app",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func RefreshToken(c *gin.Context) {
	tokenString := extractToken(c.Request)
	if tokenString == "" {
		errors.Unauthorized(c.Writer, "Missing token")
		return
	}

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		errors.Unauthorized(c.Writer, "Invalid token")
		return
	}

	if claims.ExpiresAt != nil && time.Until(claims.ExpiresAt.Time) > 30*time.Minute {
		errors.BadRequest(c.Writer, "Token not near expiration")
		return
	}

	newToken, err := GenerateToken(claims.UserID)
	if err != nil {
		errors.InternalServerError(c.Writer, "Could not generate token")
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}

// --- Device Token Logic ---

var (
	deviceTokenExpiration = time.Minute * 10
)

type DeviceClaims struct {
	DeviceID string `json:"device_id"`
	jwt.RegisteredClaims
}

func GenerateDeviceToken(deviceID string) (string, error) {
	expirationTime := time.Now().Add(deviceTokenExpiration)
	claims := &DeviceClaims{
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "appointment_app_device",
			Subject:   deviceID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ValidateDeviceToken(tokenString string) (string, error) {
	claims := &DeviceClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", fmt.Errorf("malformed token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return "", fmt.Errorf("token is expired or not yet valid")
			}
		}
		return "", fmt.Errorf("invalid token")
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if claims.DeviceID == "" {
		return "", fmt.Errorf("token missing device_id")
	}

	return claims.DeviceID, nil
}
