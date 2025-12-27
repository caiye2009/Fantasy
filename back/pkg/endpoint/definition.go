package endpoint

import (
	"github.com/gin-gonic/gin"
)

// RouteDefinition 路由定义
type RouteDefinition struct {
	Method  string          // HTTP 方法：GET, POST, PUT, DELETE 等
	Path    string          // 路由路径，如 "/client", "/client/:id"
	Handler gin.HandlerFunc // 处理函数
	Domain  string          // 业务域，如 "client", "order"
	Action  string          // 操作名称，如 "create", "update", "list"
}

// AuditMarker 审计标记函数类型（用于依赖注入，避免循环依赖）
type AuditMarker func(domain, action string) gin.HandlerFunc

var auditMarker AuditMarker

// SetAuditMarker 设置审计标记函数（在启动时注入）
func SetAuditMarker(marker AuditMarker) {
	auditMarker = marker
}

// RegisterRoutes 批量注册路由
// 自动处理：1) 添加audit中间件 2) 注册到Gin 3) 注册到GlobalRegistry
func RegisterRoutes(rg *gin.RouterGroup, routes []RouteDefinition) {
	for _, route := range routes {
		var handlers []gin.HandlerFunc

		// 1. 自动添加 audit 中间件（非GET请求才需要审计）
		if route.Domain != "" && route.Action != "" && route.Method != "GET" && auditMarker != nil {
			handlers = append(handlers, auditMarker(route.Domain, route.Action))
		}

		handlers = append(handlers, route.Handler)

		// 2. 注册到 Gin 路由
		rg.Handle(route.Method, route.Path, handlers...)

		// 3. 注册到 GlobalRegistry（用于 auth 权限检查和 permission 列表）
		if route.Domain != "" && route.Action != "" {
			endpointName := route.Domain + "." + route.Action
			fullPath := "/api/v1" + route.Path

			GlobalRegistry.Register(&Endpoint{
				Name:       endpointName,
				Path:       fullPath,
				Method:     route.Method,
				Domain:     route.Domain,
				Action:     route.Action,
				Permission: endpointName,
			})
		}
	}
}
