<template>
  <div class="material-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="openDetail"
      @edit="openDetail"
      @bulkAction="handleBulkAction"
    />

    <!-- 新增原料对话框 -->
    <AddMaterialDialog
      v-model:visible="addDialogVisible"
      @success="handleAddSuccess"
    />

    <!-- 原料详情抽屉 -->
    <MaterialDetail
      v-model:visible="detailVisible"
      :material="selectedMaterial"
      @update-material="handleMaterialUpdate"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import DataTable from '#/components/Table/index.vue'
import AddMaterialDialog from './components/AddMaterialDialog.vue'
import MaterialDetail from './components/MaterialDetail.vue'
import type { Material } from './types'
import type { PageConfig, BulkAction } from '#/components/Table/types'
import { useDataTable } from '#/composables/useDataTable'

// 数据表格配置
const { searchLoading } = useDataTable({
  index: 'material',
  pageSize: 20,
  defaultSort: [{ field: 'created_at', order: 'desc' }]
})

// 对话框和抽屉
const addDialogVisible = ref(false)
const detailVisible = ref(false)
const selectedMaterial = ref<Material | null>(null)

// 获取分类标签类型
const getCategoryTagType = (category?: string) => {
  if (!category) return 'info'
  const typeMap: Record<string, any> = {
    '胚布': 'primary',
    '染料': 'warning',
    '助剂': 'success'
  }
  return typeMap[category] || 'info'
}

// 页面配置
const pageConfig: PageConfig = {
  pageType: 'material',
  title: '原料管理',
  index: 'material',
  pageSize: 20,
  columns: [
    {
      key: 'code',
      label: '原料编号',
      showOverflowTooltip: true,
      formatter: (v: string) => v || '-'
    },
    {
      key: 'name',
      label: '原料名称',
      showOverflowTooltip: true
    },
    {
      key: 'spec',
      label: '规格',
      showOverflowTooltip: true,
      formatter: (v: string) => v || '-'
    },
    {
      key: 'category',
      label: '分类',
      showOverflowTooltip: true,
      formatter: (v: string) => v || '-'
    },
    {
      key: 'unit',
      label: '单位',
      showOverflowTooltip: true,
      formatter: (v: string) => v || '-'
    },
    {
      key: 'currentPrice',
      label: '当前价格',
      showOverflowTooltip: true,
      formatter: (v: number, row: any) => v ? `¥${v.toFixed(2)}/${row.unit}` : '-'
    },
    {
      key: 'status',
      label: '状态',
      showOverflowTooltip: true,
      formatter: (v: string) => v === 'active' ? '在用' : '停用'
    },
    {
      key: 'created_at',
      label: '创建时间',
      showOverflowTooltip: true,
      formatter: (v: string) => v ? new Date(v).toLocaleString('zh-CN') : '-'
    }
  ],
  filters: [
    {
      key: 'category',
      label: '分类',
      type: 'select',
      placeholder: '请选择分类',
      options: [
        { label: '胚布', value: '胚布' },
        { label: '染料', value: '染料' },
        { label: '助剂', value: '助剂' }
      ]
    }
  ],
  bulkActions: [
    { key: 'create', label: '新增原料', type: 'primary' },
    { key: 'delete', label: '批量删除', type: 'danger', confirm: true, confirmMessage: '确定要删除选中的原料吗？' }
  ]
}

// 新增成功
const handleAddSuccess = (newMaterial: Material) => {
  ElMessage.success('新增成功')
  // 重新加载数据
  window.location.reload()
}

// 打开详情
const openDetail = (material: Material) => {
  selectedMaterial.value = material
  detailVisible.value = true
}

// 更新原料信息
const handleMaterialUpdate = (updatedMaterial: Material) => {
  // 重新加载数据
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
.material-management {
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
