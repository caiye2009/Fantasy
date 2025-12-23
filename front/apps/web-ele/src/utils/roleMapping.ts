/**
 * 角色映射工具
 * 将后端具体角色映射到订单管理的业务角色
 */

import type { RoleType } from '#/views/order-management/types'

/**
 * 将后端角色映射到订单管理的业务角色
 * @param backendRole 后端角色
 * @returns 订单管理业务角色 (sales | follower | warehouse | system)
 */
export function mapToOrderRole(backendRole: string): RoleType {
  // 销售相关角色映射到 sales
  if (['salesManager', 'salesAssistant'].includes(backendRole)) {
    return 'sales'
  }

  // 跟单和生产相关角色映射到 follower
  if ([
    'orderCoordinator',
    'productionAssistant',
    'productionSpecialist',
    'productionDirector'
  ].includes(backendRole)) {
    return 'follower'
  }

  // 仓管角色映射到 warehouse
  if (backendRole === 'warehouse') {
    return 'warehouse'
  }

  // 其他角色（admin, hr, finance等）默认为 system
  return 'system'
}

/**
 * 检查用户是否有订单管理权限
 * @param backendRole 后端角色
 * @returns 是否有权限
 */
export function hasOrderManagementPermission(backendRole: string): boolean {
  const orderRole = mapToOrderRole(backendRole)
  return ['sales', 'follower', 'warehouse'].includes(orderRole)
}

/**
 * 获取角色显示名称
 * @param role 订单管理角色
 * @returns 显示名称
 */
export function getOrderRoleName(role: RoleType): string {
  const nameMap: Record<RoleType, string> = {
    sales: '业务',
    follower: '跟单',
    warehouse: '仓库',
    system: '系统'
  }
  return nameMap[role] || role
}
