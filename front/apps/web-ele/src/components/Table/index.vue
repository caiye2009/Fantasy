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
              :loading="filterLoadingStates[filter.key]"
              @visible-change="(visible) => handleFilterVisibleChange(visible, filter)"
              style="width: 200px"
            >
              <el-option
                v-for="option in filterOptions[filter.key] || []"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              />
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

        <!-- 中间列 -->
        <el-table-column
          v-for="column in visibleColumns"
          :key="column.key"
          :prop="column.key"
          :label="column.label"
          :min-width="120"
          :show-overflow-tooltip="true"
        >
          <template #default="scope">
            <span v-if="column.formatter">
              {{ column.formatter(scope.row[column.key], scope.row) }}
            </span>
            <span v-else>{{ scope.row[column.key] }}</span>
          </template>
        </el-table-column>

        <!-- 操作列固定右边 -->
        <el-table-column
          label="操作"
          width="100"
          fixed="right"
        >
          <template #default="scope">
            <el-button link type="primary" @click="handleView(scope.row)">查看</el-button>
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
  sort,
  selectedRows,
  selectedCount,
  hasMore,
  hasPrevious,
  initialize,
  reload,
  slideWindowDown,
  slideWindowUp,
} = useDataTable(props.config.index, props.config.pageSize)

// 本地状态
const tableWrapperRef = ref<HTMLElement>()
const searchQuery = ref('')
const filterForm = ref<Record<string, any>>({})
const filterOptions = ref<Record<string, Array<{ label: string; value: any }>>>({})
const filterLoadingStates = ref<Record<string, boolean>>({}) // 加载状态
const filterLoadedFlags = ref<Record<string, boolean>>({}) // 是否已加载

// 可见列
const visibleColumns = computed(() => columns.value.filter((col) => col.visible !== false))

// 滚动处理
let scrollTimer: number | null = null
const handleTableScroll = (e: Event) => {
  const target = e.target as HTMLElement
  if (!target) return
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
  if (val) tableData.forEach(row => { if (!selectedRows.value.includes(row)) selectedRows.value.push(row) })
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

// 刷新 & 查看
const handleRefresh = () => reload()
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

// 一次性获取所有 filter 选项（优化：合并聚合请求）
const loadAllFilterOptions = async () => {
  if (!props.config.filters) return

  // 收集所有需要聚合的 filter
  const aggRequests: Record<string, any> = {}
  const filtersNeedAgg: any[] = []

  for (const filter of props.config.filters) {
    // 跳过静态选项
    if (filter.options) {
      filterOptions.value[filter.key] = filter.options
      filterLoadedFlags.value[filter.key] = true
      continue
    }

    // 收集需要聚合的字段
    if (filter.fetchOptions && filter.key) {
      filtersNeedAgg.push(filter)
      // 假设 fetchOptions 内部使用聚合请求
      aggRequests[filter.key] = {
        type: 'terms',
        field: filter.key,
        size: 100, // 一次性获取更多数据
      }
    }
  }

  // 如果有需要聚合的字段，一次性请求
  if (Object.keys(aggRequests).length > 0) {
    try {
      const response = await elasticsearchService.search({
        index: props.config.index,
        pagination: { offset: 0, size: 0 },
        aggRequests,
      })

      // 处理聚合结果
      for (const filter of filtersNeedAgg) {
        const buckets = response.aggregations?.[filter.key]?.buckets || []
        filterOptions.value[filter.key] = buckets.map((bucket: any) => ({
          label: String(bucket.key),
          value: bucket.key,
        }))
        filterLoadedFlags.value[filter.key] = true
      }
    } catch (error) {
      console.error('Failed to load filter options:', error)
      // 失败时回退到逐个加载
      for (const filter of filtersNeedAgg) {
        filterOptions.value[filter.key] = []
      }
    }
  }
}

// 加载单个 filter 的选项（懒加载备用）
const loadFilterOption = async (filter: any) => {
  // 如果已经加载过，不重复加载
  if (filterLoadedFlags.value[filter.key]) {
    return
  }

  // 如果有静态 options，直接使用
  if (filter.options) {
    filterOptions.value[filter.key] = filter.options
    filterLoadedFlags.value[filter.key] = true
    return
  }

  // 如果有 fetchOptions 函数，调用它获取选项
  if (filter.fetchOptions) {
    filterLoadingStates.value[filter.key] = true
    try {
      filterOptions.value[filter.key] = await filter.fetchOptions()
      filterLoadedFlags.value[filter.key] = true
    } catch (error) {
      console.error(`Failed to load options for filter ${filter.key}:`, error)
      filterOptions.value[filter.key] = []
    } finally {
      filterLoadingStates.value[filter.key] = false
    }
  } else {
    // 默认为空数组
    filterOptions.value[filter.key] = []
    filterLoadedFlags.value[filter.key] = true
  }
}

// 处理下拉框显示/隐藏
const handleFilterVisibleChange = (visible: boolean, filter: any) => {
  if (visible && !filterLoadedFlags.value[filter.key]) {
    // 如果使用懒加载模式，才在这里加载
    if (!props.config.eagerLoadFilters) {
      loadFilterOption(filter)
    }
  }
}

// 初始化
onMounted(async () => {
  // 如果配置了 eagerLoadFilters，一次性加载所有 filter 选项
  if (props.config.eagerLoadFilters) {
    await loadAllFilterOptions()
  }
  // 否则使用懒加载，在用户点击下拉框时才加载

  initialize()
  const tableBody = tableWrapperRef.value?.querySelector('.el-table__body-wrapper')
  if (tableBody) tableBody.addEventListener('scroll', handleTableScroll)
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
}

.rounded-table { border-radius: 8px; overflow: hidden; }

.loading-tip, .no-more-tip {
  text-align: center;
  padding: 20px;
  color: #909399;
  font-size: 14px;
  .el-icon { margin-right: 8px; }
}
</style>
