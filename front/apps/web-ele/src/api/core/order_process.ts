import { requestClient } from '../request';

export interface OrderProcess {
  id: number;
  orderCode: string;
  processCode: string;
  processSeq: number;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface OrderProcessListResponse {
  total: number;
  list: OrderProcess[];
}

export async function createOrderProcess(data: {
  orderCode: string;
  processCode: string;
  processSeq: number;
  createdBy: number;
}) {
  const response = await requestClient.post<OrderProcess>('/order_processs', data);
  return response;
}

export async function getOrderProcesses(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<OrderProcessListResponse>('/order_processs', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getOrderProcess(id: number) {
  const response = await requestClient.get<OrderProcess>(`/order_processs/${id}`);
  return response;
}

export async function updateOrderProcess(
  id: number,
  data: {
    orderCode?: string;
    processCode?: string;
    processSeq?: number;
  }
) {
  const response = await requestClient.post<OrderProcess>(`/order_processs/${id}`, data);
  return response;
}

export async function deleteOrderProcess(id: number) {
  const response = await requestClient.delete(`/order_processs/${id}`);
  return response;
}
