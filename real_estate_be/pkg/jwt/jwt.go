package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte("SUPER_SECRET_KEY_ACCESS")
	refreshSecret = []byte("SUPER_SECRET_KEY_REFRESH")
)

type AccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateAccessToken — hết hạn sau 15 phút
func GenerateAccessToken(email string) (string, error) {
	claims := AccessClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

// GenerateRefreshToken — hết hạn sau 7 ngày
func GenerateRefreshToken(email string) (string, error) {
	claims := RefreshClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

// ParseAccessToken — giải mã access token
func ParseAccessToken(tokenStr string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&AccessClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return accessSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return token.Claims.(*AccessClaims), nil
}

// ParseRefreshToken — giải mã refresh token
func ParseRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&RefreshClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return refreshSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	return token.Claims.(*RefreshClaims), nil
}

