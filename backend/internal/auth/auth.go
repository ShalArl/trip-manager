package auth

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// TokenClaims represents the JWT claims
type TokenClaims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// AuthManager handles JWT token generation and password hashing
type AuthManager struct {
	jwtSecret     string
	tokenDuration time.Duration
}

// NewAuthManager creates a new auth manager
func NewAuthManager(jwtSecret string, tokenDuration time.Duration) *AuthManager {
	return &AuthManager{
		jwtSecret:     jwtSecret,
		tokenDuration: tokenDuration,
	}
}

// HashPassword hashes a password using bcrypt
func (am *AuthManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against its hash
func (am *AuthManager) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateToken generates a JWT token for a user
func (am *AuthManager) GenerateToken(userID, email, name string) (string, error) {
	now := time.Now()
	expirationTime := now.Add(am.tokenDuration)

	claims := TokenClaims{
		UserID: userID,
		Email:  email,
		Name:   name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(am.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyToken verifies and parses a JWT token
func (am *AuthManager) VerifyToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(am.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GetTokenExpiresIn returns the token expiration duration in seconds
func (am *AuthManager) GetTokenExpiresIn() int {
	return int(am.tokenDuration.Seconds())
}

// HashEmail creates a deterministic hash of email (for lookups without storing plain email in certain contexts)
func HashEmail(email string) string {
	hash := sha256.Sum256([]byte(email))
	return fmt.Sprintf("%x", hash)
}
