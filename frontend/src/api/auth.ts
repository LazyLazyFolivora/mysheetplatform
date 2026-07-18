import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'
import type { LoginDTO, LoginVO, ProfileVO, RegisterDTO, SendCodeDTO } from '@/types/user'

export const authApi = {
  login: (data: LoginDTO) => {
    return request.post<ApiResponse<LoginVO>>('/auth/login', data)
  },

  register: (data: RegisterDTO) => {
    return request.post<ApiResponse<LoginVO>>('/auth/register', data)
  },

  sendCode: (data: SendCodeDTO) => {
    return request.post<ApiResponse<string>>('/auth/send-code', data)
  },

  getProfile: () => {
    return request.get<ApiResponse<ProfileVO>>('/profile')
  },
}
