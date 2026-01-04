package interfaces

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"back/internal/product/application"
	"back/internal/product/domain"
	"back/pkg/endpoint"
	"back/pkg/handler"
)

// ProductHandler 产品 Handler
type ProductHandler struct {
	service *application.ProductService
}

// NewProductHandler 创建 Handler
func NewProductHandler(service *application.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// Create 创建产品
func (h *ProductHandler) Create(c *gin.Context) {
	handler.HandleCreate(c, h.service.Create, func(resp *application.ProductResponse) interface{} {
		return resp.ID
	})
}

// Get 获取产品
func (h *ProductHandler) Get(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleGet(c, uint(id), h.service.Get, domain.ErrProductNotFound)
}

// List 列表
func (h *ProductHandler) List(c *gin.Context) {
	handler.HandleList(c, h.service.List)
}

// Update 更新产品
func (h *ProductHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleUpdate(c, uint(id), h.service.Get, h.service.Update, domain.ErrProductNotFound)
}

// Delete 删除产品
func (h *ProductHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	handler.HandleDelete(c, uint(id), h.service.Get, h.service.Delete, domain.ErrProductNotFound)
}

// GetRoutes 获取路由定义
func (h *ProductHandler) GetRoutes() []endpoint.RouteDefinition {
	return []endpoint.RouteDefinition{
		{Method: "POST", Path: "/products", Handler: h.Create, Domain: "product", Action: "create"},
		{Method: "GET", Path: "/products/:id", Handler: h.Get, Domain: "product", Action: "get"},
		{Method: "GET", Path: "/products", Handler: h.List, Domain: "product", Action: "list"},
		{Method: "PUT", Path: "/products/:id", Handler: h.Update, Domain: "product", Action: "update"},
		{Method: "DELETE", Path: "/products/:id", Handler: h.Delete, Domain: "product", Action: "delete"},
	}
}
