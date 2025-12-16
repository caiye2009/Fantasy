import { requestClient } from '#/api/request';

export interface Department {
  id: number;
  name: string;
  code: string;
  description: string;
  status: string;
  parent_id?: number;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface DepartmentListResponse {
  total: number;
  departments: Department[];
}

/**
 * 获取部门列表
 */
export async function getDepartmentListApi(params?: {
  status?: string;
  page?: number;
  page_size?: number;
}): Promise<DepartmentListResponse> {
  const response = await requestClient.get<DepartmentListResponse>('/department', { params });
  return response.data;
}

/**
 * 获取部门详情
 */
export async function getDepartmentDetailApi(id: number): Promise<Department> {
  const response = await requestClient.get<Department>(`/department/${id}`);
  return response.data;
}

/**
 * 创建部门
 */
export async function createDepartmentApi(data: {
  name: string;
  code?: string;
  description?: string;
  parent_id?: number;
}): Promise<Department> {
  const response = await requestClient.post<Department>('/department', data);
  return response.data;
}

/**
 * 更新部门
 */
export async function updateDepartmentApi(
  id: number,
  data: {
    name?: string;
    code?: string;
    description?: string;
    parent_id?: number;
  },
): Promise<Department> {
  const response = await requestClient.put<Department>(`/department/${id}`, data);
  return response.data;
}

/**
 * 停用部门
 */
export async function deactivateDepartmentApi(id: number): Promise<{ message: string }> {
  const response = await requestClient.put<{ message: string }>(
    `/department/${id}/deactivate`,
  );
  return response.data;
}

/**
 * 激活部门
 */
export async function activateDepartmentApi(id: number): Promise<{ message: string }> {
  const response = await requestClient.put<{ message: string }>(`/department/${id}/activate`);
  return response.data;
}
