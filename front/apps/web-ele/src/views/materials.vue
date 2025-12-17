<template>
  <div class="material-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="handleViewDetail"
      @edit="handleEdit"
      @quote="handleQuote"
      @bulkAction="handleBulkAction"
    />

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑原料"
      width="600px"
      @close="handleEditDialogClose"
    >
      <el-form :model="currentRow" label-width="100px">
        <el-form-item label="ID">
          <el-input v-model="currentRow.id" disabled />
        </el-form-item>
        <el-form-item label="原料名称">
          <el-input v-model="currentRow.name" />
        </el-form-item>
        <el-form-item label="规格">
          <el-input v-model="currentRow.spec" />
        </el-form-item>
        <el-form-item label="单位">
          <el-input v-model="currentRow.unit" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="currentRow.description"
            type="textarea"
            :rows="3"
          />
        </el-form-item>
        <el-form-item label="创建时间" v-if="currentRow.createdAt">
          <el-input
            :value="new Date(currentRow.createdAt).toLocaleString('zh-CN')"
            disabled
          />
        </el-form-item>
        <el-form-item label="更新时间" v-if="currentRow.updatedAt">
          <el-input
            :value="new Date(currentRow.updatedAt).toLocaleString('zh-CN')"
            disabled
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">
          保存
        </el-button>
      </template>
    </el-dialog>

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

    <!-- 详情大弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="原料详情"
      width="900px"
      @close="handleDetailDialogClose"
    >
      <!-- 静态信息 -->
      <div class="detail-section">
        <h3 class="section-title">基本信息</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">
            {{ materialDetail.id }}
          </el-descriptions-item>
          <el-descriptions-item label="原料名称">
            {{ materialDetail.name }}
          </el-descriptions-item>
          <el-descriptions-item label="规格">
            {{ materialDetail.spec || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="单位">
            {{ materialDetail.unit || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">
            {{ materialDetail.description || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(materialDetail.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            {{ formatDate(materialDetail.updatedAt) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 报价历史 -->
      <div class="detail-section">
        <div class="section-header">
          <h3 class="section-title">报价历史</h3>
          <el-button
            type="primary"
            size="small"
            @click="handleQuoteFromDetail"
          >
            添加报价
          </el-button>
        </div>
        <el-table
          :data="priceHistory"
          v-loading="loadingPriceHistory"
          stripe
          border
        >
          <el-table-column prop="supplier_name" label="供应商" width="200" />
          <el-table-column prop="price" label="报价金额" width="150">
            <template #default="{ row }">
              ¥{{ row.price.toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="quoted_at" label="更新时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.quoted_at) }}
            </template>
          </el-table-column>
        </el-table>
        <el-empty
          v-if="!loadingPriceHistory && priceHistory.length === 0"
          description="暂无报价记录"
        />
      </div>

      <!-- 库存信息 -->
      <div class="detail-section">
        <h3 class="section-title">库存信息</h3>
        <el-alert
          title="库存管理功能开发中"
          type="info"
          :closable="false"
        />
      </div>

      <template #footer>
        <el-button type="primary" @click="detailDialogVisible = false">
          关闭
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
import {
  quoteMaterialPriceApi,
  getMaterialPriceHistoryApi,
  type PriceData,
} from '#/api/core/pricing'
import { getSupplierListApi, type Supplier } from '#/api/core/supplier'
import type {
  PageConfig,
  ColumnConfig,
  FilterConfig,
  BulkAction,
} from '#/components/Table/types'

import { useDataTable } from '#/composables/useDataTable'
const { searchLoading } = useDataTable('material', 20)

// 页面配置
const pageConfig: PageConfig = {
  pageType: 'material',
  title: '原料管理',
  index: 'material',
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
      key: 'createdAt',
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
      key: 'updatedAt',
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
      key: 'unit',
      label: '单位',
      type: 'select',
      placeholder: '请选择单位',
      options: [], // 后端会提供接口返回选项
    },
    {
      key: 'spec',
      label: '规格',
      type: 'select',
      placeholder: '请选择规格',
      options: [], // 后端会提供接口返回选项
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
  actions: [
    {
      key: 'quote',
      label: '报价',
      type: 'warning',
    },
  ],
}

// 编辑对话框
const editDialogVisible = ref(false)
const currentRow = ref<any>({})
const saving = ref(false)

// 报价对话框
const quoteDialogVisible = ref(false)
const quoting = ref(false)
const quoteForm = ref({
  target_id: 0,
  supplier_id: undefined as number | undefined,
  price: 0,
})
const quoteMaterialName = ref('')
const suppliers = ref<Supplier[]>([])

// 详情对话框
const detailDialogVisible = ref(false)
const materialDetail = ref<any>({})
const priceHistory = ref<PriceData[]>([])
const loadingPriceHistory = ref(false)

// 格式化日期
const formatDate = (date: string) => {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

// 加载供应商列表
const loadSuppliers = async () => {
  try {
    const res = await getSupplierListApi({ limit: 1000, offset: 0 })
    suppliers.value = res.suppliers || []
  } catch (error) {
    console.error('加载供应商列表失败:', error)
  }
}

// 查看详情
const handleViewDetail = async (row: any) => {
  materialDetail.value = { ...row }
  detailDialogVisible.value = true

  // 加载报价历史
  if (row.id) {
    loadingPriceHistory.value = true
    try {
      priceHistory.value = await getMaterialPriceHistoryApi(row.id)
    } catch (error) {
      console.error('加载报价历史失败:', error)
      priceHistory.value = []
    } finally {
      loadingPriceHistory.value = false
    }
  }
}

// 编辑
const handleEdit = (row: any) => {
  currentRow.value = { ...row }
  editDialogVisible.value = true
}

// 保存编辑
const handleSave = async () => {
  if (!currentRow.value._id) {
    ElMessage.error('缺少必要的ID信息')
    return
  }

  saving.value = true
  try {
    const { _id, id, createdAt, updatedAt, ...updateData } = currentRow.value

    await elasticsearchService.update(_id, 'materials', updateData)

    ElMessage.success('保存成功')
    editDialogVisible.value = false
    window.location.reload()
  } catch (error) {
    console.error('保存失败:', error)
  } finally {
    saving.value = false
  }
}

// 报价（从表格）
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

// 报价（从详情弹窗）
const handleQuoteFromDetail = async () => {
  quoteForm.value = {
    target_id: materialDetail.value.id,
    supplier_id: undefined,
    price: 0,
  }
  quoteMaterialName.value = materialDetail.value.name
  await loadSuppliers()
  quoteDialogVisible.value = true
}

// 提交报价
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

    // 如果详情弹窗是打开的，刷新报价历史
    if (detailDialogVisible.value) {
      loadingPriceHistory.value = true
      try {
        priceHistory.value = await getMaterialPriceHistoryApi(
          materialDetail.value.id
        )
      } catch (error) {
        console.error('刷新报价历史失败:', error)
      } finally {
        loadingPriceHistory.value = false
      }
    }
  } catch (error: any) {
    console.error('报价失败:', error)
    ElMessage.error(error.response?.data?.error || '报价失败')
  } finally {
    quoting.value = false
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
const handleEditDialogClose = () => {
  currentRow.value = {}
  saving.value = false
}

const handleQuoteDialogClose = () => {
  quoteForm.value = {
    target_id: 0,
    supplier_id: undefined,
    price: 0,
  }
  quoteMaterialName.value = ''
}

const handleDetailDialogClose = () => {
  materialDetail.value = {}
  priceHistory.value = []
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

.detail-section {
  margin-bottom: 24px;

  &:last-child {
    margin-bottom: 0;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .section-title {
    margin: 0 0 16px 0;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }
}
</style>
