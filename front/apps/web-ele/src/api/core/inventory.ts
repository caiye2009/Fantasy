import { requestClient } from '#/api/request';

/**
 * 库存响应
 */
export interface InventoryResponse {
  id: number;
  productId: number;
  category: string;
  batchId: string;
  quantity: number;
  unit: string;
  unitCost: number;
  totalCost: number;
  remark: string;
  createdAt: string;
  updatedAt: string;
}

/**
 * 库存列表响应
 */
export interface InventoryListResponse {
  total: number;
  inventories: InventoryResponse[];
}

/**
 * 创建库存请求
 */
export interface CreateInventoryRequest {
  productId: number;
  category: string;
  batchId: string;
  quantity: number;
  unit: string;
  unitCost: number;
  remark?: string;
}

/**
 * 更新库存请求
 */
export interface UpdateInventoryRequest {
  id: number;
  productId: number;
  category: string;
  batchId: string;
  quantity: number;
  unit: string;
  unitCost: number;
  remark?: string;
}

/**
 * 根据产品ID获取库存列表
 */
export async function getInventoriesByProductIdApi(productId: number): Promise<InventoryListResponse> {
  const response = await requestClient.get<InventoryListResponse>('/inventory/product', {
    params: { productId }
  });
  return response.data;
}

/**
 * 获取库存详情
 */
export async function getInventoryApi(id: number): Promise<InventoryResponse> {
  const response = await requestClient.get<InventoryResponse>(`/inventory/${id}`);
  return response.data;
}

/**
 * 根据批次ID获取库存
 */
export async function getInventoryByBatchIdApi(batchId: string): Promise<InventoryResponse> {
  const response = await requestClient.get<InventoryResponse>('/inventory/batch', {
    params: { batchId }
  });
  return response.data;
}

/**
 * 获取库存列表
 */
export async function getInventoriesApi(limit = 10, offset = 0): Promise<InventoryListResponse> {
  const response = await requestClient.get<InventoryListResponse>('/inventory/list', {
    params: { limit, offset }
  });
  return response.data;
}

/**
 * 根据类别获取库存列表
 */
export async function getInventoriesByCategoryApi(category: string, limit = 10, offset = 0): Promise<InventoryListResponse> {
  const response = await requestClient.get<InventoryListResponse>('/inventory/category', {
    params: { category, limit, offset }
  });
  return response.data;
}

/**
 * 创建库存
 */
export async function createInventoryApi(data: CreateInventoryRequest): Promise<InventoryResponse> {
  const response = await requestClient.post<InventoryResponse>('/inventory', data);
  return response.data;
}

/**
 * 更新库存
 */
export async function updateInventoryApi(data: UpdateInventoryRequest): Promise<InventoryResponse> {
  const response = await requestClient.put<InventoryResponse>('/inventory', data);
  return response.data;
}

/**
 * 删除库存
 */
export async function deleteInventoryApi(id: number): Promise<void> {
  await requestClient.delete(`/inventory/${id}`);
}

/**
 * 扣减库存
 */
export async function deductInventoryApi(id: number, quantity: number): Promise<InventoryResponse> {
  const response = await requestClient.post<InventoryResponse>('/inventory/deduct', { id, quantity });
  return response.data;
}

/**
 * 增加库存
 */
export async function addInventoryApi(id: number, quantity: number): Promise<InventoryResponse> {
  const response = await requestClient.post<InventoryResponse>('/inventory/add', { id, quantity });
  return response.data;
}
