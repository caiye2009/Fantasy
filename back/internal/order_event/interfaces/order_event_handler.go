package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/order_event/application"
	"back/internal/order_event/domain"
	"back/pkg/endpoint"
)

type OrderEventHandler struct {
	service *application.OrderEventService
}

func NewOrderEventHandler(service *application.OrderEventService) *OrderEventHandler {
	return &OrderEventHandler{service: service}
}

func (h *OrderEventHandler) Create(c *gin.Context) {
	var req application.CreateOrderEventRequest
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

func (h *OrderEventHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrOrderEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderEventHandler) List(c *gin.Context) {
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

func (h *OrderEventHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrOrderEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

func (h *OrderEventHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrOrderEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func (h *OrderEventHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/order_events", Handler: h.Create, Domain: "order_event", Action: "create"},
		{Method: "GET", Path: "/order_events/:id", Handler: h.Get, Domain: "order_event", Action: "get"},
		{Method: "GET", Path: "/order_events", Handler: h.List, Domain: "order_event", Action: "list"},
		{Method: "PUT", Path: "/order_events/:id", Handler: h.Update, Domain: "order_event", Action: "update"},
		{Method: "DELETE", Path: "/order_events/:id", Handler: h.Delete, Domain: "order_event", Action: "delete"},
	}
}
