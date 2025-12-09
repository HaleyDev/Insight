import type { RouteRecordRaw } from 'vue-router';

import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    meta: {
      icon: 'ic:baseline-view-in-ar',
      keepAlive: true,
      order: 1000,
      title: "个人数字人",
    },
    name: 'DigitalHuman ',
    path: '/dgman',
    children: [
      {
        meta: {
          title: "个人数字人",
        },
        name: 'DigitalHuman01',
        path: '/dgman/digitalhuman01',
        component: () => import('#/views/digital_human/index.vue'),
      },
    ],
  },
];

export default routes;
