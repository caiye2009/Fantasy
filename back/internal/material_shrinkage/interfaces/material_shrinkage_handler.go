package interfaces

import (
	"github.com/gin-gonic/gin"

	"back/internal/material_shrinkage/application"
	"back/pkg/endpoint"
	"back/pkg/handler"
)

type MaterialShrinkageHandler struct {
	service *application.MaterialShrinkageService
}

func NewMaterialShrinkageHandler(service *application.MaterialShrinkageService) *MaterialShrinkageHandler {
	return &MaterialShrinkageHandler{service: service}
}

func (h *MaterialShrinkageHandler) Create(c *gin.Context) {
	handler.HandleCreate(c, h.service.Create, func(resp *application.MaterialShrinkageResponse) interface{} {
		return resp.ShrinkageID
	})
}

func (h *MaterialShrinkageHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/material_shrinkages", Handler: h.Create, Domain: "material_shrinkage", Action: "create"},
		// {Method: "GET", Path: "/material_shrinkages/:id", Handler: h.Get, Domain: "material_shrinkage", Action: "get"},
		// {Method: "GET", Path: "/material_shrinkages", Handler: h.List, Domain: "material_shrinkage", Action: "list"},
		// {Method: "PUT", Path: "/material_shrinkages/:id", Handler: h.Update, Domain: "material_shrinkage", Action: "update"},
		// {Method: "DELETE", Path: "/material_shrinkages/:id", Handler: h.Delete, Domain: "material_shrinkage", Action: "delete"},
	}
}