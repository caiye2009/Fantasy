package vendor

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(rg *gin.RouterGroup, service *VendorService, authWang *auth.AuthWang) {
	handler := NewVendorHandler(service)

	vendor := rg.Group("/vendor")

	vendor.Use(authWang.AuthMiddleware())
	{
		vendor.GET("/list", handler.List)
		vendor.POST("/create", handler.Create)
		vendor.GET("/:id", handler.Get)
		vendor.PUT("/:id", handler.Update)
		vendor.DELETE("/:id", handler.Delete)
	}
}