package config

import (
	_ "back/docs" // 导入生成的 docs
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	// Auth
	authInterfaces "back/internal/auth/interfaces"

	// 基础数据模块
	clientInterfaces "back/internal/client/interfaces"
	supplierInterfaces "back/internal/supplier/interfaces"
	materialInterfaces "back/internal/material/interfaces"
	processInterfaces "back/internal/process/interfaces"
	productInterfaces "back/internal/product/interfaces"
	userInterfaces "back/internal/user/interfaces"
	userDeptInterfaces "back/internal/user_dept/interfaces"
	userRoleInterfaces "back/internal/user_role/interfaces"

	// 关联配置模块
	productFabricInterfaces "back/internal/product_fabric/interfaces"
	productProcessInterfaces "back/internal/product_process/interfaces"
	materialQuoteInterfaces "back/internal/material_quote/interfaces"
	processQuoteInterfaces "back/internal/process_quote/interfaces"

	// 订单主模块
	orderInterfaces "back/internal/order/interfaces"
	orderProductInterfaces "back/internal/order_product/interfaces"

	// 订单执行模块
	orderMaterialInterfaces "back/internal/order_material/interfaces"
	orderProcessInterfaces "back/internal/order_process/interfaces"
	purchaseOrderInterfaces "back/internal/purchase_order/interfaces"
	materialAllocationInterfaces "back/internal/material_allocation/interfaces"
	productionOrderInterfaces "back/internal/production_order/interfaces"
	productAllocationInterfaces "back/internal/product_allocation/interfaces"

	// 订单辅助模块
	orderParticipantInterfaces "back/internal/order_participant/interfaces"
	orderEventInterfaces "back/internal/order_event/interfaces"

	searchInterfaces "back/internal/search/interfaces"
	materialShrinkageInterfaces "back/internal/material_shrinkage/interfaces"

	fabricInterfaces "back/internal/fabric/interfaces"
	fabricMaterialInterfaces "back/internal/fabric_material/interfaces"
	fabricProcessinterfaces "back/internal/fabric_process/interfaces"

	"back/pkg/audit"
	"back/pkg/auth"
	"back/pkg/endpoint"
)

func InitRoutes(authWang *auth.AuthWang, services *Services, db *gorm.DB) *gin.Engine {
	// 注入 audit.Mark 函数以避免循环依赖
	endpoint.SetAuditMarker(audit.Mark)

	router := gin.Default()

	// CORS 配置
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

	// ========== 公开路由（不需要认证）==========
	public := api.Group("")
	{
		authHandler := authInterfaces.NewAuthHandler(services.Auth)
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/refresh", authHandler.RefreshToken)
	}

	// ========== 受保护路由 ==========
	protected := api.Group("")
	protected.Use(authWang.AuthMiddleware())
	protected.Use(audit.AuditMiddleware(db))
	{
		// Auth - 登出
		authHandler := authInterfaces.NewAuthHandler(services.Auth)
		protected.POST("/auth/logout", authHandler.Logout)

		// ========== 基础数据模块（6个）==========
		clientHandler := clientInterfaces.NewClientHandler(services.Client)
		endpoint.RegisterRoutes(protected, clientHandler.GetRoutes())

		supplierHandler := supplierInterfaces.NewSupplierHandler(services.Supplier)
		endpoint.RegisterRoutes(protected, supplierHandler.GetRoutes())

		materialHandler := materialInterfaces.NewMaterialHandler(services.Material)
		endpoint.RegisterRoutes(protected, materialHandler.GetRoutes())

		processHandler := processInterfaces.NewProcessHandler(services.Process)
		endpoint.RegisterRoutes(protected, processHandler.GetRoutes())

		productHandler := productInterfaces.NewProductHandler(services.Product)
		endpoint.RegisterRoutes(protected, productHandler.GetRoutes())

		userHandler := userInterfaces.NewUserHandler(services.User)
		endpoint.RegisterRoutes(protected, userHandler.GetRoutes())

		userDeptHandler := userDeptInterfaces.NewUserDeptHandler(services.UserDept)
		endpoint.RegisterRoutes(protected, userDeptHandler.GetRoutes())

		userRoleHandler := userRoleInterfaces.NewUserRoleHandler(services.UserRole)
		endpoint.RegisterRoutes(protected, userRoleHandler.GetRoutes())

		// ========== 关联配置模块（4个）==========
		productFabricHandler := productFabricInterfaces.NewProductFabricHandler(services.ProductFabric)
		endpoint.RegisterRoutes(protected, productFabricHandler.GetRoutes())

		productProcessHandler := productProcessInterfaces.NewProductProcessHandler(services.ProductProcess)
		endpoint.RegisterRoutes(protected, productProcessHandler.GetRoutes())

		materialQuoteHandler := materialQuoteInterfaces.NewMaterialQuoteHandler(services.MaterialQuote)
		endpoint.RegisterRoutes(protected, materialQuoteHandler.GetRoutes())

		processQuoteHandler := processQuoteInterfaces.NewProcessQuoteHandler(services.ProcessQuote)
		endpoint.RegisterRoutes(protected, processQuoteHandler.GetRoutes())

		// ========== 订单主模块（2个）==========
		orderHandler := orderInterfaces.NewOrderHandler(services.Order)
		endpoint.RegisterRoutes(protected, orderHandler.GetRoutes())

		orderProductHandler := orderProductInterfaces.NewOrderProductHandler(services.OrderProduct)
		endpoint.RegisterRoutes(protected, orderProductHandler.GetRoutes())

		// ========== 订单执行模块（6个）==========
		orderMaterialHandler := orderMaterialInterfaces.NewOrderMaterialHandler(services.OrderMaterial)
		endpoint.RegisterRoutes(protected, orderMaterialHandler.GetRoutes())

		orderProcessHandler := orderProcessInterfaces.NewOrderProcessHandler(services.OrderProcess)
		endpoint.RegisterRoutes(protected, orderProcessHandler.GetRoutes())

		purchaseOrderHandler := purchaseOrderInterfaces.NewPurchaseOrderHandler(services.PurchaseOrder)
		endpoint.RegisterRoutes(protected, purchaseOrderHandler.GetRoutes())

		materialAllocationHandler := materialAllocationInterfaces.NewMaterialAllocationHandler(services.MaterialAllocation)
		endpoint.RegisterRoutes(protected, materialAllocationHandler.GetRoutes())

		productionOrderHandler := productionOrderInterfaces.NewProductionOrderHandler(services.ProductionOrder)
		endpoint.RegisterRoutes(protected, productionOrderHandler.GetRoutes())

		productAllocationHandler := productAllocationInterfaces.NewProductAllocationHandler(services.ProductAllocation)
		endpoint.RegisterRoutes(protected, productAllocationHandler.GetRoutes())

		// ========== 订单辅助模块（2个）==========
		orderParticipantHandler := orderParticipantInterfaces.NewOrderParticipantHandler(services.OrderParticipant)
		endpoint.RegisterRoutes(protected, orderParticipantHandler.GetRoutes())

		orderEventHandler := orderEventInterfaces.NewOrderEventHandler(services.OrderEvent)
		endpoint.RegisterRoutes(protected, orderEventHandler.GetRoutes())

		searchHandler := searchInterfaces.NewSearchHandler(services.Search)
		endpoint.RegisterRoutes(protected, searchHandler.GetRoutes())

		materialShrinkageHandler := materialShrinkageInterfaces.NewMaterialShrinkageHandler(services.MaterialShrinkage)
		endpoint.RegisterRoutes(protected, materialShrinkageHandler.GetRoutes())

		fabricHandler := fabricInterfaces.NewFabricHandler(services.Fabric)
        endpoint.RegisterRoutes(protected, fabricHandler.GetRoutes())

        fabricMaterialHandler := fabricMaterialInterfaces.NewFabricMaterialHandler(services.FabricMaterial)
        endpoint.RegisterRoutes(protected, fabricMaterialHandler.GetRoutes())

        fabricProcessHandler := fabricProcessinterfaces.NewFabricProcessHandler(services.FabricProcess)
        endpoint.RegisterRoutes(protected, fabricProcessHandler.GetRoutes())

	}

	return router
}
