package interfaces

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"back/pkg/audit"
	"back/internal/search/application"
)

// SearchHandler 搜索 Handler（全新实现）
type SearchHandler struct {
	service *application.SearchService
}

// NewSearchHandler 创建 Handler
func NewSearchHandler(service *application.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// Search 搜索
// @Summary      搜索
// @Description  基于配置的高级搜索，支持全文搜索、筛选、聚合
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

	// 执行搜索
	result, err := h.service.Search(c.Request.Context(), &req)
	if err != nil {
		// 区分客户端错误和服务器错误
		if isClientError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "搜索失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetIndices 获取所有支持搜索的索引
// @Summary      获取索引列表
// @Description  获取所有支持搜索的索引列表
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

// isClientError 判断是否为客户端错误
func isClientError(err error) bool {
	errMsg := err.Error()
	return strings.Contains(errMsg, "unsupported index") ||
		strings.Contains(errMsg, "is not filterable") ||
		strings.Contains(errMsg, "is not aggregable") ||
		strings.Contains(errMsg, "must be") ||
		strings.Contains(errMsg, "required")
}

// RegisterSearchHandlers 注册搜索路由
func RegisterSearchHandlers(router *gin.RouterGroup, service *application.SearchService) {
	handler := NewSearchHandler(service)

	searchGroup := router.Group("/search")
	{
		// 搜索是查询操作，跳过审计
		searchGroup.POST("", audit.Skip(), handler.Search)
		searchGroup.GET("/indices", handler.GetIndices)
	}
}