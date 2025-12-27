// back/config/services.go
package config

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"back/pkg/auth"
	"back/pkg/es"

	// Auth
	authApp "back/internal/auth/application"

	// Supplier
	supplierApp "back/internal/supplier/application"
	supplierInfra "back/internal/supplier/infra"

	// Client
	clientApp "back/internal/client/application"
	clientInfra "back/internal/client/infra"

	// User
	userApp "back/internal/user/application"
	userInfra "back/internal/user/infra"

	// Material
	materialApp "back/internal/material/application"
	materialInfra "back/internal/material/infra"

	// Process
	processApp "back/internal/process/application"
	processInfra "back/internal/process/infra"

	// Pricing
	pricingApp "back/internal/pricing/application"
	pricingInfra "back/internal/pricing/infra"

	// Product
	productApp "back/internal/product/application"
	productInfra "back/internal/product/infra"

	// Plan
	planApp "back/internal/plan/application"
	planInfra "back/internal/plan/infra"

	// Order
	orderApp "back/internal/order/application"
	orderInfra "back/internal/order/infra"

	// Search
	searchApp "back/internal/search/application"
	searchInfra "back/internal/search/infra"

	// Analytics
	analyticsApp "back/internal/analytics/application"
	analyticsInfra "back/internal/analytics/infra"

	// Permission
	permissionApp "back/internal/user/permission/application"

	// Inventory
	inventoryApp "back/internal/inventory/application"
	inventoryInfra "back/internal/inventory/infra"

	// Casbin
	casbinPkg "back/pkg/casbin"
)

type Services struct {
	// Auth
	Auth *authApp.AuthService

	// Core Entities
	Supplier   *supplierApp.SupplierService
	Client     *clientApp.ClientService
	User       *userApp.UserService
	Department *userApp.DepartmentService
	Role       *userApp.RoleService
	Material   *materialApp.MaterialService
	Process    *processApp.ProcessService

	// Pricing
	MaterialPrice *pricingApp.MaterialPriceService
	ProcessPrice  *pricingApp.ProcessPriceService

	// Product
	Product               *productApp.ProductService
	ProductCostCalculator *productApp.CostCalculator
	ProductPrice          *productApp.ProductPriceService

	// Plan & Order
	Plan  *planApp.PlanService
	Order *orderApp.OrderService

	// Search
	Search *searchApp.SearchService

	// Analytics
	ReturnAnalysis *analyticsApp.ReturnAnalysisService

	// Permission
	Permission *permissionApp.PermissionService

	// Inventory
	Inventory *inventoryApp.InventoryService
}

func InitServices(db *gorm.DB, rdb *redis.Client, esClient *elasticsearch.Client, jwtWang *auth.JWTWang, whitelistManager *auth.WhitelistManager, casbinManager *casbinPkg.Manager, esSync *es.ESSync, searchRegistry *searchInfra.DomainAwareRegistry) *Services {
	// ========== Supplier ==========
	supplierRepo := supplierInfra.NewSupplierRepo(db)
	supplierService := supplierApp.NewSupplierService(supplierRepo, esSync)

	// ========== Client ==========
	clientRepo := clientInfra.NewClientRepo(db)
	clientService := clientApp.NewClientService(clientRepo, esSync)

	// ========== User ==========
	userRepo := userInfra.NewUserRepo(db)
	userService := userApp.NewUserService(userRepo, whitelistManager)

	// ========== Department ==========
	departmentRepo := userInfra.NewDepartmentRepo(db)
	departmentService := userApp.NewDepartmentService(departmentRepo)

	// ========== Role ==========
	roleRepo := userInfra.NewRoleRepo(db)
	roleService := userApp.NewRoleService(roleRepo)

	// ========== Auth ==========
	authService := authApp.NewAuthService(userService, jwtWang, whitelistManager)

	// ========== Material ==========
	materialRepo := materialInfra.NewMaterialRepo(db)
	materialService := materialApp.NewMaterialService(materialRepo, esSync)

	// ========== Process ==========
	processRepo := processInfra.NewProcessRepo(db)
	processService := processApp.NewProcessService(processRepo, esSync)

	// ========== Pricing ==========
	supplierPriceRepo := pricingInfra.NewSupplierPriceRepo(db)
	priceCache := pricingInfra.NewPriceCacheImpl(rdb)

	materialPriceService := pricingApp.NewMaterialPriceService(
		supplierPriceRepo,
		priceCache,
		materialService,
		supplierService,
	)

	processPriceService := pricingApp.NewProcessPriceService(
		supplierPriceRepo,
		priceCache,
		processService,
		supplierService,
	)

	// ========== Product ==========
	productRepo := productInfra.NewProductRepo(db)
	productService := productApp.NewProductService(productRepo, esSync)

	productCostCalculator := productApp.NewCostCalculator(
		productRepo,
		materialService,
		processService,
		materialPriceService,
		processPriceService,
	)

	productPriceService := productApp.NewProductPriceService(
		productRepo,
		materialPriceService,
		processPriceService,
	)

	// ========== Plan ==========
	planRepo := planInfra.NewPlanRepo(db)
	planService := planApp.NewPlanService(planRepo, esSync)

	// ========== Order ==========
	orderRepo := orderInfra.NewOrderRepo(db)
	orderService := orderApp.NewOrderService(orderRepo, esSync)

	// ========== Search ==========
	searchRepo := searchInfra.NewESSearchRepository(esClient)
	searchService := searchApp.NewSearchService(searchRegistry, searchRepo)

	// ========== Analytics ==========
	returnAnalysisRepo := analyticsInfra.NewReturnAnalysisRepository(db)
	returnAnalysisService := analyticsApp.NewReturnAnalysisService(returnAnalysisRepo)

	// ========== Permission ==========
	permissionService := permissionApp.NewPermissionService(casbinManager)

	// ========== Inventory ==========
	inventoryRepo := inventoryInfra.NewInventoryRepo(db)
	inventoryService := inventoryApp.NewInventoryService(inventoryRepo)

	return &Services{
		Auth:                  authService,
		Supplier:              supplierService,
		Client:                clientService,
		User:                  userService,
		Department:            departmentService,
		Role:                  roleService,
		Material:              materialService,
		Process:               processService,
		MaterialPrice:         materialPriceService,
		ProcessPrice:          processPriceService,
		Product:               productService,
		ProductCostCalculator: productCostCalculator,
		ProductPrice:          productPriceService,
		Plan:                  planService,
		Order:                 orderService,
		Search:                searchService,
		ReturnAnalysis:        returnAnalysisService,
		Permission:            permissionService,
		Inventory:             inventoryService,
	}
}
