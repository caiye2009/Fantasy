import { requestClient } from '../request';

export interface User {
  id: number;
  loginId: string;
  username: string;
  password: string;
  realName: string;
  email: string;
  role: string;
  department: string;
  isActive: boolean;
  hasInitPass: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UserListResponse {
  total: number;
  list: User[];
}

export async function createUser(data: {
  loginId: string;
  username: string;
  password: string;
  realName: string;
  email?: string;
  role: string;
  department?: string;
  isActive?: boolean;
  hasInitPass?: boolean;
}) {
  const response = await requestClient.post<User>('/users', data);
  return response;
}

export async function getUser(id: number) {
  const response = await requestClient.get<User>(`/users/${id}`);
  return response;
}

export async function updateUser(
  id: number,
  data: {
    loginId?: string;
    username?: string;
    password?: string;
    realName?: string;
    email?: string;
    role?: string;
    department?: string;
    isActive?: boolean;
    hasInitPass?: boolean;
  }
) {
  const response = await requestClient.post<User>(`/users/${id}`, data);
  return response;
}

export async function deleteUser(id: number) {
  const response = await requestClient.delete(`/users/${id}`);
  return response;
}

// List queries - User module doesn't use search, uses direct endpoint
export async function getUsers(params?: { limit?: number; offset?: number }) {
  const response = await requestClient.get<UserListResponse>('/users', {
    params: {
      limit: params?.limit || 100,
      offset: params?.offset || 0,
    },
  });
  return response;
}

// Get current user info from token or endpoint
export async function getUserInfoApi() {
  const response = await requestClient.get('/users/info');
  return response;
}

// Get departments list (hardcoded as department module was removed from backend)
export async function getDepartmentsApi(): Promise<{ departments: string[] }> {
  // TODO: 如果需要从后端获取,可以添加相应的接口
  return {
    departments: [
      '行政部',
      '人力资源部',
      '财务部',
      '生产部',
      '销售部',
      '采购部',
      '仓储部',
      '质检部',
      '研发部',
    ],
  };
}

// Get roles list (hardcoded as role module was removed from backend)
export async function getRolesApi(): Promise<{ roles: string[] }> {
  // TODO: 如果需要从后端获取,可以添加相应的接口
  return {
    roles: [
      'admin',
      'hr',
      'financeDirector',
      'finance',
      'productionDirector',
      'productionSpecialist',
      'orderCoordinator',
      'salesManager',
      'salesAssistant',
      'user',
    ],
  };
}

// Backward compatibility aliases
export const getUserListApi = getUsers;
export const getUserDetailApi = getUser;
export const createUserApi = createUser;
export const updateUserApi = updateUser;
export const deleteUserApi = deleteUser;
