// 工艺类型
export interface Process {
  id: string | number
  _id?: string // ES document ID
  name: string // 工艺名称
  description?: string // 描述
  createdAt: string
  updatedAt: string
}
