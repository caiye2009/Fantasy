package application

import "time"

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	OrderNo   string  `json:"order_no" binding:"required,max=50"`
	ClientID  uint    `json:"client_id" binding:"required"`
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"required,gte=0"`
	CreatedBy uint    `json:"created_by" binding:"required"`
}

// UpdateOrderRequest 更新订单请求
type UpdateOrderRequest struct {
	Status    string  `json:"status" binding:"omitempty,oneof=pending confirmed production completed cancelled"`
	Quantity  float64 `json:"quantity" binding:"omitempty,gt=0"`
	UnitPrice float64 `json:"unit_price" binding:"omitempty,gte=0"`
}

// OrderResponse 订单响应
type OrderResponse struct {
	ID         uint      `json:"id"`
	OrderNo    string    `json:"order_no"`
	ClientID   uint      `json:"client_id"`
	ProductID  uint      `json:"product_id"`
	Quantity   float64   `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedBy  uint      `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// OrderListResponse 订单列表响应
type OrderListResponse struct {
	Total  int64            `json:"total"`
	Orders []*OrderResponse `json:"orders"`
}

// ==================== 新增DTO ====================

// CreateOrderRequestV2 创建订单请求（新版）
type CreateOrderRequestV2 struct {
	OrderNo                 string  `json:"order_no" binding:"required,max=50"`
	ClientID                uint    `json:"client_id" binding:"required"`
	ProductID               uint    `json:"product_id" binding:"required"`
	RequiredQuantity        float64 `json:"required_quantity" binding:"required,gt=0"`         // 成品需求数量
	ProductHistoryShrinkage float64 `json:"product_history_shrinkage" binding:"omitempty,gte=0"` // 历史缩率
	UnitPrice               float64 `json:"unit_price" binding:"required,gte=0"`
}

// AssignDepartmentRequest 分配部门请求
type AssignDepartmentRequest struct {
	Department string `json:"department" binding:"required"`
}

// AssignPersonnelRequest 分配人员请求
type AssignPersonnelRequest struct {
	ProductionSpecialistID uint    `json:"production_specialist_id" binding:"required"`
	OrderCoordinatorID     uint    `json:"order_coordinator_id" binding:"required"`
	FabricTargetQuantity   float64 `json:"fabric_target_quantity" binding:"required,gt=0"` // 胚布目标数量
}

// UpdateProgressRequest 更新进度请求
type UpdateProgressRequest struct {
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	Remark   string  `json:"remark" binding:"omitempty"`
}

// AddDefectRequest 录入次品请求
type AddDefectRequest struct {
	DefectQuantity float64 `json:"defect_quantity" binding:"required,gt=0"`
	Remark         string  `json:"remark" binding:"required"`
}

// ParticipantResponse 参与者响应
type ParticipantResponse struct {
	ID         uint      `json:"id"`
	OrderID    uint      `json:"order_id"`
	UserID     uint      `json:"user_id"`
	UserName   string    `json:"user_name"`
	Role       string    `json:"role"`
	AssignedAt time.Time `json:"assigned_at"`
	AssignedBy uint      `json:"assigned_by"`
	IsActive   bool      `json:"is_active"`
}

// ProgressResponse 进度响应
type ProgressResponse struct {
	ID                uint      `json:"id"`
	OrderID           uint      `json:"order_id"`
	Type              string    `json:"type"`
	Name              string    `json:"name"` // 前端显示名称
	TargetQuantity    float64   `json:"target_quantity"`
	CompletedQuantity float64   `json:"completed_quantity"`
	Progress          int       `json:"progress"`
	Exists            bool      `json:"exists"`
	Icon              string    `json:"icon,omitempty"`
	Color             string    `json:"color,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// EventResponse 事件响应
type EventResponse struct {
	ID           uint      `json:"id"`
	OrderID      uint      `json:"order_id"`
	EventType    string    `json:"event_type"`
	OperatorID   uint      `json:"operator_id"`
	OperatorName string    `json:"operator_name"`
	OperatorRole string    `json:"operator_role"`
	BeforeData   string    `json:"before_data,omitempty"`
	AfterData    string    `json:"after_data,omitempty"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

// OrderDetailResponse 订单详情响应（包含客户名、产品名、进度项、事件等完整信息）
type OrderDetailResponse struct {
	ID                      uint                `json:"id"`
	OrderNo                 string              `json:"order_no"`
	ClientID                uint                `json:"client_id"`
	ClientName              string              `json:"client_name"`
	ProductID               uint                `json:"product_id"`
	ProductName             string              `json:"product_name"`
	ProductCode             string              `json:"product_code"`
	ProductHistoryShrinkage float64             `json:"product_history_shrinkage"`
	RequiredQuantity        float64             `json:"required_quantity"`
	UnitPrice               float64             `json:"unit_price"`
	TotalPrice              float64             `json:"total_price"`
	Status                  string              `json:"status"`
	AssignedDepartment      string              `json:"assigned_department,omitempty"`
	CreatedAt               time.Time           `json:"created_at"`
	UpdatedAt               time.Time           `json:"updated_at"`
	ProgressItems           []ProgressResponse  `json:"progress_items"`
	OperationLogs           []EventResponse     `json:"operation_logs"`
	OverallProgress         int                 `json:"overall_progress"`
}

// OrderListDetailResponse 订单列表响应（带详细信息）
type OrderListDetailResponse struct {
	Total  int64                  `json:"total"`
	Orders []*OrderDetailResponse `json:"orders"`
}