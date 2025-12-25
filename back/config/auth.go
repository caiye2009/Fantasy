package config

import (
	"log"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"github.com/casbin/casbin/v2"
	"back/pkg/auth"
)

func InitAuth(db *gorm.DB, rdb *redis.Client, cfg *Config, enforcer *casbin.Enforcer) (*auth.JWTWang, *auth.AuthWang, *auth.WhitelistManager) {
	jwtWang := auth.NewJWTWang(cfg.JWTSecret)
	whitelistManager := auth.NewWhitelistManager(rdb)
	authWang := auth.NewAuthWang(jwtWang, enforcer, whitelistManager)

	log.Println("✓ Auth initialized with Casbin and JWT whitelist")
	return jwtWang, authWang, whitelistManager
}