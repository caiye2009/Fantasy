export type DrawerMode = 'create' | 'edit' | 'view'

export type FieldType = 
  | 'text'
  | 'textarea'
  | 'number'
  | 'select'
  | 'date'
  | 'datetime'
  | 'daterange'
  | 'switch'
  | 'radio'
  | 'checkbox'
  | 'custom'

export interface FieldConfig {
  key: string
  label: string
  type: FieldType
  required?: boolean
  disabled?: boolean
  placeholder?: string
  defaultValue?: any
  requiredMessage?: string
  
  // 数字类型
  min?: number
  max?: number
  precision?: number
  step?: number
  
  // 文本域
  rows?: number
  
  // 选择框
  options?: Array<{ label: string; value: any }>
  multiple?: boolean
  filterable?: boolean
  
  // 自定义验证
  validator?: (rule: any, value: any, callback: any) => void
  
  // 异步加载选项
  fetchOptions?: () => Promise<Array<{ label: string; value: any }>>
}

export interface DrawerConfig {
  title: string
  fields: FieldConfig[]
  size?: string | number
  onSubmit?: (data: Record<string, any>, mode: DrawerMode) => Promise<void>
  beforeOpen?: (entity: any, mode: DrawerMode) => any
}