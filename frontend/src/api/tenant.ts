import apiClient from './client'

export interface Tenant {
  id: string
  name: string
  code: string
  status: number
  createdAt: string
  updatedAt: string
}

export interface ApiResponse<T = any> {
  success: boolean
  result: T
  message: string
}

export const tenantApi = {
  list: () => {
    return apiClient.get<ApiResponse<Tenant[]>>('/api/v1/tenants')
  },

  getCurrent: () => {
    return apiClient.get<ApiResponse<Tenant>>('/api/v1/tenants/current')
  }
}
