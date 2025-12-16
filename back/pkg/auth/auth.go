package auth

import (
	"context"
	"strings"
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
	jwtWang *JWTWang
}

func NewAuthWang(jwtWang *JWTWang) *AuthWang {
	return &AuthWang{
		jwtWang: jwtWang,
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

		// 5. TODO: 权限校验将在此处添加

		// 6. 放行
		c.Next()
	}
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