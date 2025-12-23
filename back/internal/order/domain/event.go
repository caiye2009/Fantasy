package domain

import (
	"time"
)

// 事件类型
const (
	EventTypeCreateOrder         = "create_order"          // 创建订单
	EventTypeAssignDepartment    = "assign_department"     // 分配部门
	EventTypeAssignPersonnel     = "assign_personnel"      // 分配人员
	EventTypeSetFabricTarget     = "set_fabric_target"     // 设定胚布目标数量
	EventTypeUpdateFabricInput   = "update_fabric_input"   // 更新胚布投入进度
	EventTypeUpdateProduction    = "update_production"     // 更新生产进度
	EventTypeUpdateWarehouseCheck = "update_warehouse_check" // 更新验货进度
	EventTypeAddDefect           = "add_defect"            // 录入次品
	EventTypeUpdateRework        = "update_rework"         // 更新回修进度
	EventTypeChangeParticipant   = "change_participant"    // 变更参与者
)

// OrderEvent 订单事件实体（核心）
type OrderEvent struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OrderID      uint      `gorm:"not null;index" json:"order_id"`                  // 订单ID
	EventType    string    `gorm:"size:50;not null;index" json:"event_type"`        // 事件类型
	OperatorID   uint      `gorm:"not null" json:"operator_id"`                     // 操作人ID
	OperatorName string    `gorm:"size:100;not null" json:"operator_name"`          // 操作人姓名
	OperatorRole string    `gorm:"size:50;not null" json:"operator_role"`           // 操作人角色
	BeforeData   string    `gorm:"type:jsonb" json:"before_data,omitempty"`         // 变更前数据（JSON）
	AfterData    string    `gorm:"type:jsonb" json:"after_data,omitempty"`          // 变更后数据（JSON）
	Description  string    `gorm:"type:text;not null" json:"description"`           // 事件描述
	CreatedAt    time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 表名
func (OrderEvent) TableName() string {
	return "order_events"
}

// IsSystemEvent 是否是系统事件
func (e *OrderEvent) IsSystemEvent() bool {
	return e.OperatorRole == "system"
}

// GetEventTypeName 获取事件类型的中文名称
func GetEventTypeName(eventType string) string {
	names := map[string]string{
		EventTypeCreateOrder:          "创建订单",
		EventTypeAssignDepartment:     "分配部门",
		EventTypeAssignPersonnel:      "分配人员",
		EventTypeSetFabricTarget:      "设定胚布目标",
		EventTypeUpdateFabricInput:    "更新胚布投入",
		EventTypeUpdateProduction:     "更新生产进度",
		EventTypeUpdateWarehouseCheck: "更新验货进度",
		EventTypeAddDefect:            "录入次品",
		EventTypeUpdateRework:         "更新回修进度",
		EventTypeChangeParticipant:    "变更参与者",
	}

	if name, ok := names[eventType]; ok {
		return name
	}
	return eventType
}
