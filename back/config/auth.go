package config

import (
	"log"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"back/pkg/auth"
)

func InitAuth(db *gorm.DB, rdb *redis.Client, cfg *Config) *auth.AuthWang {
	jwtWang := auth.NewJWTWang(cfg.JWTSecret)

	casbinWang, err := auth.NewCasbinWang(db, rdb, cfg.CasbinModelPath)
	if err != nil {
		log.Fatalf("Failed to initialize Casbin: %v", err)
	}

	// 初始化 Casbin 策略
	if err := InitCasbinPolicies(casbinWang.GetEnforcer()); err != nil {
		log.Fatalf("Failed to initialize Casbin policies: %v", err)
	}

	authWang := auth.NewAuthWang(jwtWang, casbinWang, db)

	log.Println("✓ Auth initialized")
	return authWang
}