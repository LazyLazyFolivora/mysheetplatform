import { useRouter } from 'vue-router'
import type { RouteName, RoutePath } from '@/types/router'

export function useRouterUtils() {
  const router = useRouter()

  const navigateTo = (name: RouteName, params?: Record<string, any>) => {
    return router.push({ name, params })
  }

  const navigateToPath = (path: RoutePath) => {
    return router.push(path)
  }

  const replaceTo = (name: RouteName, params?: Record<string, any>) => {
    return router.replace({ name, params })
  }

  const replaceToPath = (path: RoutePath) => {
    return router.replace(path)
  }

  const goBack = () => {
    return router.go(-1)
  }

  const getCurrentRoute = () => {
    return router.currentRoute.value
  }

  const isCurrentRoute = (name: RouteName) => {
    return getCurrentRoute().name === name
  }

  const isCurrentPath = (path: RoutePath) => {
    return getCurrentRoute().path === path
  }

  return {
    navigateTo,
    navigateToPath,
    replaceTo,
    replaceToPath,
    goBack,
    getCurrentRoute,
    isCurrentRoute,
    isCurrentPath,
  }
}

export function getBreadcrumbs(route: any) {
  const breadcrumbs = []
  if (route.meta?.title) {
    breadcrumbs.push({ name: route.meta.title, path: route.path })
  }
  return breadcrumbs
}

export function requiresAuth(route: any): boolean {
  return route.meta?.requiresAuth === true
}

export function requiresAdmin(route: any): boolean {
  return route.meta?.requiresAdmin === true
}

export function getRouteTitle(route: any): string {
  return route.meta?.title || 'Folivora懒懒'
}
