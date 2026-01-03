import { requestClient } from '../request';

export interface MaterialAllocation {
  id: number;
  allocationCode: string;
  orderCode: string;
  materialCode: string;
  allocationQty: number;
  actualQty: number;
  allocationStatus: string;
  allocatedBy: number;
  allocatedAt: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface MaterialAllocationListResponse {
  total: number;
  list: MaterialAllocation[];
}

export async function createMaterialAllocation(data: {
  allocationCode: string;
  orderCode: string;
  materialCode: string;
  allocationQty: number;
  actualQty?: number;
  allocationStatus: string;
  allocatedBy: number;
  allocatedAt?: string;
  createdBy: number;
}) {
  const response = await requestClient.post<MaterialAllocation>('/material_allocations', data);
  return response;
}

export async function getMaterialAllocations(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<MaterialAllocationListResponse>('/material_allocations', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function getMaterialAllocation(id: number) {
  const response = await requestClient.get<MaterialAllocation>(`/material_allocations/${id}`);
  return response;
}

export async function updateMaterialAllocation(
  id: number,
  data: {
    allocationCode?: string;
    orderCode?: string;
    materialCode?: string;
    allocationQty?: number;
    actualQty?: number;
    allocationStatus?: string;
    allocatedBy?: number;
    allocatedAt?: string;
  }
) {
  const response = await requestClient.post<MaterialAllocation>(`/material_allocations/${id}`, data);
  return response;
}

export async function deleteMaterialAllocation(id: number) {
  const response = await requestClient.delete(`/material_allocations/${id}`);
  return response;
}
