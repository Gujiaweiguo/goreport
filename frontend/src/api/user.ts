import apiClient from './client'

export interface UserInfo {
  id: string
  username: string
  role: string
  tenantId: string
}

export interface ApiResponse<T = any> {
  success: boolean
  result: T
  message: string
}

export const userApi = {
  getMe: () => {
    return apiClient.get<ApiResponse<UserInfo>>('/api/v1/users/me')
  }
}
