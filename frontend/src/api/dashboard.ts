import axios from 'axios'

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
    return axios.get<ApiResponse<Dashboard[]>>('/api/v1/dashboard/list')
  },

  create: (data: CreateDashboardRequest) => {
    return axios.post<ApiResponse<Dashboard>>('/api/v1/dashboard/create', data)
  },

  update: (id: string, data: UpdateDashboardRequest) => {
    return axios.put<ApiResponse<Dashboard>>(`/api/v1/dashboard/${id}`, data)
  },

  delete: (id: string) => {
    return axios.delete<ApiResponse<null>>(`/api/v1/dashboard/${id}`)
  },

  get: (id: string) => {
    return axios.get<ApiResponse<Dashboard>>(`/api/v1/dashboard/${id}`)
  }
}
