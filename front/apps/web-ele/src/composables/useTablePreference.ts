import { ref, watch } from 'vue'
import type { ColumnConfig } from '#/components/Table/types'

interface ColumnPreference {
  visible: boolean
  order: number
  width?: number
}

interface TablePreference {
  columns: Record<string, ColumnPreference>
  defaultSort?: { field: string; order: 'asc' | 'desc' }
}

export function useTablePreference(pageType: string, defaultColumns: ColumnConfig[]) {
  const storageKey = `table-preference-${pageType}`

  // 加载偏好
  const loadPreference = (): TablePreference | null => {
    const stored = localStorage.getItem(storageKey)
    return stored ? JSON.parse(stored) : null
  }

  // 保存偏好
  const savePreference = (preference: TablePreference) => {
    localStorage.setItem(storageKey, JSON.stringify(preference))
  }

  // 合并默认配置和用户偏好
  const mergeColumns = (columns: ColumnConfig[]): ColumnConfig[] => {
    const preference = loadPreference()
    if (!preference) return columns

    return columns
      .map((col) => {
        const pref = preference.columns[col.key]
        if (pref) {
          return {
            ...col,
            visible: pref.visible,
            order: pref.order,
            width: pref.width || col.width,
          }
        }
        return col
      })
      .sort((a, b) => (a.order || 0) - (b.order || 0))
  }

  const columns = ref<ColumnConfig[]>(mergeColumns(defaultColumns))

  // 更新列配置
  const updateColumns = (newColumns: ColumnConfig[]) => {
    columns.value = newColumns
    const preference: TablePreference = {
      columns: {},
    }
    newColumns.forEach((col, index) => {
      preference.columns[col.key] = {
        visible: col.visible !== false,
        order: index,
        width: col.width,
      }
    })
    savePreference(preference)
  }

  // 重置为默认
  const resetToDefault = () => {
    localStorage.removeItem(storageKey)
    columns.value = [...defaultColumns]
  }

  return {
    columns,
    updateColumns,
    resetToDefault,
  }
}