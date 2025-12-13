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
];

export default routes;