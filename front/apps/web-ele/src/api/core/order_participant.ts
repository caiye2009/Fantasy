import { requestClient } from '../request';

export interface OrderParticipant {
  id: number;
  orderCode: string;
  username: string;
  participantRole: string;
  createdAt: string;
  updatedAt: string;
}

export interface OrderParticipantListResponse {
  total: number;
  list: OrderParticipant[];
}

export async function createOrderParticipant(data: {
  orderCode: string;
  username: string;
  participantRole: string;
}) {
  const response = await requestClient.post<OrderParticipant>('/order_participants', data);
  return response;
}

export async function getOrderParticipants(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<OrderParticipantListResponse>('/order_participants', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getOrderParticipant(id: number) {
  const response = await requestClient.get<OrderParticipant>(`/order_participants/${id}`);
  return response;
}

export async function updateOrderParticipant(
  id: number,
  data: {
    orderCode?: string;
    username?: string;
    participantRole?: string;
  }
) {
  const response = await requestClient.post<OrderParticipant>(`/order_participants/${id}`, data);
  return response;
}

export async function deleteOrderParticipant(id: number) {
  const response = await requestClient.delete(`/order_participants/${id}`);
  return response;
}
