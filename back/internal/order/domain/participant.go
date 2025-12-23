package domain

import (
	"time"
)

// 参与角色类型
const (
	ParticipantRoleCreator               = "creator"                // 创建者（Sales）
	ParticipantRoleProductionDirector    = "production_director"    // 生产总监
	ParticipantRoleProductionAssistant   = "production_assistant"   // 生产助理
	ParticipantRoleProductionSpecialist  = "production_specialist"  // 生产专员
	ParticipantRoleOrderCoordinator      = "order_coordinator"      // 跟单
	ParticipantRoleWarehouse             = "warehouse"              // 仓管
)

// OrderParticipant 订单参与者实体
type OrderParticipant struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `gorm:"not null;index" json:"order_id"`                     // 订单ID
	UserID     uint      `gorm:"not null;index" json:"user_id"`                      // 参与者用户ID
	UserName   string    `gorm:"size:100;not null" json:"user_name"`                 // 参与者姓名（冗余，便于查询）
	Role       string    `gorm:"size:50;not null" json:"role"`                       // 参与角色
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`                  // 分配时间
	AssignedBy uint      `gorm:"not null" json:"assigned_by"`                        // 分配者ID
	IsActive   bool      `gorm:"default:true;index" json:"is_active"`                // 是否当前激活（支持人员替换）
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 表名
func (OrderParticipant) TableName() string {
	return "order_participants"
}

// Deactivate 停用参与者
func (p *OrderParticipant) Deactivate() {
	p.IsActive = false
}

// Activate 激活参与者
func (p *OrderParticipant) Activate() {
	p.IsActive = true
}
