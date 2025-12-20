<template>
  <div class="progress-item-card" :style="{ borderLeftColor: progressItem.color }">
    <div class="card-header">
      <div class="header-left">
        <div class="icon-wrapper" :style="{ backgroundColor: progressItem.color + '20', color: progressItem.color }">
          <el-icon :size="24">
            <component :is="getIconComponent(progressItem.icon)" />
          </el-icon>
        </div>
        <div class="header-info">
          <div class="item-name">{{ progressItem.name }}</div>
          <div class="item-type-tag">
            <el-tag :type="getTagType(progressItem.type)" size="small">
              {{ getTypeName(progressItem.type) }}
            </el-tag>
          </div>
        </div>
      </div>
      <div class="progress-percentage" :style="{ color: progressItem.color }">
        {{ progressItem.progress }}%
      </div>
    </div>

    <div class="card-content">
      <div class="quantity-info">
        <div class="quantity-item">
          <span class="quantity-label">已完成</span>
          <span class="quantity-value completed">
            {{ progressItem.completedQuantity.toLocaleString() }}
          </span>
        </div>
        <div class="quantity-divider">/</div>
        <div class="quantity-item">
          <span class="quantity-label">目标</span>
          <span class="quantity-value target">
            {{ progressItem.targetQuantity.toLocaleString() }}
          </span>
        </div>
      </div>

      <div class="progress-bar-wrapper">
        <el-progress
          :percentage="progressItem.progress"
          :color="progressItem.color"
          :stroke-width="12"
          :show-text="false"
        />
      </div>

      <div class="card-footer">
        <div class="footer-info">
          <el-icon><Clock /></el-icon>
          <span>完成度：{{ getProgressStatus(progressItem.progress) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Clock } from '@element-plus/icons-vue'
import type { ProgressItem } from '../types'

interface Props {
  progressItem: ProgressItem
}

defineProps<Props>()

// 获取图标组件（这里简化处理，实际项目中可以用动态导入）
const getIconComponent = (icon?: string) => {
  // Element Plus 图标
  if (!icon) return 'Box'

  // 简化处理：直接返回图标名称的最后一部分
  const iconName = icon.split(':')[1] || icon
  const iconMap: Record<string, string> = {
    'package-open': 'Box',
    'factory': 'OfficeBuilding',
    'clipboard-check': 'Checked',
    'wrench': 'Tools'
  }
  return iconMap[iconName] || 'Box'
}

// 获取类型标签
const getTagType = (type: string) => {
  const typeMap: Record<string, string> = {
    fabric_input: '',
    production: 'success',
    warehouse_check: 'warning',
    rework: 'danger'
  }
  return typeMap[type] || 'info'
}

// 获取类型名称
const getTypeName = (type: string) => {
  const nameMap: Record<string, string> = {
    fabric_input: '投入',
    production: '生产',
    warehouse_check: '验货',
    rework: '回修'
  }
  return nameMap[type] || type
}

// 获取进度状态描述
const getProgressStatus = (progress: number): string => {
  if (progress === 100) return '已完成'
  if (progress >= 75) return '接近完成'
  if (progress >= 50) return '进行中'
  if (progress >= 25) return '刚开始'
  return '未开始'
}
</script>

<style scoped>
.progress-item-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  border-left: 4px solid;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  transition: all 0.3s;
}

.progress-item-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 12px;
}

.header-info {
  flex: 1;
}

.item-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 4px;
}

.item-type-tag {
  display: flex;
  align-items: center;
}

.progress-percentage {
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
}

.card-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.quantity-info {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
}

.quantity-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.quantity-label {
  font-size: 12px;
  color: #909399;
}

.quantity-value {
  font-size: 20px;
  font-weight: 600;
}

.quantity-value.completed {
  color: #67C23A;
}

.quantity-value.target {
  color: #409EFF;
}

.quantity-divider {
  font-size: 20px;
  font-weight: 300;
  color: #DCDFE6;
}

.progress-bar-wrapper {
  margin: 4px 0;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 8px;
  border-top: 1px solid #E4E7ED;
}

.footer-info {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #606266;
}
</style>
