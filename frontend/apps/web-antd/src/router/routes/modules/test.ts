import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'ic:baseline-view-in-ar',
      keepAlive: true,
      order: 1000,
      title: "测试页面",
    },
    name: 'Test',
    path: '/test',
    children: [
      {
        meta: {
          title: "测试页面 01",
        },
        name: 'Test01',
        path: '/test/test01',
        component: () => import('#/views/test/test01.vue'),
      },
    ],
  },
];

export default routes;
