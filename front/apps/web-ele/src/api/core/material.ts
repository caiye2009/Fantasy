import { requestClient } from '../request';

export interface Material {
  id: number;
  materialCode: string;
  materialName: string;
  materialCategory?: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface MaterialListResponse {
  total: number;
  list: Material[];
}

export async function createMaterial(data: {
  materialCode: string;
  materialName: string;
  materialCategory?: string;
  createdBy: number;
}) {
  const response = await requestClient.post<Material>('/materials', data);
  return response;
}

export async function getMaterial(id: number) {
  const response = await requestClient.get<Material>(`/materials/${id}`);
  return response;
}

export async function getMaterials(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<MaterialListResponse>('/materials', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

export async function updateMaterial(
  id: number,
  data: {
    materialCode?: string;
    materialName?: string;
    materialCategory?: string;
  }
) {
  const response = await requestClient.put<Material>(`/materials/${id}`, data);
  return response;
}

export async function deleteMaterial(id: number) {
  const response = await requestClient.delete(`/materials/${id}`);
  return response;
}
