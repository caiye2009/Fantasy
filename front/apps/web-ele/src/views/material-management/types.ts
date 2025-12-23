// 原料类型
export interface Material {
  id: string | number
  _id?: string // ES document ID
  code?: string // 原料编号（唯一）
  name: string // 原料名称
  category?: string // 分类（胚布、染料、助剂等）
  spec?: string // 规格
  unit: string // 单位（米、kg）
  description?: string // 描述
  currentPrice?: number // 当前价格（最新供应商报价）
  status?: 'active' | 'inactive' // 状态
  updatedBy?: string // 最新更新人
  createdAt: string
  updatedAt: string
}

// 供应商报价
export interface SupplierQuote {
  id: string
  materialId: string // 原料ID
  supplierName: string // 供应商名称
  price: number // 价格
  quotedBy: string // 报价人
  quotedAt: string // 报价时间
  remark?: string // 备注
}

// 历史价格记录（保留用于其他功能）
export interface PriceHistory {
  id: string
  materialId: string
  price: number
  effectiveDate: string // 生效日期
  supplier: string
  purchaseQuantity?: number // 采购数量
  remark?: string
}

// 库存记录
export interface StockRecord {
  id: string
  materialId: string
  type: 'in' | 'out' // 入库/出库
  quantity: number
  unitPrice: number // 单价
  totalValue: number // 总价值
  relatedOrder?: string // 关联订单
  operator: string
  remark?: string
  createdAt: string
}

// 新增原料表单
export interface MaterialForm {
  code: string
  name: string
  category: string
  unit: 'kg' | 'm' // 限定为 kg 或 m
  currentPrice: number
}
