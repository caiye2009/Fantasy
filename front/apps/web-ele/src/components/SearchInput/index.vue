<template>
  <el-autocomplete
    v-model="innerValue"
    :fetch-suggestions="query"
    :placeholder="placeholder"
    :loading="loading"
    :trigger-on-focus="false"
    clearable
    @select="onSelect"
  >
    <template #default="{ item }">
      <div class="search-item" :class="{ 'search-item--empty': item.__empty }">
        <template v-if="item.__empty">
          {{ item.value }}
        </template>
        <template v-else>
          <span v-if="item.customNo || item.id" class="search-item__code">
            {{ item.customNo || item.id }}
          </span>
          <span class="search-item__name">
            {{ item.customName || item.name || item.value }}
          </span>
        </template>
      </div>
    </template>
  </el-autocomplete>
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
  searchFields?: string[]
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
        // 如果传入了 searchFields 就使用，否则不传（使用后端配置的所有字段）
        ...(props.searchFields ? { SearchFields: props.searchFields } : {}),
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

    // 兼容 data 包装和直接返回两种格式
    const data = res?.data ?? res
    const items = data?.items ?? []

    if (items.length === 0) {
      cb([{ value: "尝试其他搜索", __empty: true }])
      loading.value = false
      return
    }

    cb(
      items.map((x: any) => {
        let value = ''
        if (props.index === 'client') {
          value = x.customName || x.customNameEn || x.customNo
        } else if (props.index === 'product') {
          value = x.name || x.id
        } else {
          value = x.name || x.value || x.id
        }

        return {
          value,
          ...x,
        }
      })
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

<style scoped>
.search-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  height: 32px;
  line-height: 32px;
  overflow: hidden;
}

.search-item__code {
  flex-shrink: 0;
  font-weight: 500;
  color: #606266;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 40%;
}

.search-item__name {
  flex: 1;
  text-align: right;
  color: #909399;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.search-item--empty {
  color: #c0c4cc;
  cursor: default;
  pointer-events: none;
  justify-content: center;
  font-size: 14px;
}

/* 禁用 element-plus autocomplete 的 hover 特效 */
:deep(.el-autocomplete-suggestion__list .search-item--empty:hover) {
  background-color: transparent !important;
}
</style>
