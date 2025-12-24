package audit

import (
	"time"
)

// AuditLog 审计日志实体
type AuditLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	LoginID       uint      `gorm:"not null;index" json:"login_id"`                 // 操作人ID
	Username      string    `gorm:"size:100;not null" json:"username"`              // 操作人姓名
	Domain        string    `gorm:"size:50;not null;index" json:"domain"`           // 业务域（order, user, product等）
	Action        string    `gorm:"size:100;not null;index" json:"action"`          // 操作动作（create_order, update_user等）
	ResourceID    string    `gorm:"size:100;index" json:"resource_id"`              // 被操作的资源ID
	HTTPMethod    string    `gorm:"size:10;not null" json:"http_method"`            // HTTP方法（POST/PUT/DELETE等）
	RequestPath   string    `gorm:"size:500;not null" json:"request_path"`          // 请求路径
	IPAddress     string    `gorm:"size:50" json:"ip_address"`                      // 客户端IP
	StatusCode    int       `gorm:"not null" json:"status_code"`                    // 响应状态码
	ErrorMessage  string    `gorm:"type:text" json:"error_message,omitempty"`       // 错误信息（如果失败）
	DurationMs    int64     `gorm:"not null" json:"duration_ms"`                    // 操作耗时（毫秒）
	UserAgent     string    `gorm:"size:500" json:"user_agent,omitempty"`           // 客户端信息
	RequestID     string    `gorm:"size:100;index" json:"request_id,omitempty"`     // 请求追踪ID
	OldData       string    `gorm:"type:jsonb" json:"old_data,omitempty"`           // 变更前数据（JSON）
	NewData       string    `gorm:"type:jsonb" json:"new_data,omitempty"`           // 变更后数据（JSON）
	CreatedAt     time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// 常用的业务域
const (
	DomainOrder    = "order"
	DomainUser     = "user"
	DomainProduct  = "product"
	DomainClient   = "client"
	DomainSupplier = "supplier"
	DomainMaterial = "material"
	DomainProcess  = "process"
	DomainPricing  = "pricing"
	DomainPlan     = "plan"
)

// 常用的操作动作
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionAssign = "assign"
)
