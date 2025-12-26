<template>
  <div style="padding: 24px; max-width: 480px">
    <RemoteIndexSearch
      v-model="keyword"
      index="client"              
      placeholder="搜索客户"
      @select="handleSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import RemoteIndexSearch from "#/components/SearchInput/index.vue"   // 👈 路径按你项目改

const keyword = ref("")

const handleSelect = (item: any) => {
  console.log("选中客户：", item)
}
</script>
