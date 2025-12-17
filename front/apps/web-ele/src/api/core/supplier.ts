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
  createdAt?: string;
  updatedAt?: string;
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
}): Promise<SupplierListResponse> {
  const response = await requestClient.get<SupplierListResponse>('/supplier', { params });
  return response.data;
}

/**
 * 获取供应商详情
 */
export async function getSupplierDetailApi(id: number): Promise<Supplier> {
  const response = await requestClient.get<Supplier>(`/supplier/${id}`);
  return response.data;
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
}): Promise<Supplier> {
  const response = await requestClient.post<Supplier>('/supplier', data);
  return response.data;
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
): Promise<any> {
  const response = await requestClient.put(`/supplier/${id}`, data);
  return response.data;
}

/**
 * 删除供应商
 */
export async function deleteSupplierApi(id: number): Promise<any> {
  const response = await requestClient.delete(`/supplier/${id}`);
  return response.data;
}
