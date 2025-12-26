<template>
  <div class="material-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="openDetail"
      @edit="openDetail"
      @bulkAction="handleBulkAction"
      @topAction="handleTopAction"
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
import AddMaterialDialog from '#/components/Material/AddMaterialDialog.vue'
import MaterialDetail from '#/components/Material/MaterialDetail.vue'
import type { Material } from '#/components/Material/types'
import type { PageConfig, BulkAction } from '#/components/Table/types'
import { useDataTable } from '#/composables/useDataTable'
import { elasticsearchService } from '#/api/core/es'

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

// 页面配置
const pageConfig: PageConfig = {
  pageType: 'material',
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
      fetchOptions: async () => {
        try {
          const response = await elasticsearchService.search({
            index: 'material',
            pagination: { offset: 0, size: 0 },
            aggRequests: {
              category: {
                type: 'terms',
                field: 'category',
                size: 20,
              },
            },
          })
          const buckets = response.aggregations?.category?.buckets || []
          return buckets.map((bucket: any) => ({
            label: String(bucket.key),
            value: bucket.key,
          }))
        } catch (error) {
          console.error('加载分类选项失败:', error)
          return []
        }
      },
    },
    {
      key: 'status',
      label: '状态',
      type: 'select',
      placeholder: '请选择状态',
      fetchOptions: async () => {
        try {
          const response = await elasticsearchService.search({
            index: 'material',
            pagination: { offset: 0, size: 0 },
            aggRequests: {
              status: {
                type: 'terms',
                field: 'status',
                size: 20,
              },
            },
          })
          const buckets = response.aggregations?.status?.buckets || []
          return buckets.map((bucket: any) => ({
            label: String(bucket.key),
            value: bucket.key,
          }))
        } catch (error) {
          console.error('加载状态选项失败:', error)
          return []
        }
      },
    },
    {
      key: 'unit',
      label: '单位',
      type: 'select',
      placeholder: '请选择单位',
      fetchOptions: async () => {
        try {
          const response = await elasticsearchService.search({
            index: 'material',
            pagination: { offset: 0, size: 0 },
            aggRequests: {
              unit: {
                type: 'terms',
                field: 'unit',
                size: 20,
              },
            },
          })
          const buckets = response.aggregations?.unit?.buckets || []
          return buckets.map((bucket: any) => ({
            label: String(bucket.key),
            value: bucket.key,
          }))
        } catch (error) {
          console.error('加载单位选项失败:', error)
          return []
        }
      },
    },
  ],
  topActions: [
    { key: 'create', label: '新增原料', type: 'primary' },
  ],
  bulkActions: [
    { key: 'delete', label: '批量删除', type: 'danger', confirm: true, confirmMessage: '确定要删除选中的原料吗?' }
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

// 顶部操作
const handleTopAction = ({ action }: { action: string }) => {
  if (action === 'create') {
    addDialogVisible.value = true
  }
}

// 批量操作
const handleBulkAction = ({ action }: { action: string; rows: any[] }) => {
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
