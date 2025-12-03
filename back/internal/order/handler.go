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
// @Description  创建新的订单
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        request body Order true "订单信息"
// @Success      200 {object} Order "创建成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order [post]
func (h *OrderHandler) Create(c *gin.Context) {
	var o Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.Create(&o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, o)
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
	o, err := h.orderService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, o)
}

// List godoc
// @Summary      获取订单列表
// @Description  获取所有订单列表
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Success      200 {array} Order "获取成功"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order [get]
func (h *OrderHandler) List(c *gin.Context) {
	list, err := h.orderService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Update godoc
// @Summary      更新订单信息
// @Description  根据订单ID更新订单信息
// @Tags         订单管理
// @Accept       json
// @Produce      json
// @Param        id path int true "订单ID"
// @Param        request body map[string]interface{} true "更新的订单信息"
// @Success      200 {object} map[string]string "更新成功"
// @Failure      400 {object} map[string]string "请求参数错误"
// @Failure      500 {object} map[string]string "服务器错误"
// @Security     Bearer
// @Router       /order/{id} [put]
func (h *OrderHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	data := make(map[string]interface{})
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.orderService.Update(uint(id), data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
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

	if err := h.orderService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}