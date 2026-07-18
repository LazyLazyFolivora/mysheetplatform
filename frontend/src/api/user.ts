import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

// Go backend: user management via admin routes
export function getUsers(page = 1, pageSize = 20) {
  return request.get<ApiResponse<any>>(`/admin/users?page=${page}&page_size=${pageSize}`)
}

export function banUser(id: number) {
  return request.post<ApiResponse<any>>(`/admin/users/${id}/ban`)
}

export function unbanUser(id: number) {
  return request.post<ApiResponse<any>>(`/admin/users/${id}/unban`)
}
