package price

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(
	rg *gin.RouterGroup,
	priceService *ProcessPriceService,
	authWang *auth.AuthWang,
) {
	priceHandler := NewProcessPriceHandler(priceService)

	process := rg.Group("/process")
		process.POST("/price", priceHandler.Quote)
		process.GET("/:id/price", priceHandler.GetPrice)
		process.GET("/:id/price/history", priceHandler.GetHistory)
	process.Use(authWang.AuthMiddleware())
	{

	}
}
