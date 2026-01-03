import { requestClient } from '../request';

export interface ProductProcess {
  id: number;
  productCode: string;
  processCode: string;
  processSeq: number;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProductProcessListResponse {
  total: number;
  list: ProductProcess[];
}

export async function createProductProcess(data: {
  productCode: string;
  processCode: string;
  processSeq: number;
  createdBy: number;
}) {
  const response = await requestClient.post<ProductProcess>('/product_processs', data);
  return response;
}

export async function getProductProcesses(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<ProductProcessListResponse>('/product_processs', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getProductProcess(id: number) {
  const response = await requestClient.get<ProductProcess>(`/product_processs/${id}`);
  return response;
}

export async function updateProductProcess(
  id: number,
  data: {
    productCode?: string;
    processCode?: string;
    processSeq?: number;
  }
) {
  const response = await requestClient.post<ProductProcess>(`/product_processs/${id}`, data);
  return response;
}

export async function deleteProductProcess(id: number) {
  const response = await requestClient.delete(`/product_processs/${id}`);
  return response;
}
