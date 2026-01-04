package audit

import (
	"github.com/gin-gonic/gin"
)

// RecordCreate 记录创建操作
// resource: 创建后的完整数据
// id: 资源ID
func RecordCreate(c *gin.Context, resource interface{}, id interface{}) {
	recorder := Get(c)
	if recorder == nil {
		return
	}

	recorder.SetResourceID(id)
	recorder.SetNew(resource)
}

// RecordUpdate 记录更新操作
// old: 更新前的数据
// new: 更新后的数据
// id: 资源ID
func RecordUpdate(c *gin.Context, old, new interface{}, id interface{}) {
	recorder := Get(c)
	if recorder == nil {
		return
	}

	recorder.SetResourceID(id)
	recorder.SetOld(old)
	recorder.SetNew(new)
}

// RecordDelete 记录删除操作
// old: 删除前的数据
// id: 资源ID
func RecordDelete(c *gin.Context, old interface{}, id interface{}) {
	recorder := Get(c)
	if recorder == nil {
		return
	}

	recorder.SetResourceID(id)
	recorder.SetOld(old)
}

// RecordResourceID 仅记录资源ID（用于不需要记录详细数据的场景）
func RecordResourceID(c *gin.Context, id interface{}) {
	recorder := Get(c)
	if recorder == nil {
		return
	}

	recorder.SetResourceID(id)
}
