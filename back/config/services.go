// back/config/services.go
package config

import (
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"github.com/elastic/go-elasticsearch/v8"

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

	// Plan & Order
	Plan  *planApp.PlanService
	Order *orderApp.OrderService

	// Search
	Search *searchApp.SearchService

	// Analytics
	ReturnAnalysis *analyticsApp.ReturnAnalysisService
}

func InitServices(db *gorm.DB, rdb *redis.Client, esClient *elasticsearch.Client, jwtWang *auth.JWTWang, esSync *es.ESSync) *Services {
	// ========== Supplier ==========
	supplierRepo := supplierInfra.NewSupplierRepo(db)
	supplierService := supplierApp.NewSupplierService(supplierRepo, esSync)

	// ========== Client ==========
	clientRepo := clientInfra.NewClientRepo(db)
	clientService := clientApp.NewClientService(clientRepo, esSync)

	// ========== User ==========
	userRepo := userInfra.NewUserRepo(db)
	userService := userApp.NewUserService(userRepo)

	// ========== Department ==========
	departmentRepo := userInfra.NewDepartmentRepo(db)
	departmentService := userApp.NewDepartmentService(departmentRepo)

	// ========== Role ==========
	roleRepo := userInfra.NewRoleRepo(db)
	roleService := userApp.NewRoleService(roleRepo)

	// ========== Auth ==========
	jwtAdapter := authApp.NewJWTWangAdapter(jwtWang)
	authService := authApp.NewAuthService(userService, jwtAdapter)

	// ========== Material ==========
	materialRepo := materialInfra.NewMaterialRepo(db)
	materialService := materialApp.NewMaterialService(materialRepo, esSync)

	// ========== Process ==========
	processRepo := processInfra.NewProcessRepo(db)
	processService := processApp.NewProcessService(processRepo, esSync)

	// ========== Pricing（使用适配器） ==========
	supplierPriceRepo := pricingInfra.NewSupplierPriceRepo(db)
	priceCache := pricingInfra.NewPriceCacheImpl(rdb)

	// 创建 Supplier 适配器
	supplierAdapter := supplierApp.NewSupplierServiceAdapter(supplierService)

	materialPriceService := pricingApp.NewMaterialPriceService(
		supplierPriceRepo,
		priceCache,
		materialService,
		supplierAdapter, // 使用适配器
	)

	processPriceService := pricingApp.NewProcessPriceService(
		supplierPriceRepo,
		priceCache,
		processService,
		supplierAdapter, // 使用适配器
	)

	// ========== Product（使用适配器） ==========
	productRepo := productInfra.NewProductRepo(db)
	productService := productApp.NewProductService(productRepo, esSync)

	// 创建 Material 和 Process 适配器
	materialAdapter := materialApp.NewMaterialServiceAdapter(materialService)
	processAdapter := processApp.NewProcessServiceAdapter(processService)

	productCostCalculator := productApp.NewCostCalculator(
		productRepo,
		materialAdapter,  // 使用适配器
		processAdapter,   // 使用适配器
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
	searchRepo := searchInfra.NewESSearchRepo(esClient)
	searchService := searchApp.NewSearchService(searchRepo)
	
	// ========== Analytics ==========
	returnAnalysisRepo := analyticsInfra.NewReturnAnalysisRepository(db)
	returnAnalysisService := analyticsApp.NewReturnAnalysisService(returnAnalysisRepo)
	
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
		Plan:                  planService,
		Order:                 orderService,
		Search:                searchService,
		ReturnAnalysis:        returnAnalysisService,
	}
}