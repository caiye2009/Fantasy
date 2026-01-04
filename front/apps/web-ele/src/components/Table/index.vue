<template>
  <div class="data-table-container">
    <!-- 筛选栏和顶部操作 -->
    <div class="filter-bar">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="搜索">
          <el-input
            v-model="searchQuery"
            placeholder="输入关键词搜索"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
            class="search-input"
          >
            <template #append>
              <el-button @click="handleSearch" :icon="Search" />
            </template>
          </el-input>
        </el-form-item>

        <div class="other-filters" v-if="config.filters && config.filters.length > 0">
          <el-form-item
            v-for="filter in config.filters"
            :key="filter.key"
            :label="filter.label"
          >
            <el-input
              v-if="filter.type === 'text'"
              v-model="filterForm[filter.key]"
              :placeholder="filter.placeholder"
              clearable
              style="width: 200px"
            />
            <el-select
              v-else-if="filter.type === 'select'"
              v-model="filterForm[filter.key]"
              :placeholder="filter.placeholder"
              clearable
              filterable
              remote
              :remote-method="(query) => handleFilterSearch(filter.key, query)"
              :loading="filterStates[filter.key]?.loading || false"
              @visible-change="(visible) => handleFilterVisibleChange(visible, filter)"
              @focus="() => handleFilterFocus(filter)"
              style="width: 200px"
            >
              <el-option
                v-for="option in filterStates[filter.key]?.options || []"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
              <template #footer>
                <div
                  v-if="filterStates[filter.key]?.hasMore"
                  class="select-footer"
                  @click="loadMoreFilterOptions(filter.key)"
                >
                  <el-button
                    link
                    type="primary"
                    :loading="filterStates[filter.key]?.loading"
                  >
                    {{ filterStates[filter.key]?.loading ? '加载中...' : '加载更多' }}
                  </el-button>
                </div>
              </template>
            </el-select>
            <el-date-picker
              v-else-if="filter.type === 'date'"
              v-model="filterForm[filter.key]"
              type="date"
              :placeholder="filter.placeholder"
              style="width: 200px"
            />
            <el-date-picker
              v-else-if="filter.type === 'daterange'"
              v-model="filterForm[filter.key]"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 280px"
            />
          </el-form-item>
        </div>

        <el-form-item class="filter-buttons" v-if="config.filters && config.filters.length > 0">
          <el-button type="primary" @click="handleFilter">查询</el-button>
          <el-button @click="handleResetFilter">重置</el-button>
        </el-form-item>

        <div class="top-actions" v-if="config.topActions && config.topActions.length > 0">
          <el-button
            v-for="action in config.topActions"
            :key="action.key"
            :type="action.type || 'default'"
            @click="handleTopAction(action)"
          >
            <el-icon v-if="action.icon"><component :is="action.icon" /></el-icon>
            {{ action.label }}
          </el-button>
        </div>
      </el-form>
    </div>

    <!-- 批量操作悬浮栏 -->
    <transition name="el-fade-in">
      <div
        class="bulk-action-bar-floating"
        v-if="selectedCount > 0"
      >
        <div class="selected-info">
          已选择 <strong>{{ selectedCount }}</strong> 项
          <el-button link @click="handleClearSelection">清空</el-button>
        </div>
        <div class="bulk-actions">
          <el-button
            v-for="action in config.bulkActions"
            :key="action.key"
            :type="action.type || 'default'"
            @click="handleBulkAction(action)"
          >
            <el-icon v-if="action.icon"><component :is="action.icon" /></el-icon>
            {{ action.label }}
          </el-button>
        </div>
      </div>
    </transition>

    <!-- 表格 -->
    <div class="table-wrapper" ref="tableWrapperRef">
      <el-table
        v-loading="props.loading"
        :data="tableData"
        stripe
        style="width: 100%"
        height="calc(100vh - 180px)"
        :header-cell-style="{ background: '#f5f7fa', textAlign: 'left' }"
        @selection-change="handleSelectionChange"
        :reserve-selection="true"
        @select-all="handleSelectAll"
        class="rounded-table"
      >
        <!-- 多选列固定左边 -->
        <el-table-column
          type="selection"
          width="55"
          fixed="left"
        />

        <!-- 中间列 - 统一列宽 -->
        <el-table-column
          v-for="column in visibleColumns"
          :key="column.key"
          :prop="column.key"
          :label="column.label"
          :width="column.width || 150"
          :show-overflow-tooltip="true"
        >
          <template #default="scope">
            <span v-if="column.formatter">
              {{ column.formatter(scope.row[column.key], scope.row) }}
            </span>
            <span v-else>{{ scope.row[column.key] || '-' }}</span>
          </template>
        </el-table-column>

        <!-- 操作列固定右边 -->
        <el-table-column
          label="操作"
          width="120"
          fixed="right"
          align="center"
        >
          <template #default="scope">
            <el-button link type="primary" size="small" @click="handleView(scope.row)">
              查看
            </el-button>
          </template>
        </el-table-column>

        <template #append>
          <div v-if="loading" class="loading-tip">
            <el-icon class="is-loading"><Loading /></el-icon>
            加载中...
          </div>
          <div v-else-if="!hasMore && tableData.length > 0" class="no-more-tip">
            没有更多数据了
          </div>
        </template>

        <!-- 空状态 -->
        <template #empty>
          <div class="empty-state">
            <el-empty :description="getEmptyDescription()" />
          </div>
        </template>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { Search, Loading } from '@element-plus/icons-vue'
import { useDataTable } from '#/composables/useDataTable'
import { useTablePreference } from '#/composables/useTablePreference'
import { elasticsearchService } from '#/api/core/es'
import type { PageConfig, BulkAction } from './types'

interface Props {
  config: PageConfig
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), { loading: false })
const emit = defineEmits(['view', 'bulkAction', 'topAction'])

// 用户偏好
const { columns } = useTablePreference(
  props.config.pageType,
  props.config.columns
)

// 数据表格
const {
  loading,
  tableData,
  query,
  filters,
  selectedRows,
  selectedCount,
  hasMore,
  hasPrevious,
  scrollPosition,
  initialize,
  reload,
  slideWindowDown,
  slideWindowUp,
  saveScrollPosition,
} = useDataTable(props.config.index, props.config.pageSize)

// 本地状态
const tableWrapperRef = ref<HTMLElement>()
const searchQuery = ref('')
const filterForm = ref<Record<string, any>>({})

// 筛选器状态（支持分页和搜索）
interface FilterState {
  options: Array<{ label: string; value: any }>
  searchQuery: string
  page: number
  pageSize: number
  hasMore: boolean
  loading: boolean
  initialized: boolean
}
const filterStates = ref<Record<string, FilterState>>({})

// 可见列
const visibleColumns = computed(() => columns.value.filter((col) => col.visible !== false))

// 滚动处理
let scrollTimer: number | null = null
const handleTableScroll = (e: Event) => {
  const target = e.target as HTMLElement
  if (!target) return

  // 保存滚动位置
  saveScrollPosition(target.scrollTop)

  if (scrollTimer) clearTimeout(scrollTimer)
  scrollTimer = window.setTimeout(() => {
    const scrollTop = target.scrollTop
    const scrollHeight = target.scrollHeight
    const clientHeight = target.clientHeight
    if (scrollHeight - scrollTop - clientHeight < 100 && hasMore.value && !loading.value) slideWindowDown()
    if (scrollTop < 100 && hasPrevious.value && !loading.value && scrollTop > 0) slideWindowUp()
  }, 150)
}

// 搜索 & 筛选
const handleSearch = () => { query.value = searchQuery.value }
const handleFilter = () => {
  const activeFilters: Record<string, any> = {}
  Object.keys(filterForm.value).forEach((key) => {
    if (filterForm.value[key] !== null && filterForm.value[key] !== '') activeFilters[key] = filterForm.value[key]
  })
  filters.value = activeFilters
}
const handleResetFilter = () => { filterForm.value = {}; searchQuery.value = ''; query.value = ''; filters.value = {} }

// 选择
const handleSelectionChange = (selection: any[]) => {
  selectedRows.value = selection
}
const handleClearSelection = () => { selectedRows.value = [] }
const handleSelectAll = (val: boolean) => {
  if (val) tableData.value.forEach((row: any) => { if (!selectedRows.value.includes(row)) selectedRows.value.push(row) })
  else selectedRows.value = []
}

// 批量操作
const handleBulkAction = async (action: BulkAction) => {
  if (selectedRows.value.length === 0) return
  if (action.confirm) {
    try {
      await ElMessageBox.confirm(
        action.confirmMessage || `确定要${action.label}选中的 ${selectedCount.value} 项吗？`,
        '确认操作', { confirmButtonText: '确定', cancelButtonText: '取消', type: 'warning' }
      )
    } catch { return }
  }
  emit('bulkAction', { action: action.key, rows: selectedRows.value })
}

// 查看
const handleView = (row: any) => emit('view', row)

// 顶部操作
const handleTopAction = (action: any) => {
  emit('topAction', { action: action.key })
}

// 获取空状态描述
const getEmptyDescription = () => {
  const pageTypeMap: Record<string, string> = {
    material: '原料',
    client: '客户',
    order: '订单',
    process: '工序',
    product: '产品',
    inventory: '库存'
  }

  const pageName = pageTypeMap[props.config.pageType || ''] || '数据'
  return `请新增${pageName}`
}

// 初始化筛选器状态
const initFilterState = (filterKey: string) => {
  if (!filterStates.value[filterKey]) {
    filterStates.value[filterKey] = {
      options: [],
      searchQuery: '',
      page: 0,
      pageSize: 20, // 每页20个选项
      hasMore: true,
      loading: false,
      initialized: false,
    }
  }
}

// 获取筛选器配置
const getFilterConfig = (filterKey: string) => {
  return props.config.filters?.find((f) => f.key === filterKey)
}

// 加载筛选器选项（支持分页和搜索）
const loadFilterOptions = async (
  filterKey: string,
  searchQuery: string = '',
  page: number = 1,
  append: boolean = false
) => {
  const state = filterStates.value[filterKey]
  if (!state || state.loading) return

  const filterConfig = getFilterConfig(filterKey)
  if (!filterConfig) return

  state.loading = true

  try {
    let newOptions: Array<{ label: string; value: any }> = []

    // 如果配置了 fetchOptions，优先使用它（向后兼容）
    if (filterConfig.fetchOptions) {
      newOptions = await filterConfig.fetchOptions()

      // 如果有搜索词，在客户端过滤
      if (searchQuery) {
        newOptions = newOptions.filter((opt) =>
          String(opt.label).toLowerCase().includes(searchQuery.toLowerCase())
        )
      }

      // 模拟分页
      if (!append) {
        newOptions = newOptions.slice(0, state.pageSize)
      }
    } else {
      // 否则使用默认的 ES 聚合方式
      const aggRequests: Record<string, any> = {
        [filterKey]: {
          type: 'terms',
          field: filterKey,
          size: state.pageSize,
          // 如果有搜索词，添加搜索过滤
          ...(searchQuery ? { include: `.*${searchQuery}.*` } : {}),
        }
      }

      const response = await elasticsearchService.search({
        index: props.config.index,
        pagination: { offset: 0, size: 0 },
        aggRequests,
      })

      const buckets = response.aggregations?.[filterKey]?.buckets || []
      newOptions = buckets.map((bucket: any) => ({
        label: String(bucket.key),
        value: bucket.key,
      }))
    }

    // 更新状态
    if (append) {
      state.options = [...state.options, ...newOptions]
    } else {
      state.options = newOptions
    }

    state.page = page
    state.searchQuery = searchQuery
    state.hasMore = newOptions.length === state.pageSize
    state.initialized = true
  } catch (error) {
    console.error(`Failed to load options for filter ${filterKey}:`, error)
    state.options = []
    state.hasMore = false
  } finally {
    state.loading = false
  }
}

// 处理筛选器搜索
const handleFilterSearch = async (filterKey: string, searchQuery: string) => {
  initFilterState(filterKey)
  await loadFilterOptions(filterKey, searchQuery, 1, false)
}

// 加载更多筛选器选项
const loadMoreFilterOptions = async (filterKey: string) => {
  const state = filterStates.value[filterKey]
  if (!state) return
  if (!state.hasMore || state.loading) return

  await loadFilterOptions(filterKey, state.searchQuery, state.page + 1, true)
}

// 处理下拉框聚焦
const handleFilterFocus = async (filter: any) => {
  initFilterState(filter.key)
  const state = filterStates.value[filter.key]

  // 如果未初始化，加载第一页
  if (!state.initialized) {
    await loadFilterOptions(filter.key, '', 1, false)
  }
}

// 处理下拉框显示/隐藏
const handleFilterVisibleChange = (_visible: boolean, _filter: any) => {
  // 可以在这里添加其他逻辑，比如清空搜索等
}

// 初始化
onMounted(async () => {
  // 初始化所有筛选器状态（懒加载，在用户聚焦时才加载数据）
  if (props.config.filters) {
    for (const filter of props.config.filters) {
      if (filter.type === 'select') {
        initFilterState(filter.key)
      }
    }
  }

  await initialize()

  const tableBody = tableWrapperRef.value?.querySelector('.el-table__body-wrapper')
  if (tableBody) {
    tableBody.addEventListener('scroll', handleTableScroll)

    // 恢复滚动位置
    if (scrollPosition.value > 0) {
      setTimeout(() => {
        tableBody.scrollTop = scrollPosition.value
      }, 100)
    }
  }
})
onUnmounted(() => {
  const tableBody = tableWrapperRef.value?.querySelector('.el-table__body-wrapper')
  if (tableBody) tableBody.removeEventListener('scroll', handleTableScroll)
  if (scrollTimer) clearTimeout(scrollTimer)
})
</script>

<style scoped lang="scss">
.data-table-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-radius: 8px;
  position: relative;
}

.filter-bar {
  padding: 16px 20px;
  background: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.filter-form {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;

  .search-input { width: 300px; }
  .other-filters { display: flex; gap: 12px; flex-wrap: wrap; }
  .filter-buttons { display: flex; gap: 8px; }
  .top-actions {
    margin-left: auto;
    display: flex;
    gap: 8px;
  }
}

/* 悬浮批量操作栏靠右，不超过表格一半宽度 */
.bulk-action-bar-floating {
  position: absolute;
  top: 16px;
  right: 20px;
  max-width: 50%;
  z-index: 10;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  background: #ecf5ff;
  padding: 10px 16px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.selected-info { color: #409eff; font-size: 14px; margin-right: 12px; strong { font-size: 16px; margin: 0 4px; } }
.bulk-actions { display: flex; gap: 8px; }

.table-wrapper {
  flex: 1;
  overflow: hidden;
  padding: 0 20px 20px;
  position: relative;

  // 隐藏 el-table 内部的滚动条，但保持滚动功能
  :deep(.el-table__body-wrapper) {
    // 隐藏滚动条但保持滚动功能
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE and Edge */

    &::-webkit-scrollbar {
      display: none; /* Chrome, Safari, Opera */
    }
  }

  // 确保外部不产生滚动条
  :deep(.el-table) {
    overflow: hidden;
  }

  // 优化表格样式
  :deep(.el-table__row) {
    transition: background-color 0.2s;

    &:hover {
      background-color: #f5f7fa !important;
    }
  }

  // 优化表头样式
  :deep(.el-table__header-wrapper) {
    .el-table__header th {
      font-weight: 600;
      color: #606266;
    }
  }

  // 优化单元格内容
  :deep(.el-table__cell) {
    padding: 12px 0;

    .cell {
      padding: 0 12px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
}

.rounded-table { border-radius: 8px; overflow: hidden; }

.loading-tip, .no-more-tip {
  text-align: center;
  padding: 20px;
  color: #909399;
  font-size: 14px;
  .el-icon { margin-right: 8px; }
}

.select-footer {
  text-align: center;
  padding: 8px 0;
  cursor: pointer;

  &:hover {
    background-color: #f5f7fa;
  }
}
</style>
