package interfaces

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/process_quote/application"
	"back/internal/process_quote/domain"
	"back/pkg/endpoint"
	"back/pkg/handler"
)

type ProcessQuoteHandler struct {
	service *application.ProcessQuoteService
}

func NewProcessQuoteHandler(service *application.ProcessQuoteService) *ProcessQuoteHandler {
	return &ProcessQuoteHandler{service: service}
}

func (h *ProcessQuoteHandler) Create(c *gin.Context) {
	handler.HandleCreate(c, h.service.Create, func(resp *application.ProcessQuoteResponse) interface{} {
		return resp.ID
	})
}

func (h *ProcessQuoteHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleGet(c, uint(id), h.service.Get, domain.ErrProcessQuoteNotFound)
}

func (h *ProcessQuoteHandler) List(c *gin.Context) {
	handler.HandleList(c, h.service.List)
}

func (h *ProcessQuoteHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleUpdate(c, uint(id), h.service.Get, h.service.Update, domain.ErrProcessQuoteNotFound)
}

func (h *ProcessQuoteHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleDelete(c, uint(id), h.service.Get, h.service.Delete, domain.ErrProcessQuoteNotFound)
}

func (h *ProcessQuoteHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/process_quotes", Handler: h.Create, Domain: "process_quote", Action: "create"},
		{Method: "GET", Path: "/process_quotes/:id", Handler: h.Get, Domain: "process_quote", Action: "get"},
		{Method: "GET", Path: "/process_quotes", Handler: h.List, Domain: "process_quote", Action: "list"},
		{Method: "PUT", Path: "/process_quotes/:id", Handler: h.Update, Domain: "process_quote", Action: "update"},
		{Method: "DELETE", Path: "/process_quotes/:id", Handler: h.Delete, Domain: "process_quote", Action: "delete"},
	}
}
