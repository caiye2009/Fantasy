package endpoint

import (
	"github.com/gin-gonic/gin"
)

// RouteRegistrar 统一的路由和Endpoint注册器
// 负责同时注册Gin路由和Endpoint元数据
type RouteRegistrar struct {
	rg       *gin.RouterGroup
	registry *Registry
	domain   string // 当前业务域
}

// NewRegistrar 创建注册器
func NewRegistrar(rg *gin.RouterGroup, domain string) *RouteRegistrar {
	return &RouteRegistrar{
		rg:       rg,
		registry: GlobalRegistry,
		domain:   domain,
	}
}

// POST 注册POST路由
// path: 相对路径，如 "/client" 或 "/client/:id/detail"
// action: 操作名称，如 "create", "update"
// desc: 描述信息
// handlers: Gin处理器，可以包含 audit.Mark() 等中间件
func (r *RouteRegistrar) POST(path, action, desc string, handlers ...gin.HandlerFunc) {
	r.register("POST", path, action, desc, handlers...)
}

// GET 注册GET路由
func (r *RouteRegistrar) GET(path, action, desc string, handlers ...gin.HandlerFunc) {
	r.register("GET", path, action, desc, handlers...)
}

// PUT 注册PUT路由
func (r *RouteRegistrar) PUT(path, action, desc string, handlers ...gin.HandlerFunc) {
	r.register("PUT", path, action, desc, handlers...)
}

// DELETE 注册DELETE路由
func (r *RouteRegistrar) DELETE(path, action, desc string, handlers ...gin.HandlerFunc) {
	r.register("DELETE", path, action, desc, handlers...)
}

// PATCH 注册PATCH路由
func (r *RouteRegistrar) PATCH(path, action, desc string, handlers ...gin.HandlerFunc) {
	r.register("PATCH", path, action, desc, handlers...)
}

// register 统一注册逻辑
func (r *RouteRegistrar) register(method, path, action, desc string, handlers ...gin.HandlerFunc) {
	name := r.domain + "." + action
	fullPath := "/api/v1" + path

	// 1. 注册Gin路由
	switch method {
	case "GET":
		r.rg.GET(path, handlers...)
	case "POST":
		r.rg.POST(path, handlers...)
	case "PUT":
		r.rg.PUT(path, handlers...)
	case "DELETE":
		r.rg.DELETE(path, handlers...)
	case "PATCH":
		r.rg.PATCH(path, handlers...)
	}

	// 2. 注册Endpoint元数据到全局注册表
	r.registry.Register(&Endpoint{
		Name:       name,
		Path:       fullPath,
		Method:     method,
		Domain:     r.domain,
		Action:     action,
		Permission: name, // 权限标识等于name
		Desc:       desc,
	})
}

// GetDomain 获取当前domain (供外部创建audit.Mark时使用)
func (r *RouteRegistrar) GetDomain() string {
	return r.domain
}
