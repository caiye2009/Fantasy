package config

import (
	"gorm.io/gorm"

	"back/pkg/auth"
	"back/pkg/es"

	// Auth
	authApp "back/internal/auth/application"

	// 基础数据模块
	clientApp "back/internal/client/application"
	clientInfra "back/internal/client/infra"
	supplierApp "back/internal/supplier/application"
	supplierInfra "back/internal/supplier/infra"
	materialApp "back/internal/material/application"
	materialInfra "back/internal/material/infra"
	processApp "back/internal/process/application"
	processInfra "back/internal/process/infra"
	productApp "back/internal/product/application"
	productInfra "back/internal/product/infra"
	userApp "back/internal/user/application"
	userInfra "back/internal/user/infra"
	userDeptApp "back/internal/user_dept/application"
	userDeptInfra "back/internal/user_dept/infra"
	userRoleApp "back/internal/user_role/application"
	userRoleInfra "back/internal/user_role/infra"

	// 关联配置模块
	productMaterialApp "back/internal/product_material/application"
	productMaterialInfra "back/internal/product_material/infra"
	productProcessApp "back/internal/product_process/application"
	productProcessInfra "back/internal/product_process/infra"
	materialQuoteApp "back/internal/material_quote/application"
	materialQuoteInfra "back/internal/material_quote/infra"
	processQuoteApp "back/internal/process_quote/application"
	processQuoteInfra "back/internal/process_quote/infra"

	// 订单主模块
	orderApp "back/internal/order/application"
	orderInfra "back/internal/order/infra"
	orderProductApp "back/internal/order_product/application"
	orderProductInfra "back/internal/order_product/infra"

	// 订单执行模块
	orderMaterialApp "back/internal/order_material/application"
	orderMaterialInfra "back/internal/order_material/infra"
	orderProcessApp "back/internal/order_process/application"
	orderProcessInfra "back/internal/order_process/infra"
	purchaseOrderApp "back/internal/purchase_order/application"
	purchaseOrderInfra "back/internal/purchase_order/infra"
	materialAllocationApp "back/internal/material_allocation/application"
	materialAllocationInfra "back/internal/material_allocation/infra"
	productionOrderApp "back/internal/production_order/application"
	productionOrderInfra "back/internal/production_order/infra"
	productAllocationApp "back/internal/product_allocation/application"
	productAllocationInfra "back/internal/product_allocation/infra"

	// 订单辅助模块
	orderParticipantApp "back/internal/order_participant/application"
	orderParticipantInfra "back/internal/order_participant/infra"
	orderEventApp "back/internal/order_event/application"
	orderEventInfra "back/internal/order_event/infra"

	//搜索
	"github.com/elastic/go-elasticsearch/v8"
	searchApp "back/internal/search/application"
	searchInfra "back/internal/search/infra"
)

type Services struct {
	// Auth
	Auth *authApp.AuthService

	// 基础数据
	Client   *clientApp.ClientService
	Supplier *supplierApp.SupplierService
	Material *materialApp.MaterialService
	Process  *processApp.ProcessService
	Product  *productApp.ProductService
	User     *userApp.UserService
	UserDept *userDeptApp.UserDeptService
	UserRole *userRoleApp.UserRoleService

	// 关联配置
	ProductMaterial *productMaterialApp.ProductMaterialService
	ProductProcess  *productProcessApp.ProductProcessService
	MaterialQuote   *materialQuoteApp.MaterialQuoteService
	ProcessQuote    *processQuoteApp.ProcessQuoteService

	// 订单主表
	Order        *orderApp.OrderService
	OrderProduct *orderProductApp.OrderProductService

	// 订单执行
	OrderMaterial      *orderMaterialApp.OrderMaterialService
	OrderProcess       *orderProcessApp.OrderProcessService
	PurchaseOrder      *purchaseOrderApp.PurchaseOrderService
	MaterialAllocation *materialAllocationApp.MaterialAllocationService
	ProductionOrder    *productionOrderApp.ProductionOrderService
	ProductAllocation  *productAllocationApp.ProductAllocationService

	// 订单辅助
	OrderParticipant *orderParticipantApp.OrderParticipantService
	OrderEvent       *orderEventApp.OrderEventService

	//搜索服务
	Search *searchApp.SearchService
}

func InitServices(db *gorm.DB, esClient *elasticsearch.Client, esSync *es.ESSync, searchRegistry *searchInfra.DomainAwareRegistry, jwtWang *auth.JWTWang, whitelistManager *auth.WhitelistManager) *Services {
	// ========== 基础数据模块 ==========
	clientRepo := clientInfra.NewClientRepo(db)
	clientService := clientApp.NewClientService(clientRepo)

	supplierRepo := supplierInfra.NewSupplierRepo(db)
	supplierService := supplierApp.NewSupplierService(supplierRepo)

	materialRepo := materialInfra.NewMaterialRepo(db)
	materialService := materialApp.NewMaterialService(materialRepo)

	processRepo := processInfra.NewProcessRepo(db)
	processService := processApp.NewProcessService(processRepo)

	productRepo := productInfra.NewProductRepo(db)
	productService := productApp.NewProductService(productRepo)

	userRepo := userInfra.NewUserRepo(db)
	userService := userApp.NewUserService(userRepo)

	userDeptRepo := userDeptInfra.NewUserDeptRepo(db)
	userDeptService := userDeptApp.NewUserDeptService(userDeptRepo)

	userRoleRepo := userRoleInfra.NewUserRoleRepo(db)
	userRoleService := userRoleApp.NewUserRoleService(userRoleRepo)

	// ========== Auth ==========
	authService := authApp.NewAuthService(userService, jwtWang, whitelistManager)

	// ========== 关联配置模块 ==========
	productMaterialRepo := productMaterialInfra.NewProductMaterialRepo(db)
	productMaterialService := productMaterialApp.NewProductMaterialService(productMaterialRepo)

	productProcessRepo := productProcessInfra.NewProductProcessRepo(db)
	productProcessService := productProcessApp.NewProductProcessService(productProcessRepo)

	materialQuoteRepo := materialQuoteInfra.NewMaterialQuoteRepo(db)
	materialQuoteService := materialQuoteApp.NewMaterialQuoteService(materialQuoteRepo)

	processQuoteRepo := processQuoteInfra.NewProcessQuoteRepo(db)
	processQuoteService := processQuoteApp.NewProcessQuoteService(processQuoteRepo)

	// ========== 订单主模块 ==========
	orderRepo := orderInfra.NewOrderRepo(db)
	orderService := orderApp.NewOrderService(orderRepo)

	orderProductRepo := orderProductInfra.NewOrderProductRepo(db)
	orderProductService := orderProductApp.NewOrderProductService(orderProductRepo)

	// ========== 订单执行模块 ==========
	orderMaterialRepo := orderMaterialInfra.NewOrderMaterialRepo(db)
	orderMaterialService := orderMaterialApp.NewOrderMaterialService(orderMaterialRepo)

	orderProcessRepo := orderProcessInfra.NewOrderProcessRepo(db)
	orderProcessService := orderProcessApp.NewOrderProcessService(orderProcessRepo)

	purchaseOrderRepo := purchaseOrderInfra.NewPurchaseOrderRepo(db)
	purchaseOrderService := purchaseOrderApp.NewPurchaseOrderService(purchaseOrderRepo)

	materialAllocationRepo := materialAllocationInfra.NewMaterialAllocationRepo(db)
	materialAllocationService := materialAllocationApp.NewMaterialAllocationService(materialAllocationRepo)

	productionOrderRepo := productionOrderInfra.NewProductionOrderRepo(db)
	productionOrderService := productionOrderApp.NewProductionOrderService(productionOrderRepo)

	productAllocationRepo := productAllocationInfra.NewProductAllocationRepo(db)
	productAllocationService := productAllocationApp.NewProductAllocationService(productAllocationRepo)

	// ========== 订单辅助模块 ==========
	orderParticipantRepo := orderParticipantInfra.NewOrderParticipantRepo(db)
	orderParticipantService := orderParticipantApp.NewOrderParticipantService(orderParticipantRepo)

	orderEventRepo := orderEventInfra.NewOrderEventRepo(db)
	orderEventService := orderEventApp.NewOrderEventService(orderEventRepo)

	//搜索
	searchRepo := searchInfra.NewESSearchRepository(esClient)
	searchService := searchApp.NewSearchService(searchRegistry, searchRepo)

	return &Services{
		Auth:               authService,
		Client:             clientService,
		Supplier:           supplierService,
		Material:           materialService,
		Process:            processService,
		Product:            productService,
		User:               userService,
		UserDept:           userDeptService,
		UserRole:           userRoleService,
		ProductMaterial:    productMaterialService,
		ProductProcess:     productProcessService,
		MaterialQuote:      materialQuoteService,
		ProcessQuote:       processQuoteService,
		Order:              orderService,
		OrderProduct:       orderProductService,
		OrderMaterial:      orderMaterialService,
		OrderProcess:       orderProcessService,
		PurchaseOrder:      purchaseOrderService,
		MaterialAllocation: materialAllocationService,
		ProductionOrder:    productionOrderService,
		ProductAllocation:  productAllocationService,
		OrderParticipant:   orderParticipantService,
		OrderEvent:         orderEventService,
		Search:                searchService,
	}
}
