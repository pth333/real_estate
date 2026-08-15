package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte("SUPER_SECRET_KEY_ACCESS")
	refreshSecret = []byte("SUPER_SECRET_KEY_REFRESH")
)

type AccessClaims struct {
	Email  string `json:"email"`
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateAccessToken — hết hạn sau 15 phút, kèm user_id
func GenerateAccessToken(email string, userID uint64) (string, error) {
	claims := AccessClaims{
		Email:  email,
		UserID: userID,
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

// ExtractClaimsFromHeader trích xuất và giải mã AccessClaims từ chuỗi "Authorization" header
func ExtractClaimsFromHeader(authHeader string) (*AccessClaims, error) {
	if authHeader == "" {
		return nil, errors.New("missing authorization header")
	}

	tokenStr := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenStr = authHeader[7:]
	}

	return ParseAccessToken(tokenStr)
}

// ExtractUserIDFromHeader trích xuất thẳng UserID từ chuỗi "Authorization" header (trả về 0 nếu lỗi/khách)
func ExtractUserIDFromHeader(authHeader string) uint64 {
	claims, err := ExtractClaimsFromHeader(authHeader)
	if err != nil {
		return 0
	}
	return claims.UserID
}

// ExtractEmailFromHeader trích xuất thẳng Email từ chuỗi "Authorization" header (trả về rỗng nếu lỗi/khách)
func ExtractEmailFromHeader(authHeader string) string {
	claims, err := ExtractClaimsFromHeader(authHeader)
	if err != nil {
		return ""
	}
	return claims.Email
}
