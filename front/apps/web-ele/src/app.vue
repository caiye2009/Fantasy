<script lang="ts" setup>
import { useElementPlusDesignTokens } from '@vben/hooks';
import { ElConfigProvider } from 'element-plus';
import { elementLocale } from '#/locales';

defineOptions({ name: 'App' });

useElementPlusDesignTokens();
</script>

<template>
  <ElConfigProvider :locale="elementLocale">
    <RouterView />
  </ElConfigProvider>
</template>

<style>
html, body {
  overflow: hidden;
  height: 100%;
  margin: 0;
  padding: 0;
}
</style>