package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/order_process/application"
	"back/internal/order_process/domain"
	"back/pkg/endpoint"
)

type OrderProcessHandler struct {
	service *application.OrderProcessService
}

func NewOrderProcessHandler(service *application.OrderProcessService) *OrderProcessHandler {
	return &OrderProcessHandler{service: service}
}

func (h *OrderProcessHandler) Create(c *gin.Context) {
	var req application.CreateOrderProcessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderProcessHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrOrderProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderProcessHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	list, total, err := h.service.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"list":  list,
	})
}

func (h *OrderProcessHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrOrderProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

func (h *OrderProcessHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrOrderProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func (h *OrderProcessHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/order_processs", Handler: h.Create, Domain: "order_process", Action: "create"},
		{Method: "GET", Path: "/order_processs/:id", Handler: h.Get, Domain: "order_process", Action: "get"},
		{Method: "GET", Path: "/order_processs", Handler: h.List, Domain: "order_process", Action: "list"},
		{Method: "PUT", Path: "/order_processs/:id", Handler: h.Update, Domain: "order_process", Action: "update"},
		{Method: "DELETE", Path: "/order_processs/:id", Handler: h.Delete, Domain: "order_process", Action: "delete"},
	}
}
