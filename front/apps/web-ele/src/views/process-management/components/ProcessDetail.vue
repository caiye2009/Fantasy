<template>
  <el-drawer
    :model-value="visible"
    title="工艺详情"
    size="60%"
    @close="handleClose"
  >
    <div v-if="process" class="process-detail">
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
          <el-form-item label="工艺名称">
            <el-input v-model="editForm.name" />
          </el-form-item>
          <el-form-item label="描述">
            <el-input v-model="editForm.description" type="textarea" :rows="3" />
          </el-form-item>
        </el-form>

        <el-descriptions v-else :column="2" border>
          <el-descriptions-item label="ID">
            {{ process.id }}
          </el-descriptions-item>
          <el-descriptions-item label="工艺名称">
            {{ process.name }}
          </el-descriptions-item>
          <el-descriptions-item label="描述" :span="2">
            {{ process.description || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(process.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            {{ formatDate(process.updatedAt) }}
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 价格趋势图 -->
      <el-card v-if="priceHistory.length > 0" class="detail-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>价格趋势</span>
          </div>
        </template>
        <div v-loading="loadingPriceHistory" ref="priceChartRef" class="price-chart"></div>
      </el-card>

      <!-- 报价历史 -->
      <el-card class="detail-card" shadow="never">
        <template #header>
          <div class="card-header">
            <div>
              <span>报价历史</span>
            </div>
            <el-button type="primary" size="small" @click="openAddQuoteDialog">
              新增报价
            </el-button>
          </div>
        </template>
        <el-table :data="priceHistory" v-loading="loadingPriceHistory" border stripe>
          <el-table-column prop="supplier_name" label="供应商" width="200" show-overflow-tooltip />
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
      </el-card>
    </div>

    <!-- 新增报价对话框 -->
    <AddQuoteDialog
      v-model:visible="addQuoteDialogVisible"
      :process="process"
      @success="handleAddQuoteSuccess"
    />
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'
import type { Process } from '../types'
import AddQuoteDialog from './AddQuoteDialog.vue'
import { getProcessPriceHistoryApi, type PriceData } from '#/api/core/pricing'
import { elasticsearchService } from '#/api/core/es'

interface Props {
  visible: boolean
  process: Process | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'updateProcess', process: Process): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 价格历史
const priceHistory = ref<PriceData[]>([])
const loadingPriceHistory = ref(false)
const priceChartRef = ref<HTMLElement>()
const addQuoteDialogVisible = ref(false)

// 编辑模式
const isEditing = ref(false)
const saving = ref(false)
const editForm = ref<Partial<Process>>({})

// 开始编辑
const startEdit = () => {
  if (props.process) {
    editForm.value = { ...props.process }
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
  if (!props.process?._id) {
    ElMessage.error('缺少必要的ID信息')
    return
  }

  saving.value = true
  try {
    const { _id, id, createdAt, updatedAt, ...updateData } = editForm.value as any

    await elasticsearchService.update(props.process._id, 'process', updateData)

    ElMessage.success('保存成功')
    isEditing.value = false

    // 通知父组件更新
    emit('updateProcess', { ...props.process, ...updateData } as Process)
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败')
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

// 监听 process 变化，加载报价数据
watch(() => props.process, async (newProcess) => {
  if (newProcess && typeof newProcess.id === 'number') {
    loadingPriceHistory.value = true
    try {
      priceHistory.value = await getProcessPriceHistoryApi(newProcess.id)
      // 等待DOM更新后渲染图表
      await nextTick()
      renderPriceChart()
    } catch (error) {
      console.error('加载报价历史失败:', error)
      priceHistory.value = []
    } finally {
      loadingPriceHistory.value = false
    }
  } else {
    priceHistory.value = []
  }
}, { immediate: true })

const handleClose = () => {
  isEditing.value = false
  emit('update:visible', false)
}

// 打开新增报价对话框
const openAddQuoteDialog = () => {
  addQuoteDialogVisible.value = true
}

// 处理新增报价成功
const handleAddQuoteSuccess = async () => {
  // 刷新价格历史图表
  if (props.process && typeof props.process.id === 'number') {
    loadingPriceHistory.value = true
    try {
      priceHistory.value = await getProcessPriceHistoryApi(props.process.id)
      await nextTick()
      renderPriceChart()
    } catch (error) {
      console.error('刷新报价历史失败:', error)
    } finally {
      loadingPriceHistory.value = false
    }
  }
}
</script>

<style scoped lang="scss">
.process-detail {
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

  .price-chart {
    width: 100%;
    height: 350px;
  }
}
</style>
