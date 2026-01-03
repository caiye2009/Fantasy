package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/user/application"
	"back/internal/user/domain"
	"back/pkg/endpoint"
)

// UserHandler ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ· Handler
type UserHandler struct {
	service *application.UserService
}

// NewUserHandler ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ»ÃÂÃÂº Handler
func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Create ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¥ÃÂÃÂ»ÃÂÃÂºÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·
func (h *UserHandler) Create(c *gin.Context) {
	var req application.UserCreateRequest
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
func (h *UserHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// List ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¨ÃÂÃÂ¡ÃÂÃÂ¨
func (h *UserHandler) List(c *gin.Context) {
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
func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uint(id), updates); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "updated successfully"})
}

// Delete ÃÂÃÂ¥ÃÂÃÂÃÂÃÂ ÃÂÃÂ©ÃÂÃÂÃÂÃÂ¤ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂ¢ÃÂÃÂ¦ÃÂÃÂÃÂÃÂ·
func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

// GetRoutes ÃÂÃÂ¨ÃÂÃÂÃÂÃÂ·ÃÂÃÂ¥ÃÂÃÂÃÂÃÂÃÂÃÂ¨ÃÂÃÂ·ÃÂÃÂ¯ÃÂÃÂ§ÃÂÃÂÃÂÃÂ±ÃÂÃÂ¥ÃÂÃÂ®ÃÂÃÂÃÂÃÂ¤ÃÂÃÂ¹ÃÂÃÂ
func (h *UserHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/users", Handler: h.Create, Domain: "user", Action: "create"},
		{Method: "GET", Path: "/users/:id", Handler: h.Get, Domain: "user", Action: "get"},
		{Method: "GET", Path: "/users", Handler: h.List, Domain: "user", Action: "list"},
		{Method: "PUT", Path: "/users/:id", Handler: h.Update, Domain: "user", Action: "update"},
		{Method: "DELETE", Path: "/users/:id", Handler: h.Delete, Domain: "user", Action: "delete"},
	}
}
