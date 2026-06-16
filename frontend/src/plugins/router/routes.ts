export const routes = [
  { path: '/', redirect: '/dashboard' },
  {
    path: '/',
    component: () => import('@/layouts/default.vue'),
    children: [
      {
        path: 'dashboard',
        component: () => import('@/pages/dashboard.vue'),
      },
      {
        path: 'account-settings',
        component: () => import('@/pages/account-settings.vue'),
        // 示例：仅管理员可访问账号设置
        meta: { roles: ['admin'] },
      },
      {
        path: 'typography',
        component: () => import('@/pages/typography.vue'),
      },
      {
        path: 'icons',
        component: () => import('@/pages/icons.vue'),
      },
      {
        path: 'cards',
        component: () => import('@/pages/cards.vue'),
      },
      {
        path: 'tables',
        component: () => import('@/pages/tables.vue'),
        meta: { roles: ['admin'] },
      },
      {
        path: 'admin/users',
        component: () => import('@/pages/admin/users.vue'),
        meta: { roles: ['admin'] },
      },
      {
        path: 'form-layouts',
        component: () => import('@/pages/form-layouts.vue'),
      },
      {
        path: 'forbidden',
        component: () => import('@/pages/forbidden.vue'),
        meta: { public: true },
      },
    ],
  },
  {
    path: '/',
    component: () => import('@/layouts/blank.vue'),
    children: [
      {
        path: 'login',
        component: () => import('@/pages/login.vue'),
        meta: { public: true },
      },
      // 公开注册入口已关闭：路由仍保留，仅 admin 可访问（不删除页面代码，也不提供入口跳转）
      {
        path: 'register',
        component: () => import('@/pages/register.vue'),
        meta: { roles: ['admin'] },
      },
      {
        path: '/:pathMatch(.*)*',
        component: () => import('@/pages/[...error].vue'),
        meta: { public: true },
      },
    ],
  },
]
