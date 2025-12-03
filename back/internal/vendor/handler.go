package vendor

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VendorHandler struct {
	vendorService *VendorService
}

func NewVendorHandler(vendorService *VendorService) *VendorHandler {
	return &VendorHandler{vendorService: vendorService}
}

func (h *VendorHandler) Create(c *gin.Context) {
	var req CreateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vendor, err := h.vendorService.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, vendor)
}

func (h *VendorHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	vendor, err := h.vendorService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, vendor)
}

func (h *VendorHandler) List(c *gin.Context) {
	list, err := h.vendorService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *VendorHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req UpdateVendorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.vendorService.Update(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *VendorHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.vendorService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}