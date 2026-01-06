import { requestClient } from '../request';

export interface SearchFilter {
  field: string;
  operator: string;
  value: any;
}

export interface SearchAggregation {
  field: string;
  type: string;
}

export interface SearchRequest {
  index: string;
  query?: string;
  filters?: SearchFilter[];
  aggregations?: SearchAggregation[];
  from?: number;
  size?: number;
  sortField?: string;
  sortOrder?: string;
}

export interface SearchResponse {
  total: number;
  items: any[];  // 后端返回的是 items，不是 hits
  aggregations?: Record<string, any>;
}

export interface SearchIndicesResponse {
  indices: string[];
}

export async function search(data: SearchRequest) {
  const response = await requestClient.post<SearchResponse>('/search', data);
  return response;
}

export async function getSearchIndices() {
  const response = await requestClient.get<SearchIndicesResponse>('/search/indices');
  return response;
}

// List query functions for all modules using search API

export interface ListResponse<T> {
  total: number;
  list: T[];
}

// Client list queries
export async function getClientListApi(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'client',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.items || [],
  };
}

// Material list queries
export async function getMaterialListApi(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'material',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.items || [],
  };
}

// Order list queries
export async function getOrderListApi(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'order',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.items || [],
  };
}

// Product list queries
export async function getProductListApi(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'product',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.items || [],
  };
}

// Process list queries
export async function getProcessListApi(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'process',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.items || [],
  };
}

// Supplier list queries
export async function getSupplierListApi(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'supplier',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.items || [],
  };
}
