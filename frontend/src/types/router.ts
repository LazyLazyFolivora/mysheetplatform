import 'vue-router'

declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    requiresAuth?: boolean
    requiresAdmin?: boolean
    icon?: string
    hidden?: boolean
  }
}

export type RouteName =
  | 'Home'
  | 'Login'
  | 'Register'
  | 'Profile'
  | 'SheetMusic'
  | 'SheetMusicDetail'
  | 'About'
  | 'Contact'
  | 'FriendLinks'
  | 'AdminLogin'
  | 'admin-dashboard'
  | 'admin-users'
  | 'admin-sheet-music'
  | 'admin-permissions'
  | 'admin-roles'
  | 'admin-contact-messages'
  | 'admin-friend-links'
  | 'NotFound'

export type RoutePath =
  | '/'
  | '/login'
  | '/register'
  | '/profile'
  | '/sheet-music'
  | '/about'
  | '/contact'
  | '/friend-links'
  | '/admin/login'
  | '/admin'
  | '/admin/users'
  | '/admin/sheet-music'
  | '/admin/permissions'
  | '/admin/roles'
  | '/admin/contact-messages'
  | '/admin/friend-links'

export interface NavigationItem {
  name: string
  path: string
  icon?: string
  children?: NavigationItem[]
  meta?: {
    requiresAuth?: boolean
    requiresAdmin?: boolean
    hidden?: boolean
  }
}

export interface BreadcrumbItem {
  name: string
  path?: string
  disabled?: boolean
}
