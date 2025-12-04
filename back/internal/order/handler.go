package order

import (
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *OrderService
}

func NewOrderHandler(orderService *OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// Create godoc
// @Summary      创建订单
// @Description  创建新的订单信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        request body CreateOrderRequest true "订单信息"
// @Success      200 {object} Order "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

// Get godoc
// @Summary      获取订单详情
// @Description  根据订单ID获取订单详细信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Success      200 {object} Order "获取成功"
// @Failure      404 {object} map[string]string "订单不存在"
// @Security     Bearer
// @Router       /order/{id} [get]
func (h *OrderHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	order, err := h.orderService.Get(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// List godoc
// @Summary      获取订单列表
// @Description  获取所有订单列表
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} map[string]interface{} "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order [get]
func (h *OrderHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	orders, err := h.orderService.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  len(orders),
		"orders": orders,
	})
}

// Update godoc
// @Summary      更新订单
// @Description  根据订单ID更新订单信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body UpdateOrderRequest true "更新的订单信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id} [put]
func (h *OrderHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.Update(c.Request.Context(), uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete godoc
// @Summary      删除订单
// @Description  根据订单ID删除订单
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Success      200 {object} map[string]string "删除成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id} [delete]
func (h *OrderHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.orderService.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}