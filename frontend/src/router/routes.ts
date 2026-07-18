import { RouteRecordRaw } from 'vue-router'
import DefaultLayout from '../layouts/DefaultLayout.vue'

export const routes: RouteRecordRaw[] = [
  {
    path: '/',
    component: DefaultLayout,
    children: [
      {
        path: '',
        redirect: '/sheet-music',
      },
      {
        path: 'sheet-music',
        name: 'SheetMusic',
        component: () => import('../views/SheetMusicView.vue'),
        meta: { title: '乐谱列表' },
      },
      {
        path: 'sheet-music/:id',
        name: 'SheetMusicDetail',
        component: () => import('../views/SheetMusicDetail.vue'),
        meta: { title: '乐谱详情' },
      },
      {
        path: 'about',
        name: 'About',
        component: () => import('../views/AboutView.vue'),
        meta: { title: '关于作者' },
      },
      {
        path: 'contact',
        name: 'Contact',
        component: () => import('../views/ContactView.vue'),
        meta: { title: '联系作者' },
      },
      {
        path: 'friend-links',
        name: 'FriendLinks',
        component: () => import('../views/FriendLinksView.vue'),
        meta: { title: '友情链接' },
      },
      {
        path: 'login',
        name: 'Login',
        component: () => import('../views/LoginView.vue'),
        meta: { title: '用户登录' },
      },
      {
        path: 'register',
        name: 'Register',
        component: () => import('../views/RegisterView.vue'),
        meta: { title: '用户注册' },
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('../views/ProfileView.vue'),
        meta: { requiresAuth: true, title: '个人资料' },
      },
    ],
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('../views/admin/AdminLoginView.vue'),
    meta: { title: '管理员登录' },
  },
  {
    path: '/admin',
    component: () => import('../views/admin/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        name: 'admin-dashboard',
        component: () => import('../views/admin/DashboardView.vue'),
        meta: { title: '管理员仪表板' },
      },
      {
        path: 'users',
        name: 'admin-users',
        component: () => import('../views/admin/UsersView.vue'),
        meta: { title: '用户管理' },
      },
      {
        path: 'sheet-music',
        name: 'admin-sheet-music',
        component: () => import('../views/admin/SheetMusicView.vue'),
        meta: { title: '乐谱管理' },
      },
      {
        path: 'permissions',
        name: 'admin-permissions',
        component: () => import('../views/admin/PermissionManagement.vue'),
        meta: { title: '权限管理' },
      },
      {
        path: 'roles',
        name: 'admin-roles',
        component: () => import('../views/admin/RoleManagement.vue'),
        meta: { title: '角色管理' },
      },
      {
        path: 'contact-messages',
        name: 'admin-contact-messages',
        component: () => import('../views/admin/ContactMessages.vue'),
        meta: { title: '联系消息' },
      },
      {
        path: 'friend-links',
        name: 'admin-friend-links',
        component: () => import('../views/admin/FriendLinks.vue'),
        meta: { title: '友情链接管理' },
      },
      {
        path: 'banned-ip',
        name: 'admin-banned-ip',
        component: () => import('../views/admin/BannedIp.vue'),
        meta: { title: 'IP封禁管理' },
      },
      {
        path: 'system-config',
        name: 'admin-system-config',
        component: () => import('../views/admin/SystemConfig.vue'),
        meta: { title: '系统配置管理' },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('../views/NotFound.vue'),
    meta: { title: '页面未找到' },
  },
]
