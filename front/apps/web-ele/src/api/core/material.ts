import { requestClient } from '../request';
import { search } from './search';

export interface Material {
  id: number;
  materialCode: string;
  materialName: string;
  materialCountry: string;
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
  materialCountry?: string;
  createdBy: number;
}) {
  const response = await requestClient.post<Material>('/materials', data);
  return response;
}

export async function getMaterial(id: number) {
  const response = await requestClient.get<Material>(`/materials/${id}`);
  return response;
}

export async function updateMaterial(
  id: number,
  data: {
    materialCode?: string;
    materialName?: string;
    materialCountry?: string;
  }
) {
  const response = await requestClient.post<Material>(`/materials/${id}`, data);
  return response;
}

export async function deleteMaterial(id: number) {
  const response = await requestClient.delete(`/materials/${id}`);
  return response;
}

// List queries use search API
export async function getMaterials(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'material',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.hits || [],
  };
}

// Backward compatibility aliases
export const getMaterialListApi = getMaterials;
export const getMaterialDetailApi = getMaterial;
export const createMaterialApi = createMaterial;
export const updateMaterialApi = updateMaterial;
export const deleteMaterialApi = deleteMaterial;
