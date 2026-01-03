package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/process_quote/application"
	"back/internal/process_quote/domain"
	"back/pkg/endpoint"
)

// ProcessQuoteHandler ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ· Handler
type ProcessQuoteHandler struct {
	service *application.ProcessQuoteService
}

// NewProcessQuoteHandler ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ»ÃÂÃÂº Handler
func NewProcessQuoteHandler(service *application.ProcessQuoteService) *ProcessQuoteHandler {
	return &ProcessQuoteHandler{service: service}
}

// Create ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ»ÃÂÃÂºÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·
func (h *ProcessQuoteHandler) Create(c *gin.Context) {
	var req application.CreateProcessQuoteRequest
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
func (h *ProcessQuoteHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrProcessQuoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "process_quote not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¨ÃÂÃÂ¡ÃÂÃÂ¨
func (h *ProcessQuoteHandler) List(c *gin.Context) {
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
func (h *ProcessQuoteHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrProcessQuoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "process_quote not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete ÃÂÃÂ¥ÃÂÃÂÃÂÃÂ ÃÂÃÂ©ÃÂÃÂÃÂÃÂ¤ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·
func (h *ProcessQuoteHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrProcessQuoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "process_quote not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

// GetRoutes ÃÂÃÂ¨ÃÂÃÂÃÂÃÂ·ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¨ÃÂÃÂ·ÃÂÃÂ¯ÃÂÃÂ§ÃÂÃÂÃÂÃÂ±ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂÃÂÃÂ¤ÃÂÃÂ¹ÃÂÃÂ
func (h *ProcessQuoteHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/process_quotes", Handler: h.Create, Domain: "process_quote", Action: "create"},
		{Method: "GET", Path: "/process_quotes/:id", Handler: h.Get, Domain: "process_quote", Action: "get"},
		{Method: "GET", Path: "/process_quotes", Handler: h.List, Domain: "process_quote", Action: "list"},
		{Method: "PUT", Path: "/process_quotes/:id", Handler: h.Update, Domain: "process_quote", Action: "update"},
		{Method: "DELETE", Path: "/process_quotes/:id", Handler: h.Delete, Domain: "process_quote", Action: "delete"},
	}
}
