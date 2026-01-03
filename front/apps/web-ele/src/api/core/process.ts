import { requestClient } from '../request';
import { search } from './search';

export interface Process {
  id: number;
  processCode: string;
  processName: string;
  processCountry: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface ProcessListResponse {
  total: number;
  list: Process[];
}

export async function createProcess(data: {
  processCode: string;
  processName: string;
  processCountry?: string;
  createdBy: number;
}) {
  const response = await requestClient.post<Process>('/processs', data);
  return response;
}

export async function getProcess(id: number) {
  const response = await requestClient.get<Process>(`/processs/${id}`);
  return response;
}

export async function updateProcess(
  id: number,
  data: {
    processCode?: string;
    processName?: string;
    processCountry?: string;
  }
) {
  const response = await requestClient.post<Process>(`/processs/${id}`, data);
  return response;
}

export async function deleteProcess(id: number) {
  const response = await requestClient.delete(`/processs/${id}`);
  return response;
}

// List queries use search API
export async function getProcesses(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'process',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.hits || [],
  };
}

// Backward compatibility aliases
export const getProcessListApi = getProcesses;
export const getProcessDetailApi = getProcess;
export const createProcessApi = createProcess;
export const updateProcessApi = updateProcess;
export const deleteProcessApi = deleteProcess;
