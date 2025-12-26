<template>
  <el-dialog
    :model-value="visible"
    title="添加报价"
    width="500px"
    @close="handleClose"
  >
    <el-form :model="form" label-width="100px">
      <el-form-item label="供应商" required>
        <el-input v-model="form.supplierName" placeholder="请输入供应商名称" />
      </el-form-item>
      <el-form-item label="报价金额" required>
        <el-input-number
          v-model="form.price"
          :min="0"
          :precision="2"
          :step="0.01"
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item label="备注">
        <el-input
          v-model="form.remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注"
        />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">
        确定
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { SupplierQuote } from './types'

interface Props {
  visible: boolean
  materialId?: number
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success', quote: SupplierQuote): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const submitting = ref(false)
const form = ref({
  supplierName: '',
  price: 0,
  remark: ''
})

const handleClose = () => {
  emit('update:visible', false)
  // 重置表单
  form.value = {
    supplierName: '',
    price: 0,
    remark: ''
  }
}

const handleSubmit = async () => {
  if (!form.value.supplierName) {
    ElMessage.warning('请输入供应商名称')
    return
  }

  if (form.value.price <= 0) {
    ElMessage.warning('请输入有效的报价金额')
    return
  }

  submitting.value = true
  try {
    // 创建报价对象
    const quote: SupplierQuote = {
      id: Date.now(),
      materialId: props.materialId || 0,
      supplierName: form.value.supplierName,
      price: form.value.price,
      quotedAt: new Date().toISOString(),
      remark: form.value.remark
    }

    // 触发成功事件
    emit('success', quote)
    ElMessage.success('添加报价成功')
    handleClose()
  } catch (error) {
    console.error('添加报价失败:', error)
    ElMessage.error('添加报价失败')
  } finally {
    submitting.value = false
  }
}
</script>
