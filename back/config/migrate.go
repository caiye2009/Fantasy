package config

import (
	"gorm.io/gorm"
	"log"

	// 基础数据表
	clientDomain "back/internal/client/domain"
	supplierDomain "back/internal/supplier/domain"
	materialDomain "back/internal/material/domain"
	processDomain "back/internal/process/domain"
	productDomain "back/internal/product/domain"
	userDomain "back/internal/user/domain"
	userDeptDomain "back/internal/user_dept/domain"
	userRoleDomain "back/internal/user_role/domain"

	// 关联配置表
	productFabricDomain "back/internal/product_fabric/domain"
	productProcessDomain "back/internal/product_process/domain"
	materialQuoteDomain "back/internal/material_quote/domain"
	processQuoteDomain "back/internal/process_quote/domain"

	// 订单主表
	orderDomain "back/internal/order/domain"
	orderProductDomain "back/internal/order_product/domain"

	// 订单执行表
	orderMaterialDomain "back/internal/order_material/domain"
	orderProcessDomain "back/internal/order_process/domain"
	purchaseOrderDomain "back/internal/purchase_order/domain"
	materialAllocationDomain "back/internal/material_allocation/domain"
	productionOrderDomain "back/internal/production_order/domain"
	productAllocationDomain "back/internal/product_allocation/domain"

	// 订单辅助表
	orderParticipantDomain "back/internal/order_participant/domain"
	orderEventDomain "back/internal/order_event/domain"

	// material_shrinkage
	materialShrinkageDomain "back/internal/material_shrinkage/domain"

	"back/pkg/audit"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("=== Starting Database Migration ===")

	err := db.AutoMigrate(
		// 审计日志表
		&audit.Audit{},

		// 基础数据表（8个）
		&clientDomain.Client{},
		&supplierDomain.Supplier{},
		&materialDomain.Material{},
		&processDomain.Process{},
		&productDomain.Product{},
		&userDomain.User{},
		&userDeptDomain.UserDept{},
		&userRoleDomain.UserRole{},

		// 关联配置表（4个）
		&productFabricDomain.ProductFabric{},
		&productProcessDomain.ProductProcess{},
		&materialQuoteDomain.MaterialQuote{},
		&processQuoteDomain.ProcessQuote{},

		// 订单主表（2个）
		&orderDomain.Order{},
		&orderProductDomain.OrderProduct{},

		// 订单执行表（6个）
		&orderMaterialDomain.OrderMaterial{},
		&orderProcessDomain.OrderProcess{},
		&purchaseOrderDomain.PurchaseOrder{},
		&materialAllocationDomain.MaterialAllocation{},
		&productionOrderDomain.ProductionOrder{},
		&productAllocationDomain.ProductAllocation{},

		// 订单辅助表（2个）
		&orderParticipantDomain.OrderParticipant{},
		&orderEventDomain.OrderEvent{},

		// material_shrinkage
		&materialShrinkageDomain.MaterialShrinkage{},
	)

	if err != nil {
		log.Printf("✗ Migration failed: %v", err)
		return err
	}

	log.Println("✓ All tables migrated successfully")
	return nil
}
