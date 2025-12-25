package config

import (
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	_ "back/docs" // 导入生成的 docs

	"back/pkg/auth"
	"back/pkg/audit"
	"back/pkg/endpoint"
	authInterfaces "back/internal/auth/interfaces"
	supplierInterfaces "back/internal/supplier/interfaces"
	clientInterfaces "back/internal/client/interfaces"
	userInterfaces "back/internal/user/interfaces"
	materialInterfaces "back/internal/material/interfaces"
	processInterfaces "back/internal/process/interfaces"
	pricingInterfaces "back/internal/pricing/interfaces"
	productInterfaces "back/internal/product/interfaces"
	planInterfaces "back/internal/plan/interfaces"
	orderInterfaces "back/internal/order/interfaces"
	searchInterfaces "back/internal/search/interfaces"
	analyticsInterfaces "back/internal/analytics/interfaces"
	permissionInterfaces "back/internal/permission/interfaces"
)

func InitRoutes(authWang *auth.AuthWang, services *Services, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Swagger 文档路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api/v1")

	// 公开路由（不需要认证）
	public := api.Group("")
	{
		// 登录和刷新token不需要认证
		authHandler := authInterfaces.NewAuthHandler(services.Auth)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", authHandler.RefreshToken)
	}

	// 受保护路由
	protected := api.Group("")
	protected.Use(authWang.AuthMiddleware())  // 1. 先认证鉴权
	protected.Use(audit.AuditMiddleware(db))  // 2. 再审计（auth > handler > audit）
	{
		// 登出需要认证
		authHandler := authInterfaces.NewAuthHandler(services.Auth)
		protected.POST("/auth/logout", authHandler.Logout)

		// User
		userHandler := userInterfaces.NewUserHandler(services.User)
		endpoint.RegisterRoutes(protected, userHandler.GetRoutes())

		// Department
		departmentHandler := userInterfaces.NewDepartmentHandler(services.Department)
		endpoint.RegisterRoutes(protected, departmentHandler.GetRoutes())

		// Role
		roleHandler := userInterfaces.NewRoleHandler(services.Role)
		endpoint.RegisterRoutes(protected, roleHandler.GetRoutes())

		// Supplier
		supplierHandler := supplierInterfaces.NewSupplierHandler(services.Supplier)
		endpoint.RegisterRoutes(protected, supplierHandler.GetRoutes())

		// Client
		clientHandler := clientInterfaces.NewClientHandler(services.Client)
		endpoint.RegisterRoutes(protected, clientHandler.GetRoutes())

		// Material
		materialHandler := materialInterfaces.NewMaterialHandler(services.Material)
		endpoint.RegisterRoutes(protected, materialHandler.GetRoutes())

		// Process
		processHandler := processInterfaces.NewProcessHandler(services.Process)
		endpoint.RegisterRoutes(protected, processHandler.GetRoutes())

		// Material Price
		materialPriceHandler := pricingInterfaces.NewMaterialPriceHandler(services.MaterialPrice)
		endpoint.RegisterRoutes(protected, materialPriceHandler.GetRoutes())

		// Process Price
		processPriceHandler := pricingInterfaces.NewProcessPriceHandler(services.ProcessPrice)
		endpoint.RegisterRoutes(protected, processPriceHandler.GetRoutes())

		// Product
		productHandler := productInterfaces.NewProductHandler(services.Product, services.ProductCostCalculator, services.ProductPrice)
		endpoint.RegisterRoutes(protected, productHandler.GetRoutes())

		// Plan
		planHandler := planInterfaces.NewPlanHandler(services.Plan)
		endpoint.RegisterRoutes(protected, planHandler.GetRoutes())

		// Order
		orderHandler := orderInterfaces.NewOrderHandler(services.Order)
		endpoint.RegisterRoutes(protected, orderHandler.GetRoutes())

		// Search
		searchHandler := searchInterfaces.NewSearchHandler(services.Search)
		endpoint.RegisterRoutes(protected, searchHandler.GetRoutes())

		// Analytics
		returnAnalysisHandler := analyticsInterfaces.NewReturnAnalysisHandler(services.ReturnAnalysis)
		endpoint.RegisterRoutes(protected, returnAnalysisHandler.GetRoutes())

		// Permission
		permissionHandler := permissionInterfaces.NewPermissionHandler(services.Permission)
		endpoint.RegisterRoutes(protected, permissionHandler.GetRoutes())
	}

	return router
}