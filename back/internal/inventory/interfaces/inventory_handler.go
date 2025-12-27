package interfaces

import (
	"net/http"
	"strconv"

	"back/internal/inventory/application"
	"back/pkg/endpoint"
	"github.com/gin-gonic/gin"
)

// InventoryHandler 库存 Handler
type InventoryHandler struct {
	service *application.InventoryService
}

// NewInventoryHandler 创建 Handler
func NewInventoryHandler(service *application.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

// CreateInventory 创建库存
// @Summary      创建库存
// @Description  创建新的库存记录
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateInventoryRequest true "创建库存请求"
// @Success      200 {object} application.InventoryResponse "创建成功"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Security     Bearer
// @Router       /inventory [post]
func (h *InventoryHandler) CreateInventory(c *gin.Context) {
	var req application.CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	resp, err := h.service.CreateInventory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetInventory 获取库存详情
// @Summary      获取库存详情
// @Description  根据ID获取库存详情
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        id path int true "库存ID"
// @Success      200 {object} application.InventoryResponse "库存详情"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Failure      404 {object} map[string]interface{} "库存不存在"
// @Security     Bearer
// @Router       /inventory/{id} [get]
func (h *InventoryHandler) GetInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	resp, err := h.service.GetInventory(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetInventoriesByProductID 根据产品ID获取库存列表
// @Summary      根据产品ID获取库存列表
// @Description  根据产品ID获取该产品的所有库存记录
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        productId query int true "产品ID"
// @Success      200 {object} application.InventoryListResponse "库存列表"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Security     Bearer
// @Router       /inventory/product [get]
func (h *InventoryHandler) GetInventoriesByProductID(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Query("productId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的产品ID"})
		return
	}

	resp, err := h.service.GetInventoriesByProductID(c.Request.Context(), uint(productID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetInventoryByBatchID 根据批次ID获取库存
// @Summary      根据批次ID获取库存
// @Description  根据批次ID获取库存记录
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        batchId query string true "批次ID"
// @Success      200 {object} application.InventoryResponse "库存详情"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Failure      404 {object} map[string]interface{} "库存不存在"
// @Security     Bearer
// @Router       /inventory/batch [get]
func (h *InventoryHandler) GetInventoryByBatchID(c *gin.Context) {
	batchID := c.Query("batchId")
	if batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "批次ID不能为空"})
		return
	}

	resp, err := h.service.GetInventoryByBatchID(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListInventories 获取库存列表
// @Summary      获取库存列表
// @Description  获取所有库存列表
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} application.InventoryListResponse "库存列表"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Security     Bearer
// @Router       /inventory/list [get]
func (h *InventoryHandler) ListInventories(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.service.ListInventories(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// ListInventoriesByCategory 根据类别获取库存列表
// @Summary      根据类别获取库存列表
// @Description  根据类别获取库存列表
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        category query string true "类别"
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} application.InventoryListResponse "库存列表"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Security     Bearer
// @Router       /inventory/category [get]
func (h *InventoryHandler) ListInventoriesByCategory(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "类别不能为空"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	resp, err := h.service.ListInventoriesByCategory(c.Request.Context(), category, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateInventory 更新库存
// @Summary      更新库存
// @Description  更新库存信息
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        request body application.UpdateInventoryRequest true "更新库存请求"
// @Success      200 {object} application.InventoryResponse "更新成功"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Failure      404 {object} map[string]interface{} "库存不存在"
// @Security     Bearer
// @Router       /inventory [put]
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	var req application.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	resp, err := h.service.UpdateInventory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateQuantity 更新数量
// @Summary      更新库存数量
// @Description  更新库存数量
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        request body application.UpdateQuantityRequest true "更新数量请求"
// @Success      200 {object} application.InventoryResponse "更新成功"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Failure      404 {object} map[string]interface{} "库存不存在"
// @Security     Bearer
// @Router       /inventory/quantity [put]
func (h *InventoryHandler) UpdateQuantity(c *gin.Context) {
	var req application.UpdateQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	resp, err := h.service.UpdateQuantity(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeductInventory 扣减库存
// @Summary      扣减库存
// @Description  扣减指定数量的库存
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        request body application.DeductInventoryRequest true "扣减库存请求"
// @Success      200 {object} application.InventoryResponse "扣减成功"
// @Failure      400 {object} map[string]interface{} "参数错误或库存不足"
// @Failure      404 {object} map[string]interface{} "库存不存在"
// @Security     Bearer
// @Router       /inventory/deduct [post]
func (h *InventoryHandler) DeductInventory(c *gin.Context) {
	var req application.DeductInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	resp, err := h.service.DeductInventory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AddInventory 增加库存
// @Summary      增加库存
// @Description  增加指定数量的库存
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        request body application.AddInventoryRequest true "增加库存请求"
// @Success      200 {object} application.InventoryResponse "增加成功"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Failure      404 {object} map[string]interface{} "库存不存在"
// @Security     Bearer
// @Router       /inventory/add [post]
func (h *InventoryHandler) AddInventory(c *gin.Context) {
	var req application.AddInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	resp, err := h.service.AddInventory(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteInventory 删除库存
// @Summary      删除库存
// @Description  删除库存记录
// @Tags         库存管理
// @Accept       json
// @Produce      json
// @Param        id path int true "库存ID"
// @Success      200 {object} map[string]interface{} "删除成功"
// @Failure      400 {object} map[string]interface{} "参数错误"
// @Failure      404 {object} map[string]interface{} "库存不存在"
// @Security     Bearer
// @Router       /inventory/{id} [delete]
func (h *InventoryHandler) DeleteInventory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.service.DeleteInventory(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetRoutes 返回路由定义
func (h *InventoryHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/inventory", Handler: h.CreateInventory, Domain: "inventory", Action: "create"},
		{Method: "GET", Path: "/inventory/:id", Handler: h.GetInventory, Domain: "inventory", Action: "read"},
		{Method: "GET", Path: "/inventory/product", Handler: h.GetInventoriesByProductID, Domain: "inventory", Action: "read"},
		{Method: "GET", Path: "/inventory/batch", Handler: h.GetInventoryByBatchID, Domain: "inventory", Action: "read"},
		{Method: "GET", Path: "/inventory/list", Handler: h.ListInventories, Domain: "inventory", Action: "read"},
		{Method: "GET", Path: "/inventory/category", Handler: h.ListInventoriesByCategory, Domain: "inventory", Action: "read"},
		{Method: "PUT", Path: "/inventory", Handler: h.UpdateInventory, Domain: "inventory", Action: "update"},
		{Method: "PUT", Path: "/inventory/quantity", Handler: h.UpdateQuantity, Domain: "inventory", Action: "update"},
		{Method: "POST", Path: "/inventory/deduct", Handler: h.DeductInventory, Domain: "inventory", Action: "update"},
		{Method: "POST", Path: "/inventory/add", Handler: h.AddInventory, Domain: "inventory", Action: "update"},
		{Method: "DELETE", Path: "/inventory/:id", Handler: h.DeleteInventory, Domain: "inventory", Action: "delete"},
	}
}
