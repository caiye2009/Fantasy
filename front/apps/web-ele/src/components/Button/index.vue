<template>
  <el-button type="primary" @click="showModal">点击弹出Modal</el-button>
</template>

<script setup>
import { defineEmits } from 'vue';

const emit = defineEmits();

const showModal = () => {
  emit('openModal');
};
</script>
