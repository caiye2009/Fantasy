<template>
  <el-autocomplete
    v-model="innerValue"
    :fetch-suggestions="query"
    :placeholder="placeholder"
    :loading="loading"
    :trigger-on-focus="false"
    clearable
    @select="onSelect"
  />
</template>

<script setup lang="ts">
import { ref, watch } from "vue"
import { requestClient } from "#/api/request"   // 👈 按你项目路径调整

interface Item {
  value: string
  __empty?: boolean
  [k: string]: any
}

const props = defineProps<{
  modelValue: string
  index: string
  placeholder?: string
}>()

const emit = defineEmits<{
  (e: "update:modelValue", v: string): void
  (e: "select", item: Item): void
}>()

const innerValue = ref(props.modelValue)
watch(() => props.modelValue, v => (innerValue.value = v))
watch(innerValue, v => emit("update:modelValue", v))

const loading = ref(false)

const canSearch = (text: string) => text.trim().length >= 2

let timer: any
const debounce = (fn: Function, delay = 350) => {
  return (...args: any[]) => {
    clearTimeout(timer)
    timer = setTimeout(() => fn(...args), delay)
  }
}

// ⭐ 使用 requestClient.post
const query = debounce(async (keyword: string, cb: (x: Item[]) => void) => {
  if (!canSearch(keyword)) {
    cb([])
    return
  }

  loading.value = true

  try {
    const res = await requestClient.post(
      "/search",
      {
        index: props.index,
        query: keyword,
        fields: ["customName","customNo"],
        pagination: {
          offset: 0,
          size: 20,
        },
      },
      {
        // ⛔️ 不弹 toast
        errorMessageMode: "none",
      }
    )

    const items = res?.items ?? []

    if (items.length === 0) {
      cb([{ value: "尝试其他搜索", __empty: true }])
      loading.value = false
      return
    }

    cb(
      items.map((x: any) => ({
        value: x.customName || x.customNameEn || x.customNo,
        ...x,
      }))
    )
  } catch (e) {
    // 静默失败，只给兜底提示
    cb([{ value: "尝试其他搜索", __empty: true }])
  } finally {
    loading.value = false
  }
}, 350)

const onSelect = (item: Item) => {
  if (item.__empty) return
  emit("select", item)
}
</script>
