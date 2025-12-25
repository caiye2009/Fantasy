import { requestClient } from '../request'

export interface AssignDepartmentRequest {
  department: string
}

export interface AssignPersonnelRequest {
  productionAssistant: number
  productionSpecialist: number
  follower: number
}

export interface SetFabricTargetRequest {
  targetQuantity: number
}

export interface UpdateProgressRequest {
  quantity: number
  remark?: string
}

export interface AddDefectRequest {
  quantity: number
  remark?: string
}

export interface OrderParticipant {
  id: number
  orderId: number
  userId: number
  userName: string
  role: string
  canEdit: boolean
  createdAt: string
}

export interface OrderProgress {
  id: number
  orderId: number
  progressType: string
  targetQuantity: number
  completedQuantity: number
  progress: number
  exists: boolean
  createdAt: string
  updatedAt: string
}

export interface OrderEvent {
  id: number
  orderId: number
  eventType: string
  operatorId: number
  operatorName: string
  operatorRole: string
  beforeData?: string
  afterData?: string
  description: string
  createdAt: string
}

// 分配部门
export async function assignDepartment(orderId: number, data: AssignDepartmentRequest) {
  return requestClient.post(`/api/v1/orders/${orderId}/assign-department`, data)
}

// 分配人员（初始化进度）
export async function assignPersonnel(orderId: number, data: AssignPersonnelRequest) {
  return requestClient.post(`/api/v1/orders/${orderId}/assign-personnel`, data)
}

// 设置胚布目标数量
export async function setFabricTarget(orderId: number, data: SetFabricTargetRequest) {
  return requestClient.post(`/api/v1/orders/${orderId}/fabric-target`, data)
}

// 更新胚布投入进度
export async function updateFabricInput(orderId: number, data: UpdateProgressRequest) {
  return requestClient.post(`/api/v1/orders/${orderId}/fabric-input`, data)
}

// 更新生产进度
export async function updateProduction(orderId: number, data: UpdateProgressRequest) {
  return requestClient.post(`/api/v1/orders/${orderId}/production`, data)
}

// 更新验货进度
export async function updateWarehouseCheck(orderId: number, data: UpdateProgressRequest) {
  return requestClient.post(`/api/v1/orders/${orderId}/warehouse-check`, data)
}

// 录入次品（自动生成回修进度）
export async function addDefect(orderId: number, data: AddDefectRequest) {
  return requestClient.post(`/api/v1/orders/${orderId}/defects`, data)
}

// 更新回修进度
export async function updateRework(orderId: number, data: UpdateProgressRequest) {
  return requestClient.post(`/api/v1/orders/${orderId}/rework`, data)
}

// 获取订单参与者
export async function getOrderParticipants(orderId: number) {
  return requestClient.get<OrderParticipant[]>(`/api/v1/orders/${orderId}/participants`)
}

// 获取订单进度
export async function getOrderProgress(orderId: number) {
  return requestClient.get<OrderProgress[]>(`/api/v1/orders/${orderId}/progress`)
}

// 获取订单事件
export async function getOrderEvents(orderId: number) {
  return requestClient.get<OrderEvent[]>(`/api/v1/orders/${orderId}/events`)
}

// 订单详情响应（完整信息）
export interface OrderDetailResponse {
  id: number
  order_no: string
  client_id: number
  client_name: string
  product_id: number
  product_name: string
  product_code: string
  product_history_shrinkage: number
  required_quantity: number
  unit_price: number
  total_price: number
  status: string
  assigned_department?: string
  created_at: string
  updated_at: string
  progress_items: OrderProgress[]
  operation_logs: OrderEvent[]
  overall_progress: number
}

// 订单列表响应（含详情）
export interface OrderListDetailResponse {
  total: number
  orders: OrderDetailResponse[]
}

// 获取订单详情（含完整信息）
export async function getOrderDetail(orderId: number) {
  return requestClient.get<OrderDetailResponse>(`/order/${orderId}/detail`)
}

// 获取订单列表（含完整信息）
export async function getOrderList(params?: { limit?: number; offset?: number }) {
  return requestClient.get<OrderListDetailResponse>('/order/list-detail', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0
    }
  })
}
