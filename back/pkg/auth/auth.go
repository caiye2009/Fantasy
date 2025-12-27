package auth

import (
	"context"
	"strings"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"back/pkg/endpoint"
)

const (
	SourceContextKey     contextKey = "source"
	LoginIDContextKey    contextKey = "loginId"
	UsernameContextKey   contextKey = "username"
	RoleContextKey       contextKey = "role"
	DepartmentContextKey contextKey = "department"
)

type contextKey string

type RequestContext struct {
	Source     string
	LoginID    string
	Username   string
	Role       string
	Department string
}

type AuthWang struct {
	jwtWang          *JWTWang
	enforcer         *casbin.Enforcer
	whitelistManager *WhitelistManager
}

func NewAuthWang(jwtWang *JWTWang, enforcer *casbin.Enforcer, whitelistManager *WhitelistManager) *AuthWang {
	return &AuthWang{
		jwtWang:          jwtWang,
		enforcer:         enforcer,
		whitelistManager: whitelistManager,
	}
}

func (aw *AuthWang) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		source := c.GetHeader("Source")
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if authHeader == "" || tokenString == authHeader {
			c.JSON(401, gin.H{"error": "未提供认证信息"})
			c.Abort()
			return
		}

		claims, err := aw.jwtWang.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的token或token已过期"})
			c.Abort()
			return
		}

		// 检查 JWT 白名单
		// if aw.whitelistManager != nil {
		// 	jti := claims.ID // JWT ID from RegisteredClaims
		// 	inWhitelist, err := aw.whitelistManager.IsInWhitelist(claims.LoginID, source, jti)
		// 	if err != nil {
		// 		println("whitelist check error:", err.Error())
		// 		c.JSON(500, gin.H{"error": "验证失败"})
		// 		c.Abort()
		// 		return
		// 	}
		// 	if !inWhitelist {
		// 		c.JSON(401, gin.H{"error": "token已失效，请重新登录"})
		// 		c.Abort()
		// 		return
		// 	}
		// }

		// 设置完整的 context 到 gin.Context（供 Audit 等中间件使用）
		c.Set("source", source)
		c.Set("loginId", claims.LoginID)
		c.Set("username", claims.UserName)
		c.Set("role", claims.Role)
		c.Set("department", claims.Department)

		// 设置到 Request Context（供 Service 层使用）
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, SourceContextKey, source)
		ctx = context.WithValue(ctx, LoginIDContextKey, claims.LoginID)
		ctx = context.WithValue(ctx, UsernameContextKey, claims.UserName)
		ctx = context.WithValue(ctx, RoleContextKey, claims.Role)
		ctx = context.WithValue(ctx, DepartmentContextKey, claims.Department)
		c.Request = c.Request.WithContext(ctx)

		// 5. 根据 HTTP route 查找 endpoint
		ep := endpoint.GlobalRegistry.FindByRoute(c.Request.Method, c.Request.URL.Path)
		if ep == nil {
			// 未注册的接口，拒绝访问
			c.JSON(404, gin.H{"error": "接口未注册"})
			c.Abort()
			return
		}

		// 将 endpoint 信息存到 context（供 Audit 使用）
		c.Set("endpoint", ep)

		// 6. Casbin 权限检查（使用接口名称）
		if aw.enforcer != nil {
			loginID := claims.LoginID
			role := claims.Role
			permission := ep.GetName() // 如 "user.create", "order.list"

			// 先检查用户个性化权限
			allowed, err := aw.enforcer.Enforce(loginID, permission, "*")
			if err != nil {
				c.JSON(500, gin.H{"error": "权限检查失败"})
				c.Abort()
				return
			}

			// 如果用户没有个性化权限，检查角色权限
			if !allowed {
				allowed, err = aw.enforcer.Enforce(role, permission, "*")
				if err != nil {
					c.JSON(500, gin.H{"error": "权限检查失败"})
					c.Abort()
					return
				}
			}

			if !allowed {
				c.JSON(403, gin.H{"error": "无权限访问: " + permission})
				c.Abort()
				return
			}
		}

		// 7. 放行
		c.Next()
	}
}

func GetRequestContext(c *gin.Context) *RequestContext {
	source, _ := c.Get("source")
	loginId, _ := c.Get("loginId")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	department, _ := c.Get("department")

	return &RequestContext{
		Source:     source.(string),
		LoginID:    loginId.(string),
		Username:   username.(string),
		Role:       role.(string),
		Department: department.(string),
	}
}