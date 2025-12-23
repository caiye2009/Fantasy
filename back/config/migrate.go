package config

import (
	"log"
	"gorm.io/gorm"

	userDomain "back/internal/user/domain"
	supplierDomain "back/internal/supplier/domain"
	clientDomain "back/internal/client/domain"
	materialDomain "back/internal/material/domain"
	processDomain "back/internal/process/domain"
	pricingDomain "back/internal/pricing/domain"
	productDomain "back/internal/product/domain"
	planDomain "back/internal/plan/domain"
	orderDomain "back/internal/order/domain"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("=== Starting Database Migration ===")

	err := db.AutoMigrate(
		&userDomain.User{},
		&userDomain.Department{},
		&userDomain.Role{},
		&supplierDomain.Supplier{},
		&clientDomain.Client{},
		&materialDomain.Material{},
		&processDomain.Process{},
		&pricingDomain.SupplierPrice{},
		&productDomain.Product{},
		&planDomain.Plan{},
		&orderDomain.Order{},
		&orderDomain.OrderParticipant{},
		&orderDomain.OrderProgress{},
		&orderDomain.OrderEvent{},
	)

	if err != nil {
		log.Printf("✗ Migration failed: %v", err)
		return err
	}

	log.Println("✓ All tables migrated successfully")
	return nil
}