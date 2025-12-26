<template>
  <div class="process-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="handleViewDetail"
      @edit="handleEdit"
      @quote="handleQuote"
      @bulkAction="handleBulkAction"
      @topAction="handleTopAction"
    />

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      title="编辑工序"
      width="600px"
      @close="handleEditDialogClose"
    >
      <el-form :model="currentRow" label-width="100px">
        <el-form-item label="ID">
          <el-input v-model="currentRow.id" disabled />
        </el-form-item>
        <el-form-item label="工序名称">
          <el-input v-model="currentRow.name" />
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
        <el-form-item label="工序名称">
          <el-input :value="quoteProcessName" disabled />
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
      title="工序详情"
      width="900px"
      @close="handleDetailDialogClose"
    >
      <!-- 静态信息 -->
      <div class="detail-section">
        <h3 class="section-title">基本信息</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="ID">
            {{ processDetail.id }}
          </el-descriptions-item>
          <el-descriptions-item label="工序名称">
            {{ processDetail.name }}
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">
            {{ processDetail.description || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(processDetail.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            {{ formatDate(processDetail.updatedAt) }}
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

      <!-- 价格趋势图 -->
      <div class="detail-section" v-if="priceHistory.length > 0">
        <h3 class="section-title">价格趋势</h3>
        <div ref="priceChartRef" class="price-chart"></div>
      </div>

      <template #footer>
        <el-button type="primary" @click="detailDialogVisible = false">
          关闭
        </el-button>
      </template>
    </el-dialog>

    <!-- 新增工序对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="新增工序"
      width="600px"
      @close="handleCreateClose"
    >
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="工序名称" required>
          <el-input v-model="createForm.name" placeholder="请输入工序名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="4"
            placeholder="请输入工序描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="handleCreateSubmit">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'
import DataTable from '#/components/Table/index.vue'
import { elasticsearchService } from '#/api/core/es'
import {
  quoteProcessPriceApi,
  getProcessPriceHistoryApi,
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
const { searchLoading } = useDataTable('process', 20)

// 配置
const pageConfig: PageConfig = {
  pageType: 'process',
  index: 'process',
  pageSize: 20,
  columns: [
    { key: 'id', label: 'ID', width: 80, visible: true, sortable: true, order: 0 },
    { key: 'name', label: '工序名称', width: 200, visible: true, sortable: true, order: 1 },
    { key: 'description', label: '描述', width: 250, visible: true, order: 2 },
    {
      key: 'createdAt',
      label: '创建时间',
      width: 180,
      visible: true,
      sortable: true,
      order: 3,
      formatter: (value: string) => value ? new Date(value).toLocaleString('zh-CN') : '-'
    },
    {
      key: 'updatedAt',
      label: '更新时间',
      width: 180,
      visible: false,
      sortable: true,
      order: 4,
      formatter: (value: string) => value ? new Date(value).toLocaleString('zh-CN') : '-'
    }
  ] as ColumnConfig[],
  filters: [
    {
      key: 'type',
      label: '工序类型',
      type: 'select',
      placeholder: '请选择工序类型',
      options: [],
    },
    {
      key: 'category',
      label: '工序类别',
      type: 'select',
      placeholder: '请选择工序类别',
      options: [],
    },
  ] as FilterConfig[],
  topActions: [
    { key: 'create', label: '新增工序', type: 'primary' },
  ],
  bulkActions: [
    { key: 'delete', label: '批量删除', type: 'danger', confirm: true, confirmMessage: '确定要删除选中的工序吗？此操作不可恢复！' },
    { key: 'export', label: '导出数据', type: 'primary' },
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
const quoteProcessName = ref('')
const suppliers = ref<Supplier[]>([])

// 详情对话框
const detailDialogVisible = ref(false)
const processDetail = ref<any>({})
const priceHistory = ref<PriceData[]>([])
const loadingPriceHistory = ref(false)
const priceChartRef = ref<HTMLElement>()

// 格式化日期
const formatDate = (date: string) => {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

// 渲染价格趋势图
const renderPriceChart = () => {
  if (!priceChartRef.value || priceHistory.value.length === 0) return

  const chart = echarts.init(priceChartRef.value)

  // 按时间排序（从旧到新）
  const sortedHistory = [...priceHistory.value].sort((a, b) =>
    new Date(a.quoted_at).getTime() - new Date(b.quoted_at).getTime()
  )

  const option: EChartsOption = {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const data = params[0]
        return `${data.name}<br/>价格: ¥${data.value.toFixed(2)}`
      }
    },
    xAxis: {
      type: 'category',
      data: sortedHistory.map(item => formatDate(item.quoted_at)),
      axisLabel: {
        rotate: 45,
        interval: 0,
        fontSize: 10
      }
    },
    yAxis: {
      type: 'value',
      name: '价格（元）',
      axisLabel: {
        formatter: '¥{value}'
      }
    },
    series: [
      {
        name: '价格',
        type: 'line',
        data: sortedHistory.map(item => item.price),
        smooth: true,
        itemStyle: {
          color: '#67c23a'
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              {
                offset: 0,
                color: 'rgba(103, 194, 58, 0.3)'
              },
              {
                offset: 1,
                color: 'rgba(103, 194, 58, 0.05)'
              }
            ]
          }
        }
      }
    ],
    grid: {
      left: '10%',
      right: '5%',
      bottom: '15%',
      top: '10%'
    }
  }

  chart.setOption(option)

  // 响应式调整
  window.addEventListener('resize', () => chart.resize())
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
  processDetail.value = { ...row }
  detailDialogVisible.value = true

  // 加载报价历史
  if (row.id) {
    loadingPriceHistory.value = true
    try {
      priceHistory.value = await getProcessPriceHistoryApi(row.id)
      // 等待DOM更新后渲染图表
      await nextTick()
      renderPriceChart()
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

// 保存
const handleSave = async () => {
  if (!currentRow.value._id) {
    ElMessage.error('缺少必要的 ID 信息')
    return
  }

  saving.value = true
  try {
    const { _id, createdAt, updatedAt, ...updateData } = currentRow.value

    await elasticsearchService.update(_id, 'process', updateData)
    ElMessage.success('保存成功')

    editDialogVisible.value = false
    window.location.reload()
  } catch (error) {
    console.error('保存出错:', error)
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
  quoteProcessName.value = row.name
  await loadSuppliers()
  quoteDialogVisible.value = true
}

// 报价（从详情弹窗）
const handleQuoteFromDetail = async () => {
  quoteForm.value = {
    target_id: processDetail.value.id,
    supplier_id: undefined,
    price: 0,
  }
  quoteProcessName.value = processDetail.value.name
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
    await quoteProcessPriceApi(quoteForm.value)
    ElMessage.success('报价成功')
    quoteDialogVisible.value = false

    // 如果详情弹窗是打开的，刷新报价历史
    if (detailDialogVisible.value) {
      loadingPriceHistory.value = true
      try {
        priceHistory.value = await getProcessPriceHistoryApi(
          processDetail.value.id
        )
        // 刷新图表
        await nextTick()
        renderPriceChart()
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

// 顶部操作
const handleTopAction = ({ action }: { action: string }) => {
  if (action === 'create') {
    handleCreate()
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

// 新增工序
const createDialogVisible = ref(false)
const createForm = ref({
  name: '',
  description: '',
})
const creating = ref(false)

const handleCreate = () => {
  createDialogVisible.value = true
}

const handleCreateSubmit = async () => {
  if (!createForm.value.name) {
    ElMessage.error('请输入工序名称')
    return
  }

  creating.value = true
  try {
    await elasticsearchService.create('process', createForm.value)
    ElMessage.success('新增成功')
    createDialogVisible.value = false
    setTimeout(() => window.location.reload(), 500)
  } catch (error) {
    console.error(error)
    ElMessage.error('新增失败')
  } finally {
    creating.value = false
  }
}

const handleCreateClose = () => {
  createForm.value = {
    name: '',
    description: '',
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
  quoteProcessName.value = ''
}

const handleDetailDialogClose = () => {
  processDetail.value = {}
  priceHistory.value = []
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

  .price-chart {
    width: 100%;
    height: 350px;
  }
}
</style>
