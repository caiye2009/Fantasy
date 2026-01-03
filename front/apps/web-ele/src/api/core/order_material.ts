import { requestClient } from '../request';

export interface OrderMaterial {
  id: number;
  orderCode: string;
  materialCode: string;
  requiredQty: number;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface OrderMaterialListResponse {
  total: number;
  list: OrderMaterial[];
}

export async function createOrderMaterial(data: {
  orderCode: string;
  materialCode: string;
  requiredQty: number;
  createdBy: number;
}) {
  const response = await requestClient.post<OrderMaterial>('/order_materials', data);
  return response;
}

export async function getOrderMaterials(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<OrderMaterialListResponse>('/order_materials', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getOrderMaterial(id: number) {
  const response = await requestClient.get<OrderMaterial>(`/order_materials/${id}`);
  return response;
}

export async function updateOrderMaterial(
  id: number,
  data: {
    orderCode?: string;
    materialCode?: string;
    requiredQty?: number;
  }
) {
  const response = await requestClient.post<OrderMaterial>(`/order_materials/${id}`, data);
  return response;
}

export async function deleteOrderMaterial(id: number) {
  const response = await requestClient.delete(`/order_materials/${id}`);
  return response;
}
