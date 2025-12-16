package auth

import (
	"time"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type JWTWang struct {
	secretKey string
}

type AccessTokenClaims struct {
	LoginID string `json:"loginId"`
	UserName string `json:"userName"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	LoginID string `json:"loginId"`
	jwt.RegisteredClaims
}

func NewJWTWang(secretKey string) *JWTWang {
	return &JWTWang{secretKey: secretKey}
}

func (j *JWTWang) GenerateAccessToken(loginID string, userName string, role string) (string, error) {
	claims := &AccessTokenClaims{
		LoginID: loginID,
		UserName: userName,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTWang) GenerateRefreshToken(loginID string) (string, error) {
	claims := &RefreshTokenClaims{
		LoginID: loginID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

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