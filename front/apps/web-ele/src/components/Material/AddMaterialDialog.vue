<template>
  <el-dialog
    :model-value="visible"
    title="新增原料"
    width="600px"
    @close="handleClose"
  >
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="原料编码" prop="code">
        <el-input v-model="form.code" placeholder="请输入原料编码（唯一）" />
      </el-form-item>
      <el-form-item label="原料名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入原料名称" />
      </el-form-item>
      <el-form-item label="分类" prop="category">
        <el-select v-model="form.category" placeholder="请选择分类" style="width: 100%">
          <el-option label="胚布" value="胚布" />
          <el-option label="染料" value="染料" />
          <el-option label="助剂" value="助剂" />
        </el-select>
      </el-form-item>
      <el-form-item label="单位" prop="unit">
        <el-select v-model="form.unit" placeholder="请选择单位" style="width: 100%">
          <el-option label="米（m）" value="m" />
          <el-option label="千克（kg）" value="kg" />
        </el-select>
      </el-form-item>
      <el-form-item label="当前价格" prop="currentPrice">
        <el-input-number
          v-model="form.currentPrice"
          :min="0"
          :precision="2"
          :step="0.01"
          style="width: 100%"
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
import type { Material } from '../types'

interface Props {
  visible: boolean
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success', material: Material): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()
const submitting = ref(false)

const form = ref({
  code: '',
  name: '',
  category: '',
  unit: '' as 'kg' | 'm' | '',
  currentPrice: 0
})

const rules: FormRules = {
  code: [{ required: true, message: '请输入原料编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入原料名称', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  unit: [{ required: true, message: '请选择单位', trigger: 'change' }],
  currentPrice: [{ required: true, message: '请输入当前价格', trigger: 'blur' }]
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
  if (!valid) return

  submitting.value = true
  try {
    // 模拟提交
    await new Promise(resolve => setTimeout(resolve, 500))

    const now = new Date().toLocaleString('zh-CN')
    const newMaterial: Material = {
      id: `mat-${Date.now()}`,
      code: form.value.code,
      name: form.value.name,
      category: form.value.category,
      unit: form.value.unit as 'kg' | 'm',
      currentPrice: form.value.currentPrice,
      status: 'active',
      updatedBy: '当前用户', // TODO: 实际项目中应从用户状态获取
      createdAt: now,
      updatedAt: now
    }

    emit('success', newMaterial)
    emit('update:visible', false)
    ElMessage.success('添加成功')
  } catch (error) {
    ElMessage.error('添加失败')
  } finally {
    submitting.value = false
  }
}
</script>
