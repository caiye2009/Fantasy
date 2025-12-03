package auth

import (
    "crypto/rand"
    "encoding/hex"
    
    "golang.org/x/crypto/bcrypt"
)

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

// VerifyPassword 密码验证
func VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateActivationKey 生成激活码(32位hex)
func GenerateActivationKey() (string, error) {
    bytes := make([]byte, 16)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}