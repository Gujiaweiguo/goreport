import axios from 'axios'

export interface DataSource {
  id: string
  name: string
  type: string
  host: string
  port: number
  database: string
  username: string
  tenantId: string
  createdAt: string
  updatedAt: string
}

export interface CreateDataSourceRequest {
  name: string
  type: string
  host: string
  port: number
  database: string
  username: string
  password: string
}

export interface UpdateDataSourceRequest {
  name?: string
  type?: string
  host?: string
  port?: number
  database?: string
  username?: string
  password?: string
}

export interface ApiResponse<T = any> {
  success: boolean
  result: T
  message: string
}

export const datasourceApi = {
  list: () => {
    return axios.get<ApiResponse<DataSource[]>>('/api/v1/datasource/list')
  },

  create: (data: CreateDataSourceRequest) => {
    return axios.post<ApiResponse<DataSource>>('/api/v1/datasource/create', data)
  },

  update: (id: string, data: UpdateDataSourceRequest) => {
    return axios.put<ApiResponse<DataSource>>(`/api/v1/datasource/${id}`, data)
  },

  delete: (id: string) => {
    return axios.delete<ApiResponse<null>>(`/api/v1/datasource/${id}`)
  },

  test: (data: CreateDataSourceRequest) => {
    return axios.post<ApiResponse<null>>('/api/v1/datasource/test', data)
  },

  getTables: (id: string) => {
    return axios.get<ApiResponse<string[]>>(`/api/v1/datasource/${id}/tables`)
  },

  getFields: (id: string, table: string) => {
    return axios.get<ApiResponse<any[]>>(`/api/v1/datasource/${id}/tables/${table}/fields`)
  }
}
