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

// Create godoc
// @Summary      创建供应商
// @Description  创建新的供应商信息
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Param        request body CreateVendorRequest true "供应商信息"
// @Success      200 {object} Vendor "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /vendor [post]
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

// Get godoc
// @Summary      获取供应商详情
// @Description  根据供应商ID获取供应商详细信息
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Param        id path int true "供应商ID"
// @Success      200 {object} Vendor "获取成功"
// @Failure      404 {object} map[string]string "供应商不存在"
// @Security     Bearer
// @Router       /vendor/{id} [get]
func (h *VendorHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	vendor, err := h.vendorService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, vendor)
}

// List godoc
// @Summary      获取供应商列表
// @Description  获取所有供应商列表
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Success      200 {array} Vendor "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /vendor [get]
func (h *VendorHandler) List(c *gin.Context) {
	list, err := h.vendorService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Update godoc
// @Summary      更新供应商信息
// @Description  根据供应商ID更新供应商信息
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Param        id path int true "供应商ID"
// @Param        request body UpdateVendorRequest true "更新的供应商信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /vendor/{id} [put]
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

// Delete godoc
// @Summary      删除供应商
// @Description  根据供应商ID删除供应商
// @Tags         供应商管理
// @Accept       json
// @Produce      json
// @Param        id path int true "供应商ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /vendor/{id} [delete]
func (h *VendorHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.vendorService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}