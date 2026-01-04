package interfaces

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/order/application"
	"back/internal/order/domain"
	"back/pkg/endpoint"
	"back/pkg/handler"
)

type OrderHandler struct {
	service *application.OrderService
}

func NewOrderHandler(service *application.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(c *gin.Context) {
	handler.HandleCreate(c, h.service.Create, func(resp *application.OrderResponse) interface{} {
		return resp.ID
	})
}

func (h *OrderHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleGet(c, uint(id), h.service.Get, domain.ErrOrderNotFound)
}

func (h *OrderHandler) List(c *gin.Context) {
	handler.HandleList(c, h.service.List)
}

func (h *OrderHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleUpdate(c, uint(id), h.service.Get, h.service.Update, domain.ErrOrderNotFound)
}

func (h *OrderHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleDelete(c, uint(id), h.service.Get, h.service.Delete, domain.ErrOrderNotFound)
}

func (h *OrderHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/orders", Handler: h.Create, Domain: "order", Action: "create"},
		{Method: "GET", Path: "/orders/:id", Handler: h.Get, Domain: "order", Action: "get"},
		{Method: "GET", Path: "/orders", Handler: h.List, Domain: "order", Action: "list"},
		{Method: "PUT", Path: "/orders/:id", Handler: h.Update, Domain: "order", Action: "update"},
		{Method: "DELETE", Path: "/orders/:id", Handler: h.Delete, Domain: "order", Action: "delete"},
	}
}
