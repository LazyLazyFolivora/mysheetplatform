export interface User {
  id: number
  username: string
  nickname: string
  email: string
  avatar: string
  points: number
  role: number
  createTime: string
  updateTime: string
}

export interface UserDTO {
  id: number
  username: string
  nickname: string
  email: string
  role: number
  points: number
}

export interface CurrentUser {
  id: number
  username: string
  nickname: string
  email: string
  avatar: string
  role: number
  points: number
  createTime: string
}

export interface LoginDTO {
  email: string
  password: string
}

export interface LoginVO {
  token: string
  user_id: number
  username: string
  nickname: string
  is_admin: boolean
}

export interface ProfileVO {
  user_id: number
  username: string
  nickname: string
  email: string
  avatar: string
  is_admin: boolean
  points: number
  created_at: string
}

export interface RegisterDTO {
  username: string
  password: string
  email: string
  code: string
  nickname?: string
}

export interface SendCodeDTO {
  email: string
}
