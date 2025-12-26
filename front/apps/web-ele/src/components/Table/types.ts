export interface ColumnConfig {
  key: string
  label: string
  width?: number
  sortable?: boolean
  filterable?: boolean
  visible?: boolean
  order?: number
  formatter?: (value: any, row: any) => string
}

export interface SortConfig {
  field: string
  order: 'asc' | 'desc'
}

export interface FilterConfig {
  key: string
  label: string
  type: 'text' | 'select' | 'date' | 'daterange' | 'number'
  options?: Array<{ label: string; value: any }>
  placeholder?: string
  // 动态获取选项的函数（用于从 ES 聚合或后端 API 获取）
  fetchOptions?: () => Promise<Array<{ label: string; value: any }>>
}

export interface BulkAction {
  label: string
  key: string
  type?: 'primary' | 'success' | 'warning' | 'danger'
  icon?: string
  confirm?: boolean
  confirmMessage?: string
}

export interface TopAction {
  label: string
  key: string
  type?: 'primary' | 'success' | 'warning' | 'danger'
  icon?: string
}

export interface PageConfig {
  pageType: string
  title?: string
  index: string // 单个实体类型，如 'material', 'order', 'client'
  columns: ColumnConfig[]
  filters?: FilterConfig[]
  bulkActions?: BulkAction[]
  topActions?: TopAction[]
  eagerLoadFilters?: boolean // 是否在页面加载时一次性获取所有 filter 选项（默认 false，使用懒加载）
  pageSize: number
  actions?: any[] // 可选的自定义操作
}

export interface ESRequest {
  index: string // material, order, client, etc. (required)
  query?: string
  filters?: Record<string, any>
  aggRequests?: Record<string, any>
  pagination: {
    offset: number
    size: number
  }
  sort?: SortConfig[]
}

export interface ESResponse {
  items: any[] // 搜索结果
  total: number
  took: number
  aggregations?: Record<string, any> // 聚合结果
}

export interface CacheWindow {
  startPage: number
  endPage: number
  data: Map<number, any[]>
  currentPage: number
  totalPages: number
  totalCount: number
}