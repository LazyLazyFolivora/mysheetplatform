import axios from 'axios'
import type { AxiosRequestConfig } from 'axios'
import type { ApiResponse } from '@/types/api'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => response.data,
  (error) => Promise.reject(error)
)

export function get<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
  return api.get(url, config).then((data) => data as unknown as ApiResponse<T>)
}

export function post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
  return api.post(url, data, config).then((data) => data as unknown as ApiResponse<T>)
}

export function put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
  return api.put(url, data, config).then((data) => data as unknown as ApiResponse<T>)
}

export function del<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
  return api.delete(url, config).then((data) => data as unknown as ApiResponse<T>)
}
