<template>
  <div class="order-management">
    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="handleView"
      @bulkAction="handleBulkAction"
      @topAction="handleTopAction"
    />

    <!-- 创建订单对话框 -->
    <el-dialog
      v-model="createDialogVisible"
      title="创建订单"
      width="90%"
      :close-on-click-modal="false"
      @close="handleCreateClose"
      style="max-width: 1400px"
    >
      <div class="order-creation-wrapper">
        <!-- 步骤选择区域 -->
        <div class="steps-container" :class="{ 'slide-left': showPreview }">
          <div class="order-creation-flow">
            <!-- 1. 选择客户 -->
            <div class="step-section">
              <div
                class="step-header"
                :class="{ 'active': currentStep === 1, 'completed': selectedCustomer }"
                @click="goToStep(1)"
              >
                <div v-if="currentStep === 1 || !selectedCustomer" class="step-title">
                  <span class="step-number">1</span>
                  <span class="step-name">选择客户</span>
                </div>
                <div v-else class="step-summary-inline">
                  <span class="step-number-small">1</span>
                  <span class="summary-value">{{ selectedCustomer.name }} - {{ selectedCustomer.phone }}</span>
                </div>
              </div>

              <div v-if="currentStep === 1" class="step-content">
                <el-form :model="form" label-width="100px">
                  <el-form-item label="客户名称">
                    <el-select
                      v-model="form.customerId"
                      placeholder="请选择客户"
                      @change="handleCustomerChange"
                      style="width: 100%"
                      size="large"
                    >
                      <el-option
                        v-for="customer in customers"
                        :key="customer.id"
                        :label="customer.name"
                        :value="customer.id"
                      >
                        <div style="display: flex; justify-content: space-between">
                          <span>{{ customer.name }}</span>
                          <span style="color: #8492a6; font-size: 13px">{{ customer.phone }}</span>
                        </div>
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-form>
              </div>
            </div>

            <!-- 2. 选择产品 -->
            <div class="step-section">
              <div
                class="step-header"
                :class="{
                  'active': currentStep === 2,
                  'completed': selectedProduct,
                  'disabled': !selectedCustomer
                }"
                @click="goToStep(2)"
              >
                <div v-if="currentStep === 2 || !selectedProduct" class="step-title">
                  <span class="step-number">2</span>
                  <span class="step-name">选择产品</span>
                </div>
                <div v-else class="step-summary-inline">
                  <span class="step-number-small">2</span>
                  <span class="summary-value">{{ selectedProduct.name }} × {{ form.quantity }} (¥{{ selectedProduct.price }}/件)</span>
                </div>
              </div>

              <div v-if="currentStep === 2" class="step-content">
                <el-form :model="form" label-width="100px">
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <el-form-item label="产品名称">
                        <el-select
                          v-model="form.productId"
                          placeholder="请选择产品"
                          @change="handleProductChange"
                          style="width: 100%"
                          size="large"
                        >
                          <el-option
                            v-for="product in products"
                            :key="product.id"
                            :label="product.name"
                            :value="product.id"
                          >
                            <div style="display: flex; justify-content: space-between">
                              <span>{{ product.name }}</span>
                              <span style="color: #8492a6; font-size: 13px">¥{{ product.price }}</span>
                            </div>
                          </el-option>
                        </el-select>
                      </el-form-item>
                    </el-col>

                    <el-col :span="12">
                      <el-form-item label="购买数量">
                        <el-input-number
                          v-model="form.quantity"
                          :min="1"
                          :max="100"
                          style="width: 100%"
                          size="large"
                          @change="handleQuantityChange"
                        />
                      </el-form-item>
                    </el-col>
                  </el-row>
                </el-form>
              </div>
            </div>

            <!-- 3. 选择库存 -->
            <div class="step-section">
              <div
                class="step-header"
                :class="{
                  'active': currentStep === 3,
                  'completed': selectedWarehouse,
                  'disabled': !selectedProduct
                }"
                @click="goToStep(3)"
              >
                <div v-if="currentStep === 3 || !selectedWarehouse" class="step-title">
                  <span class="step-number">3</span>
                  <span class="step-name">选择库存</span>
                </div>
                <div v-else class="step-summary-inline">
                  <span class="step-number-small">3</span>
                  <span class="summary-value">{{ selectedWarehouse.name }} (库存: {{ selectedWarehouse.stock }})</span>
                </div>
              </div>

              <div v-if="currentStep === 3" class="step-content">
                <el-form :model="form" label-width="100px">
                  <el-form-item label="发货仓库">
                    <el-select
                      v-model="form.warehouseId"
                      placeholder="请选择仓库"
                      @change="handleWarehouseChange"
                      style="width: 100%"
                      size="large"
                    >
                      <el-option
                        v-for="warehouse in warehouses"
                        :key="warehouse.id"
                        :label="warehouse.name"
                        :value="warehouse.id"
                      >
                        <div style="display: flex; justify-content: space-between">
                          <span>{{ warehouse.name }}</span>
                          <span style="color: #8492a6; font-size: 13px">库存: {{ warehouse.stock }}</span>
                        </div>
                      </el-option>
                    </el-select>
                  </el-form-item>
                </el-form>
              </div>
            </div>
          </div>
        </div>

        <!-- 预览区域 -->
        <div class="preview-container" :class="{ 'show-preview': showPreview }">
          <div class="preview-content-full">
            <h3 style="margin-bottom: 24px; color: #303133; font-size: 20px;">订单信息预览</h3>

            <el-row :gutter="40">
              <el-col :span="8">
                <div class="info-card">
                  <div class="info-card-title">
                    <el-icon style="margin-right: 8px"><User /></el-icon>
                    客户信息
                  </div>
                  <div class="info-card-content">
                    <div class="info-item">
                      <span class="info-label">客户姓名：</span>
                      <span class="info-value">{{ selectedCustomer?.name }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">联系电话：</span>
                      <span class="info-value">{{ selectedCustomer?.phone }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">收货地址：</span>
                      <span class="info-value">{{ selectedCustomer?.address }}</span>
                    </div>
                  </div>
                </div>
              </el-col>

              <el-col :span="8">
                <div class="info-card">
                  <div class="info-card-title">
                    <el-icon style="margin-right: 8px"><ShoppingCart /></el-icon>
                    产品信息
                  </div>
                  <div class="info-card-content">
                    <div class="info-item">
                      <span class="info-label">产品名称：</span>
                      <span class="info-value">{{ selectedProduct?.name }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">产品编号：</span>
                      <span class="info-value">{{ selectedProduct?.code }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">产品单价：</span>
                      <span class="info-value">¥{{ selectedProduct?.price }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">购买数量：</span>
                      <span class="info-value">{{ form.quantity }}</span>
                    </div>
                  </div>
                </div>
              </el-col>

              <el-col :span="8">
                <div class="info-card">
                  <div class="info-card-title">
                    <el-icon style="margin-right: 8px"><Box /></el-icon>
                    库存信息
                  </div>
                  <div class="info-card-content">
                    <div class="info-item">
                      <span class="info-label">发货仓库：</span>
                      <span class="info-value">{{ selectedWarehouse?.name }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">仓库地址：</span>
                      <span class="info-value">{{ selectedWarehouse?.address }}</span>
                    </div>
                    <div class="info-item">
                      <span class="info-label">当前库存：</span>
                      <span class="info-value">{{ selectedWarehouse?.stock }}</span>
                    </div>
                  </div>
                </div>
              </el-col>
            </el-row>

            <div class="total-amount-card">
              <div class="total-label">订单总金额</div>
              <div class="total-value">¥{{ totalAmount }}</div>
            </div>

            <div class="preview-actions">
              <el-button size="large" @click="backToSteps">返回</el-button>
              <el-button type="primary" size="large" @click="submitOrder" :loading="submitting">
                确认下单
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 订单详情抽屉 -->
    <el-drawer
      v-model="detailDrawerVisible"
      title="订单详情"
      size="60%"
      @close="handleDetailClose"
    >
      <div v-if="currentOrder" class="order-detail">
        <el-card class="detail-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>订单信息</span>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="订单ID">{{ currentOrder.id || '-' }}</el-descriptions-item>
            <el-descriptions-item label="订单编号">{{ currentOrder.orderNo || '-' }}</el-descriptions-item>
            <el-descriptions-item label="客户名称">{{ currentOrder.customerName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="产品名称">{{ currentOrder.productName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="购买数量">{{ currentOrder.quantity || '-' }}</el-descriptions-item>
            <el-descriptions-item label="总金额">{{ currentOrder.totalAmount ? `¥${currentOrder.totalAmount}` : '-' }}</el-descriptions-item>
            <el-descriptions-item label="订单状态">{{ currentOrder.status || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ currentOrder.createdAt ? new Date(currentOrder.createdAt).toLocaleString('zh-CN') : '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { User, ShoppingCart, Box } from '@element-plus/icons-vue'
import DataTable from '#/components/Table/index.vue'
import { elasticsearchService } from '#/api/core/es'
import { useDataTable } from '#/composables/useDataTable'
import type { PageConfig, BulkAction } from '#/components/Table/types'

const { searchLoading } = useDataTable('order', 20)

// 页面配置
const pageConfig: PageConfig = {
  pageType: 'order',
  index: 'order',
  pageSize: 20,
  columns: [
    { key: 'id', label: 'ID', width: 80, visible: true, sortable: true, order: 0 },
    { key: 'orderNo', label: '订单编号', width: 150, visible: true, order: 1 },
    { key: 'customerName', label: '客户名称', width: 150, visible: true, order: 2 },
    { key: 'productName', label: '产品名称', width: 200, visible: true, order: 3 },
    { key: 'quantity', label: '数量', width: 100, visible: true, order: 4 },
    {
      key: 'totalAmount',
      label: '总金额',
      width: 120,
      visible: true,
      order: 5,
      formatter: (v: number) => v ? `¥${v.toFixed(2)}` : '-'
    },
    { key: 'status', label: '状态', width: 100, visible: true, order: 6 },
    {
      key: 'createdAt',
      label: '创建时间',
      width: 180,
      visible: true,
      sortable: true,
      order: 7,
      formatter: (v: string) => v ? new Date(v).toLocaleString('zh-CN') : '-'
    }
  ],
  filters: [
    {
      key: 'status',
      label: '订单状态',
      type: 'select',
      placeholder: '请选择订单状态',
      fetchOptions: async () => {
        try {
          const response = await elasticsearchService.search({
            index: 'order',
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
          console.error('加载订单状态选项失败:', error)
          return []
        }
      },
    },
    {
      key: 'customerName',
      label: '客户名称',
      type: 'select',
      placeholder: '请选择客户',
      fetchOptions: async () => {
        try {
          const response = await elasticsearchService.search({
            index: 'order',
            pagination: { offset: 0, size: 0 },
            aggRequests: {
              customerName: {
                type: 'terms',
                field: 'customerName',
                size: 50,
              },
            },
          })
          const buckets = response.aggregations?.customerName?.buckets || []
          return buckets.map((bucket: any) => ({
            label: String(bucket.key),
            value: bucket.key,
          }))
        } catch (error) {
          console.error('加载客户选项失败:', error)
          return []
        }
      },
    },
  ],
  topActions: [
    { key: 'create', label: '创建订单', type: 'primary' },
  ],
  bulkActions: [
    { key: 'delete', label: '批量删除', type: 'danger', confirm: true },
    { key: 'export', label: '导出数据', type: 'primary' },
  ] as BulkAction[],
}

// 创建订单对话框状态
const createDialogVisible = ref(false)
const currentStep = ref(1)
const submitting = ref(false)
const showPreview = ref(false)

// 表单数据
const form = ref({
  customerId: '',
  productId: '',
  quantity: 1,
  warehouseId: ''
})

// 模拟客户数据
const customers = ref([
  { id: 1, name: '张三', phone: '13800138001', address: '北京市朝阳区xxx街道xxx号' },
  { id: 2, name: '李四', phone: '13800138002', address: '上海市浦东新区xxx路xxx号' },
  { id: 3, name: '王五', phone: '13800138003', address: '广州市天河区xxx大道xxx号' },
  { id: 4, name: '赵六', phone: '13800138004', address: '深圳市南山区xxx街xxx号' }
])

// 模拟产品数据
const products = ref([
  { id: 1, name: 'iPhone 15 Pro', code: 'IP15P-001', price: 7999, category: '手机' },
  { id: 2, name: 'MacBook Pro 14', code: 'MBP14-001', price: 14999, category: '电脑' },
  { id: 3, name: 'AirPods Pro 2', code: 'APP2-001', price: 1899, category: '耳机' },
  { id: 4, name: 'iPad Air', code: 'IPA-001', price: 4799, category: '平板' }
])

// 模拟仓库数据
const warehouses = ref([
  { id: 1, name: '北京仓库', address: '北京市顺义区物流园区A区', stock: 150 },
  { id: 2, name: '上海仓库', address: '上海市青浦区工业园区B区', stock: 200 },
  { id: 3, name: '广州仓库', address: '广州市白云区物流中心C区', stock: 180 },
  { id: 4, name: '深圳仓库', address: '深圳市龙岗区科技园D区', stock: 120 }
])

// 已选择的数据
const selectedCustomer = computed(() =>
  customers.value.find(c => c.id === form.value.customerId)
)

const selectedProduct = computed(() =>
  products.value.find(p => p.id === form.value.productId)
)

const selectedWarehouse = computed(() =>
  warehouses.value.find(w => w.id === form.value.warehouseId)
)

// 计算总金额
const totalAmount = computed(() => {
  if (selectedProduct.value && form.value.quantity) {
    return (selectedProduct.value.price * form.value.quantity).toFixed(2)
  }
  return 0
})

// 处理客户选择
const handleCustomerChange = () => {
  if (form.value.customerId) {
    currentStep.value = 2
  }
}

// 处理产品选择
const handleProductChange = () => {
  if (form.value.productId && form.value.quantity) {
    currentStep.value = 3
  }
}

// 处理数量变化
const handleQuantityChange = () => {
  if (form.value.productId && form.value.quantity) {
    currentStep.value = 3
  }
}

// 处理仓库选择
const handleWarehouseChange = () => {
  if (form.value.warehouseId) {
    showPreview.value = true
  }
}

// 跳转到指定步骤
const goToStep = (step: number) => {
  if (step === 1) {
    currentStep.value = step
    showPreview.value = false
  } else if (step === 2 && selectedCustomer.value) {
    currentStep.value = step
    showPreview.value = false
  } else if (step === 3 && selectedProduct.value) {
    currentStep.value = step
    showPreview.value = false
  }
}

// 返回步骤选择
const backToSteps = () => {
  showPreview.value = false
}

// 提交订单
const submitOrder = async () => {
  submitting.value = true

  setTimeout(() => {
    const orderData = {
      customer: selectedCustomer.value,
      product: selectedProduct.value,
      quantity: form.value.quantity,
      warehouse: selectedWarehouse.value,
      totalAmount: totalAmount.value,
      orderTime: new Date().toLocaleString()
    }

    console.log('订单数据:', orderData)

    ElMessage.success('订单创建成功！')
    submitting.value = false
    createDialogVisible.value = false
    handleCreateClose()
  }, 1000)
}

// 关闭创建对话框
const handleCreateClose = () => {
  form.value = {
    customerId: '',
    productId: '',
    quantity: 1,
    warehouseId: ''
  }
  currentStep.value = 1
  showPreview.value = false
}

// 订单详情
const detailDrawerVisible = ref(false)
const currentOrder = ref<any>(null)

const handleView = (row: any) => {
  currentOrder.value = { ...row }
  detailDrawerVisible.value = true
}

const handleDetailClose = () => {
  currentOrder.value = null
}

// 顶部操作
const handleTopAction = ({ action }: { action: string }) => {
  if (action === 'create') {
    createDialogVisible.value = true
  }
}

// 批量操作
const handleBulkAction = async ({ action, rows }: { action: string; rows: any[] }) => {
  if (!rows.length) return ElMessage.warning('请选择数据')

  const ids = rows.map((r) => r._id).filter(Boolean)
  if (!ids.length) return ElMessage.error('缺少 ID')

  switch (action) {
    case 'delete':
      await handleBulkDelete(ids)
      break
    case 'export':
      await handleExport(rows)
      break
  }
}

// 批量删除
const handleBulkDelete = async (ids: string[]) => {
  const result = await elasticsearchService.bulkDelete(ids, 'order')

  if (result.success) {
    ElMessage.success(`已删除 ${result.successCount} 条`)
    if (result.failedCount > 0) {
      ElMessage.warning(`${result.failedCount} 条删除失败`)
      console.error(result.errors)
    }
    setTimeout(() => window.location.reload(), 500)
  }
}

// 导出
const handleExport = async (rows: any[]) => {
  const data = rows.map(({ _id, ...rest }) => rest)
  const headers = Object.keys(data[0] || {})
  const csvContent = [
    headers.join(','),
    ...data.map((row) =>
      headers.map((h) => JSON.stringify(row[h] || '')).join(',')
    ),
  ].join('\n')

  const blob = new Blob(['\ufeff' + csvContent], {
    type: 'text/csv;charset=utf-8;',
  })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')

  link.href = url
  link.download = `orders_${Date.now()}.csv`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)

  ElMessage.success('导出成功')
}
</script>

<style scoped lang="scss">
.order-management {
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

.order-creation-wrapper {
  position: relative;
  overflow: hidden;
  min-height: 60vh;
}

.steps-container {
  transition: transform 0.4s ease-in-out;
  padding: 10px 0;
  min-height: 60vh;
}

.steps-container.slide-left {
  transform: translateX(-100%);
}

.preview-container {
  position: absolute;
  top: 0;
  left: 100%;
  width: 100%;
  height: 100%;
  transition: transform 0.4s ease-in-out;
  padding: 20px 0;
}

.preview-container.show-preview {
  transform: translateX(-100%);
}

.preview-content-full {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.preview-actions {
  margin-top: 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20px;
  border-top: 1px solid #e4e7ed;
}

.order-creation-flow {
  padding: 10px 0;
  min-height: 60vh;
}

.step-section {
  margin-bottom: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
}

.step-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 18px 24px;
  background-color: #f5f7fa;
  cursor: pointer;
  transition: all 0.3s;
}

.step-header:hover:not(.disabled) {
  background-color: #ecf5ff;
}

.step-header.active {
  background-color: #409eff;
  color: white;
  padding: 20px 24px;
}

.step-header.completed {
  background-color: #f0f9ff;
  border-left: 4px solid #67c23a;
}

.step-header.disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.step-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.step-number {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: #409eff;
  color: white;
  font-weight: bold;
  font-size: 16px;
}

.step-header.active .step-number {
  background-color: white;
  color: #409eff;
}

.step-header.completed .step-number {
  background-color: #67c23a;
}

.step-name {
  font-size: 18px;
  font-weight: 500;
}

.step-summary-inline {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
}

.step-number-small {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background-color: #67c23a;
  color: white;
  font-weight: bold;
  font-size: 14px;
  flex-shrink: 0;
}

.summary-value {
  color: #303133;
  font-weight: 500;
  font-size: 16px;
}

.step-content {
  padding: 32px 24px;
  background-color: white;
}

.info-card {
  background: linear-gradient(135deg, #f5f7fa 0%, #ffffff 100%);
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  padding: 24px;
  height: 100%;
  transition: all 0.3s;
}

.info-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.info-card-title {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  color: #409eff;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid #409eff;
}

.info-card-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.info-item {
  display: flex;
  line-height: 1.8;
}

.info-label {
  color: #909399;
  min-width: 80px;
  font-size: 14px;
}

.info-value {
  color: #303133;
  font-weight: 500;
  flex: 1;
  font-size: 14px;
}

.total-amount-card {
  margin-top: 40px;
  padding: 32px;
  background: linear-gradient(135deg, #fff3e0 0%, #ffe0b2 100%);
  border: 2px solid #ff9800;
  border-radius: 12px;
  text-align: center;
}

.total-label {
  font-size: 16px;
  color: #666;
  margin-bottom: 12px;
}

.total-value {
  font-size: 36px;
  font-weight: bold;
  color: #f56c6c;
}

.order-detail {
  .detail-card {
    margin-bottom: 20px;

    &:last-child {
      margin-bottom: 0;
    }

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: 600;
      font-size: 16px;
    }
  }
}
</style>
