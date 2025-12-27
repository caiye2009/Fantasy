package config

import (
	"gorm.io/gorm"
	"log"

	clientDomain "back/internal/client/domain"
	inventoryDomain "back/internal/inventory/domain"
	materialDomain "back/internal/material/domain"
	orderDomain "back/internal/order/domain"
	planDomain "back/internal/plan/domain"
	pricingDomain "back/internal/pricing/domain"
	processDomain "back/internal/process/domain"
	productDomain "back/internal/product/domain"
	supplierDomain "back/internal/supplier/domain"
	userDomain "back/internal/user/domain"
	"back/pkg/audit"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("=== Starting Database Migration ===")

	err := db.AutoMigrate(
		// 审计日志表
		&audit.AuditLog{},
		// 业务表
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
		&inventoryDomain.Inventory{},
	)

	if err != nil {
		log.Printf("✗ Migration failed: %v", err)
		return err
	}

	log.Println("✓ All tables migrated successfully")
	return nil
}
