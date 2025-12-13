import { requestClient } from '#/api/request';

/**
 * 供应商信息
 */
export interface Supplier {
  id: number;
  name: string;
  contact?: string;
  phone?: string;
  address?: string;
  email?: string;
  type?: string;
  status?: string;
  created_at?: string;
  updated_at?: string;
}

/**
 * 供应商列表响应
 */
export interface SupplierListResponse {
  total: number;
  suppliers: Supplier[];
}

/**
 * 获取供应商列表
 */
export async function getSupplierListApi(params?: {
  limit?: number;
  offset?: number;
}) {
  return requestClient.get<SupplierListResponse>('/supplier', { params });
}

/**
 * 获取供应商详情
 */
export async function getSupplierDetailApi(id: number) {
  return requestClient.get<Supplier>(`/supplier/${id}`);
}

/**
 * 创建供应商
 */
export async function createSupplierApi(data: {
  name: string;
  contact?: string;
  phone?: string;
  address?: string;
  email?: string;
  type?: string;
}) {
  return requestClient.post<Supplier>('/supplier', data);
}

/**
 * 更新供应商
 */
export async function updateSupplierApi(
  id: number,
  data: {
    name?: string;
    contact?: string;
    phone?: string;
    address?: string;
    email?: string;
    type?: string;
    status?: string;
  }
) {
  return requestClient.put(`/supplier/${id}`, data);
}

/**
 * 删除供应商
 */
export async function deleteSupplierApi(id: number) {
  return requestClient.delete(`/supplier/${id}`);
}
