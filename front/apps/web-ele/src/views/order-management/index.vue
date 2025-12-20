<template>
  <div class="order-management-page">
    <!-- 顶部角色切换器 -->
    <div class="page-header">
      <div class="header-left">
        <h2>订单管理</h2>
        <el-tag type="info" size="large">共 {{ orders.length }} 个订单</el-tag>
      </div>
      <div class="header-right">
        <el-button type="primary" size="large" @click="openCreateOrderDialog">
          创建订单
        </el-button>
        <div class="role-switcher">
          <span class="role-label">当前角色：</span>
          <el-select v-model="currentRole" size="large" style="width: 200px">
            <el-option label="业务（Sales）" value="sales" />
            <el-option label="跟单（Follower）" value="follower" />
            <el-option label="仓库（Warehouse）" value="warehouse" />
          </el-select>
        </div>
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
import { ref } from 'vue'
import OrderList from './components/OrderList.vue'
import OrderDetail from './components/OrderDetail.vue'
import CreateOrderDialog from './components/CreateOrderDialog.vue'
import { mockOrders } from './mockData'
import type { Order } from './types'

// 当前角色
const currentRole = ref<'sales' | 'follower' | 'warehouse'>('sales')

// 订单数据
const orders = ref<Order[]>(mockOrders)

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
  // TODO: 刷新订单列表或添加新订单
  // 这里可以调用 API 重新获取订单列表
  console.log('订单创建成功，刷新列表')
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

.role-switcher {
  display: flex;
  align-items: center;
  gap: 12px;
}

.role-label {
  font-size: 14px;
  font-weight: 500;
  color: #606266;
}
</style>
