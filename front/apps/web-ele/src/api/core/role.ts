import { requestClient } from '#/api/request';

export interface Role {
  id: number;
  name: string;
  code: string;
  description: string;
  status: string;
  level: number;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface RoleListResponse {
  total: number;
  roles: Role[];
}

/**
 * 获取职位列表
 */
export async function getRoleListApi(params?: {
  status?: string;
  page?: number;
  page_size?: number;
}): Promise<RoleListResponse> {
  const response = await requestClient.get<RoleListResponse>('/role', { params });
  return response.data;
}

/**
 * 获取职位详情
 */
export async function getRoleDetailApi(id: number): Promise<Role> {
  const response = await requestClient.get<Role>(`/role/${id}`);
  return response.data;
}

/**
 * 创建职位
 */
export async function createRoleApi(data: {
  name: string;
  code: string;
  description?: string;
  level?: number;
}): Promise<Role> {
  const response = await requestClient.post<Role>('/role', data);
  return response.data;
}

/**
 * 更新职位
 */
export async function updateRoleApi(
  id: number,
  data: {
    name?: string;
    code?: string;
    description?: string;
    level?: number;
  },
): Promise<Role> {
  const response = await requestClient.put<Role>(`/role/${id}`, data);
  return response.data;
}

/**
 * 停用职位
 */
export async function deactivateRoleApi(id: number): Promise<{ message: string }> {
  const response = await requestClient.put<{ message: string }>(`/role/${id}/deactivate`);
  return response.data;
}

/**
 * 激活职位
 */
export async function activateRoleApi(id: number): Promise<{ message: string }> {
  const response = await requestClient.put<{ message: string }>(`/role/${id}/activate`);
  return response.data;
}
