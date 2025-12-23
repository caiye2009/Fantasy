package auth

import (
	"context"
	"regexp"
	"strings"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

const (
	SourceContextKey  contextKey = "source"
	LoginIDContextKey contextKey = "loginId"
	RoleContextKey    contextKey = "role"
)

type contextKey string

type RequestContext struct {
	Source  string
	LoginID string
	Role    string
}

type AuthWang struct {
	jwtWang  *JWTWang
	enforcer *casbin.Enforcer
}

func NewAuthWang(jwtWang *JWTWang, enforcer *casbin.Enforcer) *AuthWang {
	return &AuthWang{
		jwtWang:  jwtWang,
		enforcer: enforcer,
	}
}

func (aw *AuthWang) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		source := c.GetHeader("Source")
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// if source == "" || authHeader == "" || tokenString == authHeader {
		if authHeader == "" || tokenString == authHeader {
			c.JSON(401, gin.H{"error": "菜叶"})
			c.Abort()
			return
		}

		claims, err := aw.jwtWang.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的token或token已过期"})
			c.Abort()
			return
		}

		c.Set("source", source)
		c.Set("loginId", claims.LoginID)
		c.Set("role", claims.Role)

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, SourceContextKey, source)
		ctx = context.WithValue(ctx, LoginIDContextKey, claims.LoginID)
		ctx = context.WithValue(ctx, RoleContextKey, claims.Role)
		c.Request = c.Request.WithContext(ctx)

		println("===== Auth Context =====")
		println("source:", source)
		println("loginId:", claims.LoginID)
		println("role:", claims.Role)
		println("========================")

		// 5. Casbin权限检查
		if aw.enforcer != nil {
			role := claims.Role
			path := c.Request.URL.Path
			method := c.Request.Method

			// 标准化路径（将数字ID替换为:id）
			normalizedPath := normalizePath(path)

			println("===== Permission Check =====")
			println("role:", role)
			println("path:", path)
			println("normalizedPath:", normalizedPath)
			println("method:", method)

			allowed, err := aw.enforcer.Enforce(role, normalizedPath, method)
			if err != nil {
				println("permission check error:", err.Error())
				c.JSON(500, gin.H{"error": "权限检查失败"})
				c.Abort()
				return
			}

			println("allowed:", allowed)
			println("===========================")

			if !allowed {
				c.JSON(403, gin.H{"error": "无权限访问该资源"})
				c.Abort()
				return
			}
		}

		// 6. 放行
		c.Next()
	}
}

// normalizePath 路径标准化（将数字ID替换为:id）
// 例如：/api/v1/order/123/detail -> /api/v1/order/:id/detail
func normalizePath(path string) string {
	// 使用正则表达式匹配路径中的数字ID
	re := regexp.MustCompile(`/\d+(/|$)`)
	return re.ReplaceAllString(path, "/:id$1")
}

func GetRequestContext(c *gin.Context) *RequestContext {
	source, _ := c.Get("source")
	loginId, _ := c.Get("loginId")
	role, _ := c.Get("role")

	return &RequestContext{
		Source:  source.(string),
		LoginID: loginId.(string),
		Role:    role.(string),
	}
}