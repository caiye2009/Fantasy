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
		// 1. 提取 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "缺少认证信息"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(401, gin.H{"error": "认证格式错误"})
			c.Abort()
			return
		}

		// 2. 验证 JWT
		claims, err := aw.jwtWang.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}

		// 3. 查询用户状态 (使用 login_id)
		var user struct {
			ID     uint   `gorm:"column:id"`
			Status string `gorm:"column:status"`
		}
		if err := aw.db.Table("users").Select("id, status").Where("login_id = ?", claims.LoginID).First(&user).Error; err != nil {
			c.JSON(401, gin.H{"error": "用户不存在"})
			c.Abort()
			return
		}

		if user.Status != "active" {
			c.JSON(403, gin.H{"error": "账号已被停用"})
			c.Abort()
			return
		}

		// 4. 特殊路径:修改密码
		path := c.Request.URL.Path
		method := c.Request.Method

		if path == "/api/v1/user/password" && method == "PUT" {
			c.Set("login_id", claims.LoginID)
			c.Set("role", claims.Role)
			c.Next()
			return
		}

		// 5. Casbin 权限检查
		ok, err := aw.casbinWang.CheckPermission(claims.Role, path, method)
		if err != nil {
			c.JSON(500, gin.H{"error": "权限检查失败"})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(403, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}

		// 6. 设置上下文
		c.Set("login_id", claims.LoginID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// GetJWTWang 返回 JWTWang (供 config 使用)
func (aw *AuthWang) GetJWTWang() *JWTWang {
	return aw.jwtWang
}