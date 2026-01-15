<script setup lang="ts">
defineProps<{ modelValue: string }>()
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()
</script>

<template>
  <input 
    :value="modelValue" 
    @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    placeholder="搜索..."
  />
</template>