package config

import (
	"log"
	"gorm.io/gorm"

	supplierInfra "back/internal/supplier/infra"
	clientInfra "back/internal/client/infra"
	userInfra "back/internal/user/infra"
	materialInfra "back/internal/material/infra"
	processInfra "back/internal/process/infra"
	pricingInfra "back/internal/pricing/infra"
	productInfra "back/internal/product/infra"
	planInfra "back/internal/plan/infra"
	orderInfra "back/internal/order/infra"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("=== Starting Database Migration ===")

	err := db.AutoMigrate(
		&userInfra.UserPO{},
		&supplierInfra.SupplierPO{},
		&clientInfra.ClientPO{},
		&materialInfra.MaterialPO{},
		&processInfra.ProcessPO{},
		&pricingInfra.SupplierPricePO{},
		&productInfra.ProductPO{},
		&planInfra.PlanPO{},
		&orderInfra.OrderPO{},
	)

	if err != nil {
		log.Printf("✗ Migration failed: %v", err)
		return err
	}

	log.Println("✓ All tables migrated successfully")
	return nil
}