package endpoint

import (
	"github.com/gin-gonic/gin"
)

// RouteDefinition 路由定义
type RouteDefinition struct {
	Method      string            // HTTP 方法：GET, POST, PUT, DELETE 等
	Path        string            // 路由路径，如 "/client", "/client/:id"
	Handler     gin.HandlerFunc   // 处理函数
	Middlewares []gin.HandlerFunc // 中间件列表
	Name        string            // 接口名称，用于 audit 和日志
}

// RegisterRoutes 批量注册路由
func RegisterRoutes(rg *gin.RouterGroup, routes []RouteDefinition) {
	for _, route := range routes {
		handlers := make([]gin.HandlerFunc, 0, len(route.Middlewares)+1)
		handlers = append(handlers, route.Middlewares...)
		handlers = append(handlers, route.Handler)

		switch route.Method {
		case "GET":
			rg.GET(route.Path, handlers...)
		case "POST":
			rg.POST(route.Path, handlers...)
		case "PUT":
			rg.PUT(route.Path, handlers...)
		case "DELETE":
			rg.DELETE(route.Path, handlers...)
		case "PATCH":
			rg.PATCH(route.Path, handlers...)
		}
	}
}
