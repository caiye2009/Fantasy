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

	vendor, err := h.vendorService.Create(c.Request.Context(), &req)
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
	
	vendor, err := h.vendorService.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
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
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} map[string]interface{} "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /vendor [get]
func (h *VendorHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	vendors, err := h.vendorService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   len(vendors),
		"vendors": vendors,
	})
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

	if err := h.vendorService.Update(c.Request.Context(), uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
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

	if err := h.vendorService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}