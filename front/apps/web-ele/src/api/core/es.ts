// 向后兼容的 ES API，实际使用 search API
export * from './search';
import { search as searchApi, type SearchResponse } from './search';

// 导出 elasticsearchService 作为向后兼容
export const elasticsearchService = {
  search: async (params: any): Promise<any> => {
    // 转换参数格式
    const searchParams: any = {
      index: params.index,
      query: params.query,
      filters: params.filters,
      from: params.pagination?.offset || 0,
      size: params.pagination?.size || 20,
      sortField: params.sort?.[0]?.field,
      sortOrder: params.sort?.[0]?.order,
    };

    // 处理聚合请求
    if (params.aggRequests) {
      searchParams.aggregations = Object.entries(params.aggRequests).map(([_key, value]: [string, any]) => ({
        field: value.field,
        type: value.type,
      }));
    }

    const result = await searchApi(searchParams);
    return {
      ...result,
      items: result.hits,
    };
  },

  update: async (_id: string, index: string, data: any) => {
    // TODO: 实现更新逻辑,可能需要后端提供相应接口
    console.warn(`elasticsearchService.update not fully implemented for index: ${index}, _id: ${_id}`, data);
    return { success: true };
  },

  create: async (index: string, data: any) => {
    // TODO: 实现创建逻辑,可能需要后端提供相应接口
    console.warn(`elasticsearchService.create not fully implemented for index: ${index}`, data);
    return { success: true, id: Date.now() };
  },

  bulkDelete: async (ids: string[], index: string) => {
    // TODO: 实现批量删除逻辑,可能需要后端提供相应接口
    console.warn(`elasticsearchService.bulkDelete not fully implemented for index: ${index}`, ids);
    return {
      success: true,
      successCount: ids.length,
      failedCount: 0,
      errors: [],
    };
  },
};
