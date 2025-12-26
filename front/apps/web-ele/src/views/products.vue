<template>
  <div class="product-management">

    <DataTable
      :config="pageConfig"
      :loading="searchLoading"
      @view="handleView"
      @bulkAction="handleBulkAction"
      @topAction="handleTopAction"
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

    <!-- 产品详情抽屉 -->
    <el-drawer
      v-model="detailDialogVisible"
      title="产品详情"
      size="60%"
      @close="handleDetailClose"
    >
      <div v-if="currentProduct.id" class="product-detail">
        <!-- 基本信息 -->
        <el-card class="detail-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
              <el-button v-if="!isEditing" type="primary" size="small" @click="startEdit">
                修改
              </el-button>
              <div v-else>
                <el-button size="small" @click="cancelEdit">取消</el-button>
                <el-button type="primary" size="small" :loading="saving" @click="saveEdit">
                  保存
                </el-button>
              </div>
            </div>
          </template>

          <el-form v-if="isEditing" :model="editForm" label-width="100px">
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="产品名称">
                  <el-input v-model="editForm.name" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="状态">
                  <el-select v-model="editForm.status" placeholder="请选择状态" style="width: 100%">
                    <el-option label="草稿" value="draft" />
                    <el-option label="已提交" value="submitted" />
                    <el-option label="已审批" value="approved" />
                    <el-option label="已拒绝" value="rejected" />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <!-- 原料配置 -->
            <el-form-item label="原料配置">
              <el-button
                size="small"
                type="primary"
                @click="addMaterial"
                style="margin-bottom: 8px"
              >
                添加原料
              </el-button>

              <el-table :data="editForm.materials" border>
                <el-table-column label="原料">
                  <template #default="{ row }">
                    <el-select
                      v-model="row.material_id"
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
                      :min="0"
                      :max="100"
                    />
                  </template>
                </el-table-column>

                <el-table-column label="操作">
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
                :class="{ error: totalRatioEdit !== 100 && editForm.materials.length }"
              >
                总占比：{{ totalRatioEdit }}%
              </div>
            </el-form-item>

            <!-- 工艺配置 -->
            <el-form-item label="工艺配置">
              <el-select
                v-model="selectedProcessIdsEdit"
                multiple
                filterable
                placeholder="选择工艺"
                @visible-change="onProcessSelectVisible"
                @change="onProcessChangeEdit"
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

          <el-descriptions v-else :column="2" border>
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
            <el-descriptions-item label="更新时间" :span="2">
              {{ formatDate(currentProduct.updated_at) }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- 价格信息 -->
        <el-card class="detail-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>价格信息</span>
            </div>
          </template>
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
        </el-card>

        <!-- 产品公式 -->
        <el-card class="detail-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span>产品公式</span>
            </div>
          </template>

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
        </el-card>

        <!-- 操作区域 -->
        <div class="action-area">
          <el-button type="success" @click="handleGenerateQuote" size="large">
            生成报价单
          </el-button>
        </div>
      </div>
    </el-drawer>

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
import DataTable from '#/components/Table/index.vue'
import { elasticsearchService } from '#/api/core/es'
import { createProductApi, updateProductApi, getProductPriceApi, type ProductPriceResponse } from '#/api/core/product'
import { useDataTable } from '#/composables/useDataTable'

/* ================= 基础 ================= */

const { searchLoading } = useDataTable('product', 20)

const pageConfig = {
  pageType: 'product',
  index: 'product',
  rowKey: 'id',
  pageSize: 20,
  columns: [
    { key: 'id', label: 'ID', width: 80, showOverflowTooltip: true },
    { key: 'name', label: '产品名称', width: 200, showOverflowTooltip: true },
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
      showOverflowTooltip: true,
      formatter: (val: string) => val ? new Date(val).toLocaleString('zh-CN') : '-'
    },
    {
      key: 'updated_at',
      label: '更新时间',
      width: 180,
      showOverflowTooltip: true,
      formatter: (val: string) => val ? new Date(val).toLocaleString('zh-CN') : '-'
    },
  ],
  filters: [
    {
      key: 'status',
      label: '状态',
      type: 'select',
      placeholder: '请选择状态',
      fetchOptions: async () => {
        try {
          const res = await elasticsearchService.search({
            index: 'product',
            pagination: { offset: 0, size: 0 },
            aggRequests: {
              status: {
                type: 'terms',
                field: 'status',
                size: 20,
              },
            },
          })
          const buckets = res.aggregations?.status?.buckets || []
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
  ],
  topActions: [
    { label: '新建产品', key: 'create', type: 'primary' as const },
  ],
  bulkActions: [
    { label: '批量删除', key: 'delete', type: 'danger' as const, confirm: true }
  ],
  actions: ['view'],
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

/* ================= 详情 Drawer ================= */

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

// 编辑模式
const isEditing = ref(false)
const saving = ref(false)
const editForm = ref<any>({
  name: '',
  status: 'draft',
  materials: [],
  processes: []
})
const selectedProcessIdsEdit = ref<number[]>([])

// 编辑时的总占比
const totalRatioEdit = computed(() =>
  editForm.value.materials.reduce((s: number, m: any) => s + (m.ratioPercent || 0), 0)
)

// 开始编辑
const startEdit = () => {
  if (currentProduct.value.id) {
    editForm.value = {
      name: currentProduct.value.name,
      status: currentProduct.value.status,
      materials: currentProduct.value.materials.map((m: any) => ({
        material_id: m.material_id,
        ratioPercent: m.ratio * 100,
      })),
      processes: currentProduct.value.processes,
    }
    selectedProcessIdsEdit.value = currentProduct.value.processes.map((p: any) => p.process_id)
    isEditing.value = true
  }
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
  editForm.value = {
    name: '',
    status: 'draft',
    materials: [],
    processes: []
  }
  selectedProcessIdsEdit.value = []
}

// 工艺变更（编辑模式）
const onProcessChangeEdit = (ids: number[]) => {
  editForm.value.processes = ids.map(id => ({ process_id: id }))
}

// 保存编辑
const saveEdit = async () => {
  if (!currentProduct.value.id) {
    ElMessage.error('缺少必要的ID信息')
    return
  }

  // 验证总占比
  if (editForm.value.materials.length > 0 && totalRatioEdit.value !== 100) {
    ElMessage.error('原料占比必须等于100%')
    return
  }

  if (!editForm.value.name) {
    ElMessage.error('请输入产品名称')
    return
  }

  saving.value = true
  try {
    const payload = {
      name: editForm.value.name,
      status: editForm.value.status,
      materials: editForm.value.materials.map((m: any) => ({
        material_id: m.material_id,
        ratio: m.ratioPercent / 100,
      })),
      processes: editForm.value.processes,
    }

    await updateProductApi(currentProduct.value.id, payload)

    ElMessage.success('保存成功')
    isEditing.value = false

    // 更新当前产品信息
    currentProduct.value = {
      ...currentProduct.value,
      ...payload,
      materials: payload.materials,
      processes: payload.processes
    }

    // 重新加载页面
    setTimeout(() => window.location.reload(), 500)
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

/* ================= 报价单 Dialog ================= */

const quoteDialogVisible = ref(false)
const quoteForm = ref({
  client_id: undefined as number | undefined,
  quantity: 1,
  unit_price: 0
})
const clientsList = ref<any[]>([])
const generatingQuote = ref(false)

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

const handleView = async (row: any) => {
  currentProduct.value = { ...row }
  isEditing.value = false
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

/* ================= 顶部操作 & 批量 ================= */

const handleTopAction = ({ action }: any) => {
  if (action === 'create') handleCreate()
}

const handleBulkAction = ({ action }: any) => {
  // TODO: 实现批量删除
}

/* ================= reset ================= */

const handleDialogClose = () => {
  formRef.value?.clearValidate()
  formData.value = { name: '', materials: [], processes: [] }
  selectedProcessIds.value = []
}

const handleDetailClose = () => {
  isEditing.value = false
  editForm.value = {
    name: '',
    status: 'draft',
    materials: [],
    processes: []
  }
  selectedProcessIdsEdit.value = []
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

.ratio-summary {
  margin-top: 8px;
  font-weight: 600;
}

.ratio-summary.error {
  color: #f56c6c;
}

.product-detail {
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

  .action-area {
    display: flex;
    justify-content: center;
    padding: 20px 0;
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
}
</style>
