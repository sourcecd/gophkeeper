// Package auth for authorization support
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const tokenDuration = 12 * time.Hour

// JWTManagerIface interface
type JWTManagerIface interface {
	PrepareUser(username string, password string) (*User, error)
	Generate(userid int64) (string, error)
	Verify(accessToken string) (*UserClaims, error)
}

// User struct for work with username and password
type User struct {
	Username       string
	HashedPassword string
}

// UserClaims main jwt token structure
type UserClaims struct {
	jwt.RegisteredClaims
	UserID int64
}

// JWTManager configuration for jwt
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// PrepareUser method for prepare password hash
func (manager *JWTManager) PrepareUser(username string, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	user := &User{
		Username:       username,
		HashedPassword: string(hashedPassword),
	}

	return user, nil
}

// IsCorrectPassword compare password with hash
func (user *User) IsCorrectPassword(passwordrealhashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordrealhashed), []byte(user.HashedPassword))
	return err == nil
}

// NewJWTManager create jwt instance
func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

// Generate jwt token
func (manager *JWTManager) Generate(userid int64) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
		},
		UserID: userid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// Verify jwt, unpack and extract parameters
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
