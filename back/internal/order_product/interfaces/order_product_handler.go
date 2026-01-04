package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/order_product/application"
	"back/internal/order_product/domain"
	"back/pkg/endpoint"
)

type OrderProductHandler struct {
	service *application.OrderProductService
}

func NewOrderProductHandler(service *application.OrderProductService) *OrderProductHandler {
	return &OrderProductHandler{service: service}
}

func (h *OrderProductHandler) Create(c *gin.Context) {
	var req application.CreateOrderProductRequest
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

func (h *OrderProductHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrOrderProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *OrderProductHandler) List(c *gin.Context) {
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

func (h *OrderProductHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrOrderProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

func (h *OrderProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrOrderProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func (h *OrderProductHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/order_products", Handler: h.Create, Domain: "order_product", Action: "create"},
		{Method: "GET", Path: "/order_products/:id", Handler: h.Get, Domain: "order_product", Action: "get"},
		{Method: "GET", Path: "/order_products", Handler: h.List, Domain: "order_product", Action: "list"},
		{Method: "PUT", Path: "/order_products/:id", Handler: h.Update, Domain: "order_product", Action: "update"},
		{Method: "DELETE", Path: "/order_products/:id", Handler: h.Delete, Domain: "order_product", Action: "delete"},
	}
}
