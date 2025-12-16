import type { RouteRecordStringComponent } from '@vben/types';

import { requestClient } from '#/api/request';

/**
 * 获取用户所有菜单
 */
export async function getAllMenusApi(): Promise<RouteRecordStringComponent[]> {
  const response = await requestClient.get<RouteRecordStringComponent[]>('/menu/all');
  return response.data;
}
