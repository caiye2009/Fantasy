<script setup lang="ts">
import { ref, watch } from 'vue'
import { fetchAggOptions } from '#/api/core/search'
import type { AggOption } from '#/api/core/search'

const props = defineProps<{
  filterKey: string
  modelValue: any
  index: string
  size?: number
  labelFormatter?: (key: any) => string
  placeholder?: string
}>()

const emit = defineEmits<{ 'update:modelValue': [value: any] }>()

const options = ref<AggOption[]>([])
const loading = ref(false)

// 组件挂载时或首次聚焦时获取选项
const fetchOptions = async () => {
  if (options.value.length > 0) return // 已加载过就不重复加载
  
  loading.value = true
  try {
    options.value = await fetchAggOptions({
      index: props.index,
      field: props.filterKey,
      size: props.size || 20,
      labelFormatter: props.labelFormatter
    })
  } catch (error) {
    console.error(`获取选项失败: ${props.filterKey}`, error)
    options.value = []
  } finally {
    loading.value = false
  }
}

// 聚焦时加载
const handleFocus = () => {
  fetchOptions()
}
</script>

<template>
  <el-select
    :model-value="modelValue"
    @update:model-value="emit('update:modelValue', $event)"
    :loading="loading"
    :placeholder="filterKey"
    clearable
    filterable
    style="width: 240px"
    @focus="handleFocus"
  >
    <el-option
      v-for="opt in options"
      :key="opt.value"
      :label="opt.label"
      :value="opt.value"
    />
  </el-select>
</template>