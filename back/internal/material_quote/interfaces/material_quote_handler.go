package interfaces

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/material_quote/application"
	"back/internal/material_quote/domain"
	"back/pkg/endpoint"
	"back/pkg/handler"
)

type MaterialQuoteHandler struct {
	service *application.MaterialQuoteService
}

func NewMaterialQuoteHandler(service *application.MaterialQuoteService) *MaterialQuoteHandler {
	return &MaterialQuoteHandler{service: service}
}

func (h *MaterialQuoteHandler) Create(c *gin.Context) {
	handler.HandleCreate(c, h.service.Create, func(resp *application.MaterialQuoteResponse) interface{} {
		return resp.ID
	})
}

func (h *MaterialQuoteHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleGet(c, uint(id), h.service.Get, domain.ErrMaterialQuoteNotFound)
}

func (h *MaterialQuoteHandler) List(c *gin.Context) {
	handler.HandleList(c, h.service.List)
}

func (h *MaterialQuoteHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleUpdate(c, uint(id), h.service.Get, h.service.Update, domain.ErrMaterialQuoteNotFound)
}

func (h *MaterialQuoteHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleDelete(c, uint(id), h.service.Get, h.service.Delete, domain.ErrMaterialQuoteNotFound)
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
