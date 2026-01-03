import { requestClient } from '../request';

export interface OrderEvent {
  id: number;
  eventId: string;
  orderCode: string;
  eventType: string;
  eventContent: string;
  operator: string;
  operatedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface OrderEventListResponse {
  total: number;
  list: OrderEvent[];
}

export async function createOrderEvent(data: {
  eventId: string;
  orderCode: string;
  eventType: string;
  eventContent?: string;
  operator: string;
  operatedAt: string;
}) {
  const response = await requestClient.post<OrderEvent>('/order_events', data);
  return response;
}

export async function getOrderEvents(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<OrderEventListResponse>('/order_events', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getOrderEvent(id: number) {
  const response = await requestClient.get<OrderEvent>(`/order_events/${id}`);
  return response;
}

export async function updateOrderEvent(
  id: number,
  data: {
    eventId?: string;
    orderCode?: string;
    eventType?: string;
    eventContent?: string;
    operator?: string;
    operatedAt?: string;
  }
) {
  const response = await requestClient.post<OrderEvent>(`/order_events/${id}`, data);
  return response;
}

export async function deleteOrderEvent(id: number) {
  const response = await requestClient.delete(`/order_events/${id}`);
  return response;
}
