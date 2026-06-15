import { defineStore } from 'pinia'
import http, { clearToken, getToken, setToken } from '@/utils/http'

export interface UserInfo {
  id?: number
  username?: string
  email?: string
  avatar?: string
  /** 用户角色，用于路由权限控制，例如 'admin' | 'user' */
  role?: string
  [key: string]: unknown
}

interface LoginPayload {
  email: string
  password: string
}

interface RegisterPayload {
  username: string
  email: string
  password: string
  confirmPassword: string
}

interface LoginResponseData {
  Token?: string
  token?: string
}

const USER_KEY = 'insight_user'

function loadUser(): UserInfo | null {
  try {
    const raw = localStorage.getItem(USER_KEY)
    return raw ? (JSON.parse(raw) as UserInfo) : null
  }
  catch {
    return null
  }
}

function saveUser(user: UserInfo | null) {
  if (user)
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  else
    localStorage.removeItem(USER_KEY)
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: getToken() || '',
    user: loadUser(),
  }),

  getters: {
    isAuthenticated: state => Boolean(state.token),
    role: (state): string => state.user?.role || '',
    /**
     * 权限判断：传入路由所需角色数组，返回是否允许访问
     * 空数组或 undefined 表示不限角色（仅需登录）
     */
    hasRole: (state) => (roles?: string[]): boolean => {
      if (!roles || roles.length === 0)
        return true
      return roles.includes(state.user?.role || '')
    },
  },

  actions: {
    /**
     * 邮箱登录：POST /api/v1/login
     */
    async login(payload: LoginPayload): Promise<void> {
      const data = await http.post<unknown, LoginResponseData>('/login', payload)
      const token = data?.Token || data?.token
      if (!token)
        throw new Error('登录失败：未获取到 Token')

      this.token = token
      setToken(token)

      // 可选：登录后拉取用户信息
      await this.fetchProfile().catch(() => undefined)
    },

    /**
     * 用户注册：POST /api/v1/Register
     */
    async register(payload: RegisterPayload): Promise<void> {
      await http.post('/Register', payload)
    },

    /**
     * 拉取当前用户信息（按需对接后端 /users/me 等接口）
     */
    async fetchProfile(): Promise<UserInfo | null> {
      try {
        const data = await http.get<unknown, UserInfo>('/users/me')
        this.user = data
        saveUser(data)
        return data
      }
      catch {
        return null
      }
    },

    /**
     * 登出：清理本地状态
     */
    logout(): void {
      this.token = ''
      this.user = null
      clearToken()
      saveUser(null)
    },
  },
})
