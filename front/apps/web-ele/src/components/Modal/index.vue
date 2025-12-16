<template>
  <el-dialog :visible.sync="dialogVisible" title="这是一个Modal">
    <p>这里是Modal内容</p>
    <span slot="footer" class="dialog-footer">
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="dialogVisible = false">确认</el-button>
    </span>
  </el-dialog>
</template>

<script setup>
import { ref } from 'vue';

const dialogVisible = ref(false);

const openModal = () => {
  dialogVisible.value = true;
};
</script>
