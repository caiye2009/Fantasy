<template>
  <div style="padding: 40px">
    <h2>Vue 3 Test Page</h2>

    <p>count: {{ count }}</p>
    <el-button @click="count++">count +1</el-button>

    <el-divider />

    <el-button type="primary" @click="openModal">
      打开 Modal
    </el-button>

    <BaseModal
      v-model="dialogVisible"
      title="测试 Modal"
      @confirm="onConfirm"
    >
      <p>这是从 test.vue 传进来的内容</p>
    </BaseModal>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import BaseModal from '#/components/Modal/index.vue'

const count = ref(0)

const dialogVisible = ref(false)

const openModal = () => {
  dialogVisible.value = true
}

const onConfirm = () => {
  console.log('点击了确认')
}
</script>
