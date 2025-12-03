package process

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(rg *gin.RouterGroup, service *ProcessService, authWang *auth.AuthWang) {
	handler := NewProcessHandler(service)

	process := rg.Group("/process")

	process.Use(authWang.AuthMiddleware())
	{
		process.GET("/list", handler.List)
		process.POST("/create", handler.Create)
		process.GET("/:id", handler.Get)
		process.PUT("/:id", handler.Update)
		process.DELETE("/:id", handler.Delete)
	}
}