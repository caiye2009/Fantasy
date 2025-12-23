<template>
  <el-dialog
    :model-value="visible"
    title="添加报价"
    width="500px"
    @close="handleClose"
  >
    <el-form :model="quoteForm" label-width="100px">
      <el-form-item label="工艺名称">
        <el-input :value="process?.name" disabled />
      </el-form-item>
      <el-form-item label="供应商" required>
        <el-select
          v-model="quoteForm.supplier_id"
          placeholder="请选择供应商"
          filterable
          style="width: 100%"
        >
          <el-option
            v-for="supplier in suppliers"
            :key="supplier.id"
            :label="supplier.name"
            :value="supplier.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="报价金额" required>
        <el-input-number
          v-model="quoteForm.price"
          :min="0"
          :precision="2"
          :step="0.01"
          style="width: 100%"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="quoting">
        提交报价
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { Process } from '../types'
import { quoteProcessPriceApi } from '#/api/core/pricing'
import { getSupplierListApi, type Supplier } from '#/api/core/supplier'

interface Props {
  visible: boolean
  process: Process | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const quoteForm = ref({
  target_id: 0,
  supplier_id: undefined as number | undefined,
  price: 0,
})

const suppliers = ref<Supplier[]>([])
const quoting = ref(false)

// 加载供应商列表
const loadSuppliers = async () => {
  try {
    const res = await getSupplierListApi({ limit: 1000, offset: 0 })
    suppliers.value = res.suppliers || []
  } catch (error) {
    console.error('加载供应商列表失败:', error)
  }
}

// 监听弹窗打开
watch(() => props.visible, async (visible) => {
  if (visible) {
    await loadSuppliers()
    if (props.process) {
      quoteForm.value.target_id = Number(props.process.id)
    }
  }
})

const handleClose = () => {
  emit('update:visible', false)
  quoteForm.value = {
    target_id: 0,
    supplier_id: undefined,
    price: 0,
  }
}

const handleSubmit = async () => {
  if (!quoteForm.value.supplier_id) {
    ElMessage.warning('请选择供应商')
    return
  }

  if (quoteForm.value.price <= 0) {
    ElMessage.warning('请输入有效的报价金额')
    return
  }

  quoting.value = true
  try {
    await quoteProcessPriceApi(quoteForm.value)
    ElMessage.success('报价成功')
    emit('success')
    handleClose()
  } catch (error: any) {
    console.error('报价失败:', error)
    ElMessage.error(error.response?.data?.error || '报价失败')
  } finally {
    quoting.value = false
  }
}
</script>
