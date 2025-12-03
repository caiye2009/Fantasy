package order

import (
	"back/pkg/auth"
	"github.com/gin-gonic/gin"
)

func Register(rg *gin.RouterGroup, service *OrderService, authWang *auth.AuthWang) {
	handler := NewOrderHandler(service)

	order := rg.Group("/order")


	order.Use(authWang.AuthMiddleware())
	{
		order.GET("/list", handler.List)
		order.POST("/create", handler.Create)
		order.GET("/:id", handler.Get)
		order.PUT("/:id", handler.Update)
		order.DELETE("/:id", handler.Delete)
	}
}