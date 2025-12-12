<template>
  <div class="process-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="handleView"
      @edit="handleEdit"
      @bulkAction="handleBulkAction"
    />

    <!-- 详情/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="600px"
      @close="handleDialogClose"
    >
      <el-form :model="currentRow" label-width="100px">
        <el-form-item label="ID">
          <el-input v-model="currentRow.id" disabled />
        </el-form-item>

        <el-form-item label="工序名称">
          <el-input v-model="currentRow.name" :disabled="dialogMode === 'view'" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            v-model="currentRow.description"
            type="textarea"
            :rows="3"
            :disabled="dialogMode === 'view'"
          />
        </el-form-item>

        <el-form-item label="创建时间" v-if="currentRow.created_at">
          <el-input
            :value="new Date(currentRow.created_at).toLocaleString('zh-CN')"
            disabled
          />
        </el-form-item>

        <el-form-item label="更新时间" v-if="currentRow.updated_at">
          <el-input
            :value="new Date(currentRow.updated_at).toLocaleString('zh-CN')"
            disabled
          />
        </el-form-item>
      </el-form>

      <template #footer v-if="dialogMode === 'edit'">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import DataTable from '#/components/Table/index.vue'
import { elasticsearchService } from '#/api/core/es'

import type {
  PageConfig,
  ColumnConfig,
  FilterConfig,
  BulkAction,
} from '#/components/Table/types'

// 复用 useDataTable 的 searchLoading
import { useDataTable } from '#/composables/useDataTable'
const { searchLoading } = useDataTable(['process'], 20)

// 配置
const pageConfig: PageConfig = {
  pageType: 'process',
  title: '工序管理',
  indices: ['process'],
  pageSize: 20,
  columns: [
    { key: 'id', label: 'ID', width: 80, visible: true, sortable: true, order: 0 },
    { key: 'name', label: '工序名称', width: 200, visible: true, sortable: true, order: 1 },
    { key: 'description', label: '描述', width: 250, visible: true, order: 2 },
    {
      key: 'created_at',
      label: '创建时间',
      width: 180,
      visible: true,
      sortable: true,
      order: 3,
      formatter: (value: string) => value ? new Date(value).toLocaleString('zh-CN') : '-'
    },
    {
      key: 'updated_at',
      label: '更新时间',
      width: 180,
      visible: false,
      sortable: true,
      order: 4,
      formatter: (value: string) => value ? new Date(value).toLocaleString('zh-CN') : '-'
    }
  ] as ColumnConfig[],
  filters: [
    { key: 'name', label: '工序名称', type: 'text', placeholder: '请输入工序名称' },
  ] as FilterConfig[],
  bulkActions: [
    { key: 'delete', label: '批量删除', type: 'danger', confirm: true, confirmMessage: '确定要删除选中的工序吗？此操作不可恢复！' },
    { key: 'export', label: '导出数据', type: 'primary' },
  ] as BulkAction[],
}

// dialog 状态
const dialogVisible = ref(false)
const dialogMode = ref<'view' | 'edit'>('view')
const currentRow = ref<any>({})
const saving = ref(false)

// 动态标题
const dialogTitle = computed(() =>
  dialogMode.value === 'view' ? '查看工序' : '编辑工序'
)

// 查看
const handleView = (row: any) => {
  dialogMode.value = 'view'
  currentRow.value = { ...row }
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogMode.value = 'edit'
  currentRow.value = { ...row }
  dialogVisible.value = true
}

// 保存
const handleSave = async () => {
  if (!currentRow.value._id) {
    ElMessage.error('缺少必要的 ID 信息')
    return
  }

  saving.value = true
  try {
    const { _id, created_at, updated_at, ...updateData } = currentRow.value

    await elasticsearchService.update(_id, 'process', updateData)
    ElMessage.success('保存成功')

    dialogVisible.value = false
    window.location.reload()
  } catch (error) {
    console.error('保存出错:', error)
  } finally {
    saving.value = false
  }
}

// 批量操作入口
const handleBulkAction = async ({
  action,
  rows,
}: {
  action: string
  rows: any[]
}) => {
  if (rows.length === 0) return ElMessage.warning('请先选择数据')

  const ids = rows.map(r => r._id).filter(Boolean)
  if (ids.length === 0) return ElMessage.error('选中数据缺少 _id')

  switch (action) {
    case 'delete':
      await handleBulkDelete(ids)
      break
    case 'export':
      await handleExport(rows)
      break
    default:
      ElMessage.warning(`未知操作：${action}`)
  }
}

// 批量删除
const handleBulkDelete = async (ids: string[]) => {
  try {
    const result = await elasticsearchService.bulkDelete(ids, 'process')

    if (result.success) {
      ElMessage.success(`成功删除 ${result.successCount} 条数据`)
      if (result.failedCount > 0) {
        ElMessage.warning(`${result.failedCount} 条数据删除失败`)
      }
      setTimeout(() => window.location.reload(), 400)
    }
  } catch (error) {
    console.error(error)
  }
}

// 导出
const handleExport = async (rows: any[]) => {
  try {
    const data = rows.map(({ _id, ...rest }) => rest)
    const headers = Object.keys(data[0] || {})

    const csvContent = [
      headers.join(','),
      ...data.map(row =>
        headers.map(h => JSON.stringify(row[h] || '')).join(',')
      ),
    ].join('\n')

    const blob = new Blob(['\ufeff' + csvContent], {
      type: 'text/csv;charset=utf-8;',
    })

    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `process_${Date.now()}.csv`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)

    ElMessage.success('导出成功')
  } catch (e) {
    console.error(e)
    ElMessage.error('导出失败')
  }
}

const handleDialogClose = () => {
  currentRow.value = {}
  saving.value = false
}
</script>

<style scoped lang="scss">
.process-management {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
  overflow: hidden;

  :deep(.data-table-container) {
    flex: 1;
    margin: 20px;
    overflow: hidden;
  }
}
</style>
