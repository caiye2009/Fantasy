package fields

import "github.com/gin-gonic/gin"

const (
	CodeSuccess = 0
	CodeError   = 1
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code: CodeSuccess,
		Data: data,
	})
}

func Error(c *gin.Context, msg string) {
	c.JSON(200, Response{
		Code: CodeError,
		Msg:  msg,
	})
}

func ErrorWithCode(c *gin.Context, code int, msg string) {
	c.JSON(200, Response{
		Code: code,
		Msg:  msg,
	})
}