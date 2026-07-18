export interface PermissionModule {
  id: number
  name: string
  parent_id: number | null
  path: string
  icon: string
  sort_order: number
  created_at: string
}

export interface ModuleTreeNode extends PermissionModule {
  children: ModuleTreeNode[] | null
}

export interface Role {
  id: number
  name: string
  description: string
  created_at: string
}

export interface Permission {
  id: number
  name: string
  code: string
  url: string
  method: string
  created_at: string
}
