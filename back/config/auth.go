package config

import (
	"log"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"back/pkg/auth"
)

func InitAuth(db *gorm.DB, rdb *redis.Client, cfg *Config) (*auth.JWTWang, *auth.AuthWang) {
	jwtWang := auth.NewJWTWang(cfg.JWTSecret)
	authWang := auth.NewAuthWang(jwtWang)

	log.Println("✓ Auth initialized")
	return jwtWang, authWang
}