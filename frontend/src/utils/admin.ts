import http from '@/utils/http'
import type { UserInfo } from '@/stores/auth'

/** 创建用户的请求体（与后端 RegisterRequest 对齐） */
export interface CreateUserPayload {
  username: string
  email: string
  password: string
  confirmPassword: string
}

/**
 * 管理员获取用户列表
 * GET /api/v1/admin/users
 */
export function listUsers(): Promise<UserInfo[]> {
  return http.get<unknown, UserInfo[]>('/admin/users')
}

/**
 * 管理员创建用户
 * POST /api/v1/admin/users
 * 复用后端 Register handler，新增账号默认 role=user
 */
export function createUser(payload: CreateUserPayload): Promise<void> {
  return http.post<unknown, void>('/admin/users', payload)
}

/** 管理员修改用户的请求体（除头像外的字段），空字符串 / 未传表示不修改 */
export interface UpdateUserPayload {
  username?: string
  email?: string
  password?: string
  role?: string
}

/**
 * 管理员修改用户
 * PUT /api/v1/admin/users/:id
 */
export function updateUser(id: number, payload: UpdateUserPayload): Promise<void> {
  return http.put<unknown, void>(`/admin/users/${id}`, payload)
}

/**
 * 管理员删除用户
 * DELETE /api/v1/admin/users/:id
 */
export function deleteUser(id: number): Promise<void> {
  return http.delete<unknown, void>(`/admin/users/${id}`)
}
