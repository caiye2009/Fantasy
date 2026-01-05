package interfaces

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/fabric_material/application"
	"back/internal/fabric_material/domain"
	"back/pkg/endpoint"
)

type FabricMaterialHandler struct {
	service *application.FabricMaterialService
}

func NewFabricMaterialHandler(service *application.FabricMaterialService) *FabricMaterialHandler {
	return &FabricMaterialHandler{service: service}
}

func (h *FabricMaterialHandler) Create(c *gin.Context) {
	var req application.CreateFabricMaterialRequest
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

func (h *FabricMaterialHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	resp, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, domain.ErrFabricMaterialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "fabric_material not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}


func (h *FabricMaterialHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.service.Delete(c.Request.Context(), uint(id)); err != nil {
		if errors.Is(err, domain.ErrFabricMaterialNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "fabric_material not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func (h *FabricMaterialHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/fabric_materials", Handler: h.Create, Domain: "fabric_material", Action: "create"},
		{Method: "GET", Path: "/fabric_materials/:id", Handler: h.Get, Domain: "fabric_material", Action: "get"},
		// {Method: "GET", Path: "/fabric_materials", Handler: h.List, Domain: "fabric_material", Action: "list"},
		// {Method: "PUT", Path: "/fabric_materials/:id", Handler: h.Update, Domain: "fabric_material", Action: "update"},
		{Method: "DELETE", Path: "/fabric_materials/:id", Handler: h.Delete, Domain: "fabric_material", Action: "delete"},
	}
}