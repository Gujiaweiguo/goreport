import apiClient from './client'

export interface Chart {
  id: string
  tenantId: string
  name: string
  code: string
  type: string
  config: string
  status: number
  createdAt: string
  updatedAt: string
}

export interface CreateChartRequest {
  name: string
  code?: string
  type: string
  config: any
}

export interface UpdateChartRequest {
  id: string
  name?: string
  code?: string
  type?: string
  config?: any
}

export interface ApiResponse<T = any> {
  success: boolean
  result: T
  message: string
}

export const chartApi = {
  list: () => {
    return apiClient.get<ApiResponse<Chart[]>>('/api/v1/charts')
  },

  get: (id: string) => {
    return apiClient.get<ApiResponse<Chart>>('/api/v1/charts/get', {
      params: { id }
    })
  },

  create: (data: CreateChartRequest) => {
    return apiClient.post<ApiResponse<Chart>>('/api/v1/charts/create', data)
  },

  update: (data: UpdateChartRequest) => {
    return apiClient.post<ApiResponse<Chart>>('/api/v1/charts/update', data)
  },

  delete: (id: string) => {
    return apiClient.delete<ApiResponse<null>>('/api/v1/charts/delete', {
      params: { id }
    })
  },

  render: (id: string) => {
    return apiClient.get<ApiResponse<any>>('/api/v1/charts/render', {
      params: { id }
    })
  }
}
