import type { RouteRecordRaw } from 'vue-router';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'lucide:layout-dashboard',
      order: -1,
      title: $t('page.dashboard.title'),
    },
    name: 'Dashboard',
    path: '/dashboard',
    children: [
      {
        name: 'Analytics',
        path: '/analytics',
        component: () => import('#/views/dashboard/analytics/index.vue'),
        meta: {
          affixTab: true,
          icon: 'lucide:area-chart',
          title: $t('page.dashboard.analytics'),
        },
      },
      {
        name: 'Workspace',
        path: '/workspace',
        component: () => import('#/views/dashboard/workspace/index.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: $t('page.dashboard.workspace'),
        },
      },
      {
        name: 'list',
        path: '/list',
        component: () => import('#/views/dashboard/list/index.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: $t('page.dashboard.list'),
        },
      },
      {
        name: 'material',
        path: '/material',
        component: () => import('#/views/dashboard/material/index.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: $t('page.dashboard.material'),
        },
      },
      {
        name: 'client',
        path: '/client',
        component: () => import('#/views/dashboard/client/index.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: $t('page.dashboard.client'),
        },
      },
      {
        name: 'product',
        path: '/product',
        component: () => import('#/views/dashboard/product/index.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: $t('page.dashboard.product'),
        },
      },
      {
        name: 'process',
        path: '/process',
        component: () => import('#/views/dashboard/process/index.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: $t('page.dashboard.process'),
        },
      },
    ],
  },
];

export default routes;
