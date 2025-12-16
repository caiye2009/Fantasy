import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { elasticsearchService } from '#/api/core/es'
import type { ESResponse } from '#/components/DataTable/types'

export function useDataTable(entityType: string, pageSize: number = 20) {
  const router = useRouter()
  const route = useRoute()

  // 加载状态
  const loading = ref(false)          // 滚动加载
  const searchLoading = ref(false)    // 搜索/筛选加载（用于显示Table遮罩层）

  // 缓存相关（不删除任何数据）
  const cache = ref<Map<number, any[]>>(new Map())
  const currentPage = ref(1)
  const totalCount = ref(0)

  // 搜索 / 筛选 / 排序
  const query = ref<string>('')
  const filters = ref<Record<string, any>>({})
  const sort = ref<Array<{ field: string; order: string }>>([])

  // selection
  const selectedRows = ref<any[]>([])
  const selectedCount = computed(() => selectedRows.value.length)

  // 渲染到界面的总数据
  const tableData = computed(() => {
    const all: any[] = []
    const sortedPages = Array.from(cache.value.keys()).sort((a, b) => a - b)
    for (const p of sortedPages) {
      const d = cache.value.get(p)
      if (d) all.push(...d)
    }
    return all
  })

  const hasMore = computed(() => {
    const keys = Array.from(cache.value.keys())
    if (keys.length === 0) return false
    const maxPage = Math.max(...keys)
    const last = cache.value.get(maxPage)
    return last && last.length === pageSize
  })

  const hasPrevious = computed(() => {
    const keys = Array.from(cache.value.keys())
    if (keys.length === 0) return false
    const minPage = Math.min(...keys)
    return minPage > 1
  })

  /** URL → 内存 */
  const syncFromURL = () => {
    if (route.query.q) query.value = route.query.q as string
    if (route.query.filters) {
      try {
        filters.value = JSON.parse(route.query.filters as string)
      } catch {
        filters.value = {}
      }
    }
    if (route.query.sort) {
      try {
        sort.value = JSON.parse(route.query.sort as string)
      } catch {
        sort.value = []
      }
    }
  }

  /** 内存 → URL */
  const syncToURL = () => {
    const params: Record<string, string> = {}
    if (query.value) params.q = query.value
    if (Object.keys(filters.value).length > 0) params.filters = JSON.stringify(filters.value)
    if (sort.value.length > 0) params.sort = JSON.stringify(sort.value)

    router.replace({ query: params })
  }

  /** 构建请求体（POST body） */
  const buildSearchRequest = (page: number) => ({
    entityType,
    query: query.value || undefined,
    filters: Object.keys(filters.value).length > 0 ? filters.value : undefined,
    sort: sort.value.length > 0 ? sort.value : undefined,
    pagination: {
      offset: (page - 1) * pageSize,
      size: pageSize,
    },
  })

  /** 加载页（如果已缓存直接返回，否则请求） */
  const loadPage = async (page: number, isSearch = false): Promise<any[] | null> => {
    // 已缓存
    if (cache.value.has(page)) {
      return cache.value.get(page) || null
    }

    // 判断加载状态
    if (isSearch) searchLoading.value = true
    else loading.value = true

    try {
      const requestBody = buildSearchRequest(page)
      const response: ESResponse = await elasticsearchService.search(requestBody)

      totalCount.value = response.total || 0

      const data = response.items?.map((item) => ({
        _id: item.id || item._id,
        ...item,
      })) || []

      cache.value.set(page, data)
      return data
    } catch (e: any) {
      ElMessage.error(e.message || '加载失败')
      return null
    } finally {
      if (isSearch) searchLoading.value = false
      else loading.value = false
    }
  }

  /** 清空缓存 */
  const clearAllCache = () => {
    cache.value.clear()
    currentPage.value = 1
    selectedRows.value = []
  }

  /** 初始化 */
  const initialize = async () => {
    syncFromURL()
    clearAllCache()
    await loadPage(1, true)
  }

  /** 滚动加载下一页 */
  const slideWindowDown = async () => {
    if (loading.value) return
    const nextPage = currentPage.value + 1
    if (!cache.value.has(nextPage)) {
      currentPage.value = nextPage
      await loadPage(nextPage)
    }
  }

  /** 滚动加载上一页 */
  const slideWindowUp = async () => {
    if (loading.value) return
    const prevPage = currentPage.value - 1
    if (prevPage > 0 && !cache.value.has(prevPage)) {
      currentPage.value = prevPage
      await loadPage(prevPage)
    }
  }

  /** 重新加载（刷新） */
  const reload = async () => {
    clearAllCache()
    await loadPage(1, true)
  }

  /** 监听搜索条件 */
  watch([query, filters, sort], () => {
    clearAllCache()
    syncToURL()
    loadPage(1, true)
  }, { deep: true })

  /** 滚动事件（用于表格滚动） */
  const handleScroll = (scrollTop: number, scrollHeight: number, clientHeight: number) => {
    if (loading.value || tableData.value.length === 0) return

    const itemHeight = scrollHeight / tableData.value.length
    const visibleStart = Math.floor(scrollTop / itemHeight)
    const visibleEnd = Math.floor((scrollTop + clientHeight) / itemHeight)

    const visibleIndex = Math.floor((visibleStart + visibleEnd) / 2)
    const newPage = Math.floor(visibleIndex / pageSize) + 1

    if (newPage !== currentPage.value && newPage > 0) {
      currentPage.value = newPage
    }

    // 预加载下一页
    if (visibleEnd >= tableData.value.length - 5) {
      slideWindowDown()
    }
    // 预加载上一页
    if (visibleStart <= 5) {
      slideWindowUp()
    }
  }

  return {
    loading,
    searchLoading,
    tableData,
    query,
    filters,
    sort,
    selectedRows,
    selectedCount,
    hasMore,
    hasPrevious,
    totalCount,
    currentPage,

    initialize,
    handleScroll,
    loadPage,
    slideWindowDown,
    slideWindowUp,
    reload,
  }
}
