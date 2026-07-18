import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { CurrentUser, LoginVO, RegisterDTO } from '@/types/user'
import { authApi } from '@/api/auth'
import { ElMessage } from 'element-plus'

// Go 后端返回扁平的 LoginVO（is_admin 布尔值），前端统一用 role 数字表示
function toCurrentUser(vo: LoginVO): CurrentUser {
  return {
    id: vo.user_id,
    username: vo.username,
    nickname: vo.nickname,
    email: '',
    avatar: '',
    role: vo.is_admin ? 1 : 0,
    points: 0,
    createTime: '',
  }
}

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<CurrentUser | null>(null)

  const isLoggedIn = computed(() => !!token.value)
  const points = computed(() => userInfo.value?.points || 0)
  const isAdmin = computed(() => userInfo.value?.role === 1)

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info: CurrentUser) => {
    userInfo.value = info
    try {
      localStorage.setItem('userInfo', JSON.stringify(info))
    } catch (error) {
      console.error('保存用户信息到localStorage失败:', error)
    }
  }

  const login = async (email: string, password: string) => {
    try {
      const response = await authApi.login({ email, password })
      if (response.code === 200 && response.result) {
        setToken(response.result.token)
        setUserInfo(toCurrentUser(response.result))
        return true
      }
      return false
    } catch (error: any) {
      if (error.response?.data?.message) {
        throw new Error(error.response.data.message)
      }
      throw new Error(error.message || '登录失败，请稍后重试')
    }
  }

  const register = async (data: RegisterDTO) => {
    try {
      const response = await authApi.register(data)
      if (response.code === 200 && response.result) {
        setToken(response.result.token)
        setUserInfo(toCurrentUser(response.result))
        ElMessage.success('注册成功')
        return true
      }
      ElMessage.error(response.message || '注册失败')
      return false
    } catch (error: any) {
      if (error.response?.data?.message) {
        ElMessage.error(error.response.data.message)
      } else {
        ElMessage.error(error.message || '注册失败，请稍后重试')
      }
      return false
    }
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('userInfo')
  }

  const initUserInfo = () => {
    const storedUserInfo = localStorage.getItem('userInfo')
    if (storedUserInfo && storedUserInfo !== 'undefined') {
      try {
        userInfo.value = JSON.parse(storedUserInfo)
      } catch (error) {
        console.error('解析localStorage中的用户信息失败:', error)
        localStorage.removeItem('userInfo')
      }
    }
  }

  const fetchProfile = async () => {
    try {
      const res = await authApi.getProfile()
      if (res.code === 200 && res.result) {
        const vo = res.result
        const updated: CurrentUser = {
          id: vo.user_id,
          username: vo.username,
          nickname: vo.nickname,
          email: vo.email || '',
          avatar: vo.avatar || '',
          role: vo.is_admin ? 1 : 0,
          points: vo.points,
          createTime: vo.created_at,
        }
        setUserInfo(updated)
      }
    } catch {
      // 静默失败，不影响下载流程
    }
  }

  initUserInfo()

  return {
    token,
    userInfo,
    isLoggedIn,
    isAdmin,
    points,
    setToken,
    setUserInfo,
    login,
    register,
    logout,
    initUserInfo,
    fetchProfile,
  }
})
