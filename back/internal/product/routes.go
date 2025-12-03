package product

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(rg *gin.RouterGroup, service *ProductService, authWang *auth.AuthWang) {
	handler := NewProductHandler(service)

	product := rg.Group("/product")

	product.Use(authWang.AuthMiddleware())
	{
		product.GET("/list", handler.List)
		product.POST("/create", handler.Create)
		product.GET("/:id", handler.Get)
		product.PUT("/:id", handler.Update)
		product.DELETE("/:id", handler.Delete)
		product.POST("/cost", handler.CalculateCost)
	}
}