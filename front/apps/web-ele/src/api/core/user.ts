import type { UserInfo } from '@vben/types';

import { requestClient } from '#/api/request';

/**
 * 获取用户信息
 */
export async function getUserInfoApi() {
  return requestClient.get<UserInfo>('/user/info');
}

/**
 * 获取用户列表
 */
export async function getUserListApi(params?: { limit?: number; offset?: number }) {
  return requestClient.get<{ total: number; users: any[] }>('/user', { params });
}

/**
 * 创建用户
 */
export async function createUserApi(data: {
  username: string;
  department: string;
  role: string;
  email?: string;
}) {
  return requestClient.post<{ login_id: string; password: string }>('/user', data);
}

/**
 * 获取部门列表
 */
export async function getDepartmentsApi() {
  return requestClient.get<{ departments: string[] }>('/user/departments');
}

/**
 * 获取角色列表
 */
export async function getRolesApi() {
  return requestClient.get<{ roles: string[] }>('/user/roles');
}

/**
 * 获取用户详情
 */
export async function getUserDetailApi(id: number) {
  return requestClient.get<any>(`/user/${id}`);
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
}) {
  return requestClient.put(`/user/${id}`, data);
}

/**
 * 删除用户
 */
export async function deleteUserApi(id: number) {
  return requestClient.delete(`/user/${id}`);
}
