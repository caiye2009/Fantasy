import type { UserInfo } from '@vben/types';

import { requestClient } from '#/api/request';

/**
 * 获取用户信息
 */
export async function getUserInfoApi(): Promise<UserInfo> {
  const response = await requestClient.get<UserInfo>('/user/info');
  return response.data;
}

/**
 * 获取用户列表
 */
export async function getUserListApi(params?: { limit?: number; offset?: number }): Promise<{ total: number; users: any[] }> {
  const response = await requestClient.get<{ total: number; users: any[] }>('/user', { params });
  return response.data;
}

/**
 * 创建用户
 */
export async function createUserApi(data: {
  username: string;
  department: string;
  role: string;
  email?: string;
}): Promise<{ login_id: string; password: string }> {
  const response = await requestClient.post<{ login_id: string; password: string }>('/user', data);
  return response.data;
}

/**
 * 获取部门列表
 */
export async function getDepartmentsApi(): Promise<{ departments: string[] }> {
  const response = await requestClient.get<{ departments: string[] }>('/user/departments');
  return response.data;
}

/**
 * 获取角色列表
 */
export async function getRolesApi(): Promise<{ roles: string[] }> {
  const response = await requestClient.get<{ roles: string[] }>('/user/roles');
  return response.data;
}

/**
 * 获取用户详情
 */
export async function getUserDetailApi(id: number): Promise<any> {
  const response = await requestClient.get<any>(`/user/${id}`);
  return response.data;
}

/**
 * 更新用户
 */
export async function updateUserApi(id: number, data: {
  username?: string;
  department?: string;
  email?: string;
  role?: string;
  status?: string;
}): Promise<any> {
  const response = await requestClient.put(`/user/${id}`, data);
  return response.data;
}

/**
 * 删除用户
 */
export async function deleteUserApi(id: number): Promise<any> {
  const response = await requestClient.delete(`/user/${id}`);
  return response.data;
}
