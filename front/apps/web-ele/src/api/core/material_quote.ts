import { requestClient } from '../request';

export interface MaterialQuote {
  id: number;
  quoteId: string;
  materialCode: string;
  supplierCode: string;
  quotePrice: number;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface MaterialQuoteListResponse {
  total: number;
  list: MaterialQuote[];
}

export async function createMaterialQuote(data: {
  quoteId: string;
  materialCode: string;
  supplierCode: string;
  quotePrice: number;
  createdBy: number;
}) {
  const response = await requestClient.post<MaterialQuote>('/material_quotes', data);
  return response;
}

export async function getMaterialQuotes(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<MaterialQuoteListResponse>('/material_quotes', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getMaterialQuote(id: number) {
  const response = await requestClient.get<MaterialQuote>(`/material_quotes/${id}`);
  return response;
}

export async function updateMaterialQuote(
  id: number,
  data: {
    quoteId?: string;
    materialCode?: string;
    supplierCode?: string;
    quotePrice?: number;
  }
) {
  const response = await requestClient.post<MaterialQuote>(`/material_quotes/${id}`, data);
  return response;
}

export async function deleteMaterialQuote(id: number) {
  const response = await requestClient.delete(`/material_quotes/${id}`);
  return response;
}

// 获取原料报价历史 (按material_id筛选)
export interface MaterialPriceHistory {
  supplier_id: number;
  supplier_name: string;
  price: number;
  quoted_at: string;
}

export async function getMaterialPriceHistoryApi(materialId: number): Promise<MaterialPriceHistory[]> {
  // TODO: 实现获取原料报价历史的逻辑
  // 可能需要后端添加按material_id筛选的接口
  console.warn('getMaterialPriceHistoryApi not implemented yet, materialId:', materialId);
  return [];
}
