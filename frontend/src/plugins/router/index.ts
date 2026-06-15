import type { App } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { routes } from './routes'
import { useAuthStore } from '@/stores/auth'
import { getToken } from '@/utils/http'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

router.beforeEach((to) => {
  const isPublic = to.meta?.public === true
  const requiredRoles = to.meta?.roles

  // store 在 router 安装到 app 之前可能尚未就绪，这里直接读 token 做兜底
  let authed = false
  let role = ''
  try {
    const auth = useAuthStore()
    authed = auth.isAuthenticated
    role = auth.role
  }
  catch {
    authed = Boolean(getToken())
  }

  // 1) 未登录访问受保护路由 → 跳登录
  if (!authed && !isPublic)
    return { path: '/login', query: { redirect: to.fullPath } }

  // 2) 已登录访问登录/注册页 → 回主页
  if (authed && (to.path === '/login' || to.path === '/register'))
    return { path: '/' }

  // 3) 角色权限校验：路由声明 meta.roles 时，仅匹配角色的用户可访问
  if (authed && requiredRoles && requiredRoles.length > 0) {
    if (!role || !requiredRoles.includes(role))
      return { path: '/forbidden', query: { from: to.fullPath } }
  }

  return true
})

export default function (app: App) {
  app.use(router)
}

export { router }
