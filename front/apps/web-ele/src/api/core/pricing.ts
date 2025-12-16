import { requestClient } from '#/api/request';

/**
 * 报价请求
 */
export interface QuoteRequest {
  target_id: number;
  supplier_id: number;
  price: number;
}

/**
 * 价格数据
 */
export interface PriceData {
  price: number;
  supplier_id: number;
  supplier_name: string;
  quoted_at: string;
}

/**
 * 材料报价
 */
export async function quoteMaterialPriceApi(data: QuoteRequest): Promise<any> {
  const response = await requestClient.post('/material/price/quote', data);
  return response.data;
}

/**
 * 获取材料价格（最低价和最高价）
 */
export async function getMaterialPriceApi(materialId: number): Promise<{ min: number; max: number }> {
  const response = await requestClient.get<{ min: number; max: number }>(
    `/material/${materialId}/price`
  );
  return response.data;
}

/**
 * 获取材料价格历史
 */
export async function getMaterialPriceHistoryApi(materialId: number): Promise<PriceData[]> {
  const response = await requestClient.get<PriceData[]>(`/material/${materialId}/price/history`);
  return response.data;
}

/**
 * 工艺报价
 */
export async function quoteProcessPriceApi(data: QuoteRequest): Promise<any> {
  const response = await requestClient.post('/process/price/quote', data);
  return response.data;
}

/**
 * 获取工艺价格（最低价和最高价）
 */
export async function getProcessPriceApi(processId: number): Promise<{ min: number; max: number }> {
  const response = await requestClient.get<{ min: number; max: number }>(
    `/process/${processId}/price`
  );
  return response.data;
}

/**
 * 获取工艺价格历史
 */
export async function getProcessPriceHistoryApi(processId: number): Promise<PriceData[]> {
  const response = await requestClient.get<PriceData[]>(`/process/${processId}/price/history`);
  return response.data;
}
