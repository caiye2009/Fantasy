package price

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MaterialPriceHandler struct {
	priceService *MaterialPriceService
}

func NewMaterialPriceHandler(priceService *MaterialPriceService) *MaterialPriceHandler {
	return &MaterialPriceHandler{priceService: priceService}
}

// Quote godoc
// @Summary      厂商报价
// @Description  厂商为指定材料提交报价
// @Tags         材料价格管理
// @Accept       json
// @Produce      json
// @Param        request body object{vendor_id=uint,material_id=uint,price=float64} true "报价信息"
// @Success      200 {object} map[string]string "报价成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material/price/quote [post]
func (h *MaterialPriceHandler) Quote(c *gin.Context) {
	var req struct {
		VendorID   uint    `json:"vendor_id" binding:"required" example:"1"`
		MaterialID uint    `json:"material_id" binding:"required" example:"1"`
		Price      float64 `json:"price" binding:"required,gt=0" example:"99.99"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.priceService.Quote(c.Request.Context(), req.VendorID, req.MaterialID, req.Price); err != nil {
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
// @Router       /material/price/{id} [get]
func (h *MaterialPriceHandler) GetPrice(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	materialID := uint(id)

	minPrice, err := h.priceService.GetMinPrice(materialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	maxPrice, err := h.priceService.GetMaxPrice(materialID)
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
// @Success      200 {array} MaterialPrice "价格历史记录"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /material/price/{id}/history [get]
func (h *MaterialPriceHandler) GetHistory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	materialID := uint(id)

	history, err := h.priceService.GetHistory(c.Request.Context(), materialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}