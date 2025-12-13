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
}) {
  return requestClient.get<RoleListResponse>('/role', { params });
}

/**
 * 获取职位详情
 */
export async function getRoleDetailApi(id: number) {
  return requestClient.get<Role>(`/role/${id}`);
}

/**
 * 创建职位
 */
export async function createRoleApi(data: {
  name: string;
  code: string;
  description?: string;
  level?: number;
}) {
  return requestClient.post<Role>('/role', data);
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
) {
  return requestClient.put<Role>(`/role/${id}`, data);
}

/**
 * 停用职位
 */
export async function deactivateRoleApi(id: number) {
  return requestClient.put<{ message: string }>(`/role/${id}/deactivate`);
}

/**
 * 激活职位
 */
export async function activateRoleApi(id: number) {
  return requestClient.put<{ message: string }>(`/role/${id}/activate`);
}
