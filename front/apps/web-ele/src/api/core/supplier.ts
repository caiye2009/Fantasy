import { requestClient } from '../request';
import { search } from './search';

export interface Supplier {
  id: number;
  supplierCode: string;
  supplierName: string;
  supplierCategory: string;
  createdBy: number;
  createdAt: string;
  updatedAt: string;
}

export interface SupplierListResponse {
  total: number;
  list: Supplier[];
}

export async function createSupplier(data: {
  supplierCode: string;
  supplierName: string;
  supplierCategory?: string;
  createdBy: number;
}) {
  const response = await requestClient.post<Supplier>('/suppliers', data);
  return response;
}

export async function getSupplier(id: number) {
  const response = await requestClient.get<Supplier>(`/suppliers/${id}`);
  return response;
}

export async function updateSupplier(
  id: number,
  data: {
    supplierCode?: string;
    supplierName?: string;
    supplierCategory?: string;
  }
) {
  const response = await requestClient.post<Supplier>(`/suppliers/${id}`, data);
  return response;
}

export async function deleteSupplier(id: number) {
  const response = await requestClient.delete(`/suppliers/${id}`);
  return response;
}

// List queries use search API
export async function getSuppliers(params?: { limit?: number; offset?: number }) {
  const result = await search({
    index: 'supplier',
    from: params?.offset || 0,
    size: params?.limit || 20,
  });
  return {
    total: result.total || 0,
    list: result.hits || [],
  };
}

// Backward compatibility aliases
export const getSupplierListApi = getSuppliers;
export const getSupplierDetailApi = getSupplier;
export const createSupplierApi = createSupplier;
export const updateSupplierApi = updateSupplier;
export const deleteSupplierApi = deleteSupplier;
