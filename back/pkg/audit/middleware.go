package audit

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditMiddleware 审计中间件
// 在 auth 中间件之后使用，自动记录所有非 GET 请求
func AuditMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 跳过 GET 请求（只读操作不需要审计）
		if c.Request.Method == "GET" {
			c.Next()
			return
		}

		// 2. 跳过 OPTIONS 请求（CORS 预检）
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 3. 创建 Recorder 并注入到 context
		recorder := NewRecorder(c, db)
		c.Set(RecorderContextKey, recorder)

		// 4. 执行 handler
		c.Next()

		// 5. ← 统一调用点：handler 执行完毕后自动保存审计日志
		// 即使业务层没有设置 old/new，也会记录基本信息
		if err := recorder.Save(); err != nil {
			// 记录审计日志失败不影响业务流程，只打印日志
			println("Failed to save audit log:", err.Error())
		}
	}
}

// Mark 标记路由的 domain 和 action（在路由注册时使用）
// 示例：rg.POST("/order", audit.Mark("order", "orderCreation"), handler.Create)
func Mark(domain, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 预先设置 domain 和 action，避免自动推断
		c.Set("audit_domain", domain)
		c.Set("audit_action", action)
		c.Next()
	}
}

// Skip 创建一个跳过审计的中间件（用于特殊路由）
func Skip() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("skip_audit", true)
		c.Next()
	}
}

// shouldSkip 检查是否应该跳过审计
func shouldSkip(c *gin.Context) bool {
	skip, exists := c.Get("skip_audit")
	return exists && skip.(bool)
}
