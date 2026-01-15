<template>
  <el-drawer
    :model-value="visible"
    @update:model-value="$emit('update:visible', $event)"
    :title="drawerTitle"
    :size="size"
    :before-close="handleClose"
    destroy-on-close
    class="entity-drawer"
  >
    <div v-loading="loading" class="drawer-content">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
        class="entity-form"
      >
        <el-form-item
          v-for="field in fields"
          :key="field.key"
          :label="field.label"
          :prop="field.key"
          :required="field.required"
        >
          <el-input
            v-if="field.type === 'text'"
            v-model="formData[field.key]"
            :placeholder="field.placeholder || `请输入${field.label}`"
            :disabled="field.disabled || false"
            clearable
            @input="handleFieldChange"
          />

          <el-input
            v-else-if="field.type === 'textarea'"
            v-model="formData[field.key]"
            type="textarea"
            :rows="field.rows || 4"
            :placeholder="field.placeholder || `请输入${field.label}`"
            :disabled="field.disabled || false"
            @input="handleFieldChange"
          />

          <el-input-number
            v-else-if="field.type === 'number'"
            v-model="formData[field.key]"
            :min="field.min"
            :max="field.max"
            :precision="field.precision || 0"
            :step="field.step || 1"
            :disabled="field.disabled || false"
            style="width: 100%"
            @change="handleFieldChange"
          />

          <el-select
            v-else-if="field.type === 'select'"
            v-model="formData[field.key]"
            :placeholder="field.placeholder || `请选择${field.label}`"
            :disabled="field.disabled || false"
            :multiple="field.multiple || false"
            :filterable="field.filterable !== false"
            clearable
            style="width: 100%"
            @change="handleFieldChange"
          >
            <el-option
              v-for="option in field.options"
              :key="option.value"
              :label="option.label"
              :value="option.value"
            />
          </el-select>

          <el-date-picker
            v-else-if="field.type === 'date'"
            v-model="formData[field.key]"
            type="date"
            :placeholder="field.placeholder || `请选择${field.label}`"
            :disabled="field.disabled || false"
            style="width: 100%"
            @change="handleFieldChange"
          />

          <el-switch
            v-else-if="field.type === 'switch'"
            v-model="formData[field.key]"
            :disabled="field.disabled || false"
            @change="handleFieldChange"
          />

          <slot
            v-else-if="field.type === 'custom'"
            :name="`field-${field.key}`"
            :form-data="formData"
            :field="field"
            :on-change="handleFieldChange"
          />
        </el-form-item>
      </el-form>
    </div>

    <template #footer>
      <div class="drawer-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ mode === 'create' ? '创建' : '保存' }}
        </el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import type { EntityDrawerConfig, EntityDrawerMode } from './types'

interface Props {
  visible: boolean
  config: EntityDrawerConfig
  entity?: any
  mode?: EntityDrawerMode
  drawerId: string  // 用于区分不同抽屉的唯一标识
}

const props = withDefaults(defineProps<Props>(), {
  mode: 'create'
})

const emit = defineEmits(['update:visible', 'submit', 'close'])

const formRef = ref<FormInstance>()
const formData = ref<Record<string, any>>({})
const loading = ref(false)
const submitting = ref(false)

// 缓存key = drawerId + mode
const cacheKey = computed(() => `drawer_cache_${props.drawerId}_${props.mode}`)

const drawerTitle = computed(() => {
  if (props.mode === 'view') return `查看${props.config.title}`
  if (props.mode === 'edit') return `编辑${props.config.title}`
  return `新增${props.config.title}`
})

const size = computed(() => props.config.size || '50%')

const fields = computed(() => {
  if (props.mode === 'view') {
    return props.config.fields.map(field => ({ ...field, disabled: true }))
  }
  return props.config.fields
})

const rules = computed<FormRules>(() => {
  const formRules: FormRules = {}
  props.config.fields.forEach(field => {
    if (field.required) {
      formRules[field.key] = [{
        required: true,
        message: field.requiredMessage || `请输入${field.label}`,
        trigger: field.type === 'select' ? 'change' : 'blur'
      }]
    }
    if (field.validator) {
      if (!formRules[field.key]) formRules[field.key] = []
      formRules[field.key].push({ validator: field.validator, trigger: 'blur' })
    }
  })
  return formRules
})

const initFormData = () => {
  const data: Record<string, any> = {}
  
  // 1. 设置默认值
  props.config.fields.forEach(field => {
    if (field.defaultValue !== undefined) {
      data[field.key] = field.defaultValue
    } else {
      if (field.type === 'number') data[field.key] = 0
      else if (field.type === 'switch') data[field.key] = false
      else if (field.type === 'checkbox') data[field.key] = []
      else data[field.key] = ''
    }
  })

  // 2. 编辑/查看模式：使用实体数据
  if (props.entity && (props.mode === 'edit' || props.mode === 'view')) {
    props.config.fields.forEach(field => {
      if (props.entity[field.key] !== undefined) {
        data[field.key] = props.entity[field.key]
      }
    })
  }

  // 3. 新增模式：尝试从缓存恢复
  if (props.mode === 'create') {
    const cached = loadCache()
    if (cached) {
      Object.assign(data, cached)
    }
  }

  formData.value = data
}

const saveCache = () => {
  try {
    sessionStorage.setItem(cacheKey.value, JSON.stringify(formData.value))
  } catch (e) {
    console.error('缓存失败:', e)
  }
}

const loadCache = () => {
  try {
    const cached = sessionStorage.getItem(cacheKey.value)
    return cached ? JSON.parse(cached) : null
  } catch (e) {
    return null
  }
}

const clearCache = () => {
  try {
    sessionStorage.removeItem(cacheKey.value)
  } catch (e) {}
}

const handleFieldChange = () => {
  if (props.mode === 'create') {
    saveCache()
  }
}

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    submitting.value = true

    if (props.config.onSubmit) {
      await props.config.onSubmit(formData.value, props.mode)
    }

    emit('submit', formData.value, props.mode)
    ElMessage.success(props.mode === 'create' ? '创建成功' : '保存成功')
    
    if (props.mode === 'create') {
      clearCache()
    }
    
    handleClose()
  } catch (error: any) {
    console.error('提交失败:', error)
    if (error.response?.data?.error) {
      ElMessage.error(error.response.data.error)
    }
  } finally {
    submitting.value = false
  }
}

watch(() => props.visible, (val) => {
  if (val) {
    initFormData()
    setTimeout(() => formRef.value?.clearValidate(), 0)
  }
})

defineExpose({ formRef, formData, clearCache })
</script>

<style scoped lang="scss">
.entity-drawer {
  :deep(.el-drawer__header) {
    margin-bottom: 20px;
    padding-bottom: 20px;
    border-bottom: 1px solid #ebeef5;
  }

  :deep(.el-drawer__body) {
    padding: 20px;
    display: flex;
    flex-direction: column;
  }

  .drawer-content {
    flex: 1;
    overflow-y: auto;
    padding-right: 10px;

    &::-webkit-scrollbar { width: 6px; }
    &::-webkit-scrollbar-thumb {
      background-color: #dcdfe6;
      border-radius: 3px;
      &:hover { background-color: #c0c4cc; }
    }
  }

  .entity-form .el-form-item { margin-bottom: 22px; }

  .drawer-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding-top: 20px;
    border-top: 1px solid #ebeef5;
  }
}
</style>