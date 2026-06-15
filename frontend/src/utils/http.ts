import axios, { AxiosError } from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios'

/**
 * 统一后端响应结构（与 backend/pkg/app/response.go 对齐）
 */
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
  details?: string[]
}

const TOKEN_KEY = 'insight_token'

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

const http: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api/v1',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：注入 Token
http.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getToken()
    if (token && config.headers)
      config.headers.Authorization = `Bearer ${token}`

    return config
  },
  (error: AxiosError) => Promise.reject(error),
)

// 响应拦截器：解包业务结构 + 401 跳登录
// 注意：成功时直接返回 body.data，调用方 `await http.post()` 直接拿到业务数据
http.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const body = response.data

    // 兼容直接返回非标准结构（例如静态资源、二进制流）
    if (body == null || typeof body !== 'object' || !('code' in body))
      return response.data as unknown as AxiosResponse

    if (body.code === 0)
      return body.data as unknown as AxiosResponse

    // 业务错误：rejection 给上层
    return Promise.reject(Object.assign(new Error(body.message || '请求失败'), {
      code: body.code,
      details: body.details,
    }))
  },
  (error: AxiosError<ApiResponse>) => {
    const status = error.response?.status
    const body = error.response?.data

    if (status === 401) {
      clearToken()
      // 避免在登录页再次触发跳转
      if (!location.pathname.startsWith('/login'))
        location.replace('/login')
    }

    const message = body?.message || error.message || '网络错误'

    return Promise.reject(Object.assign(new Error(message), {
      code: body?.code ?? status,
      details: body?.details,
    }))
  },
)

/**
 * 类型化请求工具：返回 data 部分
 */
export function request<T = unknown>(config: AxiosRequestConfig): Promise<T> {
  return http.request<unknown, T>(config)
}

export default http
