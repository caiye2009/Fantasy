package material

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(rg *gin.RouterGroup, service *MaterialService, authWang *auth.AuthWang) {
	handler := NewMaterialHandler(service)

	material := rg.Group("/material")

	material.Use(authWang.AuthMiddleware())
	{
		material.GET("/list", handler.List)
		material.POST("/create", handler.Create)
		material.GET("/:id", handler.Get)
		material.PUT("/:id", handler.Update)
		material.DELETE("/:id", handler.Delete)
	}
}