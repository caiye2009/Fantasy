package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/product_fabric/application"
	"back/internal/product_fabric/domain"
	"back/pkg/endpoint"
)

type ProductFabricHandler struct {
	service *application.ProductFabricService
}

func NewProductFabricHandler(service *application.ProductFabricService) *ProductFabricHandler {
	return &ProductFabricHandler{service: service}
}

func (h *ProductFabricHandler) Create(c *gin.Context) {
	var req application.CreateProductFabricRequest
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

func (h *ProductFabricHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrProductFabricNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product_fabric not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProductFabricHandler) List(c *gin.Context) {
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

func (h *ProductFabricHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrProductFabricNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product_fabric not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

func (h *ProductFabricHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrProductFabricNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product_fabric not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func (h *ProductFabricHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/product_fabrics", Handler: h.Create, Domain: "product_fabric", Action: "create"},
		{Method: "GET", Path: "/product_fabrics/:id", Handler: h.Get, Domain: "product_fabric", Action: "get"},
		{Method: "GET", Path: "/product_fabrics", Handler: h.List, Domain: "product_fabric", Action: "list"},
		{Method: "PUT", Path: "/product_fabrics/:id", Handler: h.Update, Domain: "product_fabric", Action: "update"},
		{Method: "DELETE", Path: "/product_fabrics/:id", Handler: h.Delete, Domain: "product_fabric", Action: "delete"},
	}
}