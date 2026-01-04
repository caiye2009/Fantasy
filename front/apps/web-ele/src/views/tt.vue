<template>
  <div style="padding: 24px; max-width: 480px">
    <RemoteIndexSearch
      v-model="keyword"
      index="client"              
      placeholder="搜索客户"
      @select="handleSelect"
    />
    <Button @openModal="visible = true" />

    <el-dialog v-model="visible" title="新增客户">
      这里是内容
    </el-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref } from "vue"
import RemoteIndexSearch from "#/components/SearchInput/index.vue"   // 👈 路径按你项目改
import Button from "#/components/Button/index.vue"

const keyword = ref("")

const handleSelect = (item: any) => {
  console.log("选中客户：", item)
}
</script>
