<template>
  <div class="material-management">
    <div class="page-header">
      <h2>原料管理</h2>
      <el-button type="primary" @click="openAddDialog">
        <el-icon><Plus /></el-icon>
        新增原料
      </el-button>
    </div>

    <!-- 筛选器 -->
    <div class="filter-bar">
      <el-input
        v-model="searchKeyword"
        placeholder="搜索原料名称或编号"
        clearable
        style="width: 300px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select v-model="filterCategory" placeholder="全部分类" clearable style="width: 150px">
        <el-option label="全部" value="" />
        <el-option label="胚布" value="胚布" />
        <el-option label="染料" value="染料" />
        <el-option label="助剂" value="助剂" />
      </el-select>
    </div>

    <!-- 原料列表 -->
    <div class="material-list">
      <el-table
        :data="filteredMaterials"
        stripe
        @row-click="openDetail"
      >
        <el-table-column prop="code" label="原料编号" width="140" />
        <el-table-column prop="name" label="原料名称" width="200" />
        <el-table-column prop="category" label="分类" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getCategoryTagType(row.category)" size="small">
              {{ row.category }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="currentPrice" label="当前价格" width="150" align="right">
          <template #default="{ row }">
            <span class="price">¥{{ row.currentPrice.toFixed(2) }}/{{ row.unit }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="updatedBy" label="最新更新人" width="120" align="center" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '在用' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" size="small" link @click.stop="openDetail(row)">
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 新增原料对话框 -->
    <AddMaterialDialog
      v-model:visible="addDialogVisible"
      @success="handleAddSuccess"
    />

    <!-- 原料详情抽屉 -->
    <MaterialDetail
      v-model:visible="detailVisible"
      :material="selectedMaterial"
      @update-material="handleMaterialUpdate"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Plus, Search } from '@element-plus/icons-vue'
import { mockMaterials } from './mockData'
import AddMaterialDialog from './components/AddMaterialDialog.vue'
import MaterialDetail from './components/MaterialDetail.vue'
import type { Material } from './types'

// 数据
const materials = ref<Material[]>(mockMaterials)
const searchKeyword = ref('')
const filterCategory = ref('')

// 对话框和抽屉
const addDialogVisible = ref(false)
const detailVisible = ref(false)
const selectedMaterial = ref<Material | null>(null)

// 过滤后的原料列表
const filteredMaterials = computed(() => {
  return materials.value.filter(material => {
    const matchKeyword = !searchKeyword.value ||
      material.name.includes(searchKeyword.value) ||
      material.code.includes(searchKeyword.value)
    const matchCategory = !filterCategory.value ||
      material.category === filterCategory.value
    return matchKeyword && matchCategory
  })
})

// 获取分类标签类型
const getCategoryTagType = (category: string) => {
  const typeMap: Record<string, any> = {
    '胚布': 'primary',
    '染料': 'warning',
    '助剂': 'success'
  }
  return typeMap[category] || 'info'
}

// 打开新增对话框
const openAddDialog = () => {
  addDialogVisible.value = true
}

// 新增成功
const handleAddSuccess = (newMaterial: Material) => {
  materials.value.unshift(newMaterial)
}

// 打开详情
const openDetail = (material: Material) => {
  selectedMaterial.value = material
  detailVisible.value = true
}

// 更新原料信息
const handleMaterialUpdate = (updatedMaterial: Material) => {
  // 更新列表中的原料数据
  const index = materials.value.findIndex(m => m.id === updatedMaterial.id)
  if (index !== -1) {
    materials.value[index] = updatedMaterial
  }
  // 更新当前选中的原料
  selectedMaterial.value = updatedMaterial
}
</script>

<style scoped>
.material-management {
  padding: 24px;
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.page-header h2 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  padding: 16px 24px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.material-list {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.price {
  font-weight: 600;
  color: #F56C6C;
}

.low-stock {
  color: #F56C6C;
  font-weight: 600;
}

:deep(.el-table__row) {
  cursor: pointer;
}

:deep(.el-table__row:hover) {
  background-color: #f5f7fa !important;
}
</style>
