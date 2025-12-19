import { requestClient } from '#/api/request';

/**
 * 材料配置
 */
export interface MaterialConfig {
  material_id: number;
  ratio: number; // 占比 (0-1)
}

/**
 * 工艺配置
 */
export interface ProcessConfig {
  process_id: number;
  quantity?: number;
}

/**
 * 产品价格响应
 */
export interface ProductPriceResponse {
  current_price: number;
  historical_high: number;
  historical_low: number;
}

/**
 * 创建产品请求
 */
export interface CreateProductRequest {
  name: string;
  materials: MaterialConfig[];
  processes: ProcessConfig[];
}

/**
 * 更新产品请求
 */
export interface UpdateProductRequest {
  name?: string;
  status?: string;
  materials?: MaterialConfig[];
  processes?: ProcessConfig[];
}

/**
 * 产品响应
 */
export interface ProductResponse {
  id: number;
  name: string;
  status: string;
  materials: MaterialConfig[];
  processes: ProcessConfig[];
  created_at: string;
  updated_at: string;
}

/**
 * 获取产品价格
 */
export async function getProductPriceApi(productId: number): Promise<ProductPriceResponse> {
  const response = await requestClient.get<ProductPriceResponse>(`/product/${productId}/price`);
  return response.data;
}

/**
 * 创建产品
 */
export async function createProductApi(data: CreateProductRequest): Promise<ProductResponse> {
  const response = await requestClient.post<ProductResponse>('/product', data);
  return response.data;
}

/**
 * 获取产品详情
 */
export async function getProductApi(productId: number): Promise<ProductResponse> {
  const response = await requestClient.get<ProductResponse>(`/product/${productId}`);
  return response.data;
}

/**
 * 更新产品
 */
export async function updateProductApi(productId: number, data: UpdateProductRequest): Promise<void> {
  await requestClient.post(`/product/${productId}`, data);
}

/**
 * 删除产品
 */
export async function deleteProductApi(productId: number): Promise<void> {
  await requestClient.delete(`/product/${productId}`);
}
