<template>
  <el-drawer
    v-model="drawerVisible"
    :size="1200"
    :close-on-click-modal="false"
  >
    <template #header>
      <div class="drawer-header">
        <span class="order-title">{{ order?.orderNo }}</span>
        <el-button type="primary" size="default" @click="openOperationModal">
          操作
        </el-button>
      </div>
    </template>

    <div v-if="order" class="order-detail">
      <!-- 订单流程进度条（紧凑版） -->
      <div class="flow-progress-compact">
        <div
          v-for="(step, index) in flowSteps"
          :key="index"
          class="flow-step"
          :class="{
            'active': index === currentFlowStep,
            'finished': index < currentFlowStep,
            'clickable': step.clickable
          }"
          @click="handleStepClick(step, index)"
        >
          <div class="step-circle">
            <span v-if="index < currentFlowStep">✓</span>
            <span v-else>{{ index + 1 }}</span>
          </div>
          <div class="step-label">{{ step.label }}</div>
        </div>
      </div>

      <!-- 左右分栏：生产进度明细 + 基本信息 -->
      <el-row :gutter="16" class="main-content">
        <!-- 左：生产进度明细 -->
        <el-col :span="14">
          <div class="detail-section compact">
            <div class="section-title-simple">生产进度明细</div>
            <div class="progress-items-compact">
              <div
                v-for="item in visibleProgressItems"
                :key="item.type"
                class="progress-item-row"
              >
                <div class="item-header">
                  <span class="item-name">{{ item.name }}</span>
                  <span class="item-progress">{{ item.completedQuantity.toLocaleString() }} / {{ item.targetQuantity.toLocaleString() }}</span>
                </div>
                <div class="item-bar">
                  <el-progress
                    :percentage="item.progress"
                    :color="item.color"
                    :stroke-width="8"
                  >
                    <template #default="{ percentage }">
                      <span style="font-size: 12px">{{ percentage }}%</span>
                    </template>
                  </el-progress>
                </div>
              </div>
            </div>
          </div>
        </el-col>

        <!-- 右：基本信息 -->
        <el-col :span="10">
          <div class="detail-section compact">
            <div class="section-title-simple">基本信息</div>
            <div class="info-list">
              <div class="info-row">
                <span class="label">客户</span>
                <span class="value">{{ order.clientName }}</span>
              </div>
              <div class="info-row">
                <span class="label">产品</span>
                <span class="value">{{ order.productName }}</span>
              </div>
              <div class="info-row">
                <span class="label">产品编号</span>
                <span class="value">{{ order.productCode }}</span>
              </div>
              <div class="info-row">
                <span class="label">需求数量</span>
                <span class="value strong">{{ order.requiredQuantity.toLocaleString() }} 件</span>
              </div>
              <div class="info-row">
                <span class="label">历史缩率</span>
                <span class="value">{{ order.productHistoryShrinkage }}% <span class="muted">(参考)</span></span>
              </div>
              <div class="info-row">
                <span class="label">创建时间</span>
                <span class="value">{{ order.createdAt }}</span>
              </div>
              <div class="info-row">
                <span class="label">状态</span>
                <span class="value">
                  <el-tag :type="order.status === 'completed' ? 'success' : 'warning'" size="small">
                    {{ order.status === 'completed' ? '已完成' : '进行中' }}
                  </el-tag>
                </span>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>

      <!-- 操作日志 -->
      <div class="detail-section compact">
        <div class="section-title-simple">
          操作日志
          <span class="log-count">{{ order.operationLogs.length }} 条</span>
        </div>
        <OperationLog :logs="order.operationLogs" />
      </div>
    </div>

    <!-- 操作面板Modal -->
    <OperationModal
      v-model:visible="operationModalVisible"
      :order="order"
      :current-role="currentRole"
      @update="handleUpdate"
    />
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import OperationLog from './OperationLog.vue'
import OperationModal from './OperationModal.vue'
import type { Order, RoleType } from '../types'

interface Props {
  visible: boolean
  order: Order | null
  currentRole: RoleType
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'update', order: Order): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 抽屉显示状态
const drawerVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 操作面板Modal
const operationModalVisible = ref(false)

// 流程步骤定义
const flowSteps = [
  { label: '下单', clickable: false },
  { label: '派单', clickable: false },
  { label: '生产', clickable: true },
  { label: '核验', clickable: true },
  { label: '发货', clickable: false },
  { label: '完成', clickable: false }
]

// 只显示存在的进度项
const visibleProgressItems = computed(() => {
  if (!props.order) return []
  return props.order.progressItems.filter(item => item.exists)
})

// 计算当前流程步骤
const currentFlowStep = computed(() => {
  if (!props.order) return 0

  const fabricItem = props.order.progressItems.find(item => item.type === 'fabric_input')
  const productionItem = props.order.progressItems.find(item => item.type === 'production')
  const warehouseItem = props.order.progressItems.find(item => item.type === 'warehouse_check')
  const reworkItem = props.order.progressItems.find(item => item.type === 'rework')

  // 6. 完成 - 所有进度100%
  if (props.order.status === 'completed') return 6

  // 5. 发货 - 验货完成且（无次品或回修完成）
  if (warehouseItem && warehouseItem.progress === 100) {
    if (!reworkItem?.exists || (reworkItem.exists && reworkItem.progress === 100)) {
      return 5
    }
  }

  // 4. 核验 - 验货开始
  if (warehouseItem && warehouseItem.progress > 0) return 4

  // 3. 生产 - 生产开始
  if (productionItem && productionItem.progress > 0) return 3

  // 2. 派单 - 胚布投入开始
  if (fabricItem && fabricItem.progress > 0) return 2

  // 1. 下单 - 订单创建
  return 1
})

// 打开操作面板
const openOperationModal = () => {
  operationModalVisible.value = true
}

// 步骤点击
const handleStepClick = (step: any, index: number) => {
  if (!step.clickable) return
  console.log('点击了步骤:', step.label, index)
  // 这里可以添加跳转逻辑
}

// 处理更新
const handleUpdate = (updatedOrder: Order) => {
  emit('update', updatedOrder)
}
</script>

<style scoped>
.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.order-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.order-detail {
  padding: 0 4px;
}

/* 紧凑流程进度条 */
.flow-progress-compact {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
  position: relative;
}

.flow-progress-compact::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 60px;
  right: 60px;
  height: 2px;
  background: #e4e7ed;
  transform: translateY(-50%);
  z-index: 0;
}

.flow-step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  position: relative;
  z-index: 1;
}

.flow-step.clickable {
  cursor: pointer;
}

.flow-step.clickable:hover .step-circle {
  transform: scale(1.1);
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.3);
}

.step-circle {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: white;
  border: 2px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: #909399;
  transition: all 0.3s;
}

.flow-step.active .step-circle {
  background: #409EFF;
  border-color: #409EFF;
  color: white;
  box-shadow: 0 0 0 4px rgba(64, 158, 255, 0.1);
}

.flow-step.finished .step-circle {
  background: #67C23A;
  border-color: #67C23A;
  color: white;
}

.step-label {
  font-size: 12px;
  color: #606266;
  white-space: nowrap;
}

.flow-step.active .step-label {
  color: #409EFF;
  font-weight: 600;
}

.flow-step.finished .step-label {
  color: #67C23A;
}

/* 主内容区 */
.main-content {
  margin-bottom: 16px;
}

/* 精简的区块 */
.detail-section.compact {
  padding: 16px;
  background: white;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
  height: 100%;
}

/* 简单标题 */
.section-title-simple {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e4e7ed;
}

.log-count {
  font-size: 12px;
  color: #909399;
  font-weight: 400;
  margin-left: 8px;
}

/* 信息列表 */
.info-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-row {
  display: flex;
  align-items: center;
  font-size: 13px;
  padding: 4px 0;
}

.info-row .label {
  color: #909399;
  min-width: 80px;
  flex-shrink: 0;
}

.info-row .value {
  color: #303133;
  flex: 1;
}

.info-row .value.strong {
  font-weight: 600;
  color: #409EFF;
  font-size: 14px;
}

.info-row .muted {
  color: #909399;
  font-size: 12px;
}

/* 精简的进度项 */
.progress-items-compact {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.progress-item-row {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
}

.item-name {
  color: #303133;
  font-weight: 500;
}

.item-progress {
  color: #606266;
  font-size: 12px;
  font-family: monospace;
}

.item-bar {
  width: 100%;
}

:deep(.el-drawer__header) {
  margin-bottom: 16px;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
}

:deep(.el-drawer__body) {
  padding: 16px;
}
</style>
