<template>
  <div class="order-list">
    <el-table
      :data="orders"
      stripe
      style="width: 100%"
      @row-click="handleRowClick"
    >
      <el-table-column prop="orderNo" label="订单号" width="180" fixed>
        <template #default="{ row }">
          <div class="order-no">
            <el-icon><Document /></el-icon>
            <span>{{ row.orderNo }}</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="clientName" label="客户名称" width="180" />

      <el-table-column prop="productName" label="产品名称" width="160" />

      <el-table-column prop="productCode" label="产品编号" width="140" />

      <el-table-column prop="requiredQuantity" label="需求数量" width="120" align="right">
        <template #default="{ row }">
          <span class="quantity">{{ row.requiredQuantity.toLocaleString() }}</span>
        </template>
      </el-table-column>

      <el-table-column prop="productHistoryShrinkage" label="历史缩率" width="120" align="center">
        <template #default="{ row }">
          <el-tag type="info" size="small">{{ row.productHistoryShrinkage }}%</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="订单整体进度" width="280">
        <template #default="{ row }">
          <div class="progress-cell">
            <el-progress
              :percentage="row.overallProgress"
              :color="getProgressColor(row.overallProgress)"
              :stroke-width="16"
            >
              <template #default="{ percentage }">
                <span class="progress-text">{{ percentage }}%</span>
              </template>
            </el-progress>
          </div>
        </template>
      </el-table-column>

      <el-table-column prop="status" label="状态" width="120" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 'completed' ? 'success' : 'warning'" size="large">
            {{ row.status === 'completed' ? '已完成' : '进行中' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="createdAt" label="创建时间" width="180" />

      <el-table-column label="操作" width="120" fixed="right" align="center">
        <template #default="{ row }">
          <el-button
            type="primary"
            size="small"
            @click.stop="handleViewDetail(row)"
          >
            查看详情
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { Document } from '@element-plus/icons-vue'
import type { Order } from '../types'

interface Props {
  orders: Order[]
  currentRole: 'sales' | 'follower' | 'warehouse'
}

interface Emits {
  (e: 'view-detail', order: Order): void
}

defineProps<Props>()
const emit = defineEmits<Emits>()

// 获取进度条颜色
const getProgressColor = (percentage: number): string => {
  if (percentage === 100) return '#67C23A'
  if (percentage >= 75) return '#409EFF'
  if (percentage >= 50) return '#E6A23C'
  return '#F56C6C'
}

// 查看详情
const handleViewDetail = (order: Order) => {
  emit('view-detail', order)
}

// 点击行也可查看详情
const handleRowClick = (order: Order) => {
  emit('view-detail', order)
}
</script>

<style scoped>
.order-list {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.order-no {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #409EFF;
  font-weight: 500;
}

.quantity {
  font-weight: 600;
  color: #303133;
}

.progress-cell {
  padding: 4px 0;
}

.progress-text {
  font-size: 13px;
  font-weight: 600;
}

:deep(.el-table__row) {
  cursor: pointer;
  transition: background-color 0.2s;
}

:deep(.el-table__row:hover) {
  background-color: #f5f7fa !important;
}

:deep(.el-progress__text) {
  font-size: 13px !important;
  font-weight: 600;
}
</style>
