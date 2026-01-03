import { requestClient } from '../request';

export interface PurchaseOrder {
  id: number;
  purchaseCode: string;
  orderCode: string;
  supplierCode: string;
  purchaseDate: string;
  purchaseStatus: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface PurchaseOrderListResponse {
  total: number;
  list: PurchaseOrder[];
}

export async function createPurchaseOrder(data: {
  purchaseCode: string;
  orderCode: string;
  supplierCode: string;
  purchaseDate: string;
  purchaseStatus: string;
  createdBy: number;
}) {
  const response = await requestClient.post<PurchaseOrder>('/purchase_orders', data);
  return response;
}

export async function getPurchaseOrders(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<PurchaseOrderListResponse>('/purchase_orders', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getPurchaseOrder(id: number) {
  const response = await requestClient.get<PurchaseOrder>(`/purchase_orders/${id}`);
  return response;
}

export async function updatePurchaseOrder(
  id: number,
  data: {
    purchaseCode?: string;
    orderCode?: string;
    supplierCode?: string;
    purchaseDate?: string;
    purchaseStatus?: string;
  }
) {
  const response = await requestClient.post<PurchaseOrder>(`/purchase_orders/${id}`, data);
  return response;
}

export async function deletePurchaseOrder(id: number) {
  const response = await requestClient.delete(`/purchase_orders/${id}`);
  return response;
}
