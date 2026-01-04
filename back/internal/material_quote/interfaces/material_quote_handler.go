package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/material_quote/application"
	"back/internal/material_quote/domain"
	"back/pkg/endpoint"
)

type MaterialQuoteHandler struct {
	service *application.MaterialQuoteService
}

func NewMaterialQuoteHandler(service *application.MaterialQuoteService) *MaterialQuoteHandler {
	return &MaterialQuoteHandler{service: service}
}

func (h *MaterialQuoteHandler) Create(c *gin.Context) {
	var req application.CreateMaterialQuoteRequest
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

func (h *MaterialQuoteHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrMaterialQuoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "material_quote not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MaterialQuoteHandler) List(c *gin.Context) {
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

func (h *MaterialQuoteHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrMaterialQuoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "material_quote not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

func (h *MaterialQuoteHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrMaterialQuoteNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "material_quote not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func (h *MaterialQuoteHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/material_quotes", Handler: h.Create, Domain: "material_quote", Action: "create"},
		{Method: "GET", Path: "/material_quotes/:id", Handler: h.Get, Domain: "material_quote", Action: "get"},
		{Method: "GET", Path: "/material_quotes", Handler: h.List, Domain: "material_quote", Action: "list"},
		{Method: "PUT", Path: "/material_quotes/:id", Handler: h.Update, Domain: "material_quote", Action: "update"},
		{Method: "DELETE", Path: "/material_quotes/:id", Handler: h.Delete, Domain: "material_quote", Action: "delete"},
	}
}
