package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/product_allocation/application"
	"back/internal/product_allocation/domain"
	"back/pkg/endpoint"
)

// ProductAllocationHandler ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ· Handler
type ProductAllocationHandler struct {
	service *application.ProductAllocationService
}

// NewProductAllocationHandler ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ»ÃÂÃÂº Handler
func NewProductAllocationHandler(service *application.ProductAllocationService) *ProductAllocationHandler {
	return &ProductAllocationHandler{service: service}
}

// Create ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ»ÃÂÃÂºÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·
func (h *ProductAllocationHandler) Create(c *gin.Context) {
	var req application.CreateProductAllocationRequest
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

// Get ÃÂÃÂ¨ÃÂÃÂÃÂÃÂ·ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·
func (h *ProductAllocationHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrProductAllocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product_allocation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¨ÃÂÃÂ¡ÃÂÃÂ¨
func (h *ProductAllocationHandler) List(c *gin.Context) {
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

// Update ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ´ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ°ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·
func (h *ProductAllocationHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrProductAllocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product_allocation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete ÃÂÃÂ¥ÃÂÃÂÃÂÃÂ ÃÂÃÂ©ÃÂÃÂÃÂÃÂ¤ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·
func (h *ProductAllocationHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrProductAllocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product_allocation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

// GetRoutes ÃÂÃÂ¨ÃÂÃÂÃÂÃÂ·ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¨ÃÂÃÂ·ÃÂÃÂ¯ÃÂÃÂ§ÃÂÃÂÃÂÃÂ±ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂÃÂÃÂ¤ÃÂÃÂ¹ÃÂÃÂ
func (h *ProductAllocationHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/product_allocations", Handler: h.Create, Domain: "product_allocation", Action: "create"},
		{Method: "GET", Path: "/product_allocations/:id", Handler: h.Get, Domain: "product_allocation", Action: "get"},
		{Method: "GET", Path: "/product_allocations", Handler: h.List, Domain: "product_allocation", Action: "list"},
		{Method: "PUT", Path: "/product_allocations/:id", Handler: h.Update, Domain: "product_allocation", Action: "update"},
		{Method: "DELETE", Path: "/product_allocations/:id", Handler: h.Delete, Domain: "product_allocation", Action: "delete"},
	}
}
