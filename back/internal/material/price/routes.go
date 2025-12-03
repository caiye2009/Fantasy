package price

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(
	rg *gin.RouterGroup,
	priceService *MaterialPriceService,
	authWang *auth.AuthWang,
) {
	priceHandler := NewMaterialPriceHandler(priceService)

	material := rg.Group("/material")
	material.POST("/price", priceHandler.Quote)
	material.GET("/:id/price", priceHandler.GetPrice)
	material.GET("/:id/price/history", priceHandler.GetHistory)
	material.Use(authWang.AuthMiddleware())
	{

	}
}
