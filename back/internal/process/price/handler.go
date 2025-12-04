package price

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProcessPriceHandler struct {
	priceService *ProcessPriceService
}

func NewProcessPriceHandler(priceService *ProcessPriceService) *ProcessPriceHandler {
	return &ProcessPriceHandler{priceService: priceService}
}

// Quote godoc
// @Summary      厂商报价
// @Description  厂商为指定工序提交报价
// @Tags         工序价格管理
// @Accept       json
// @Produce      json
// @Param        request body object{vendor_id=uint,process_id=uint,price=float64} true "报价信息"
// @Success      200 {object} map[string]string "报价成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process/price/quote [post]
func (h *ProcessPriceHandler) Quote(c *gin.Context) {
	var req struct {
		VendorID  uint    `json:"vendor_id" binding:"required" example:"1"`
		ProcessID uint    `json:"process_id" binding:"required" example:"1"`
		Price     float64 `json:"price" binding:"required,gt=0" example:"199.99"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.priceService.Quote(c.Request.Context(), req.VendorID, req.ProcessID, req.Price); err != nil {
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
// @Router       /process/price/{id} [get]
func (h *ProcessPriceHandler) GetPrice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	processID := uint(id)

	minPrice, err := h.priceService.GetMinPrice(processID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	maxPrice, err := h.priceService.GetMaxPrice(processID)
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
// @Success      200 {array} ProcessPrice "价格历史记录"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /process/price/{id}/history [get]
func (h *ProcessPriceHandler) GetHistory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	processID := uint(id)

	history, err := h.priceService.GetHistory(c.Request.Context(), processID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}