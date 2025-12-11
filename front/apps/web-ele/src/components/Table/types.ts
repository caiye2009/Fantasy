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
}

export interface BulkAction {
  label: string
  key: string
  type?: 'primary' | 'success' | 'warning' | 'danger'
  icon?: string
  confirm?: boolean
  confirmMessage?: string
}

export interface PageConfig {
  pageType: string
  title: string
  indices: string[]
  columns: ColumnConfig[]
  filters: FilterConfig[]
  bulkActions: BulkAction[]
  pageSize: number
}

export interface ESRequest {
  fields?: string[]
  filters?: Record<string, any>
  from: number
  indices: string[]
  query?: string
  size: number
  sort?: SortConfig[]
}

export interface ESResponse {
  total: number
  took: number
  max_score: number
  results: Array<{
    index: string
    type: string
    id: string
    score: number
    source: any
  }>
}

export interface CacheWindow {
  startPage: number
  endPage: number
  data: Map<number, any[]>
  currentPage: number
  totalPages: number
  totalCount: number
}