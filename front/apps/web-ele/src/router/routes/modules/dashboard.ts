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
  // 管理模块
  {
    name: 'Materials',
    path: '/materials',
    component: () => import('#/views/materials.vue'),
    meta: {
      icon: 'lucide:circle-pile',
      title: '原料管理',
    },
  },
  {
    name: 'Processes',
    path: '/processes',
    component: () => import('#/views/processes.vue'),
    meta: {
      icon: 'lucide:spool',
      title: '工艺管理',
    },
  },
  {
    name: 'Products',
    path: '/products',
    component: () => import('#/views/products.vue'),
    meta: {
      icon: 'lucide:package',
      title: '产品管理',
    },
  },
  {
    name: 'Clients',
    path: '/clients',
    component: () => import('#/views/clients.vue'),
    meta: {
      icon: 'lucide:book-user',
      title: '客户管理',
    },
  },
  // 其他模块
  {
    name: 'Orders',
    path: '/orders',
    component: () => import('#/views/orders.vue'),
    meta: {
      icon: 'lucide:clipboard-list',
      title: '订单管理',
    },
  },
  {
    name: 'Inventory',
    path: '/inventory',
    component: () => import('#/views/inventory.vue'),
    meta: {
      icon: 'lucide:warehouse',
      title: '库存管理',
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
    name: 'Users',
    path: '/users',
    component: () => import('#/views/users.vue'),
    meta: {
      icon: 'lucide:users',
      title: $t('page.dashboard.users'),
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
    name: 'Tt',
    path: '/tt',
    component: () => import('#/views/tt.vue'),
    meta: {
      icon: 'lucide:area-chart',
      title: $t('page.dashboard.tt'),
    },
  },
];

export default routes;