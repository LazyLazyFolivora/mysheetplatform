import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'
import type { ModuleTreeNode, Permission, PermissionModule, Role } from '@/types/permission'

// ===== 权限模块 =====

export function getPermissionTree() {
  return request.get<ApiResponse<ModuleTreeNode[]>>('/admin/permissions/tree')
}

export function createModule(data: { name: string; parent_id?: number | null; path?: string; icon?: string }) {
  return request.post<ApiResponse<PermissionModule>>('/admin/permissions/modules', data)
}

export function updateModule(
  id: number,
  data: { name: string; parent_id?: number | null; path?: string; icon?: string }
) {
  return request.put<ApiResponse<PermissionModule>>(`/admin/permissions/modules/${id}`, data)
}

export function deleteModule(id: number) {
  return request.delete<ApiResponse<void>>(`/admin/permissions/modules/${id}`)
}

export function getModulePermissions(moduleId: number) {
  return request.get<ApiResponse<number[]>>(`/admin/permissions/modules/${moduleId}/permissions`)
}

export function assignPermissionsToModule(moduleId: number, permIds: number[]) {
  return request.post<ApiResponse<void>>(`/admin/permissions/modules/${moduleId}/permissions`, {
    perm_ids: permIds,
  })
}

// ===== 权限 =====

export interface PermissionPayload {
  name: string
  code: string
  url?: string
  method?: string
  module_id?: number | null
}

export function getAllPermissions() {
  return request.get<ApiResponse<Permission[]>>('/admin/permissions')
}

export function createPermission(data: PermissionPayload) {
  return request.post<ApiResponse<Permission>>('/admin/permissions', data)
}

export function updatePermission(id: number, data: PermissionPayload) {
  return request.put<ApiResponse<Permission>>(`/admin/permissions/${id}`, data)
}

export function deletePermission(id: number) {
  return request.delete<ApiResponse<void>>(`/admin/permissions/${id}`)
}

// ===== 角色 =====

export function getAllRoles() {
  return request.get<ApiResponse<Role[]>>('/admin/roles')
}

export function createRole(data: { name: string; description?: string }) {
  return request.post<ApiResponse<Role>>('/admin/roles', data)
}

export function updateRole(id: number, data: { name: string; description?: string }) {
  return request.put<ApiResponse<Role>>(`/admin/roles/${id}`, data)
}

export function deleteRole(id: number) {
  return request.delete<ApiResponse<void>>(`/admin/roles/${id}`)
}

export function getAssignedModules(roleId: number) {
  return request.get<ApiResponse<number[]>>(`/admin/roles/${roleId}/modules`)
}

export function assignModules(roleId: number, moduleIds: number[]) {
  return request.post<ApiResponse<void>>(`/admin/roles/${roleId}/modules`, {
    module_ids: moduleIds,
  })
}
