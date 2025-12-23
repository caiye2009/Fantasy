<template>
  <div class="process-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="openDetail"
      @edit="openDetail"
      @bulkAction="handleBulkAction"
    />

    <!-- 新增工艺对话框 -->
    <AddProcessDialog
      v-model:visible="addDialogVisible"
      @success="handleAddSuccess"
    />

    <!-- 工艺详情抽屉 -->
    <ProcessDetail
      v-model:visible="detailVisible"
      :process="selectedProcess"
      @update-process="handleProcessUpdate"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import DataTable from '#/components/Table/index.vue'
import AddProcessDialog from './components/AddProcessDialog.vue'
import ProcessDetail from './components/ProcessDetail.vue'
import type { Process } from './types'
import type { PageConfig } from '#/components/Table/types'
import { useDataTable } from '#/composables/useDataTable'

// 数据表格配置
const { searchLoading } = useDataTable({
  index: 'process',
  pageSize: 20,
  defaultSort: [{ field: 'created_at', order: 'desc' }]
})

// 对话框和抽屉
const addDialogVisible = ref(false)
const detailVisible = ref(false)
const selectedProcess = ref<Process | null>(null)

// 页面配置
const pageConfig: PageConfig = {
  pageType: 'process',
  title: '工艺管理',
  index: 'process',
  pageSize: 20,
  columns: [
    {
      key: 'id',
      label: 'ID',
      showOverflowTooltip: true
    },
    {
      key: 'name',
      label: '工艺名称',
      showOverflowTooltip: true
    },
    {
      key: 'description',
      label: '描述',
      showOverflowTooltip: true,
      formatter: (v: string) => v || '-'
    },
    {
      key: 'created_at',
      label: '创建时间',
      showOverflowTooltip: true,
      formatter: (v: string) => v ? new Date(v).toLocaleString('zh-CN') : '-'
    }
  ],
  filters: [],
  bulkActions: [
    { key: 'create', label: '新增工艺', type: 'primary' },
    { key: 'delete', label: '批量删除', type: 'danger', confirm: true, confirmMessage: '确定要删除选中的工艺吗？' }
  ]
}

// 新增成功
const handleAddSuccess = () => {
  ElMessage.success('新增成功')
  window.location.reload()
}

// 打开详情
const openDetail = (process: Process) => {
  selectedProcess.value = process
  detailVisible.value = true
}

// 更新工艺信息
const handleProcessUpdate = (updatedProcess: Process) => {
  window.location.reload()
}

// 批量操作
const handleBulkAction = ({ action }: { action: string; rows: any[] }) => {
  if (action === 'create') {
    addDialogVisible.value = true
  }
  // TODO: 实现批量删除功能
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
