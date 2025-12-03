package config

import (
	"log"
	"gorm.io/gorm"

	"back/internal/vendor"
	"back/internal/client"
	"back/internal/user"
	"back/internal/material"
	materialPrice "back/internal/material/price"
	"back/internal/process"
	processPrice "back/internal/process/price"
	"back/internal/product"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("=== Starting Database Migration ===")

	err := db.AutoMigrate(
		&user.User{},
		&vendor.Vendor{},
		&client.Client{},
		&material.Material{},
		&materialPrice.MaterialPrice{},
		&process.Process{},
		&processPrice.ProcessPrice{},
		&product.Product{},
	)

	if err != nil {
		log.Printf("✗ Migration failed: %v", err)
		return err
	}

	log.Println("✓ All tables migrated successfully")
	return nil
}