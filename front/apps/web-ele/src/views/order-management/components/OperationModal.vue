<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`操作面板 - ${getRoleName(currentRole)}`"
    width="600px"
    :close-on-click-modal="false"
  >
    <div class="operation-modal">
      <!-- 业务（Sales）操作 -->
      <div v-if="currentRole === 'sales'" class="operation-form">
        <el-form :model="salesForm" label-width="100px" size="default">
          <el-form-item label="追加数量">
            <el-input-number
              v-model="salesForm.additionalQuantity"
              :min="1"
              :max="10000"
              :step="10"
              style="width: 200px"
            />
            <span class="form-tip">当前：{{ order?.requiredQuantity }} 件</span>
          </el-form-item>
          <el-form-item label="备注">
            <el-input
              v-model="salesForm.remark"
              type="textarea"
              :rows="3"
              placeholder="请输入追加原因"
            />
          </el-form-item>
        </el-form>
      </div>

      <!-- 跟单（Follower）操作 -->
      <div v-if="currentRole === 'follower'" class="operation-form">
        <el-tabs v-model="followerActiveTab">
          <!-- 胚布投入 -->
          <el-tab-pane label="胚布投入" name="fabric">
            <el-form :model="followerFabricForm" label-width="100px" size="default">
              <el-form-item label="当前投入">
                <span class="value-display">{{ getFabricInputItem()?.completedQuantity || 0 }} / {{ getFabricInputItem()?.targetQuantity || 0 }}</span>
              </el-form-item>
              <el-form-item label="新增数量">
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
                  placeholder="可选"
                />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 生产进度 -->
          <el-tab-pane label="生产进度" name="production">
            <el-form :model="followerProductionForm" label-width="100px" size="default">
              <el-form-item label="当前生产">
                <span class="value-display">{{ getProductionItem()?.completedQuantity || 0 }} / {{ getProductionItem()?.targetQuantity || 0 }}</span>
              </el-form-item>
              <el-form-item label="新增数量">
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
                  placeholder="可选"
                />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 回修进度 -->
          <el-tab-pane v-if="getReworkItem()?.exists" label="回修进度" name="rework">
            <el-form :model="followerReworkForm" label-width="100px" size="default">
              <el-form-item label="当前回修">
                <span class="value-display">{{ getReworkItem()?.completedQuantity || 0 }} / {{ getReworkItem()?.targetQuantity || 0 }}</span>
              </el-form-item>
              <el-form-item label="新增数量">
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
                  placeholder="可选"
                />
              </el-form-item>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 仓库（Warehouse）操作 -->
      <div v-if="currentRole === 'warehouse'" class="operation-form">
        <el-tabs v-model="warehouseActiveTab">
          <!-- 验货进度 -->
          <el-tab-pane label="验货进度" name="check">
            <el-form :model="warehouseCheckForm" label-width="100px" size="default">
              <el-form-item label="当前验收">
                <span class="value-display">{{ getWarehouseCheckItem()?.completedQuantity || 0 }} / {{ getWarehouseCheckItem()?.targetQuantity || 0 }}</span>
              </el-form-item>
              <el-form-item label="新增数量">
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
                  placeholder="可选"
                />
              </el-form-item>
            </el-form>
          </el-tab-pane>

          <!-- 次品录入 -->
          <el-tab-pane label="次品录入" name="defect">
            <el-form :model="warehouseDefectForm" label-width="100px" size="default">
              <el-form-item label="当前次品">
                <span class="value-display danger">{{ getReworkItem()?.targetQuantity || 0 }} 件</span>
              </el-form-item>
              <el-form-item label="新增次品">
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
                  placeholder="请详细描述次品原因"
                />
              </el-form-item>
              <el-alert
                type="warning"
                :closable="false"
                show-icon
                style="margin-top: 8px"
              >
                录入次品后将自动生成回修进度
              </el-alert>
            </el-form>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          确认提交
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
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

const submitting = ref(false)

// Dialog 显示状态
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

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
const getFabricInputItem = () => props.order?.progressItems.find(item => item.type === 'fabric_input')
const getProductionItem = () => props.order?.progressItems.find(item => item.type === 'production')
const getWarehouseCheckItem = () => props.order?.progressItems.find(item => item.type === 'warehouse_check')
const getReworkItem = () => props.order?.progressItems.find(item => item.type === 'rework')

// 获取角色名称
const getRoleName = (role: RoleType) => {
  const nameMap = {
    sales: '业务',
    follower: '跟单',
    warehouse: '仓库',
    system: '系统'
  }
  return nameMap[role] || role
}

// 提交
const handleSubmit = () => {
  submitting.value = true

  setTimeout(() => {
    if (props.currentRole === 'sales') {
      if (!salesForm.additionalQuantity) {
        ElMessage.warning('请输入追加数量')
        submitting.value = false
        return
      }
      ElMessage.success(`追加需求数量 ${salesForm.additionalQuantity} 件（演示）`)
      resetForms()
    } else if (props.currentRole === 'follower') {
      if (followerActiveTab.value === 'fabric') {
        if (!followerFabricForm.quantity) {
          ElMessage.warning('请输入投入数量')
          submitting.value = false
          return
        }
        ElMessage.success(`更新胚布投入 ${followerFabricForm.quantity} 件（演示）`)
      } else if (followerActiveTab.value === 'production') {
        if (!followerProductionForm.quantity) {
          ElMessage.warning('请输入生产数量')
          submitting.value = false
          return
        }
        ElMessage.success(`更新生产进度 ${followerProductionForm.quantity} 件（演示）`)
      } else if (followerActiveTab.value === 'rework') {
        if (!followerReworkForm.quantity) {
          ElMessage.warning('请输入回修数量')
          submitting.value = false
          return
        }
        ElMessage.success(`更新回修进度 ${followerReworkForm.quantity} 件（演示）`)
      }
      resetForms()
    } else if (props.currentRole === 'warehouse') {
      if (warehouseActiveTab.value === 'check') {
        if (!warehouseCheckForm.quantity) {
          ElMessage.warning('请输入验收数量')
          submitting.value = false
          return
        }
        ElMessage.success(`更新验收进度 ${warehouseCheckForm.quantity} 件（演示）`)
      } else if (warehouseActiveTab.value === 'defect') {
        if (!warehouseDefectForm.defectQuantity) {
          ElMessage.warning('请输入次品数量')
          submitting.value = false
          return
        }
        ElMessage.warning(`录入次品 ${warehouseDefectForm.defectQuantity} 件（演示）`)
      }
      resetForms()
    }

    submitting.value = false
    dialogVisible.value = false
  }, 500)
}

// 取消
const handleCancel = () => {
  resetForms()
  dialogVisible.value = false
}

// 重置表单
const resetForms = () => {
  salesForm.additionalQuantity = 0
  salesForm.remark = ''
  followerFabricForm.quantity = 0
  followerFabricForm.remark = ''
  followerProductionForm.quantity = 0
  followerProductionForm.remark = ''
  followerReworkForm.quantity = 0
  followerReworkForm.remark = ''
  warehouseCheckForm.quantity = 0
  warehouseCheckForm.remark = ''
  warehouseDefectForm.defectQuantity = 0
  warehouseDefectForm.remark = ''
}
</script>

<style scoped>
.operation-modal {
  min-height: 200px;
}

.operation-form {
  padding: 8px 0;
}

.form-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}

.value-display {
  font-size: 14px;
  font-weight: 600;
  color: #409EFF;
  font-family: monospace;
}

.value-display.danger {
  color: #F56C6C;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-tabs__item) {
  font-size: 14px;
}

:deep(.el-form-item__label) {
  font-weight: 500;
  font-size: 13px;
}

:deep(.el-textarea__inner) {
  font-size: 13px;
}
</style>
