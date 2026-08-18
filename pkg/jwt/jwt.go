package jwt

import (
	"errors"
	"os"
	"strings"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"backend-challenge-golang/pkg/datetime"
)

const (
	EnvSecretKey          = "JWT_SECRET"
	FallbackSecret        = "local-dev-jwt-secret"
	DefaultExpireDuration = 7 * 24 * time.Hour
)

var ErrSecretKeyRequired = errors.New("jwt secret key is required")

type Claims struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
	jwtlib.RegisteredClaims
}

type JWTService struct {
	secretKey      []byte
	expireDuration time.Duration
}

func SecretFromEnvironment() string {
	secretKey := strings.TrimSpace(os.Getenv(EnvSecretKey))
	if secretKey == "" {
		return FallbackSecret
	}
	return secretKey
}

func NewJWTService(secretKey string, expireDuration time.Duration) (*JWTService, error) {
	secretKey = strings.TrimSpace(secretKey)
	if secretKey == "" {
		return nil, ErrSecretKeyRequired
	}
	if expireDuration <= 0 {
		expireDuration = DefaultExpireDuration
	}

	return &JWTService{
		secretKey:      []byte(secretKey),
		expireDuration: expireDuration,
	}, nil
}

func (service *JWTService) CreateToken(userID string, name string) (string, error) {
	now := datetime.GetCurrentDateTimeNow()
	claims := &Claims{
		UserID: userID,
		Name:   name,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwtlib.NewNumericDate(now.Add(service.expireDuration)),
			IssuedAt:  jwtlib.NewNumericDate(now),
			NotBefore: jwtlib.NewNumericDate(now),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(service.secretKey)
}

func (service *JWTService) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwtlib.ParseWithClaims(tokenString, claims, func(token *jwtlib.Token) (any, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return service.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
