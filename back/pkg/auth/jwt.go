package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTWang struct {
	secretKey string
}

type Claims struct {
	LoginID string `json:"login_id"` // 改为 login_id
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTWang(secretKey string) *JWTWang {
	return &JWTWang{secretKey: secretKey}
}

// GenerateToken 使用 login_id 生成 JWT
func (j *JWTWang) GenerateToken(loginID string, role string) (string, error) {
	claims := &Claims{
		LoginID: loginID,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTWang) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}