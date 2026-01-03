import { requestClient } from '../request';

export interface OrderProduct {
  id: number;
  orderCode: string;
  productCode: string;
  orderedQty: number;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface OrderProductListResponse {
  total: number;
  list: OrderProduct[];
}

export async function createOrderProduct(data: {
  orderCode: string;
  productCode: string;
  orderedQty: number;
  createdBy: number;
}) {
  const response = await requestClient.post<OrderProduct>('/order_products', data);
  return response;
}

export async function getOrderProducts(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<OrderProductListResponse>('/order_products', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getOrderProduct(id: number) {
  const response = await requestClient.get<OrderProduct>(`/order_products/${id}`);
  return response;
}

export async function updateOrderProduct(
  id: number,
  data: {
    orderCode?: string;
    productCode?: string;
    orderedQty?: number;
  }
) {
  const response = await requestClient.post<OrderProduct>(`/order_products/${id}`, data);
  return response;
}

export async function deleteOrderProduct(id: number) {
  const response = await requestClient.delete(`/order_products/${id}`);
  return response;
}
