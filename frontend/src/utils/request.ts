import axios from 'axios'
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import router from '@/router'
import type { ApiResponse } from '@/types/api'

interface CustomAxiosInstance extends Omit<AxiosInstance, 'get' | 'post' | 'put' | 'delete'> {
  get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T>
  post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>
  put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>
  delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T>
}

const baseURL = import.meta.env.VITE_API_BASE_URL || '/api'
const request = axios.create({
  baseURL,
  timeout: 15000,
  withCredentials: true,
}) as CustomAxiosInstance

request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    // FormData 必须由运行时自动带 boundary，禁止外部写死 Content-Type
    if (typeof FormData !== 'undefined' && config.data instanceof FormData) {
      if (config.headers) {
        delete (config.headers as Record<string, unknown>)['Content-Type']
        delete (config.headers as Record<string, unknown>)['content-type']
      }
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  (response: AxiosResponse) => {
    if (response.config.responseType === 'blob') {
      return response
    }

    const res = response.data
    // Go backend: code 200 = success
    if (res.code === 200) {
      return res
    }

    // Handle 401 via business error
    if (res.code === 401) {
      const userStore = useUserStore()
      userStore.logout()
      router.push('/login')
      ElMessage.warning('请先登录')
    } else {
      ElMessage.error(res.message || '请求失败')
    }
    return Promise.reject(new Error(res.message || '请求失败'))
  },
  (error) => {
    const userStore = useUserStore()

    if (error.config?.url === '/auth/logout' && error.response?.status === 401) {
      return Promise.reject(error)
    }

    if (error.response?.status === 401) {
      userStore.logout()

      const currentPath = router.currentRoute.value.path
      if (currentPath.startsWith('/admin')) {
        router.push('/admin/login')
      } else {
        router.push('/login')
      }

      ElMessage.warning('登录已过期，请重新登录')
      return Promise.reject(error)
    }

    if (error.response) {
      const { status, data } = error.response
      switch (status) {
        case 403:
          ElMessage.error('没有权限')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error(data?.message || '服务器错误')
          break
        default:
          ElMessage.error(data?.message || '请求失败')
      }
    } else if (error.code === 'ECONNABORTED') {
      ElMessage.error('请求超时，请重试')
    } else if (error.code === 'ERR_NETWORK') {
      ElMessage.error('网络错误，请检查后端服务是否启动')
    } else {
      ElMessage.error('未知错误，请稍后重试')
    }
    return Promise.reject(error)
  }
)

export default request
