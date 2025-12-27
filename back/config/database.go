package config

import (
	"gorm.io/driver/postgres" // ← 改为 postgres
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func InitDatabase(cfg *Config) *gorm.DB {
	// PostgreSQL DSN 格式:
	// host=localhost user=postgres password=postgres dbname=fantasy port=5432 sslmode=disable TimeZone=Asia/Shanghai

	db, err := gorm.Open(postgres.Open(cfg.DatabaseDSN), &gorm.Config{ // ← 改为 postgres.Open
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✓ Database connected (PostgreSQL)")
	return db
}
