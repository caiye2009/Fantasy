package interfaces

import (
	"errors"
	"net/http"
	
	"github.com/gin-gonic/gin"
	
	"back/internal/search/application"
	"back/internal/search/domain"
)

// SearchHandler 搜索 Handler
type SearchHandler struct {
	service *application.SearchService
}

// NewSearchHandler 创建 Handler
func NewSearchHandler(service *application.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// Search 通用搜索
// @Summary      通用搜索
// @Description  执行跨多个索引的通用搜索
// @Tags         搜索
// @Accept       json
// @Produce      json
// @Param        request body application.SearchRequest true "搜索请求参数"
// @Success      200 {object} application.SearchResponse "搜索结果"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /search [post]
func (h *SearchHandler) Search(c *gin.Context) {
	var req application.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	results, err := h.service.Search(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidSize) ||
			errors.Is(err, domain.ErrSizeTooLarge) ||
			errors.Is(err, domain.ErrInvalidFrom) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}
	
	c.JSON(http.StatusOK, results)
}

// SearchByIndex 按索引搜索
// @Summary      按索引搜索
// @Description  在指定索引中执行搜索
// @Tags         搜索
// @Accept       json
// @Produce      json
// @Param        index path string true "索引名称" example:"products"
// @Param        request body application.SearchRequest true "搜索请求参数"
// @Success      200 {object} application.SearchResponse "搜索结果"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /search/{index} [post]
func (h *SearchHandler) SearchByIndex(c *gin.Context) {
	index := c.Param("index")
	
	var req application.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// 强制指定索引
	req.Indices = []string{index}
	
	results, err := h.service.Search(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidSize) ||
			errors.Is(err, domain.ErrSizeTooLarge) ||
			errors.Is(err, domain.ErrInvalidFrom) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}
	
	c.JSON(http.StatusOK, results)
}

// GetIndices 获取所有可用索引
// @Summary      获取所有可用索引
// @Description  获取系统中所有可搜索的索引列表
// @Tags         搜索
// @Accept       json
// @Produce      json
// @Success      200 {object} application.IndexListResponse "索引列表"
// @Security     Bearer
// @Router       /search/indices [get]
func (h *SearchHandler) GetIndices(c *gin.Context) {
	resp := h.service.GetIndices()
	c.JSON(http.StatusOK, resp)
}

// RegisterSearchHandlers 注册路由
func RegisterSearchHandlers(rg *gin.RouterGroup, service *application.SearchService) {
	handler := NewSearchHandler(service)
	
	rg.POST("/search", handler.Search)
	rg.POST("/search/:index", handler.SearchByIndex)
	rg.GET("/search/indices", handler.GetIndices)
}