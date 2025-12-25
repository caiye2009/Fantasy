package endpoint

// Endpoint 接口元数据
type Endpoint struct {
	Name       string // 全局唯一标识：user.create, order.list
	Path       string // HTTP 路径：/api/v1/user
	Method     string // HTTP 方法：POST, GET, PUT, DELETE
	Domain     string // 审计域：user, order, client
	Action     string // 审计动作：create, update, delete, list
	Permission string // 权限标识（通常等于 Name）
	Desc       string // 描述
}

// GetName 获取完整名称 (domain.action)
func (e *Endpoint) GetName() string {
	if e.Name != "" {
		return e.Name
	}
	return e.Domain + "." + e.Action
}
