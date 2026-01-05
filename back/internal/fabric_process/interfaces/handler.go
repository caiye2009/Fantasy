package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/fabric_process/application"
	"back/internal/fabric_process/domain"
	"back/pkg/endpoint"
)

type FabricProcessHandler struct {
	service *application.FabricProcessService
}

func NewFabricProcessHandler(service *application.FabricProcessService) *FabricProcessHandler {
	return &FabricProcessHandler{service: service}
}

func (h *FabricProcessHandler) Create(c *gin.Context) {
	var req application.CreateFabricProcessRequest
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

func (h *FabricProcessHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrFabricProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "fabric_process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *FabricProcessHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrFabricProcessNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "fabric_process not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func (h *FabricProcessHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/fabric_processes", Handler: h.Create, Domain: "fabric_process", Action: "create"},
		{Method: "GET", Path: "/fabric_processes/:id", Handler: h.Get, Domain: "fabric_process", Action: "get"},
		// {Method: "GET", Path: "/fabric_processes", Handler: h.List, Domain: "fabric_process", Action: "list"},
		// {Method: "PUT", Path: "/fabric_processes/:id", Handler: h.Update, Domain: "fabric_process", Action: "update"},
		{Method: "DELETE", Path: "/fabric_processes/:id", Handler: h.Delete, Domain: "fabric_process", Action: "delete"},
	}
}