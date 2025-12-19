<template>
  <el-dialog
    v-model="visible"
    :title="title"
    width="400px"
  >
    <slot>
      <p>这里是 Modal 内容</p>
    </slot>

    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" @click="handleConfirm">确认</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  modelValue: boolean
  title?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm'): void
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val: boolean) => emit('update:modelValue', val),
})

const handleConfirm = () => {
  emit('confirm')
  visible.value = false
}
</script>
