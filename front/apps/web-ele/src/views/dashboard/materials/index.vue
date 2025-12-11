<template>
  <div class="material-management">
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
        <el-form-item label="原料名称">
          <el-input v-model="currentRow.name" :disabled="dialogMode === 'view'" />
        </el-form-item>
        <el-form-item label="规格">
          <el-input v-model="currentRow.spec" :disabled="dialogMode === 'view'" />
        </el-form-item>
        <el-form-item label="单位">
          <el-input v-model="currentRow.unit" :disabled="dialogMode === 'view'" />
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
import { ElMessage, ElMessageBox } from 'element-plus'
import DataTable from '#/components/Table/index.vue'
import { elasticsearchService } from '#/api/core/es'
import type {
  PageConfig,
  ColumnConfig,
  FilterConfig,
  BulkAction,
} from '#/components/Table/types'

// ---- 新增：从 composable 取 searchLoading，用于在 query/filters/sort 变更时显示表格处 loader
// 注意：useDataTable 并不会在此自动初始化（不影响 DataTable 内部逻辑），只是读取 searchLoading 状态
import { useDataTable } from '#/composables/useDataTable'
const { searchLoading } = useDataTable(['materials'], 20)
// ---- 新增结束

// 页面配置
const pageConfig: PageConfig = {
  pageType: 'material',
  title: '原料管理',
  indices: ['materials'],
  pageSize: 20,
  columns: [
    {
      key: 'id',
      label: 'ID',
      width: 80,
      sortable: true,
      visible: true,
      order: 0,
    },
    {
      key: 'name',
      label: '原料名称',
      width: 300,
      sortable: true,
      visible: true,
      order: 1,
    },
    {
      key: 'spec',
      label: '规格',
      width: 150,
      visible: true,
      order: 2,
    },
    {
      key: 'unit',
      label: '单位',
      width: 100,
      visible: true,
      order: 3,
    },
    {
      key: 'description',
      label: '描述',
      width: 200,
      visible: true,
      order: 4,
    },
    {
      key: 'created_at',
      label: '创建时间',
      width: 180,
      sortable: true,
      visible: true,
      order: 5,
      formatter: (value: string) => {
        return value ? new Date(value).toLocaleString('zh-CN') : '-'
      },
    },
    {
      key: 'updated_at',
      label: '更新时间',
      width: 180,
      sortable: true,
      visible: false,
      order: 6,
      formatter: (value: string) => {
        return value ? new Date(value).toLocaleString('zh-CN') : '-'
      },
    },
  ] as ColumnConfig[],
  filters: [
    {
      key: 'name',
      label: '原料名称',
      type: 'text',
      placeholder: '请输入原料名称',
    },
    {
      key: 'spec',
      label: '规格',
      type: 'text',
      placeholder: '请输入规格',
    },
  ] as FilterConfig[],
  bulkActions: [
    {
      key: 'delete',
      label: '批量删除',
      type: 'danger',
      confirm: true,
      confirmMessage: '确定要删除选中的原料吗？此操作不可恢复！',
    },
    {
      key: 'export',
      label: '导出数据',
      type: 'primary',
    },
  ] as BulkAction[],
}

// 对话框状态
const dialogVisible = ref(false)
const dialogMode = ref<'view' | 'edit'>('view')
const currentRow = ref<any>({})
const saving = ref(false)

const dialogTitle = computed(() => {
  return dialogMode.value === 'view' ? '查看原料' : '编辑原料'
})

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
    ElMessage.error('缺少必要的ID信息')
    return
  }

  saving.value = true
  try {
    const { _id, id, created_at, updated_at, ...updateData } = currentRow.value

    await elasticsearchService.update(_id, 'materials', updateData)

    ElMessage.success('保存成功')
    dialogVisible.value = false
    window.location.reload()
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    saving.value = false
  }
}

// 批量操作
const handleBulkAction = async ({
  action,
  rows,
}: {
  action: string
  rows: any[]
}) => {
  if (rows.length === 0) {
    ElMessage.warning('请先选择要操作的数据')
    return
  }

  const ids = rows.map((row) => row._id).filter(Boolean)

  if (ids.length === 0) {
    ElMessage.error('选中的数据缺少ID信息')
    return
  }

  try {
    switch (action) {
      case 'delete':
        await handleBulkDelete(ids)
        break

      case 'export':
        await handleExport(rows)
        break

      default:
        ElMessage.warning(`未知的操作: ${action}`)
    }
  } catch (error) {
    console.error('批量操作失败:', error)
  }
}

// 批量删除
const handleBulkDelete = async (ids: string[]) => {
  try {
    const result = await elasticsearchService.bulkDelete(ids, 'materials')

    if (result.success) {
      ElMessage.success(`成功删除 ${result.successCount} 条数据`)

      if (result.failedCount > 0) {
        ElMessage.warning(
          `${result.failedCount} 条数据删除失败，请查看详情`
        )
        console.error('删除失败的记录:', result.errors)
      }

      setTimeout(() => {
        window.location.reload()
      }, 500)
    }
  } catch (error) {
    console.error('批量删除失败:', error)
  }
}

// 导出数据
const handleExport = async (rows: any[]) => {
  try {
    const data = rows.map((row) => {
      const { _id, ...rest } = row
      return rest
    })

    const headers = Object.keys(data[0] || {})
    const csvContent = [
      headers.join(','),
      ...data.map((row) =>
        headers.map((header) => JSON.stringify(row[header] || '')).join(',')
      ),
    ].join('\n')

    const blob = new Blob(['\ufeff' + csvContent], {
      type: 'text/csv;charset=utf-8;',
    })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `materials_${Date.now()}.csv`
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

// 关闭对话框
const handleDialogClose = () => {
  currentRow.value = {}
  saving.value = false
}
</script>

<style scoped lang="scss">
.material-management {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
  overflow: hidden;  // 防止外部滚动
  
  :deep(.data-table-container) {
    flex: 1;
    margin: 20px;
    overflow: hidden;
  }
}
</style>