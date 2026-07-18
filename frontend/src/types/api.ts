export interface ApiResponse<T> {
  code: number
  message: string
  result: T
}

export interface PaginatedResponse<T> {
  total: number
  list: T[]
}
