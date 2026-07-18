import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

export interface SystemConfig {
  id: number
  config_key: string
  config_val: string
  remark: string
  created_at: string
  updated_at: string
}

export function getSystemConfigs() {
  return request.get<ApiResponse<SystemConfig[]>>('/admin/system-config')
}

export function createSystemConfig(key: string, value: string, remark?: string) {
  return request.post<ApiResponse<SystemConfig>>('/admin/system-config', { key, value, remark })
}

export function updateSystemConfig(key: string, value: string, remark?: string) {
  return request.put<ApiResponse<void>>('/admin/system-config', { key, value, remark })
}

export function deleteSystemConfig(key: string) {
  return request.delete<ApiResponse<void>>(`/admin/system-config/${encodeURIComponent(key)}`)
}
