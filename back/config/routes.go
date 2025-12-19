package config

import (
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "back/docs" // 导入生成的 docs

	"back/pkg/auth"
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
)

func InitRoutes(authWang *auth.AuthWang, services *Services) *gin.Engine {
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

	// 公开路由
	public := api.Group("")
	{
		authInterfaces.RegisterAuthHandlers(public, services.Auth)
	}

	// 受保护路由
	protected := api.Group("")
	protected.Use(authWang.AuthMiddleware())
	{
		userInterfaces.RegisterUserHandlers(protected, services.User)
		userInterfaces.RegisterDepartmentHandlers(protected, services.Department)
		userInterfaces.RegisterRoleHandlers(protected, services.Role)
		supplierInterfaces.RegisterSupplierHandlers(protected, services.Supplier)
		clientInterfaces.RegisterClientHandlers(protected, services.Client)
		materialInterfaces.RegisterMaterialHandlers(protected, services.Material)
		processInterfaces.RegisterProcessHandlers(protected, services.Process)
		pricingInterfaces.RegisterMaterialPriceHandlers(protected, services.MaterialPrice)
		pricingInterfaces.RegisterProcessPriceHandlers(protected, services.ProcessPrice)
		productInterfaces.RegisterProductHandlers(protected, services.Product, services.ProductCostCalculator, services.ProductPrice)
		planInterfaces.RegisterPlanHandlers(protected, services.Plan)
		orderInterfaces.RegisterOrderHandlers(protected, services.Order)
		searchInterfaces.RegisterSearchHandlers(protected, services.Search)

		// Analytics
		returnAnalysisHandler := analyticsInterfaces.NewReturnAnalysisHandler(services.ReturnAnalysis)
		returnAnalysisHandler.RegisterRoutes(protected)
	}

	return router
}