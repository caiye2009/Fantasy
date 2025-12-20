<template>
  <el-dialog
    :model-value="visible"
    title="新增供应商报价"
    width="600px"
    @close="handleClose"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="供应商名称" prop="supplierName">
        <el-input v-model="form.supplierName" placeholder="请输入供应商名称" />
      </el-form-item>
      <el-form-item label="报价" prop="price">
        <el-input-number
          v-model="form.price"
          :min="0"
          :precision="2"
          :step="0.01"
          style="width: 100%"
        />
        <span v-if="material" style="margin-left: 8px; color: #909399;">
          /{{ material.unit }}
        </span>
      </el-form-item>
      <el-form-item label="报价人" prop="quotedBy">
        <el-input v-model="form.quotedBy" placeholder="请输入报价人" />
      </el-form-item>
      <el-form-item label="备注" prop="remark">
        <el-input
          v-model="form.remark"
          type="textarea"
          :rows="3"
          placeholder="请输入备注信息（可选）"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="handleSubmit" :loading="submitting">
        确定
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import type { Material, SupplierQuote } from '../types'

interface Props {
  visible: boolean
  material: Material | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success', quote: SupplierQuote): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = ref({
  supplierName: '',
  price: 0,
  quotedBy: '',
  remark: ''
})

const rules: FormRules = {
  supplierName: [{ required: true, message: '请输入供应商名称', trigger: 'blur' }],
  price: [{ required: true, message: '请输入报价', trigger: 'blur' }],
  quotedBy: [{ required: true, message: '请输入报价人', trigger: 'blur' }]
}

watch(() => props.visible, (val) => {
  if (!val) {
    formRef.value?.resetFields()
  }
})

const handleClose = () => {
  emit('update:visible', false)
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate()
  if (!valid || !props.material) return

  submitting.value = true
  try {
    // 模拟提交
    await new Promise(resolve => setTimeout(resolve, 500))

    const now = new Date().toLocaleString('zh-CN')
    const newQuote: SupplierQuote = {
      id: `sq-${Date.now()}`,
      materialId: props.material.id,
      supplierName: form.value.supplierName,
      price: form.value.price,
      quotedBy: form.value.quotedBy,
      quotedAt: now,
      remark: form.value.remark || undefined
    }

    emit('success', newQuote)
    emit('update:visible', false)
    ElMessage.success('报价添加成功')
  } catch (error) {
    ElMessage.error('报价添加失败')
  } finally {
    submitting.value = false
  }
}
</script>
