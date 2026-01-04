<template>
  <el-drawer
    :model-value="visible"
    :title="isCreateMode ? '新增原料' : '原料详情'"
    size="60%"
    @close="handleClose"
  >
    <div v-if="isCreateMode || material" class="material-detail">
      <!-- 基本信息 -->
      <el-card class="detail-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <el-button v-if="!isCreateMode && !isEditing" type="primary" size="small" @click="startEdit">
              修改
            </el-button>
            <div v-if="isCreateMode || isEditing">
              <el-button size="small" @click="cancelEdit">取消</el-button>
              <el-button type="primary" size="small" :loading="saving" @click="saveEdit">
                {{ isCreateMode ? '创建' : '保存' }}
              </el-button>
            </div>
          </div>
        </template>

        <el-form v-if="isCreateMode || isEditing" :model="editForm" :rules="formRules" ref="formRef" label-width="100px">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="原料编码">
                <el-input v-model="editForm.code" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="原料名称">
                <el-input v-model="editForm.name" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="规格">
                <el-input v-model="editForm.spec" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="分类">
                <el-select v-model="editForm.category" placeholder="请选择分类" style="width: 100%">
                  <el-option label="胚布" value="胚布" />
                  <el-option label="染料" value="染料" />
                  <el-option label="助剂" value="助剂" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="单位">
                <el-input v-model="editForm.unit" />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="当前价格">
                <el-input-number v-model="editForm.currentPrice" :precision="2" :step="0.01" style="width: 100%" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="状态">
                <el-select v-model="editForm.status" placeholder="请选择状态" style="width: 100%">
                  <el-option label="在用" value="active" />
                  <el-option label="停用" value="inactive" />
                </el-select>
              </el-form-item>
            </el-col>
          </el-row>
          <el-form-item label="描述">
            <el-input v-model="editForm.description" type="textarea" :rows="3" />
          </el-form-item>
        </el-form>

        <el-descriptions v-else :column="2" border>
          <el-descriptions-item label="原料编码">
            {{ material.code || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="原料名称">
            {{ material.name }}
          </el-descriptions-item>
          <el-descriptions-item label="规格">
            {{ material.spec || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="分类">
            <el-tag v-if="material.category" :type="getCategoryTagType(material.category)" size="small">
              {{ material.category }}
            </el-tag>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="单位">
            {{ material.unit }}
          </el-descriptions-item>
          <el-descriptions-item label="当前价格">
            <span v-if="material.currentPrice" class="price">¥{{ material.currentPrice.toFixed(2) }}/{{ material.unit }}</span>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag v-if="material.status" :type="material.status === 'active' ? 'success' : 'info'" size="small">
              {{ material.status === 'active' ? '在用' : '停用' }}
            </el-tag>
            <span v-else>-</span>
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">
            {{ material.description || '-' }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 价格趋势图（仅在查看/编辑模式显示） -->
      <el-card v-if="!isCreateMode && priceHistory.length > 0" class="detail-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>价格趋势</span>
          </div>
        </template>
        <div v-loading="loadingPriceHistory" ref="priceChartRef" class="price-chart"></div>
      </el-card>

      <!-- 供应商报价（仅在查看/编辑模式显示） -->
      <el-card v-if="!isCreateMode" class="detail-card" shadow="never">
        <template #header>
          <div class="card-header">
            <div>
              <span>供应商报价</span>
              <span class="sub-header">（从最新到最旧）</span>
            </div>
            <el-button type="primary" size="small" @click="openAddQuoteDialog">
              新增报价
            </el-button>
          </div>
        </template>
        <el-table :data="supplierQuotes" border stripe>
          <el-table-column prop="supplierName" label="供应商名称" width="150" />
          <el-table-column label="价格" width="150" align="right">
            <template #default="{ row }">
              <span class="price">¥{{ row.price.toFixed(2) }}/{{ material.unit }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="quotedAt" label="报价时间" width="180" />
          <el-table-column prop="quotedBy" label="报价人" width="120" align="center" />
          <el-table-column prop="remark" label="备注" />
        </el-table>
      </el-card>
    </div>

    <!-- 新增报价对话框 -->
    <AddQuoteDialog
      v-model:visible="addQuoteDialogVisible"
      :material="material"
      @success="handleAddQuoteSuccess"
    />
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'
import type { Material, SupplierQuote } from './types'
import AddQuoteDialog from './AddQuoteDialog.vue'
import { elasticsearchService } from '#/api/core/es'
import { getMaterialPriceHistoryApi, type MaterialPriceHistory } from '#/api/core/material_quote'
import { createMaterial } from '#/api/core/material'

type PriceData = MaterialPriceHistory

interface Props {
  visible: boolean
  material: Material | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'updateMaterial', material: Material): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 供应商报价数据（使用ref以支持动态更新）
const supplierQuotes = ref<SupplierQuote[]>([])
const addQuoteDialogVisible = ref(false)

// 价格历史
const priceHistory = ref<PriceData[]>([])
const loadingPriceHistory = ref(false)
const priceChartRef = ref<HTMLElement>()

// 编辑模式
const isEditing = ref(false)
const saving = ref(false)
const editForm = ref<Partial<Material>>({})
const formRef = ref<FormInstance>()

// 新增模式（当 material 为 null 时）
const isCreateMode = computed(() => !props.material)

// 表单验证规则
const formRules: FormRules = {
  code: [{ required: true, message: '请输入原料编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入原料名称', trigger: 'blur' }]
}

// 开始编辑
const startEdit = () => {
  if (props.material) {
    editForm.value = { ...props.material }
    isEditing.value = true
  }
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
  editForm.value = {}
}

// 保存编辑
const saveEdit = async () => {
  // 表单验证
  const valid = await formRef.value?.validate()
  if (!valid) return

  saving.value = true
  try {
    if (isCreateMode.value) {
      // 新增模式
      const result = await createMaterial({
        materialCode: editForm.value.code!,
        materialName: editForm.value.name!,
        createdBy: 1 // TODO: 从用户状态获取实际用户ID
      })

      ElMessage.success('创建成功')

      // 通知父组件新增成功
      const newMaterial: Material = {
        id: `mat-${result.id}`,
        code: result.materialCode,
        name: result.materialName,
        category: editForm.value.category || '',
        unit: editForm.value.unit || 'kg',
        currentPrice: editForm.value.currentPrice || 0,
        status: editForm.value.status || 'active',
        updatedBy: '当前用户',
        createdAt: result.createdAt,
        updatedAt: result.updatedAt
      }
      emit('updateMaterial', newMaterial)
      emit('update:visible', false)
    } else {
      // 编辑模式
      if (!props.material?._id) {
        ElMessage.error('缺少必要的ID信息')
        return
      }

      const { _id, id, createdAt, updatedAt, ...updateData } = editForm.value as any

      await elasticsearchService.update(props.material._id, 'material', updateData)

      ElMessage.success('保存成功')
      isEditing.value = false

      // 通知父组件更新
      emit('updateMaterial', { ...props.material, ...updateData } as Material)
    }
  } catch (error) {
    console.error(isCreateMode.value ? '创建失败:' : '保存失败:', error)
    ElMessage.error(isCreateMode.value ? '创建失败，请稍后重试' : '保存失败')
  } finally {
    saving.value = false
  }
}

// 格式化日期
const formatDate = (date?: string) => {
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
          color: '#409eff'
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
                color: 'rgba(64, 158, 255, 0.3)'
              },
              {
                offset: 1,
                color: 'rgba(64, 158, 255, 0.05)'
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

// 监听 material 变化，加载报价数据
watch(() => props.material, async (newMaterial) => {
  if (newMaterial) {
    // 从后台加载供应商报价数据
    // TODO: 调用API获取报价数据
    supplierQuotes.value = []

    // 如果有 ID，加载实际价格历史
    if (typeof newMaterial.id === 'number') {
      loadingPriceHistory.value = true
      try {
        priceHistory.value = await getMaterialPriceHistoryApi(newMaterial.id)
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
  } else {
    // 新增模式：初始化表单
    supplierQuotes.value = []
    priceHistory.value = []
    editForm.value = {
      code: '',
      name: '',
      category: '',
      unit: 'kg',
      currentPrice: 0,
      status: 'active',
      description: ''
    }
  }
}, { immediate: true })

const handleClose = () => {
  emit('update:visible', false)
}

const getCategoryTagType = (category?: string) => {
  if (!category) return 'info'
  const typeMap: Record<string, any> = {
    '胚布': 'primary',
    '染料': 'warning',
    '助剂': 'success'
  }
  return typeMap[category] || 'info'
}

// 打开新增报价对话框
const openAddQuoteDialog = () => {
  addQuoteDialogVisible.value = true
}

// 处理新增报价成功
const handleAddQuoteSuccess = async (newQuote: SupplierQuote) => {
  // 将新报价添加到列表顶部（最新的）
  supplierQuotes.value.unshift(newQuote)

  // 更新原料的当前价格为最新报价
  if (props.material) {
    const updatedMaterial: Material = {
      ...props.material,
      currentPrice: newQuote.price,
      updatedBy: newQuote.quotedBy,
      updatedAt: newQuote.quotedAt
    }
    emit('updateMaterial', updatedMaterial)

    // 刷新价格历史图表
    if (typeof props.material.id === 'number') {
      loadingPriceHistory.value = true
      try {
        priceHistory.value = await getMaterialPriceHistoryApi(props.material.id)
        await nextTick()
        renderPriceChart()
      } catch (error) {
        console.error('刷新报价历史失败:', error)
      } finally {
        loadingPriceHistory.value = false
      }
    }
  }
}
</script>

<style scoped lang="scss">
.material-detail {
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

      .sub-header {
        margin-left: 8px;
        font-size: 12px;
        font-weight: 400;
        color: #909399;
      }
    }
  }

  .price {
    font-weight: 600;
    color: #f56c6c;
    font-size: 16px;
  }

  .price-chart {
    width: 100%;
    height: 350px;
  }

  :deep(.el-statistic) {
    text-align: center;

    .el-statistic__head {
      color: #909399;
      font-size: 14px;
    }

    .el-statistic__content {
      font-size: 32px;
      font-weight: 600;
      color: #409eff;
    }
  }
}
</style>
