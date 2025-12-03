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

func (h *OrderHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	o, err := h.orderService.Get(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, o)
}

func (h *OrderHandler) List(c *gin.Context) {
	list, err := h.orderService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

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

func (h *OrderHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.orderService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}