<template>
  <el-dialog
    :model-value="visible"
    title="新增工艺"
    width="500px"
    @close="handleClose"
  >
    <el-form :model="form" label-width="100px">
      <el-form-item label="工艺名称" required>
        <el-input v-model="form.name" placeholder="请输入工艺名称" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="3"
          placeholder="请输入描述"
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
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

interface Props {
  visible: boolean
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'success'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const form = ref({
  name: '',
  description: ''
})

const submitting = ref(false)

const handleClose = () => {
  emit('update:visible', false)
  form.value = { name: '', description: '' }
}

const handleSubmit = async () => {
  if (!form.value.name) {
    ElMessage.warning('请输入工艺名称')
    return
  }

  submitting.value = true
  try {
    // TODO: 调用API创建工艺
    await new Promise(resolve => setTimeout(resolve, 500))

    ElMessage.success('创建成功')
    emit('success')
    handleClose()
  } catch (error) {
    console.error('创建失败:', error)
    ElMessage.error('创建失败')
  } finally {
    submitting.value = false
  }
}
</script>
