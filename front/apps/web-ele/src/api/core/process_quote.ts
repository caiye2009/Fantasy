import { requestClient } from '../request';

export interface ProcessQuote {
  id: number;
  quoteId: string;
  processCode: string;
  supplierCode: string;
  quotePrice: number;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProcessQuoteListResponse {
  total: number;
  list: ProcessQuote[];
}

export async function createProcessQuote(data: {
  quoteId: string;
  processCode: string;
  supplierCode: string;
  quotePrice: number;
  createdBy: number;
}) {
  const response = await requestClient.post<ProcessQuote>('/process_quotes', data);
  return response;
}

export async function getProcessQuotes(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<ProcessQuoteListResponse>('/process_quotes', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getProcessQuote(id: number) {
  const response = await requestClient.get<ProcessQuote>(`/process_quotes/${id}`);
  return response;
}

export async function updateProcessQuote(
  id: number,
  data: {
    quoteId?: string;
    processCode?: string;
    supplierCode?: string;
    quotePrice?: number;
  }
) {
  const response = await requestClient.put<ProcessQuote>(`/process_quotes/${id}`, data);
  return response;
}

export async function deleteProcessQuote(id: number) {
  const response = await requestClient.delete(`/process_quotes/${id}`);
  return response;
}

// 获取工序报价历史 (按target_id筛选)
export interface ProcessPriceHistory {
  supplier_id: number;
  supplier_name: string;
  price: number;
  quoted_at: string;
}

export async function getProcessPriceHistoryApi(processId: number): Promise<ProcessPriceHistory[]> {
  // TODO: 实现获取工序报价历史的逻辑
  // 可能需要后端添加按process_id筛选的接口
  console.warn('getProcessPriceHistoryApi not implemented yet, returning empty array');
  return [];
}

// 提交工序报价
export async function quoteProcessPriceApi(data: {
  target_id: number;
  supplier_id: number;
  price: number;
}) {
  const requestData = {
    quoteId: `PQ-${Date.now()}`,
    processCode: String(data.target_id),
    supplierCode: String(data.supplier_id),
    quotePrice: data.price,
    createdBy: 1, // TODO: 从用户上下文获取
  };
  const response = await createProcessQuote(requestData);
  return response;
}
