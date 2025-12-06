package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTWang struct {
	secretKey string
}

// AccessTokenClaims Access Token 的 Claims（包含 loginID 和 role）
type AccessTokenClaims struct {
	LoginID string `json:"loginId"` // 改为小驼峰
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// RefreshTokenClaims Refresh Token 的 Claims（只包含 loginID）
type RefreshTokenClaims struct {
	LoginID string `json:"loginId"`
	jwt.RegisteredClaims
}

func NewJWTWang(secretKey string) *JWTWang {
	return &JWTWang{secretKey: secretKey}
}

// GenerateAccessToken 生成 Access Token（1小时有效期）
func (j *JWTWang) GenerateAccessToken(loginID string, role string) (string, error) {
	claims := &AccessTokenClaims{
		LoginID: loginID,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)), // 1小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// GenerateRefreshToken 生成 Refresh Token（30天有效期）
func (j *JWTWang) GenerateRefreshToken(loginID string) (string, error) {
	claims := &RefreshTokenClaims{
		LoginID: loginID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30天
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateAccessToken 验证 Access Token
func (j *JWTWang) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid access token")
}

// ValidateRefreshToken 验证 Refresh Token（返回 loginID）
func (j *JWTWang) ValidateRefreshToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims.LoginID, nil
	}

	return "", errors.New("invalid refresh token")
}