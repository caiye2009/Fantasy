<template>
  <div class="operation-log">
    <div v-if="logs.length > 0" class="log-list">
      <div
        v-for="log in sortedLogs"
        :key="log.id"
        class="log-row"
      >
        <div class="log-time">{{ formatTime(log.createdAt) }}</div>
        <div class="log-content">
          <div class="log-main">
            <el-tag :type="getRoleTagType(log.role)" size="small" effect="plain">
              {{ getRoleName(log.role) }}
            </el-tag>
            <span class="log-desc">{{ log.description }}</span>
          </div>
          <div v-if="log.operatorName !== '系统'" class="log-operator">
            {{ log.operatorName }}
          </div>
        </div>
      </div>
    </div>

    <div v-else class="empty-log">
      <span>暂无日志</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { OperationLog, RoleType } from '../types'

interface Props {
  logs: OperationLog[]
}

const props = defineProps<Props>()

// 按时间倒序排列
const sortedLogs = computed(() => {
  return [...props.logs].sort((a, b) => {
    return new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
  })
})

// 获取角色标签类型
const getRoleTagType = (role: RoleType) => {
  const typeMap: Record<RoleType, any> = {
    sales: 'primary',
    follower: 'success',
    warehouse: 'warning',
    system: 'info'
  }
  return typeMap[role] || 'info'
}

// 获取角色名称
const getRoleName = (role: RoleType) => {
  const nameMap: Record<RoleType, string> = {
    sales: '业务',
    follower: '跟单',
    warehouse: '仓库',
    system: '系统'
  }
  return nameMap[role] || role
}

// 格式化时间
const formatTime = (dateStr: string) => {
  // 简化时间显示：12-19 14:30
  const date = new Date(dateStr)
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${month}-${day} ${hour}:${minute}`
}
</script>

<style scoped>
.operation-log {
  max-height: 400px;
  overflow-y: auto;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.log-row {
  display: flex;
  gap: 12px;
  padding: 8px 0;
  font-size: 13px;
  border-bottom: 1px solid #f0f0f0;
}

.log-row:last-child {
  border-bottom: none;
}

.log-time {
  flex-shrink: 0;
  width: 80px;
  color: #909399;
  font-size: 12px;
  font-family: monospace;
}

.log-content {
  flex: 1;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.log-main {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.log-desc {
  color: #606266;
}

.log-operator {
  flex-shrink: 0;
  color: #909399;
  font-size: 12px;
}

.empty-log {
  padding: 20px;
  text-align: center;
  color: #909399;
  font-size: 13px;
}

/* 滚动条样式 */
.operation-log::-webkit-scrollbar {
  width: 6px;
}

.operation-log::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}

.operation-log::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}
</style>
