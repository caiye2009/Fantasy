package auth

import "github.com/gin-gonic/gin"

func RegisterPublic(rg *gin.RouterGroup, service *AuthService) {
	handler := NewAuthHandler(service)

	auth := rg.Group("/auth")
	{
		auth.POST("/login", handler.Login)
	}
}