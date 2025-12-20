<template>
  <el-drawer
    :model-value="visible"
    title="原料详情"
    size="60%"
    @close="handleClose"
  >
    <div v-if="material" class="material-detail">
      <!-- 基本信息 -->
      <el-card class="detail-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
          </div>
        </template>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="原料编码">
            {{ material.code }}
          </el-descriptions-item>
          <el-descriptions-item label="原料名称">
            {{ material.name }}
          </el-descriptions-item>
          <el-descriptions-item label="分类">
            <el-tag :type="getCategoryTagType(material.category)" size="small">
              {{ material.category }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="单位">
            {{ material.unit }}
          </el-descriptions-item>
          <el-descriptions-item label="当前价格">
            <span class="price">¥{{ material.currentPrice.toFixed(2) }}/{{ material.unit }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="material.status === 'active' ? 'success' : 'info'" size="small">
              {{ material.status === 'active' ? '在用' : '停用' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 供应商报价 -->
      <el-card class="detail-card" shadow="never">
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
import { ref, watch } from 'vue'
import type { Material, SupplierQuote } from '../types'
import { mockSupplierQuotes } from '../mockData'
import AddQuoteDialog from './AddQuoteDialog.vue'

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

// 监听 material 变化，加载报价数据
watch(() => props.material, (newMaterial) => {
  if (newMaterial) {
    supplierQuotes.value = mockSupplierQuotes[newMaterial.id] || []
  } else {
    supplierQuotes.value = []
  }
}, { immediate: true })

const handleClose = () => {
  emit('update:visible', false)
}

const getCategoryTagType = (category: string) => {
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
const handleAddQuoteSuccess = (newQuote: SupplierQuote) => {
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
