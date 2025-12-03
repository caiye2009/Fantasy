package config

import (
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "back/docs" // 导入生成的 docs

	"back/pkg/auth"
	"back/internal/vendor"
	"back/internal/client"
	"back/internal/user"
	internalAuth "back/internal/auth"
	"back/internal/material"
	materialPrice "back/internal/material/price"
	"back/internal/process"
	processPrice "back/internal/process/price"
	"back/internal/product"
	"back/internal/search"
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

	public := api.Group("")
	{
		internalAuth.RegisterPublic(public, services.Auth)
	}

	protected := api.Group("")
	protected.Use(authWang.AuthMiddleware())
	{
		user.Register(protected, services.User, authWang)
		vendor.Register(protected, services.Vendor, authWang)
		client.Register(protected, services.Client, authWang)
		material.Register(protected, services.Material, authWang)
		materialPrice.Register(protected, services.MaterialPrice, authWang)
		process.Register(protected, services.Process, authWang)
		processPrice.Register(protected, services.ProcessPrice, authWang)
		product.Register(protected, services.Product, authWang)
		search.Register(protected, services.Search, authWang)
	}

	return router
}