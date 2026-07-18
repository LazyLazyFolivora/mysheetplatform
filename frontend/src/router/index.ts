import { routes } from './routes'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const userInfoStr = localStorage.getItem('userInfo')
  let userInfo: any = null

  if (userInfoStr && userInfoStr !== 'undefined') {
    try {
      userInfo = JSON.parse(userInfoStr)
    } catch {
      // ignore parse errors
    }
  }

  const isLoggedIn = !!token
  const isAdmin = userInfo?.role === 1

  if (to.meta.title) {
    document.title = `${to.meta.title} - Folivora懒懒`
  }

  // Redirect logged-in users away from login pages
  if (isLoggedIn && (to.path === '/login' || to.path === '/admin/login')) {
    next('/')
    return
  }

  // Check admin access
  if (to.meta.requiresAdmin && !isAdmin) {
    next('/')
    return
  }

  // Check auth
  if (to.meta.requiresAuth && !isLoggedIn) {
    next('/login')
    return
  }

  next()
})

export default router
