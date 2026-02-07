import apiClient from './client'

export interface Dashboard {
  id: string
  name: string
  description: string
  config: {
    width: number
    height: number
    backgroundColor: string
  }
  components: any[]
  createdAt: string
  updatedAt: string
}

export interface CreateDashboardRequest {
  name: string
  description: string
  config: {
    width: number
    height: number
    backgroundColor: string
  }
  components: any[]
}

export interface UpdateDashboardRequest {
  name?: string
  description?: string
  config?: {
    width?: number
    height?: number
    backgroundColor?: string
  }
  components?: any[]
}

export interface ApiResponse<T = any> {
  success: boolean
  result: T
  message: string
}

export const dashboardApi = {
  list: () => {
    return apiClient.get<ApiResponse<Dashboard[]>>('/api/v1/dashboard/list')
  },

  create: (data: CreateDashboardRequest) => {
    return apiClient.post<ApiResponse<Dashboard>>('/api/v1/dashboard/create', data)
  },

  update: (id: string, data: UpdateDashboardRequest) => {
    return apiClient.put<ApiResponse<Dashboard>>(`/api/v1/dashboard/${id}`, data)
  },

  delete: (id: string) => {
    return apiClient.delete<ApiResponse<null>>(`/api/v1/dashboard/${id}`)
  },

  get: (id: string) => {
    return apiClient.get<ApiResponse<Dashboard>>(`/api/v1/dashboard/${id}`)
  }
}
