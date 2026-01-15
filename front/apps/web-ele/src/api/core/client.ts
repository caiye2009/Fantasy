import { requestClient } from '../request';

export interface Client {
  id: number;
  clientCode: string;
  clientName: string;
  clientCountry: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface ClientListResponse {
  total: number;
  list: Client[];
}

export async function createClient(data: {
  clientCode: string;
  clientName: string;
  clientCountry?: string;
  createdBy: number;
}) {
  const response = await requestClient.post<Client>('/clients', data);
  return response;
}

export async function getClient(id: number) {
  const response = await requestClient.get<Client>(`/clients/${id}`);
  return response;
}

export async function updateClient(
  id: number,
  data: {
    clientCode?: string;
    clientName?: string;
    clientCountry?: string;
  }
) {
  const response = await requestClient.post<Client>(`/clients/${id}`, data);
  return response;
}

export async function deleteClient(id: number) {
  const response = await requestClient.delete(`/clients/${id}`);
  return response;
}

// Backward compatibility aliases
export const getClientDetailApi = getClient;
export const createClientApi = createClient;
export const updateClientApi = updateClient;
export const deleteClientApi = deleteClient;
