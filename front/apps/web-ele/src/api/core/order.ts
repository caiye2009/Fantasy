import { requestClient } from '../request';
import { search } from './search';

export interface Order {
  id: number;
  orderCode: string;
  clientCode: string;
  orderDate: string;
  deliveryDate: string;
  orderStatus: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface OrderListResponse {
  total: number;
  list: Order[];
}

export async function createOrder(data: {
  orderCode: string;
  clientCode: string;
  orderDate: string;
  deliveryDate?: string;
  orderStatus: string;
  createdBy: number;
}) {
  const response = await requestClient.post<Order>('/orders', data);
  return response;
}

export async function getOrder(id: number) {
  const response = await requestClient.get<Order>(`/orders/${id}`);
  return response;
}

export async function updateOrder(
  id: number,
  data: {
    orderCode?: string;
    clientCode?: string;
    orderDate?: string;
    deliveryDate?: string;
    orderStatus?: string;
  }
) {
  const response = await requestClient.post<Order>(`/orders/${id}`, data);
  return response;
}

export async function deleteOrder(id: number) {
  const response = await requestClient.delete(`/orders/${id}`);
  return response;
}

// List queries use search API
export async function getOrders(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'order',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.hits || [],
  };
}

// Backward compatibility aliases
export const getOrderListApi = getOrders;
export const getOrderDetailApi = getOrder;
export const createOrderApi = createOrder;
export const updateOrderApi = updateOrder;
export const deleteOrderApi = deleteOrder;
