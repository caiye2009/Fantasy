<template>
  <div class="inventory-table">
    <div class="table-header">
      <h3>{{ title }}</h3>
      <el-space>
        <el-input
          v-model="searchText"
          placeholder="搜索..."
          clearable
          style="width: 300px"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button @click="handleExport">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </el-space>
    </div>

    <el-table
      :data="filteredData"
      stripe
      border
      height="600"
      :default-sort="{ prop: 'inboundDate', order: 'descending' }"
    >
      <el-table-column
        v-for="col in columns"
        :key="col.prop"
        :prop="col.prop"
        :label="col.label"
        :width="col.width"
        :align="col.align || 'left'"
        :sortable="col.sortable !== false"
      >
        <template #default="{ row }">
          <span v-if="col.type === 'currency'" class="currency">
            ¥{{ row[col.prop].toFixed(2) }}
          </span>
          <span
            v-else-if="col.type === 'difference'"
            :class="getDifferenceClass(row[col.prop])"
          >
            {{ row[col.prop] >= 0 ? '+' : '' }}¥{{ row[col.prop].toFixed(2) }}
          </span>
          <span v-else>
            {{ row[col.prop] }}
          </span>
        </template>
      </el-table-column>
    </el-table>

    <div class="table-footer">
      <el-text>共 {{ filteredData.length }} 条记录</el-text>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Download } from '@element-plus/icons-vue'

interface Column {
  prop: string
  label: string
  width?: number
  align?: 'left' | 'center' | 'right'
  type?: 'currency' | 'difference'
  sortable?: boolean
}

interface Props {
  data: any[]
  columns: Column[]
  title: string
}

const props = defineProps<Props>()

const searchText = ref('')

// 过滤数据
const filteredData = computed(() => {
  if (!searchText.value) {
    return props.data
  }

  const keyword = searchText.value.toLowerCase()
  return props.data.filter(item => {
    return Object.values(item).some(value => {
      if (typeof value === 'string') {
        return value.toLowerCase().includes(keyword)
      }
      return false
    })
  })
})

// 获取差价样式类
const getDifferenceClass = (value: number) => {
  if (value > 0) return 'positive-difference'
  if (value < 0) return 'negative-difference'
  return ''
}

// 导出数据
const handleExport = () => {
  try {
    const headers = props.columns.map(col => col.label)
    const csvContent = [
      headers.join(','),
      ...filteredData.value.map(row =>
        props.columns.map(col => {
          const value = row[col.prop]
          if (typeof value === 'number') {
            return value
          }
          return `"${value || ''}"`
        }).join(',')
      )
    ].join('\n')

    const blob = new Blob(['\ufeff' + csvContent], {
      type: 'text/csv;charset=utf-8;'
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${props.title}_${Date.now()}.csv`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  }
}
</script>

<style scoped lang="scss">
.inventory-table {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }
  }

  .currency {
    font-weight: 600;
    color: #409eff;
  }

  .positive-difference {
    font-weight: 600;
    color: #67c23a;
  }

  .negative-difference {
    font-weight: 600;
    color: #f56c6c;
  }

  .table-footer {
    margin-top: 16px;
    text-align: right;
    color: #909399;
  }
}
</style>
