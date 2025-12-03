package plan

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(rg *gin.RouterGroup, service *PlanService, authWang *auth.AuthWang) {
	handler := NewPlanHandler(service)
	
	plan := rg.Group("/plan")
	
	plan.Use(authWang.AuthMiddleware())
	{
		plan.GET("/list", handler.List)
		plan.POST("/create", handler.Create)
		plan.GET("/:id", handler.Get)
		plan.PUT("/:id", handler.Update)
		plan.DELETE("/:id", handler.Delete)	
	}
}