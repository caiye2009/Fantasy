// 产品类型
export interface Product {
  id: string | number
  _id?: string // ES document ID
  name: string // 产品名称
  status?: string // 状态
  materials?: any[] // 原料配置
  processes?: any[] // 工艺配置
  created_at: string
  updated_at: string
}
