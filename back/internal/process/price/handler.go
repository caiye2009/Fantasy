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

func (h *ProcessPriceHandler) Quote(c *gin.Context) {
	var req struct {
		VendorID  uint    `json:"vendor_id" binding:"required"`
		ProcessID uint    `json:"process_id" binding:"required"`
		Price     float64 `json:"price" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.priceService.Quote(req.VendorID, req.ProcessID, req.Price); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "quoted"})
}

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

func (h *ProcessPriceHandler) GetHistory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	processID := uint(id)

	history, err := h.priceService.GetHistory(processID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}