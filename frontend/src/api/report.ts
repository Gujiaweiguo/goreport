import apiClient from './client'

export interface Report {
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

export interface CreateReportRequest {
  name: string
  code?: string
  type?: string
  config: any // JSON 格式的报表配置
}

export interface UpdateReportRequest {
  id: string
  name?: string
  code?: string
  type?: string
  config?: any // JSON 格式的报表配置
}

export interface PreviewReportRequest {
  id: string
  params?: Record<string, any>
}

export interface PreviewReportResponse {
  html: string
}

export interface ApiResponse<T = any> {
  success: boolean
  result: T
  message: string
}

export const reportApi = {
  list: () => {
    return apiClient.get<ApiResponse<Report[]>>('/api/v1/jmreport/list')
  },

  get: (id: string) => {
    return apiClient.get<ApiResponse<Report>>('/api/v1/jmreport/get', {
      params: { id }
    })
  },

  create: (data: CreateReportRequest) => {
    return apiClient.post<ApiResponse<Report>>('/api/v1/jmreport/create', data)
  },

  update: (data: UpdateReportRequest) => {
    return apiClient.post<ApiResponse<Report>>('/api/v1/jmreport/update', data)
  },

  delete: (id: string) => {
    return apiClient.delete<ApiResponse<null>>('/api/v1/jmreport/delete', {
      params: { id }
    })
  },

  preview: (data: PreviewReportRequest) => {
    return apiClient.post<ApiResponse<PreviewReportResponse>>('/api/v1/jmreport/preview', data)
  }
}
