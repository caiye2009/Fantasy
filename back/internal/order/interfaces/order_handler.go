package interfaces

import (
	"errors"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"back/internal/order/application"
	"back/internal/order/domain"
)

// OrderHandler 订单 Handler
type OrderHandler struct {
	service *application.OrderService
}

// NewOrderHandler 创建 Handler
func NewOrderHandler(service *application.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// Create 创建订单
// @Summary      创建订单
// @Description  创建新的订单信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        request body application.CreateOrderRequest true "订单信息"
// @Success      200 {object} application.OrderResponse "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var req application.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNoDuplicate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "订单编号已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Get 获取订单详情
// @Summary      获取订单详情
// @Description  根据订单ID获取订单详细信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Success      200 {object} application.OrderResponse "获取成功"
// @Failure      404 {object} map[string]string "订单不存在"
// @Security     Bearer
// @Router       /order/{id} [get]
func (h *OrderHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// List 获取订单列表
// @Summary      获取订单列表
// @Description  获取所有订单列表
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        limit query int false "每页数量" default(10)
// @Param        offset query int false "偏移量" default(0)
// @Success      200 {object} application.OrderListResponse "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order [get]
func (h *OrderHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	
	resp, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, resp)
}

// Update 更新订单
// @Summary      更新订单
// @Description  根据订单ID更新订单信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body application.UpdateOrderRequest true "更新的订单信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id} [post]
func (h *OrderHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	
	var req application.UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err := h.service.Update(c.Request.Context(), uint(id), &req); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		if errors.Is(err, domain.ErrCannotConfirm) ||
			errors.Is(err, domain.ErrCannotStartProduction) ||
			errors.Is(err, domain.ErrCannotComplete) ||
			errors.Is(err, domain.ErrCannotUpdateCompleted) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除订单
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
	
	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "订单不存在"})
			return
		}
		if errors.Is(err, domain.ErrCannotDeleteCompleted) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "已完成的订单不能删除"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// RegisterOrderHandlers 注册路由
func RegisterOrderHandlers(rg *gin.RouterGroup, service *application.OrderService) {
	handler := NewOrderHandler(service)
	
	rg.POST("/order", handler.Create)
	rg.GET("/order/:id", handler.Get)
	rg.GET("/order", handler.List)
	rg.POST("/order/:id", handler.Update)
	rg.DELETE("/order/:id", handler.Delete)
}