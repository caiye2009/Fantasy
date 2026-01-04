package response

import (
	"github.com/gin-gonic/gin"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, data)
}

// Error 错误响应
// 自动将错误信息设置到context，供audit中间件记录
func Error(c *gin.Context, code int, err error) {
	// 设置错误信息到context（供audit使用）
	if err != nil {
		c.Set("error", err.Error())
	}

	c.JSON(code, gin.H{"error": err.Error()})
}

// Message 消息响应
func Message(c *gin.Context, message string) {
	c.JSON(200, gin.H{"message": message})
}

// ErrorMessage 错误消息响应
func ErrorMessage(c *gin.Context, code int, message string) {
	c.Set("error", message)
	c.JSON(code, gin.H{"error": message})
}
