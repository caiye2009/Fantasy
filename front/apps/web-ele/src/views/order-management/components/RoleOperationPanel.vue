<template>
  <div class="role-operation-panel">
    <!-- 业务（Sales）操作 -->
    <div v-if="currentRole === 'sales'" class="operation-section">
      <div class="operation-title">
        <el-icon><ShoppingCart /></el-icon>
        <span>业务操作</span>
      </div>
      <div class="operation-content">
        <el-form :model="salesForm" label-width="120px">
          <el-form-item label="追加需求数量">
            <el-input-number
              v-model="salesForm.additionalQuantity"
              :min="1"
              :max="10000"
              :step="10"
              style="width: 200px"
            />
            <span class="form-tip">当前需求：{{ order.requiredQuantity }} 件</span>
          </el-form-item>
          <el-form-item label="备注">
            <el-input
              v-model="salesForm.remark"
              type="textarea"
              :rows="3"
              placeholder="请输入追加原因或备注"
              style="width: 400px"
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSalesOperation" :loading="submitting">
              确认追加
            </el-button>
            <el-button @click="resetSalesForm">重置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>

    <!-- 跟单（Follower）操作 -->
    <div v-if="currentRole === 'follower'" class="operation-section">
      <div class="operation-title">
        <el-icon><Monitor /></el-icon>
        <span>跟单操作</span>
      </div>
      <div class="operation-content">
        <el-tabs v-model="followerActiveTab">
          <!-- 胚布投入 -->
          <el-tab-pane label="胚布投入" name="fabric">
            <el-form :model="followerFabricForm" label-width="140px">
              <el-form-item label="当前投入数量">
                <span class="current-value">{{ getFabricInputItem()?.completedQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="目标数量">
                <span class="current-value">{{ getFabricInputItem()?.targetQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="新增投入数量">
                <el-input-number
                  v-model="followerFabricForm.quantity"
                  :min="1"
                  :max="10000"
                  :step="10"
                  style="width: 200px"
                />
              </el-form-item>
              <el-form-item label="备注">
                <el-input
                  v-model="followerFabricForm.remark"
                  type="textarea"
                  :rows="2"
                  placeholder="请输入备注"
                  style="width: 400px"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleFabricUpdate" :loading="submitting">
                  更新胚布投入
                </el-button>
                <el-button @click="resetFollowerFabricForm">重置</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 生产进度 -->
          <el-tab-pane label="生产进度" name="production">
            <el-form :model="followerProductionForm" label-width="140px">
              <el-form-item label="当前生产数量">
                <span class="current-value">{{ getProductionItem()?.completedQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="目标数量">
                <span class="current-value">{{ getProductionItem()?.targetQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="新增生产数量">
                <el-input-number
                  v-model="followerProductionForm.quantity"
                  :min="1"
                  :max="10000"
                  :step="10"
                  style="width: 200px"
                />
              </el-form-item>
              <el-form-item label="备注">
                <el-input
                  v-model="followerProductionForm.remark"
                  type="textarea"
                  :rows="2"
                  placeholder="请输入备注"
                  style="width: 400px"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleProductionUpdate" :loading="submitting">
                  更新生产进度
                </el-button>
                <el-button @click="resetFollowerProductionForm">重置</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 回修进度 (仅在存在回修时显示) -->
          <el-tab-pane v-if="getReworkItem()?.exists" label="回修进度" name="rework">
            <el-form :model="followerReworkForm" label-width="140px">
              <el-form-item label="当前回修数量">
                <span class="current-value">{{ getReworkItem()?.completedQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="待回修数量">
                <span class="current-value">{{ getReworkItem()?.targetQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="新增回修数量">
                <el-input-number
                  v-model="followerReworkForm.quantity"
                  :min="1"
                  :max="getReworkItem()?.targetQuantity || 1000"
                  :step="10"
                  style="width: 200px"
                />
              </el-form-item>
              <el-form-item label="备注">
                <el-input
                  v-model="followerReworkForm.remark"
                  type="textarea"
                  :rows="2"
                  placeholder="请输入备注"
                  style="width: 400px"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleReworkUpdate" :loading="submitting">
                  更新回修进度
                </el-button>
                <el-button @click="resetFollowerReworkForm">重置</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>

    <!-- 仓库（Warehouse）操作 -->
    <div v-if="currentRole === 'warehouse'" class="operation-section">
      <div class="operation-title">
        <el-icon><Box /></el-icon>
        <span>仓库操作</span>
      </div>
      <div class="operation-content">
        <el-tabs v-model="warehouseActiveTab">
          <!-- 验货进度 -->
          <el-tab-pane label="验货进度" name="check">
            <el-form :model="warehouseCheckForm" label-width="140px">
              <el-form-item label="当前验收数量">
                <span class="current-value">{{ getWarehouseCheckItem()?.completedQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="目标数量">
                <span class="current-value">{{ getWarehouseCheckItem()?.targetQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="新增验收数量">
                <el-input-number
                  v-model="warehouseCheckForm.quantity"
                  :min="1"
                  :max="10000"
                  :step="10"
                  style="width: 200px"
                />
              </el-form-item>
              <el-form-item label="备注">
                <el-input
                  v-model="warehouseCheckForm.remark"
                  type="textarea"
                  :rows="2"
                  placeholder="请输入备注"
                  style="width: 400px"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleWarehouseCheckUpdate" :loading="submitting">
                  更新验收进度
                </el-button>
                <el-button @click="resetWarehouseCheckForm">重置</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 次品录入 -->
          <el-tab-pane label="次品录入" name="defect">
            <el-form :model="warehouseDefectForm" label-width="140px">
              <el-form-item label="当前次品数量">
                <span class="current-value danger">{{ getReworkItem()?.targetQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="新增次品数量">
                <el-input-number
                  v-model="warehouseDefectForm.defectQuantity"
                  :min="1"
                  :max="1000"
                  style="width: 200px"
                />
              </el-form-item>
              <el-form-item label="次品描述">
                <el-input
                  v-model="warehouseDefectForm.remark"
                  type="textarea"
                  :rows="3"
                  placeholder="请详细描述次品原因和情况"
                  style="width: 400px"
                />
              </el-form-item>
              <el-form-item>
                <el-alert
                  type="warning"
                  :closable="false"
                  style="margin-bottom: 12px"
                >
                  <template #title>
                    录入次品后将自动生成回修进度，目标数量为次品数量
                  </template>
                </el-alert>
              </el-form-item>
              <el-form-item>
                <el-button type="danger" @click="handleDefectAdd" :loading="submitting">
                  录入次品
                </el-button>
                <el-button @click="resetWarehouseDefectForm">重置</el-button>
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>

    <!-- 无权限提示 -->
    <div v-if="!['sales', 'follower', 'warehouse'].includes(currentRole)" class="no-permission">
      <el-empty description="当前角色无操作权限" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { ShoppingCart, Monitor, Box } from '@element-plus/icons-vue'
import type { Order, RoleType } from '../types'

interface Props {
  order: Order
  currentRole: RoleType
}

interface Emits {
  (e: 'update', order: Order): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const submitting = ref(false)

// 业务表单
const salesForm = reactive({
  additionalQuantity: 0,
  remark: ''
})

// 跟单表单
const followerActiveTab = ref('fabric')
const followerFabricForm = reactive({
  quantity: 0,
  remark: ''
})
const followerProductionForm = reactive({
  quantity: 0,
  remark: ''
})
const followerReworkForm = reactive({
  quantity: 0,
  remark: ''
})

// 仓库表单
const warehouseActiveTab = ref('check')
const warehouseCheckForm = reactive({
  quantity: 0,
  remark: ''
})
const warehouseDefectForm = reactive({
  defectQuantity: 0,
  remark: ''
})

// 获取进度项
const getFabricInputItem = () => props.order.progressItems.find(item => item.type === 'fabric_input')
const getProductionItem = () => props.order.progressItems.find(item => item.type === 'production')
const getWarehouseCheckItem = () => props.order.progressItems.find(item => item.type === 'warehouse_check')
const getReworkItem = () => props.order.progressItems.find(item => item.type === 'rework')

// 业务操作
const handleSalesOperation = () => {
  if (!salesForm.additionalQuantity) {
    ElMessage.warning('请输入追加数量')
    return
  }
  ElMessage.success(`追加需求数量 ${salesForm.additionalQuantity} 件（仅演示，未实际修改数据）`)
  resetSalesForm()
}

const resetSalesForm = () => {
  salesForm.additionalQuantity = 0
  salesForm.remark = ''
}

// 跟单操作 - 胚布
const handleFabricUpdate = () => {
  if (!followerFabricForm.quantity) {
    ElMessage.warning('请输入投入数量')
    return
  }
  ElMessage.success(`更新胚布投入 ${followerFabricForm.quantity} 件（仅演示，未实际修改数据）`)
  resetFollowerFabricForm()
}

const resetFollowerFabricForm = () => {
  followerFabricForm.quantity = 0
  followerFabricForm.remark = ''
}

// 跟单操作 - 生产
const handleProductionUpdate = () => {
  if (!followerProductionForm.quantity) {
    ElMessage.warning('请输入生产数量')
    return
  }
  ElMessage.success(`更新生产进度 ${followerProductionForm.quantity} 件（仅演示，未实际修改数据）`)
  resetFollowerProductionForm()
}

const resetFollowerProductionForm = () => {
  followerProductionForm.quantity = 0
  followerProductionForm.remark = ''
}

// 跟单操作 - 回修
const handleReworkUpdate = () => {
  if (!followerReworkForm.quantity) {
    ElMessage.warning('请输入回修数量')
    return
  }
  ElMessage.success(`更新回修进度 ${followerReworkForm.quantity} 件（仅演示，未实际修改数据）`)
  resetFollowerReworkForm()
}

const resetFollowerReworkForm = () => {
  followerReworkForm.quantity = 0
  followerReworkForm.remark = ''
}

// 仓库操作 - 验货
const handleWarehouseCheckUpdate = () => {
  if (!warehouseCheckForm.quantity) {
    ElMessage.warning('请输入验收数量')
    return
  }
  ElMessage.success(`更新验收进度 ${warehouseCheckForm.quantity} 件（仅演示，未实际修改数据）`)
  resetWarehouseCheckForm()
}

const resetWarehouseCheckForm = () => {
  warehouseCheckForm.quantity = 0
  warehouseCheckForm.remark = ''
}

// 仓库操作 - 次品
const handleDefectAdd = () => {
  if (!warehouseDefectForm.defectQuantity) {
    ElMessage.warning('请输入次品数量')
    return
  }
  ElMessage.warning(`录入次品 ${warehouseDefectForm.defectQuantity} 件，将自动生成回修进度（仅演示，未实际修改数据）`)
  resetWarehouseDefectForm()
}

const resetWarehouseDefectForm = () => {
  warehouseDefectForm.defectQuantity = 0
  warehouseDefectForm.remark = ''
}
</script>

<style scoped>
.role-operation-panel {
  background: white;
  border-radius: 8px;
  padding: 20px;
}

.operation-section {
  width: 100%;
}

.operation-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #409EFF;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid #409EFF;
}

.operation-content {
  padding: 12px 0;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}

.current-value {
  font-size: 18px;
  font-weight: 600;
  color: #409EFF;
}

.current-value.danger {
  color: #F56C6C;
}

.no-permission {
  padding: 40px 0;
  text-align: center;
}

:deep(.el-tabs__item) {
  font-size: 15px;
  font-weight: 500;
}

:deep(.el-form-item__label) {
  font-weight: 500;
}
</style>
