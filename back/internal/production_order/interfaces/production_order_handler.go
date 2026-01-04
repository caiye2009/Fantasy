package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/production_order/application"
	"back/internal/production_order/domain"
	"back/pkg/endpoint"
)

type ProductionOrderHandler struct {
	service *application.ProductionOrderService
}

func NewProductionOrderHandler(service *application.ProductionOrderService) *ProductionOrderHandler {
	return &ProductionOrderHandler{service: service}
}

func (h *ProductionOrderHandler) Create(c *gin.Context) {
	var req application.CreateProductionOrderRequest
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

func (h *ProductionOrderHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrProductionOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "production_order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductionOrderHandler) List(c *gin.Context) {
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

func (h *ProductionOrderHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrProductionOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "production_order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

func (h *ProductionOrderHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrProductionOrderNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "production_order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func (h *ProductionOrderHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/production_orders", Handler: h.Create, Domain: "production_order", Action: "create"},
		{Method: "GET", Path: "/production_orders/:id", Handler: h.Get, Domain: "production_order", Action: "get"},
		{Method: "GET", Path: "/production_orders", Handler: h.List, Domain: "production_order", Action: "list"},
		{Method: "PUT", Path: "/production_orders/:id", Handler: h.Update, Domain: "production_order", Action: "update"},
		{Method: "DELETE", Path: "/production_orders/:id", Handler: h.Delete, Domain: "production_order", Action: "delete"},
	}
}
