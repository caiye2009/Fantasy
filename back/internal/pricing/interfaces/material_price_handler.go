package interfaces

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/pkg/audit"
	"back/pkg/endpoint"
	"back/internal/pricing/application"
)

// MaterialPriceHandler Material 报价 Handler
type MaterialPriceHandler struct {
	service *application.MaterialPriceService
}

// NewMaterialPriceHandler 创建 Handler
func NewMaterialPriceHandler(service *application.MaterialPriceService) *MaterialPriceHandler {
	return &MaterialPriceHandler{service: service}
}

// Quote godoc
// @Summary      厂商报价
// @Description  厂商为指定材料提交报价
// @Tags         材料价格管理
// @Accept       json
// @Produce      json
// @Param        request body application.QuoteRequest true "报价信息"
// @Success      200 {object} map[string]string "报价成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material/price/quote [post]
func (h *MaterialPriceHandler) Quote(c *gin.Context) {
	var req application.QuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.Quote(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "报价成功"})
}

// GetPrice godoc
// @Summary      查询材料价格
// @Description  查询指定材料的最低价格和最高价格
// @Tags         材料价格管理
// @Accept       json
// @Produce      json
// @Param        id path int true "材料ID"
// @Success      200 {object} map[string]interface{} "价格信息"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material/:id/price [get]
func (h *MaterialPriceHandler) GetPrice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	materialID := uint(id)
	
	minPrice, err := h.service.GetMinPrice(c.Request.Context(), materialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	maxPrice, err := h.service.GetMaxPrice(c.Request.Context(), materialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"min": minPrice,
		"max": maxPrice,
	})
}

// GetHistory godoc
// @Summary      查询价格历史
// @Description  查询指定材料的价格历史记录
// @Tags         材料价格管理
// @Accept       json
// @Produce      json
// @Param        id path int true "材料ID"
// @Success      200 {array} domain.PriceData "价格历史记录"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material/:id/price/history [get]
func (h *MaterialPriceHandler) GetHistory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	materialID := uint(id)
	
	history, err := h.service.GetHistory(c.Request.Context(), materialID, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, history)
}

// GetRoutes 获取路由定义
func (h *MaterialPriceHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{
			Method:      "POST",
			Path:        "/pricing/material",
			Handler:     h.Quote,
			Middlewares: []gin.HandlerFunc{audit.Mark("pricing", "materialUpsert")},
			Name:        "材料报价",
		},
		{
			Method:      "GET",
			Path:        "/pricing/material/:id",
			Handler:     h.GetPrice,
			Middlewares: nil,
			Name:        "获取材料价格",
		},
		{
			Method:      "GET",
			Path:        "/pricing/material/:id/history",
			Handler:     h.GetHistory,
			Middlewares: nil,
			Name:        "获取材料价格历史",
		},
	}
}