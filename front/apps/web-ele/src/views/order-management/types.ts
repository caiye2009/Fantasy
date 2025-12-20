// 进度项类型
export type ProgressItemType = 'fabric_input' | 'production' | 'warehouse_check' | 'rework'

// 订单状态
export type OrderStatus = 'in_progress' | 'completed'

// 角色类型
export type RoleType = 'sales' | 'follower' | 'warehouse' | 'system'

// 操作类型
export type OperationType =
  | 'create_order'
  | 'add_required_quantity'
  | 'update_fabric_input'
  | 'update_production'
  | 'update_warehouse_check'
  | 'add_defect'
  | 'update_rework'
  | 'generate_rework'

// 进度项
export interface ProgressItem {
  type: ProgressItemType
  name: string
  targetQuantity: number
  completedQuantity: number
  progress: number
  exists: boolean
  icon?: string
  color?: string
}

// 订单
export interface Order {
  id: string
  orderNo: string
  clientName: string
  productName: string
  productCode: string
  productHistoryShrinkage: number // 产品历史缩率（仅展示，不参与计算）
  requiredQuantity: number // 需求产品数量
  createdAt: string
  updatedAt: string
  status: OrderStatus
  progressItems: ProgressItem[]
  overallProgress: number // 订单整体进度 0-100
  operationLogs: OperationLog[]
}

// 操作日志
export interface OperationLog {
  id: string
  orderId: string
  operator: string
  operatorName: string
  role: RoleType
  operationType: OperationType
  operationName: string
  beforeData?: any
  afterData?: any
  description: string
  createdAt: string
}

// 操作表单
export interface OperationForm {
  quantity?: number
  defectQuantity?: number
  remark?: string
}
