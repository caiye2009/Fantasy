<template>
  <div class="material-management">
    <!-- 数据表格 -->
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      :table-data="tableData"
      :has-more="hasMore"
      :has-previous="hasPrevious"
      :selected-rows="selectedRows"
      :selected-count="selectedCount"
      @view="openDetail"
      @edit="openDetail"
      @quote="handleQuote"
      @bulkAction="handleBulkAction"
      @topAction="handleTopAction"
      @scroll="handleScroll"
      @search="handleSearch"
      @filter="handleFilter"
    />

    <!-- 原料详情/新增抽屉 -->
    <MaterialDetail
      v-model:visible="detailVisible"
      :material="selectedMaterial"
      @update-material="handleMaterialUpdate"
    />

    <!-- 报价对话框 -->
    <el-dialog
      v-model="quoteDialogVisible"
      title="添加报价"
      width="500px"
      @close="handleQuoteDialogClose"
    >
      <el-form :model="quoteForm" label-width="100px">
        <el-form-item label="原料名称">
          <el-input :value="quoteMaterialName" disabled />
        </el-form-item>
        <el-form-item label="供应商" required>
          <el-select
            v-model="quoteForm.supplier_id"
            placeholder="请选择供应商"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="supplier in suppliers"
              :key="supplier.id"
              :label="supplier.name"
              :value="supplier.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="报价金额" required>
          <el-input-number
            v-model="quoteForm.price"
            :min="0"
            :precision="2"
            :step="0.01"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="quoteDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleQuoteSubmit" :loading="quoting">
          提交报价
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import DataTable from '#/components/Table/index.vue'
import MaterialDetail from '#/components/Material/MaterialDetail.vue'
import type { Material } from '#/components/Material/types'
import type { PageConfig } from '#/components/Table/types'
import { useDataTable } from '#/composables/useDataTable'
import { elasticsearchService } from '#/api/core/es'
import { getSupplierListApi, type Supplier } from '#/api/core/supplier'
import { quoteMaterialPriceApi } from '#/api/core/material_quote'

// =====================
// 数据表格逻辑
// =====================
const {
  searchLoading,
  tableData,
  query,
  filters,
  selectedRows,
  selectedCount,
  hasMore,
  hasPrevious,
  initialize,
  handleScroll,
  reload
} = useDataTable({
  index: 'material',
  pageSize: 40,
  defaultSort: [{ field: 'created_at', order: 'desc' }]
})

console.log("dllm", tableData.value)

// 搜索处理
const handleSearch = (searchQuery: string) => {
  query.value = searchQuery
}

// 筛选处理
const handleFilter = (filterValues: Record<string, any>) => {
  filters.value = filterValues
}

// 初始化数据
initialize()

// =====================
// 抽屉逻辑
// =====================
const detailVisible = ref(false)
const selectedMaterial = ref<Material | null>(null)

// =====================
// 表格配置
// =====================
const pageConfig: PageConfig = {
  pageType: 'material',
  index: 'material',
  pageSize: 40,
  columns: [
    // 只展示四个字段
    {
      key: 'materialCode',
      label: '原料编号',
      showOverflowTooltip: true,
      formatter: (v: string) => v || '-'
    },
    {
      key: 'materialName',
      label: '原料名称',
      showOverflowTooltip: true
    },
    {
      key: 'materialCategory',
      label: '规格',
      showOverflowTooltip: true,
      formatter: (v: string) => v || '-'
    },
    {
      key: 'materialStatus',
      label: '状态',
      showOverflowTooltip: true,
      formatter: (v: string) => v || '-'
    },
  ],
  filters: [
    // 保留现有筛选逻辑
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
              category: { type: 'terms', field: 'category', size: 20 },
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
              status: { type: 'terms', field: 'status', size: 20 },
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
              unit: { type: 'terms', field: 'unit', size: 20 },
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
  ],
  actions: [
    {
      key: 'quote',
      label: '报价',
      type: 'warning',
    },
  ],
}

// =====================
// 行操作逻辑
// =====================
const openDetail = (material: Material) => {
  selectedMaterial.value = material
  detailVisible.value = true
}

// 更新/新增原料信息
const handleMaterialUpdate = (updatedMaterial: Material) => {
  ElMessage.success('操作成功')
  // 重新加载数据
}

// 顶部操作
const handleTopAction = ({ action }: { action: string }) => {
  if (action === 'create') {
    selectedMaterial.value = null  // null 表示新增模式
    detailVisible.value = true
  }
}

// 批量操作
const handleBulkAction = ({ action }: { action: string; rows: any[] }) => {
  // TODO: 实现批量删除功能
}

// =====================
// 报价逻辑
// =====================
const quoteDialogVisible = ref(false)
const quoting = ref(false)
const quoteForm = ref({
  target_id: 0,
  supplier_id: undefined as number | undefined,
  price: 0,
})
const quoteMaterialName = ref('')
const suppliers = ref<Supplier[]>([])

const loadSuppliers = async () => {
  try {
    const res = await getSupplierListApi({ limit: 1000, offset: 0 })
    suppliers.value = res.list || []
  } catch (error) {
    console.error('加载供应商列表失败:', error)
  }
}

const handleQuote = async (row: any) => {
  quoteForm.value = {
    target_id: row.id,
    supplier_id: undefined,
    price: 0,
  }
  quoteMaterialName.value = row.name
  await loadSuppliers()
  quoteDialogVisible.value = true
}

const handleQuoteSubmit = async () => {
  if (!quoteForm.value.supplier_id) {
    ElMessage.warning('请选择供应商')
    return
  }

  if (quoteForm.value.price <= 0) {
    ElMessage.warning('请输入有效的报价金额')
    return
  }

  quoting.value = true
  try {
    await quoteMaterialPriceApi(quoteForm.value)
    ElMessage.success('报价成功')
    quoteDialogVisible.value = false
  } catch (error: any) {
    console.error('报价失败:', error)
    ElMessage.error(error.response?.data?.error || '报价失败')
  } finally {
    quoting.value = false
  }
}

const handleQuoteDialogClose = () => {
  quoteForm.value = {
    target_id: 0,
    supplier_id: undefined,
    price: 0,
  }
  quoteMaterialName.value = ''
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
