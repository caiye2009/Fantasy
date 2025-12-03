package user

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(rg *gin.RouterGroup, service *UserService, authWang *auth.AuthWang) {
	handler := NewUserHandler(service)

	user := rg.Group("/user")
	user.Use(authWang.AuthMiddleware())
	{
		user.POST("/create", handler.Create)
		user.GET("/list", handler.List)
		user.GET("/:id", handler.Get)
		user.PUT("/:id", handler.Update)
		user.DELETE("/:id", handler.Delete)
		user.PUT("/password", handler.ChangePassword)
	}
}