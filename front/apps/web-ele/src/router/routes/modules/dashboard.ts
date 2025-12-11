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
        name: 'Workspace',
        path: '/workspace',
        component: () => import('#/views/dashboard/workspace/index.vue'),
        meta: {
          affixTab: true,
          icon: 'carbon:workspace',
          title: $t('page.dashboard.workspace'),
        },
      },
      {
        name: 'Materials',
        path: '/materials',
        component: () => import('#/views/dashboard/materials/index.vue'),
        meta: {
          icon: 'carbon:workspace',
          title: $t('page.dashboard.materials'),
        },
      },
      {
        name: 'Analytics',
        path: '/analytics',
        component: () => import('#/views/dashboard/analytics/index.vue'),
        meta: {
          icon: 'lucide:area-chart',
          title: $t('page.dashboard.analytics'),
        },
      },
      {
        name: 'list',
        path: '/list',
        component: () => import('#/views/dashboard/list/index.vue'),
        meta: {
          icon: 'lucide:area-chart',
          title: $t('page.dashboard.list'),
        },
      },
    ],
  },
];

export default routes;