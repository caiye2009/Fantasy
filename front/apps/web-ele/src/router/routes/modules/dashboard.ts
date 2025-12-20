import type { RouteRecordRaw } from 'vue-router';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    name: 'Workspace',
    path: '/workspace',
    component: () => import('#/views/workspace.vue'),
    meta: {
      affixTab: true,
      icon: 'carbon:workspace',
      title: $t('page.dashboard.workspace'),
    },
  },
  {
    name: 'Materials',
    path: '/materials',
    component: () => import('#/views/materials.vue'),
    meta: {
      icon: 'lucide:circle-pile',
      title: $t('page.dashboard.materials'),
    },
  },
  {
    name: 'Processes',
    path: '/processes',
    component: () => import('#/views/processes.vue'),
    meta: {
      icon: 'lucide:spool',
      title: $t('page.dashboard.processes'),
    },
  },
  {
    name: 'Products',
    path: '/products',
    component: () => import('#/views/products.vue'),
    meta: {
      icon: 'lucide:package',
      title: $t('page.dashboard.products'),
    },
  },
  {
    name: 'Clients',
    path: '/clients',
    component: () => import('#/views/clients.vue'),
    meta: {
      icon: 'lucide:book-user',
      title: $t('page.dashboard.clients'),
    },
  },
  {
    name: 'Analytics',
    path: '/analytics',
    component: () => import('#/views/analytics/index.vue'),
    meta: {
      icon: 'lucide:area-chart',
      title: $t('page.dashboard.analytics'),
    },
  },
  {
    name: 'Demo',
    path: '/demo',
    component: () => import('#/views/demo.vue'),
    meta: {
      icon: 'lucide:area-chart',
      title: $t('page.dashboard.demo'),
    },
  },
  {
    name: 'Users',
    path: '/users',
    component: () => import('#/views/users.vue'),
    meta: {
      icon: 'lucide:users',
      title: $t('page.dashboard.users'),
    },
  },
  {
    name: 'OrderManagement',
    path: '/order-management',
    component: () => import('#/views/order-management/index.vue'),
    meta: {
      icon: 'lucide:clipboard-list',
      title: '订单管理',
    },
  },
  {
    name: 'MaterialManagement',
    path: '/material-management',
    component: () => import('#/views/material-management/index.vue'),
    meta: {
      icon: 'lucide:package-search',
      title: '原料管理',
    },
  },
  {
    name: 'InventoryManagement',
    path: '/inventory-management',
    component: () => import('#/views/inventory-management/index.vue'),
    meta: {
      icon: 'lucide:warehouse',
      title: '库存管理',
    },
  },
];

export default routes;