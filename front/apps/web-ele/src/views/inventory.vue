<template>
  <div class="inventory-management">
    <div class="page-header">
      <el-space>
        <el-button type="primary" @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </el-space>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-cards">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="总库存价值" :value="totalInventoryValue" :precision="2">
            <template #prefix>¥</template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="订单库存价值" :value="orderInventoryValue" :precision="2">
            <template #prefix>¥</template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="公共库存价值" :value="publicInventoryValue" :precision="2">
            <template #prefix>¥</template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="库存差价" :value="inventoryDifference" :precision="2" :value-style="differenceStyle">
            <template #prefix>¥</template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>

    <!-- Tab切换 -->
    <el-tabs v-model="activeTab" class="inventory-tabs">
      <!-- 订单库存 -->
      <el-tab-pane label="订单库存" name="order">
        <InventoryTable
          :data="orderInventory"
          :columns="orderColumns"
          title="订单库存"
        />
      </el-tab-pane>

      <!-- 公共库存 -->
      <el-tab-pane label="公共库存" name="public">
        <InventoryTable
          :data="publicInventory"
          :columns="publicColumns"
          title="公共库存"
        />
      </el-tab-pane>

      <!-- 原料库存 -->
      <el-tab-pane label="原料库存" name="material">
        <InventoryTable
          :data="materialInventory"
          :columns="materialColumns"
          title="原料库存"
        />
      </el-tab-pane>

      <!-- 半成品库存 -->
      <el-tab-pane label="半成品库存" name="semifinished">
        <InventoryTable
          :data="semifinishedInventory"
          :columns="semifinishedColumns"
          title="半成品库存"
        />
      </el-tab-pane>

      <!-- 成品库存 -->
      <el-tab-pane label="成品库存" name="finished">
        <InventoryTable
          :data="finishedInventory"
          :columns="finishedColumns"
          title="成品库存"
        />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import InventoryTable from '#/components/InventoryTable/index.vue'

// Tab状态
const activeTab = ref('order')

// 库存数据（从后台获取）
const orderInventory = ref<any[]>([])
const publicInventory = ref<any[]>([])
const materialInventory = ref<any[]>([])
const semifinishedInventory = ref<any[]>([])
const finishedInventory = ref<any[]>([])

// 表格列配置
const orderColumns = [
  { prop: 'orderNo', label: '订单编号', width: 150 },
  { prop: 'clientName', label: '客户名称', width: 150 },
  { prop: 'productName', label: '产品名称', width: 180 },
  { prop: 'quantity', label: '数量', width: 100, align: 'right' as const },
  { prop: 'unit', label: '单位', width: 80 },
  { prop: 'inboundValue', label: '入库价值', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'currentMarketValue', label: '现在市价', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'difference', label: '差价', width: 120, align: 'right' as const, type: 'difference' as const },
  { prop: 'warehouseName', label: '仓库', width: 120 },
  { prop: 'location', label: '货位', width: 120 },
  { prop: 'inboundDate', label: '入库时间', width: 180 }
]

const publicColumns = [
  { prop: 'itemCode', label: '物料编码', width: 150 },
  { prop: 'itemName', label: '物料名称', width: 180 },
  { prop: 'itemType', label: '类型', width: 100 },
  { prop: 'quantity', label: '数量', width: 100, align: 'right' as const },
  { prop: 'unit', label: '单位', width: 80 },
  { prop: 'inboundValue', label: '入库价值', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'currentMarketValue', label: '现在市价', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'difference', label: '差价', width: 120, align: 'right' as const, type: 'difference' as const },
  { prop: 'warehouseName', label: '仓库', width: 120 },
  { prop: 'location', label: '货位', width: 120 },
  { prop: 'inboundDate', label: '入库时间', width: 180 }
]

const materialColumns = [
  { prop: 'materialCode', label: '原料编码', width: 150 },
  { prop: 'materialName', label: '原料名称', width: 180 },
  { prop: 'specification', label: '规格', width: 120 },
  { prop: 'quantity', label: '数量', width: 100, align: 'right' as const },
  { prop: 'unit', label: '单位', width: 80 },
  { prop: 'inboundValue', label: '入库价值', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'currentMarketValue', label: '现在市价', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'difference', label: '差价', width: 120, align: 'right' as const, type: 'difference' as const },
  { prop: 'supplierName', label: '供应商', width: 150 },
  { prop: 'warehouseName', label: '仓库', width: 120 },
  { prop: 'inboundDate', label: '入库时间', width: 180 }
]

const semifinishedColumns = [
  { prop: 'itemCode', label: '半成品编码', width: 150 },
  { prop: 'itemName', label: '半成品名称', width: 180 },
  { prop: 'specification', label: '规格', width: 120 },
  { prop: 'quantity', label: '数量', width: 100, align: 'right' as const },
  { prop: 'unit', label: '单位', width: 80 },
  { prop: 'inboundValue', label: '入库价值', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'currentMarketValue', label: '现在市价', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'difference', label: '差价', width: 120, align: 'right' as const, type: 'difference' as const },
  { prop: 'warehouseName', label: '仓库', width: 120 },
  { prop: 'location', label: '货位', width: 120 },
  { prop: 'inboundDate', label: '入库时间', width: 180 }
]

const finishedColumns = [
  { prop: 'productCode', label: '产品编码', width: 150 },
  { prop: 'productName', label: '产品名称', width: 180 },
  { prop: 'specification', label: '规格', width: 120 },
  { prop: 'quantity', label: '数量', width: 100, align: 'right' as const },
  { prop: 'unit', label: '单位', width: 80 },
  { prop: 'inboundValue', label: '入库价值', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'currentMarketValue', label: '现在市价', width: 120, align: 'right' as const, type: 'currency' as const },
  { prop: 'difference', label: '差价', width: 120, align: 'right' as const, type: 'difference' as const },
  { prop: 'warehouseName', label: '仓库', width: 120 },
  { prop: 'location', label: '货位', width: 120 },
  { prop: 'inboundDate', label: '入库时间', width: 180 }
]

// 计算总库存价值
const totalInventoryValue = computed(() => {
  const orderValue = orderInventory.value.reduce((sum, item) => sum + item.currentMarketValue, 0)
  const publicValue = publicInventory.value.reduce((sum, item) => sum + item.currentMarketValue, 0)
  return orderValue + publicValue
})

const orderInventoryValue = computed(() => {
  return orderInventory.value.reduce((sum, item) => sum + item.currentMarketValue, 0)
})

const publicInventoryValue = computed(() => {
  return publicInventory.value.reduce((sum, item) => sum + item.currentMarketValue, 0)
})

// 库存差价(现在市价 - 入库价值)
const inventoryDifference = computed(() => {
  const allItems = [
    ...orderInventory.value,
    ...publicInventory.value
  ]
  const currentValue = allItems.reduce((sum, item) => sum + item.currentMarketValue, 0)
  const inboundValue = allItems.reduce((sum, item) => sum + item.inboundValue, 0)
  return currentValue - inboundValue
})

const differenceStyle = computed(() => {
  return inventoryDifference.value >= 0
    ? { color: '#67C23A' }
    : { color: '#F56C6C' }
})

// 刷新数据
const handleRefresh = async () => {
  // TODO: 从后台API获取库存数据
  // 目前暂时清空数据
  orderInventory.value = []
  publicInventory.value = []
  materialInventory.value = []
  semifinishedInventory.value = []
  finishedInventory.value = []
}
</script>

<style scoped lang="scss">
.inventory-management {
  padding: 24px;
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.stats-cards {
  margin-bottom: 20px;

  :deep(.el-card) {
    border-radius: 8px;

    .el-card__body {
      padding: 20px;
    }
  }

  :deep(.el-statistic) {
    text-align: center;

    .el-statistic__head {
      color: #909399;
      font-size: 14px;
      margin-bottom: 8px;
    }

    .el-statistic__content {
      font-size: 28px;
      font-weight: 600;
    }
  }
}

.inventory-tabs {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

  :deep(.el-tabs__header) {
    margin-bottom: 20px;
  }

  :deep(.el-tabs__item) {
    font-size: 16px;
    font-weight: 500;
  }
}
</style>
