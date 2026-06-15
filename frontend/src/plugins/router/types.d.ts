import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    /** 是否为公共路由（无需登录即可访问） */
    public?: boolean
    /** 允许访问该路由的用户角色列表；为空表示登录用户均可访问 */
    roles?: string[]
  }
}
