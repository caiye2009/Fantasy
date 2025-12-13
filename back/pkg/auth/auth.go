package auth

import (
	"strings"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthWang struct {
	jwtWang    *JWTWang
	casbinWang *CasbinWang
	db         *gorm.DB
}

func NewAuthWang(
	jwtWang *JWTWang,
	casbinWang *CasbinWang,
	db *gorm.DB,
) *AuthWang {
	return &AuthWang{
		jwtWang:    jwtWang,
		casbinWang: casbinWang,
		db:         db,
	}
}

func (aw *AuthWang) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method

		// 1. 提取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			println("===== Auth Log =====")
			println("jwt: <missing>")
			println("user: <none>")
			println("role: <none>")
			println("endpoint:", method, path)
			println("outcome: deny (缺少认证信息)")
			println("====================")
			c.JSON(401, gin.H{"error": "缺少认证信息"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			println("===== Auth Log =====")
			println("jwt:", authHeader[:20]+"...")
			println("user: <none>")
			println("role: <none>")
			println("endpoint:", method, path)
			println("outcome: deny (认证格式错误)")
			println("====================")
			c.JSON(401, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		// 2. 验证 Access Token（改用 ValidateAccessToken）
		claims, err := aw.jwtWang.ValidateAccessToken(tokenString)
		if err != nil {
			println("===== Auth Log =====")
			println("jwt:", tokenString[:20]+"...")
			println("user: <invalid>")
			println("role: <invalid>")
			println("endpoint:", method, path)
			println("outcome: deny (token无效或已过期)")
			println("====================")
			c.JSON(401, gin.H{"error": "无效的token或token已过期"})
			c.Abort()
			return
		}

		// 3. 查询用户状态（使用 loginId）
		var user struct {
			ID     uint   `gorm:"column:id"`
			Status string `gorm:"column:status"`
		}
		if err := aw.db.Table("users").Select("id, status").Where("login_id = ?", claims.LoginID).First(&user).Error; err != nil {
			println("===== Auth Log =====")
			println("jwt:", tokenString[:20]+"...")
			println("user:", claims.LoginID)
			println("role:", claims.Role)
			println("endpoint:", method, path)
			println("outcome: deny (用户不存在)")
			println("====================")
			c.JSON(401, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		if user.Status != "active" {
			println("===== Auth Log =====")
			println("jwt:", tokenString[:20]+"...")
			println("user:", claims.LoginID)
			println("role:", claims.Role)
			println("endpoint:", method, path)
			println("outcome: deny (账号已被停用)")
			println("====================")
			c.JSON(403, gin.H{"error": "账号已被停用"})
			c.Abort()
			return
		}

		// 4. 特殊路径：修改密码（不需要 Casbin 检查）
		if path == "/api/v1/user/password" && method == "PUT" {
			c.Set("loginId", claims.LoginID)
			c.Set("role", claims.Role)

			// 将 loginId 和 role 设置到 request context
			ctx := c.Request.Context()
			ctx = SetLoginID(ctx, claims.LoginID)
			ctx = SetRole(ctx, claims.Role)
			c.Request = c.Request.WithContext(ctx)

			println("===== Auth Log =====")
			println("jwt:", tokenString[:20]+"...")
			println("user:", claims.LoginID)
			println("role:", claims.Role)
			println("endpoint:", method, path)
			println("outcome: allow (修改密码，跳过权限检查)")
			println("====================")

			c.Next()
			return
		}

		// 5. Casbin 权限检查
		ok, err := aw.casbinWang.CheckPermission(claims.Role, path, method)
		if err != nil {
			println("===== Auth Log =====")
			println("jwt:", tokenString[:20]+"...")
			println("user:", claims.LoginID)
			println("role:", claims.Role)
			println("endpoint:", method, path)
			println("outcome: deny (权限检查失败)")
			println("====================")
			c.JSON(500, gin.H{"error": "权限检查失败"})
			c.Abort()
			return
		}

		if !ok {
			println("===== Auth Log =====")
			println("jwt:", tokenString[:20]+"...")
			println("user:", claims.LoginID)
			println("role:", claims.Role)
			println("endpoint:", method, path)
			println("outcome: deny (权限不足)")
			println("====================")
			c.JSON(403, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		// 6. 设置上下文（同时设置到 gin context 和 request context）
		c.Set("loginId", claims.LoginID)
		c.Set("role", claims.Role)

		// 将 loginId 和 role 设置到 request context，供业务层使用
		ctx := c.Request.Context()
		ctx = SetLoginID(ctx, claims.LoginID)
		ctx = SetRole(ctx, claims.Role)
		c.Request = c.Request.WithContext(ctx)

		println("===== Auth Log =====")
		println("jwt:", tokenString[:20]+"...")
		println("user:", claims.LoginID)
		println("role:", claims.Role)
		println("endpoint:", method, path)
		println("outcome: allow")
		println("====================")

		c.Next()
	}
}

// GetJWTWang 返回 JWTWang（供 config 使用）
func (aw *AuthWang) GetJWTWang() *JWTWang {
	return aw.jwtWang
}