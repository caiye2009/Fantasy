package interfaces

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"back/internal/pricing/application"
)

// ProcessPriceHandler Process 报价 Handler
type ProcessPriceHandler struct {
	service *application.ProcessPriceService
}

// NewProcessPriceHandler 创建 Handler
func NewProcessPriceHandler(service *application.ProcessPriceService) *ProcessPriceHandler {
	return &ProcessPriceHandler{service: service}
}

// Quote godoc
// @Summary      厂商报价
// @Description  厂商为指定工序提交报价
// @Tags         工序价格管理
// @Accept       json
// @Produce      json
// @Param        request body application.QuoteRequest true "报价信息"
// @Success      200 {object} map[string]string "报价成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process/price/quote [post]
func (h *ProcessPriceHandler) Quote(c *gin.Context) {
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
// @Summary      查询工序价格
// @Description  查询指定工序的最低价格和最高价格
// @Tags         工序价格管理
// @Accept       json
// @Produce      json
// @Param        id path int true "工序ID"
// @Success      200 {object} map[string]interface{} "价格信息"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process/:id/price [get]
func (h *ProcessPriceHandler) GetPrice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	processID := uint(id)
	
	minPrice, err := h.service.GetMinPrice(c.Request.Context(), processID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	maxPrice, err := h.service.GetMaxPrice(c.Request.Context(), processID)
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
// @Description  查询指定工序的价格历史记录
// @Tags         工序价格管理
// @Accept       json
// @Produce      json
// @Param        id path int true "工序ID"
// @Success      200 {array} domain.VendorPrice "价格历史记录"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process/:id/price/history [get]
func (h *ProcessPriceHandler) GetHistory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	processID := uint(id)
	
	history, err := h.service.GetHistory(c.Request.Context(), processID, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, history)
}

// RegisterProcessPriceHandlers 注册路由
func RegisterProcessPriceHandlers(rg *gin.RouterGroup, service *application.ProcessPriceService) {
	handler := NewProcessPriceHandler(service)
	
	rg.POST("/process/price/quote", handler.Quote)
	rg.GET("/process/:id/price", handler.GetPrice)
	rg.GET("/process/:id/price/history", handler.GetHistory)
}