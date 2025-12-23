<template>
  <div class="order-management-page">
    <!-- 顶部工具栏 -->
    <div class="page-header">
      <div class="header-left">
        <h2>订单管理</h2>
        <el-tag type="info" size="large">共 {{ orders.length }} 个订单</el-tag>
      </div>
      <div class="header-right">
        <el-tag :type="getRoleTagType(currentRole)" size="large">
          当前角色：{{ getRoleName(currentRole) }}
        </el-tag>
        <el-button type="primary" size="large" @click="openCreateOrderDialog">
          创建订单
        </el-button>
      </div>
    </div>

    <!-- 订单列表 -->
    <OrderList
      :orders="orders"
      :current-role="currentRole"
      @view-detail="handleViewDetail"
    />

    <!-- 订单详情抽屉 -->
    <OrderDetail
      v-model:visible="detailVisible"
      :order="selectedOrder"
      :current-role="currentRole"
      @update="handleOrderUpdate"
    />

    <!-- 创建订单对话框 -->
    <CreateOrderDialog
      v-model:visible="createOrderDialogVisible"
      @success="handleCreateOrderSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@vben/stores'
import { ElMessage } from 'element-plus'
import OrderList from './components/OrderList.vue'
import OrderDetail from './components/OrderDetail.vue'
import CreateOrderDialog from './components/CreateOrderDialog.vue'
import type { Order, RoleType, ProgressItem } from './types'
import { mapToOrderRole, getOrderRoleName } from '#/utils/roleMapping'
import { getOrderList, type OrderDetailResponse } from '#/api/core/order'

// 获取用户信息
const userStore = useUserStore()

// 从用户信息中获取当前角色
const currentRole = computed<RoleType>(() => {
  // 获取用户的第一个角色（如果有多个角色）
  const backendRole = userStore.userInfo?.role || userStore.userRoles?.[0] || ''
  // 映射到订单管理角色
  return mapToOrderRole(backendRole)
})

// 订单数据
const orders = ref<Order[]>([])
const loading = ref(false)

// 将后端数据转换为前端格式
function convertToFrontendOrder(backendOrder: OrderDetailResponse): Order {
  return {
    id: String(backendOrder.id),
    orderNo: backendOrder.order_no,
    clientName: backendOrder.client_name,
    productName: backendOrder.product_name,
    productCode: backendOrder.product_code,
    productHistoryShrinkage: backendOrder.product_history_shrinkage,
    requiredQuantity: backendOrder.required_quantity,
    createdAt: new Date(backendOrder.created_at).toLocaleString('zh-CN'),
    updatedAt: new Date(backendOrder.updated_at).toLocaleString('zh-CN'),
    status: backendOrder.status === 'completed' ? 'completed' : 'in_progress',
    progressItems: backendOrder.progress_items.map(item => ({
      type: item.type as ProgressItem['type'],
      name: item.name,
      targetQuantity: item.target_quantity,
      completedQuantity: item.completed_quantity,
      progress: item.progress,
      exists: item.exists,
      icon: item.icon,
      color: item.color
    })),
    overallProgress: backendOrder.overall_progress,
    operationLogs: backendOrder.operation_logs.map(log => ({
      id: String(log.id),
      orderId: String(log.order_id),
      operator: String(log.operator_id),
      operatorName: log.operator_name,
      role: log.operator_role as any,
      operationType: log.event_type as any,
      operationName: log.description,
      beforeData: log.before_data ? JSON.parse(log.before_data) : null,
      afterData: log.after_data ? JSON.parse(log.after_data) : null,
      description: log.description,
      createdAt: new Date(log.created_at).toLocaleString('zh-CN')
    }))
  }
}

// 加载订单数据
async function loadOrders() {
  try {
    loading.value = true
    const response = await getOrderList({ limit: 100, offset: 0 })
    orders.value = response.orders.map(convertToFrontendOrder)
  } catch (error: any) {
    console.error('Failed to load orders:', error)
    ElMessage.error('加载订单列表失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadOrders()
})

// 详情抽屉
const detailVisible = ref(false)
const selectedOrder = ref<Order | null>(null)

// 创建订单对话框
const createOrderDialogVisible = ref(false)

// 打开创建订单对话框
const openCreateOrderDialog = () => {
  createOrderDialogVisible.value = true
}

// 创建订单成功
const handleCreateOrderSuccess = () => {
  ElMessage.success('订单创建成功')
  loadOrders() // 重新加载订单列表
}

// 查看详情
const handleViewDetail = (order: Order) => {
  selectedOrder.value = order
  detailVisible.value = true
}

// 更新订单
const handleOrderUpdate = (updatedOrder: Order) => {
  const index = orders.value.findIndex(o => o.id === updatedOrder.id)
  if (index !== -1) {
    orders.value[index] = updatedOrder
  }
}

// 获取角色名称
const getRoleName = (role: RoleType) => {
  return getOrderRoleName(role)
}

// 获取角色标签类型
const getRoleTagType = (role: RoleType) => {
  const typeMap: Record<RoleType, 'success' | 'warning' | 'info' | 'danger'> = {
    sales: 'success',
    follower: 'warning',
    warehouse: 'info',
    system: 'danger'
  }
  return typeMap[role] || 'info'
}
</script>

<style scoped>
.order-management-page {
  padding: 24px;
  background-color: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 20px 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-left h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
</style>
