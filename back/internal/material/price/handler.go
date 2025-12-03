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

// 厂商报价
func (h *MaterialPriceHandler) Quote(c *gin.Context) {
	var req struct {
		VendorID   uint    `json:"vendor_id" binding:"required"`
		MaterialID uint    `json:"material_id" binding:"required"`
		Price      float64 `json:"price" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.priceService.Quote(req.VendorID, req.MaterialID, req.Price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "quoted"})
}

// 查询价格 (最低和最高)
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

// 查询价格历史
func (h *MaterialPriceHandler) GetHistory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	materialID := uint(id)

	history, err := h.priceService.GetHistory(materialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}