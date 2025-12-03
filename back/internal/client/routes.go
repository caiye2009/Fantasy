package client

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(rg *gin.RouterGroup, service *ClientService, authWang *auth.AuthWang) {
	handler := NewClientHandler(service)

	client := rg.Group("/client")

	client.Use(authWang.AuthMiddleware())
	{
		client.GET("/list", handler.List)
		client.POST("/create", handler.Create)
		client.GET("/:id", handler.Get)
		client.PUT("/:id", handler.Update)
		client.DELETE("/:id", handler.Delete)	
	}
}