package search

import (
	"github.com/gin-gonic/gin"
	"back/pkg/auth"
)

func Register(rg *gin.RouterGroup, service *SearchService, authWang *auth.AuthWang) {
	handler := NewSearchHandler(service)

	search := rg.Group("/search")

	search.Use(authWang.AuthMiddleware())
	{
		search.POST("", handler.Search)
		search.POST("/:index", handler.SearchByIndex)
		search.GET("/indices", handler.GetIndices)
	}
}