<template>
  <div class="product-management">
    <!-- 顶部操作栏 -->
    <div class="top-bar">
      <h2 class="page-title">产品管理</h2>
      <div class="top-actions">
        <el-button type="warning" @click="handleSeedData" :loading="seeding">
          初始化数据
        </el-button>
        <el-button type="primary" @click="handleCreate">
          新增产品
        </el-button>
      </div>
    </div>

    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="handleView"
      @edit="handleEdit"
      @bulkAction="handleBulkAction"
    />

    <!-- 编辑/新建对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="900px"
      @close="handleDialogClose"
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="产品名称" prop="name">
          <el-input
            v-model="formData.name"
            :disabled="dialogMode === 'view'"
          />
        </el-form-item>

        <!-- 原料配置 -->
        <el-form-item label="原料配置" prop="materials">
          <el-button
            v-if="dialogMode !== 'view'"
            size="small"
            type="primary"
            @click="addMaterial"
            style="margin-bottom: 8px"
          >
            添加原料
          </el-button>

          <el-table :data="formData.materials" border>
            <el-table-column label="原料">
              <template #default="{ row }">
                <el-select
                  v-model="row.material_id"
                  :disabled="dialogMode === 'view'"
                  filterable
                  placeholder="选择原料"
                  @visible-change="onMaterialSelectVisible"
                >
                  <el-option
                    v-for="m in materialsList"
                    :key="m.id"
                    :label="m.name"
                    :value="m.id"
                  />
                </el-select>
              </template>
            </el-table-column>

            <el-table-column label="占比 %">
              <template #default="{ row }">
                <el-input-number
                  v-model="row.ratioPercent"
                  :disabled="dialogMode === 'view'"
                  :min="0"
                  :max="100"
                />
              </template>
            </el-table-column>

            <el-table-column v-if="dialogMode !== 'view'" label="操作">
              <template #default="{ $index }">
                <el-button
                  type="danger"
                  size="small"
                  @click="removeMaterial($index)"
                >
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <div
            class="ratio-summary"
            :class="{ error: totalRatio !== 100 && formData.materials.length }"
          >
            总占比：{{ totalRatio }}%
          </div>
        </el-form-item>

        <!-- 工艺配置 -->
        <el-form-item label="工艺配置" prop="processes">
          <el-select
            v-model="selectedProcessIds"
            multiple
            filterable
            placeholder="选择工艺"
            :disabled="dialogMode === 'view'"
            @visible-change="onProcessSelectVisible"
            @change="onProcessChange"
          >
            <el-option
              v-for="p in processesList"
              :key="p.id"
              :label="p.name"
              :value="p.id"
            />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer v-if="dialogMode !== 'view'">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          保存
        </el-button>
      </template>
    </el-dialog>

    <!-- 产品详情大弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="产品详情"
      width="1000px"
      @close="handleDetailClose"
    >
      <!-- 基本信息 -->
      <div class="detail-section">
        <h3 class="section-title">基本信息</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="产品ID">
            {{ currentProduct.id }}
          </el-descriptions-item>
          <el-descriptions-item label="产品名称">
            {{ currentProduct.name }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(currentProduct.status)">
              {{ getStatusText(currentProduct.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(currentProduct.created_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 价格信息 -->
      <div class="detail-section">
        <h3 class="section-title">价格信息</h3>
        <el-row :gutter="20" v-loading="loadingPrice">
          <el-col :span="8">
            <el-statistic title="当前价格" :value="priceInfo.current_price" :precision="2">
              <template #prefix>¥</template>
            </el-statistic>
          </el-col>
          <el-col :span="8">
            <el-statistic title="历史最高" :value="priceInfo.historical_high" :precision="2">
              <template #prefix>¥</template>
            </el-statistic>
          </el-col>
          <el-col :span="8">
            <el-statistic title="历史最低" :value="priceInfo.historical_low" :precision="2">
              <template #prefix>¥</template>
            </el-statistic>
          </el-col>
        </el-row>
      </div>

      <!-- 产品公式 -->
      <div class="detail-section">
        <div class="section-header">
          <h3 class="section-title">产品公式</h3>
          <el-button v-if="canEditFormula" type="primary" size="small" @click="handleEditFormula">
            修改公式
          </el-button>
        </div>

        <div class="formula-content">
          <div class="formula-item">
            <h4>原料配置</h4>
            <el-table :data="currentProduct.materials" border>
              <el-table-column label="原料名称" width="200">
                <template #default="{ row }">
                  {{ getMaterialName(row.material_id) }}
                </template>
              </el-table-column>
              <el-table-column label="占比" align="center">
                <template #default="{ row }">
                  {{ (row.ratio * 100).toFixed(2) }}%
                </template>
              </el-table-column>
            </el-table>
          </div>

          <div class="formula-item">
            <h4>工艺配置</h4>
            <el-table :data="currentProduct.processes" border>
              <el-table-column label="工艺名称" width="200">
                <template #default="{ row }">
                  {{ getProcessName(row.process_id) }}
                </template>
              </el-table-column>
              <el-table-column label="数量" align="center">
                <template #default="{ row }">
                  {{ row.quantity || '-' }}
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
        <el-button type="success" @click="handleGenerateQuote">
          生成报价单
        </el-button>
      </template>
    </el-dialog>

    <!-- 生成报价单对话框 -->
    <el-dialog
      v-model="quoteDialogVisible"
      title="生成报价单"
      width="600px"
      @close="handleQuoteDialogClose"
    >
      <el-form :model="quoteForm" label-width="120px">
        <el-form-item label="产品">
          <el-input :value="currentProduct.name" disabled />
        </el-form-item>
        <el-form-item label="客户" required>
          <el-select
            v-model="quoteForm.client_id"
            placeholder="请选择客户"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="client in clientsList"
              :key="client.id"
              :label="client.name"
              :value="client.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="数量" required>
          <el-input-number
            v-model="quoteForm.quantity"
            :min="1"
            :step="1"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="单价" required>
          <el-input-number
            v-model="quoteForm.unit_price"
            :min="0"
            :precision="2"
            :step="0.01"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="总价">
          <el-input
            :value="`¥${(quoteForm.quantity * quoteForm.unit_price).toFixed(2)}`"
            disabled
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="quoteDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="generatingQuote" @click="handleQuoteSubmit">
          生成
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { useAccessStore } from '@vben/stores'
import DataTable from '#/components/Table/index.vue'
import { elasticsearchService } from '#/api/core/es'
import { createProductApi, updateProductApi, getProductPriceApi, type ProductPriceResponse } from '#/api/core/product'
import { useDataTable } from '#/composables/useDataTable'
import { seedData } from '#/utils/seedData'

/* ================= 基础 ================= */

const { searchLoading } = useDataTable('product', 20)
const accessStore = useAccessStore()

const pageConfig = {
  pageType: 'product',
  title: '产品管理',
  index: 'product',
  rowKey: 'id',
  pageSize: 20,
  columns: [
    { key: 'id', label: 'ID', width: 80 },
    { key: 'name', label: '产品名称', width: 200 },
    {
      key: 'status',
      label: '状态',
      width: 100,
      formatter: (val: string) => {
        const textMap: Record<string, string> = {
          draft: '草稿',
          submitted: '已提交',
          approved: '已审批',
          rejected: '已拒绝'
        }
        return textMap[val] || val
      }
    },
    {
      key: 'created_at',
      label: '创建时间',
      width: 180,
      formatter: (val: string) => val ? new Date(val).toLocaleString('zh-CN') : '-'
    },
    {
      key: 'updated_at',
      label: '更新时间',
      width: 180,
      formatter: (val: string) => val ? new Date(val).toLocaleString('zh-CN') : '-'
    },
  ],
  filters: [],
  bulkActions: [
    { label: '新建产品', key: 'create', type: 'primary' as const },
    { label: '批量删除', key: 'delete', type: 'danger' as const, confirm: true }
  ],
  actions: ['view', 'edit', 'create'],
}

/* ================= 编辑 Dialog ================= */

const dialogVisible = ref(false)
const dialogMode = ref<'create' | 'edit' | 'view'>('view')
const dialogTitle = computed(() =>
  dialogMode.value === 'create'
    ? '新建产品'
    : dialogMode.value === 'edit'
    ? '编辑产品'
    : '查看产品'
)

const submitting = ref(false)
const formRef = ref<FormInstance>()

/* ================= 详情 Dialog ================= */

const detailDialogVisible = ref(false)
const currentProduct = ref<any>({
  id: 0,
  name: '',
  status: 'draft',
  materials: [],
  processes: [],
  created_at: '',
  updated_at: ''
})

const priceInfo = ref<ProductPriceResponse>({
  current_price: 0,
  historical_high: 0,
  historical_low: 0
})

const loadingPrice = ref(false)

// 权限判断：是否可以修改公式（简单判断，可根据实际需求修改）
const canEditFormula = computed(() => {
  // 这里可以根据用户角色判断，例如只有管理员或产品经理可以修改
  const userRoles = accessStore.accessCodes || []
  return userRoles.includes('admin') || userRoles.includes('product_manager')
})

/* ================= 报价单 Dialog ================= */

const quoteDialogVisible = ref(false)
const quoteForm = ref({
  client_id: undefined as number | undefined,
  quantity: 1,
  unit_price: 0
})
const clientsList = ref<any[]>([])
const generatingQuote = ref(false)
const seeding = ref(false)

/* ================= 数据 ================= */

const formData = ref({
  name: '',
  materials: [] as { material_id: number; ratioPercent: number }[],
  processes: [] as { process_id: number }[],
})

const selectedProcessIds = ref<number[]>([])

/* ================= 校验 ================= */

const totalRatio = computed(() =>
  formData.value.materials.reduce((s, m) => s + (m.ratioPercent || 0), 0)
)

const formRules: FormRules = {
  name: [{ required: true, message: '请输入产品名称' }],
  materials: [
    {
      validator: (_, __, cb) => {
        if (!formData.value.materials.length) {
          cb(new Error('至少一个原料'))
        } else if (totalRatio.value !== 100) {
          cb(new Error('原料占比必须等于100%'))
        } else cb()
      },
    },
  ],
}

/* ================= 懒加载 ================= */

const materialsList = ref<any[]>([])
const processesList = ref<any[]>([])

const materialsLoaded = ref(false)
const processesLoaded = ref(false)

const loadMaterials = async () => {
  if (materialsLoaded.value) return
  try {
    const res = await elasticsearchService.search({
      index: 'material',
      pagination: { offset: 0, size: 100 }, // 后端最大限制100
    })
    materialsList.value = res.items || []
    materialsLoaded.value = true
  } catch (error) {
    console.error('加载原料列表失败:', error)
    materialsList.value = []
  }
}

const loadProcesses = async () => {
  if (processesLoaded.value) return
  try {
    const res = await elasticsearchService.search({
      index: 'process',
      pagination: { offset: 0, size: 100 }, // 后端最大限制100
    })
    processesList.value = res.items || []
    processesLoaded.value = true
  } catch (error) {
    console.error('加载工艺列表失败:', error)
    processesList.value = []
  }
}

const loadClients = async () => {
  try {
    const res = await elasticsearchService.search({
      index: 'client',
      pagination: { offset: 0, size: 100 }, // 后端最大限制100
    })
    clientsList.value = res.items || []
  } catch (error) {
    console.error('加载客户列表失败:', error)
  }
}

const onMaterialSelectVisible = (visible: boolean) => {
  if (visible) loadMaterials()
}

const onProcessSelectVisible = (visible: boolean) => {
  if (visible) loadProcesses()
}

/* ================= 辅助函数 ================= */

const getMaterialName = (materialId: number) => {
  const material = materialsList.value.find(m => m.id === materialId)
  return material?.name || `原料ID: ${materialId}`
}

const getProcessName = (processId: number) => {
  const process = processesList.value.find(p => p.id === processId)
  return process?.name || `工艺ID: ${processId}`
}

const getStatusType = (status: string) => {
  const typeMap: Record<string, any> = {
    draft: 'info',
    submitted: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: string) => {
  const textMap: Record<string, string> = {
    draft: '草稿',
    submitted: '已提交',
    approved: '已审批',
    rejected: '已拒绝'
  }
  return textMap[status] || status
}

const formatDate = (date: string) => {
  return date ? new Date(date).toLocaleString('zh-CN') : '-'
}

/* ================= 行为 ================= */

const addMaterial = () => {
  formData.value.materials.push({ material_id: 0, ratioPercent: 0 })
}

const removeMaterial = (i: number) => {
  formData.value.materials.splice(i, 1)
}

const onProcessChange = (ids: number[]) => {
  formData.value.processes = ids.map(id => ({ process_id: id }))
}

const handleCreate = async () => {
  dialogMode.value = 'create'
  dialogVisible.value = true
  await Promise.all([loadMaterials(), loadProcesses()])
}

const handleSeedData = async () => {
  try {
    await ElMessageBox.confirm(
      '将插入10条原料数据和10条工艺数据到数据库，确定继续吗？',
      '初始化测试数据',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    seeding.value = true
    const result = await seedData()

    ElMessage.success(`数据初始化成功！原料: ${result.materials}条，工艺: ${result.processes}条`)

    // 重置加载状态，重新加载数据
    materialsLoaded.value = false
    processesLoaded.value = false
    await Promise.all([loadMaterials(), loadProcesses()])
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('初始化数据失败:', error)
    }
  } finally {
    seeding.value = false
  }
}

const handleEdit = async (row: any) => {
  dialogMode.value = 'edit'
  formData.value = {
    name: row.name,
    materials: row.materials.map((m: any) => ({
      material_id: m.material_id,
      ratioPercent: m.ratio * 100,
    })),
    processes: row.processes,
  }
  selectedProcessIds.value = row.processes.map((p: any) => p.process_id)
  dialogVisible.value = true
  await Promise.all([loadMaterials(), loadProcesses()])
}

const handleView = async (row: any) => {
  currentProduct.value = { ...row }
  detailDialogVisible.value = true

  // 加载数据
  await Promise.all([loadMaterials(), loadProcesses()])

  // 加载价格信息
  if (row.id) {
    loadingPrice.value = true
    try {
      priceInfo.value = await getProductPriceApi(row.id)
    } catch (error) {
      console.error('加载价格信息失败:', error)
      ElMessage.warning('加载价格信息失败')
    } finally {
      loadingPrice.value = false
    }
  }
}

const handleEditFormula = () => {
  detailDialogVisible.value = false
  handleEdit(currentProduct.value)
}

const handleGenerateQuote = async () => {
  quoteForm.value = {
    client_id: undefined,
    quantity: 1,
    unit_price: priceInfo.value.current_price || 0
  }

  await loadClients()
  quoteDialogVisible.value = true
}

const handleQuoteSubmit = async () => {
  if (!quoteForm.value.client_id) {
    ElMessage.warning('请选择客户')
    return
  }

  if (quoteForm.value.quantity <= 0) {
    ElMessage.warning('请输入有效的数量')
    return
  }

  if (quoteForm.value.unit_price <= 0) {
    ElMessage.warning('请输入有效的单价')
    return
  }

  generatingQuote.value = true
  try {
    // 这里调用生成报价单的API
    // 目前先用ElMessage模拟
    ElMessage.success('报价单生成成功！')
    quoteDialogVisible.value = false

    // 实际应该是：
    // const quoteData = {
    //   product_id: currentProduct.value.id,
    //   client_id: quoteForm.value.client_id,
    //   quantity: quoteForm.value.quantity,
    //   unit_price: quoteForm.value.unit_price,
    //   total_price: quoteForm.value.quantity * quoteForm.value.unit_price
    // }
    // await createQuoteApi(quoteData)
  } catch (error: any) {
    console.error('生成报价单失败:', error)
    ElMessage.error(error.response?.data?.error || '生成报价单失败')
  } finally {
    generatingQuote.value = false
  }
}

/* ================= 提交 ================= */

const handleSubmit = async () => {
  await formRef.value?.validate()
  submitting.value = true

  const payload = {
    name: formData.value.name,
    materials: formData.value.materials.map(m => ({
      material_id: m.material_id,
      ratio: m.ratioPercent / 100,
    })),
    processes: formData.value.processes,
  }

  try {
    dialogMode.value === 'create'
      ? await createProductApi(payload)
      : await updateProductApi(currentProduct.value.id, payload)

    ElMessage.success('保存成功')
    dialogVisible.value = false
    location.reload()
  } finally {
    submitting.value = false
  }
}

/* ================= 批量 ================= */

const handleBulkAction = ({ action }: any) => {
  if (action === 'create') handleCreate()
}

/* ================= reset ================= */

const handleDialogClose = () => {
  formRef.value?.clearValidate()
  formData.value = { name: '', materials: [], processes: [] }
  selectedProcessIds.value = []
}

const handleDetailClose = () => {
  currentProduct.value = {
    id: 0,
    name: '',
    status: 'draft',
    materials: [],
    processes: [],
    created_at: '',
    updated_at: ''
  }
  priceInfo.value = {
    current_price: 0,
    historical_high: 0,
    historical_low: 0
  }
}

const handleQuoteDialogClose = () => {
  quoteForm.value = {
    client_id: undefined,
    quantity: 1,
    unit_price: 0
  }
}
</script>

<style scoped>
.product-management {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.top-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #fff;
  border-bottom: 1px solid #ebeef5;
}

.page-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.top-actions {
  display: flex;
  gap: 12px;
}

.ratio-summary {
  margin-top: 8px;
  font-weight: 600;
}

.ratio-summary.error {
  color: #f56c6c;
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

.formula-content {
  display: flex;
  flex-direction: column;
  gap: 20px;

  .formula-item {
    h4 {
      margin: 0 0 12px 0;
      font-size: 14px;
      font-weight: 600;
      color: #606266;
    }
  }
}

:deep(.el-statistic) {
  text-align: center;

  .el-statistic__head {
    color: #909399;
    font-size: 14px;
  }

  .el-statistic__content {
    font-size: 24px;
    font-weight: 600;
  }
}
</style>
