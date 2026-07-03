package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
	"fmt"

	"github.com/KuRoute/kuroute/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func getSecretKey() string {
	_ = godotenv.Load()
	return os.Getenv("JWT_SECRET")
}

func GenerateAccessToken(userID uuid.UUID, role domain.UserRole, hubID uuid.UUID) (string, error) {
	secretKey := getSecretKey()
	if secretKey == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	expiryStr := os.Getenv("JWT_ACCESS_EXPIRY")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil || expiryStr == "" {
		expiry = 15 * time.Minute
	}

	claims := &domain.JWTClaims{
		UserID: userID,
		Role:   role,
		HubID:  hubID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateRefreshToken(userID uuid.UUID, role domain.UserRole, hubID uuid.UUID) (string, error) {
	secretKey := getSecretKey()
	if secretKey == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	expiryStr := os.Getenv("JWT_REFRESH_EXPIRY")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil || expiryStr == "" {
		expiry = 720 * time.Hour
	}

	claims := &domain.JWTClaims{
		UserID: userID,
		Role:   role,
		HubID:  hubID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func GenerateSetupToken(userID uuid.UUID) (string, error) {
	secretKey := getSecretKey()
	if secretKey == "" {
		return "", errors.New("JWT_SECRET not set")
	}

	claims := &domain.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func ValidateToken(tokenString string) (*domain.JWTClaims, error) {
	secretKey := getSecretKey()
	if secretKey == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func ExtractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	fmt.Printf("Authorization: %q\n", r.Header.Get("Authorization"))

	parts := strings.Fields(authHeader)

	if len(parts) != 2 {
		return ""
	}

	if parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func ExtractUserID(r *http.Request) (uuid.UUID, error) {
	tokenString := ExtractToken(r)
	if tokenString == "" {
		return uuid.Nil, errors.New("no token provided")
	}

	claims, err := ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}

	return claims.UserID, nil
}

func HashToken(tokenString string) string {
	hash := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(hash[:])
}

func GetAccessExpiryDuration() time.Duration {
	expiryStr := os.Getenv("JWT_ACCESS_EXPIRY")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil || expiryStr == "" {
		return 15 * time.Minute
	}
	return expiry
}

func GetRefreshExpiryDuration() time.Duration {
	expiryStr := os.Getenv("JWT_REFRESH_EXPIRY")
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil || expiryStr == "" {
		return 720 * time.Hour
	}
	return expiry
}