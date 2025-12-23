package config

import (
	"log"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"github.com/casbin/casbin/v2"
	"back/pkg/auth"
)

func InitAuth(db *gorm.DB, rdb *redis.Client, cfg *Config, enforcer *casbin.Enforcer) (*auth.JWTWang, *auth.AuthWang) {
	jwtWang := auth.NewJWTWang(cfg.JWTSecret)
	authWang := auth.NewAuthWang(jwtWang, enforcer)

	log.Println("✓ Auth initialized with Casbin")
	return jwtWang, authWang
}