import axios from 'axios'

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
    return axios.get<ApiResponse<Report[]>>('/api/v1/jmreport/list')
  },

  get: (id: string) => {
    return axios.get<ApiResponse<Report>>('/api/v1/jmreport/get', {
      params: { id }
    })
  },

  create: (data: CreateReportRequest) => {
    return axios.post<ApiResponse<Report>>('/api/v1/jmreport/create', data)
  },

  update: (data: UpdateReportRequest) => {
    return axios.post<ApiResponse<Report>>('/api/v1/jmreport/update', data)
  },

  delete: (id: string) => {
    return axios.delete<ApiResponse<null>>('/api/v1/jmreport/delete', {
      params: { id }
    })
  },

  preview: (data: PreviewReportRequest) => {
    return axios.post<ApiResponse<PreviewReportResponse>>('/api/v1/jmreport/preview', data)
  }
}
