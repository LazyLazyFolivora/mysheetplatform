import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

export const adminApi = {
  // Dashboard stats
  getStats: () => request.get<ApiResponse<any>>('/admin/dashboard'),

  // Users
  getUsers: (page = 1, pageSize = 20) =>
    request.get<ApiResponse<any>>(`/admin/users?page=${page}&page_size=${pageSize}`),
  banUser: (id: number) => request.post<ApiResponse<any>>(`/admin/users/${id}/ban`),
  unbanUser: (id: number) => request.post<ApiResponse<any>>(`/admin/users/${id}/unban`),

  // Banned IPs
  getBannedIPs: (page = 1, pageSize = 20) =>
    request.get<ApiResponse<any>>(`/admin/banned-ip?page=${page}&page_size=${pageSize}`),
  createBannedIP: (data: any) => request.post<ApiResponse<any>>('/admin/banned-ip', data),
  deleteBannedIP: (id: number) => request.delete<ApiResponse<any>>(`/admin/banned-ip/${id}`),

  // Contact messages
  getMessages: (page = 1, pageSize = 20) =>
    request.get<ApiResponse<any>>(`/admin/contact-messages?page=${page}&page_size=${pageSize}`),
  markMessageRead: (id: number) => request.post<ApiResponse<any>>(`/admin/contact-messages/${id}/read`),
  deleteMessage: (id: number) => request.delete<ApiResponse<any>>(`/admin/contact-messages/${id}`),

  // Friend links
  getFriendLinks: (page = 1, pageSize = 10) =>
    request.get<ApiResponse<any>>(`/admin/friend-links?page=${page}&page_size=${pageSize}`),
  createFriendLink: (data: any) => request.post<ApiResponse<any>>('/admin/friend-links', data),
  updateFriendLink: (id: number, data: any) => request.put<ApiResponse<any>>(`/admin/friend-links/${id}`, data),
  deleteFriendLink: (id: number) => request.delete<ApiResponse<any>>(`/admin/friend-links/${id}`),
}

export default adminApi
