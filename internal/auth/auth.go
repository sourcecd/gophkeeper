package auth

import (
	// "context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/metadata"
)

const tokenDuration = 12 * time.Hour

type User struct {
	Username       string
	HashedPassword string
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int64
}

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

func (jwtm *JWTManager) PrepareUser(username string, password string) (*User, error) {
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

func (user *User) IsCorrectPassword(passwordrealhashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordrealhashed), []byte(user.HashedPassword))
	return err == nil
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

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

/*func JwtInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = metadata.NewOutgoingContext(ctx, map[string][]string{
		"token": {"test"},
	})
	return handler(ctx, req)
}*/
