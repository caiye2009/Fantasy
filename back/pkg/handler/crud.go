package handler

import (
	"context"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"back/pkg/audit"
	"back/pkg/response"
)

// HandleCreate 处理创建操作
// Req: 请求类型
// Resp: 响应类型
// createFunc: 创建函数
// getID: 从响应中提取ID的函数
func HandleCreate[Req, Resp any](
	c *gin.Context,
	createFunc func(ctx context.Context, req *Req) (*Resp, error),
	getID func(*Resp) interface{},
) {
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, err)
		return
	}

	resp, err := createFunc(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, 500, err)
		return
	}

	// 记录审计：创建操作
	audit.RecordCreate(c, resp, getID(resp))

	response.Success(c, resp)
}

// HandleGet 处理获取操作
// Resp: 响应类型
// id: 资源ID
// getFunc: 获取函数
// notFoundErr: 用于判断资源不存在的错误
func HandleGet[Resp any](
	c *gin.Context,
	id uint,
	getFunc func(ctx context.Context, id uint) (*Resp, error),
	notFoundErr error,
) {
	resp, err := getFunc(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, notFoundErr) {
			response.ErrorMessage(c, 404, "resource not found")
			return
		}
		response.Error(c, 500, err)
		return
	}

	response.Success(c, resp)
}

// HandleList 处理列表操作
// Resp: 响应类型
// listFunc: 列表查询函数
func HandleList[Resp any](
	c *gin.Context,
	listFunc func(ctx context.Context, limit, offset int) ([]*Resp, int64, error),
) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	list, total, err := listFunc(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, 500, err)
		return
	}

	response.Success(c, gin.H{
		"total": total,
		"list":  list,
	})
}

// HandleUpdate 处理更新操作
// Resp: 响应类型
// id: 资源ID
// getFunc: 获取函数（用于读取old和new）
// updateFunc: 更新函数
// notFoundErr: 用于判断资源不存在的错误
func HandleUpdate[Resp any](
	c *gin.Context,
	id uint,
	getFunc func(ctx context.Context, id uint) (*Resp, error),
	updateFunc func(ctx context.Context, id uint, updates map[string]interface{}) error,
	notFoundErr error,
) {
	// 1. 读取更新前的数据（用于audit）
	old, err := getFunc(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, notFoundErr) {
			response.ErrorMessage(c, 404, "resource not found")
			return
		}
		response.Error(c, 500, err)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.Error(c, 400, err)
		return
	}

	// 2. 执行更新
	if err := updateFunc(c.Request.Context(), id, updates); err != nil {
		if errors.Is(err, notFoundErr) {
			response.ErrorMessage(c, 404, "resource not found")
			return
		}
		response.Error(c, 500, err)
		return
	}

	// 3. 读取更新后的数据（用于audit）
	new, _ := getFunc(c.Request.Context(), id)

	// 记录审计：更新操作
	audit.RecordUpdate(c, old, new, id)

	response.Message(c, "updated successfully")
}

// HandleDelete 处理删除操作
// Resp: 响应类型
// id: 资源ID
// getFunc: 获取函数（用于读取old）
// deleteFunc: 删除函数
// notFoundErr: 用于判断资源不存在的错误
func HandleDelete[Resp any](
	c *gin.Context,
	id uint,
	getFunc func(ctx context.Context, id uint) (*Resp, error),
	deleteFunc func(ctx context.Context, id uint) error,
	notFoundErr error,
) {
	// 1. 读取删除前的数据（用于audit）
	old, err := getFunc(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, notFoundErr) {
			response.ErrorMessage(c, 404, "resource not found")
			return
		}
		response.Error(c, 500, err)
		return
	}

	// 2. 执行删除
	if err := deleteFunc(c.Request.Context(), id); err != nil {
		if errors.Is(err, notFoundErr) {
			response.ErrorMessage(c, 404, "resource not found")
			return
		}
		response.Error(c, 500, err)
		return
	}

	// 记录审计：删除操作
	audit.RecordDelete(c, old, id)

	response.Message(c, "deleted successfully")
}
