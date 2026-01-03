import { requestClient } from '../request';

export interface ProductionOrder {
  id: number;
  productionCode: string;
  orderCode: string;
  processCode: string;
  supplierCode: string;
  productionQty: number;
  productionStatus: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProductionOrderListResponse {
  total: number;
  list: ProductionOrder[];
}

export async function createProductionOrder(data: {
  productionCode: string;
  orderCode: string;
  processCode: string;
  supplierCode: string;
  productionQty: number;
  productionStatus: string;
  createdBy: number;
}) {
  const response = await requestClient.post<ProductionOrder>('/production_orders', data);
  return response;
}

export async function getProductionOrders(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<ProductionOrderListResponse>('/production_orders', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getProductionOrder(id: number) {
  const response = await requestClient.get<ProductionOrder>(`/production_orders/${id}`);
  return response;
}

export async function updateProductionOrder(
  id: number,
  data: {
    productionCode?: string;
    orderCode?: string;
    processCode?: string;
    supplierCode?: string;
    productionQty?: number;
    productionStatus?: string;
  }
) {
  const response = await requestClient.post<ProductionOrder>(`/production_orders/${id}`, data);
  return response;
}

export async function deleteProductionOrder(id: number) {
  const response = await requestClient.delete(`/production_orders/${id}`);
  return response;
}
