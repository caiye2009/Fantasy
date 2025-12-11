import { baseRequestClient } from '#/api/request'
import type { ESRequest, ESResponse } from '#/components/DataTable/types'

/**
 * 批量操作响应
 */
interface BulkOperationResponse {
  success: boolean
  successCount: number
  failedCount: number
  errors?: Array<{
    id: string
    reason: string
  }>
}

/**
 * 导出选项
 */
interface ExportOptions {
  format?: 'csv' | 'xlsx' | 'json'
  filename?: string
}

/**
 * Elasticsearch 服务
 */
export const elasticsearchService = {
  /**
   * Elasticsearch 搜索接口
   */
  async search(params: ESRequest): Promise<ESResponse> {
    const response = await baseRequestClient.post<ESResponse>('/search', params)
    return response.data
  },

  /**
   * 批量删除
   */
  async bulkDelete(
    ids: string[],
    index: string
  ): Promise<BulkOperationResponse> {
    const response = await baseRequestClient.post<BulkOperationResponse>('/bulk-delete', {
      ids,
      index,
    })
    return response.data
  },

  /**
   * 批量更新
   */
  async bulkUpdate(
    updates: Array<{ id: string; data: Record<string, any> }>,
    index: string
  ): Promise<BulkOperationResponse> {
    const response = await baseRequestClient.post<BulkOperationResponse>('/bulk-update', {
      updates,
      index,
    })
    return response.data
  },

  /**
   * 导出数据
   */
  async export(
    params: ESRequest,
    options: ExportOptions = {}
  ): Promise<Blob> {
    const { format = 'xlsx', filename } = options
    const response = await baseRequestClient.post<Blob>(
      '/export',
      { ...params, format, filename },
      {
        responseType: 'blob',
      }
    )
    return response.data
  },

  /**
   * 获取单条数据
   */
  async getById(id: string, index: string): Promise<any> {
    const response = await baseRequestClient.get(`/${index}/${id}`)
    return response.data
  },

  /**
   * 更新单条数据
   */
  async update(
    id: string,
    index: string,
    data: Record<string, any>
  ): Promise<{ success: boolean; data: any }> {
    const response = await baseRequestClient.put(`/${index}/${id}`, data)
    return response.data
  },

  /**
   * 创建单条数据
   */
  async create(
    index: string,
    data: Record<string, any>
  ): Promise<{ success: boolean; id: string; data: any }> {
    const response = await baseRequestClient.post(`/${index}`, data)
    return response.data
  },

  /**
   * 删除单条数据
   */
  async delete(
    id: string,
    index: string
  ): Promise<{ success: boolean }> {
    const response = await baseRequestClient.delete(`/${index}/${id}`)
    return response.data
  },

  /**
   * 获取聚合数据（统计）
   */
  async aggregate(params: {
    index: string
    aggregations: Record<string, any>
    query?: any
  }): Promise<any> {
    const response = await baseRequestClient.post('/aggregate', params)
    return response.data
  },
}