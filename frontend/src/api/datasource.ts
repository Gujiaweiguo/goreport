import apiClient from './client'

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

export interface SelectOption {
  label: string
  value: string
}

const normalizeDatasourceType = (type?: string) => {
  if (type === 'postgresql') {
    return 'postgres'
  }
  return type
}

const normalizeDatasourcePayload = <T extends { type?: string }>(data: T): T => {
  return {
    ...data,
    type: normalizeDatasourceType(data.type)
  }
}

export const datasourceApi = {
  list: (page?: number, pageSize?: number) => {
    return apiClient.get<ApiResponse<{datasources: DataSource[], total: number, page: number, pageSize: number}>>('/api/v1/datasources', {
      params: {
        page: page || 1,
        pageSize: pageSize || 10
      }
    })
  },

  create: (data: CreateDataSourceRequest) => {
    return apiClient.post<ApiResponse<DataSource>>('/api/v1/datasources', normalizeDatasourcePayload(data))
  },

  update: (id: string, data: UpdateDataSourceRequest) => {
    return apiClient.put<ApiResponse<DataSource>>(`/api/v1/datasources/${id}`, normalizeDatasourcePayload(data))
  },

  delete: (id: string) => {
    return apiClient.delete<ApiResponse<null>>(`/api/v1/datasources/${id}`)
  },

  test: (data: CreateDataSourceRequest) => {
    return apiClient.post<ApiResponse<null>>('/api/v1/datasources/test', normalizeDatasourcePayload(data))
  },

  testById: (id: string) => {
    return apiClient.post<ApiResponse<null>>(`/api/v1/datasources/${id}/test`, {})
  },

  getTables: (id: string) => {
    return apiClient.get<ApiResponse<string[]>>(`/api/v1/datasources/${id}/tables`)
  },

  getFields: (id: string, table: string) => {
    return apiClient.get<ApiResponse<any[]>>(`/api/v1/datasources/${id}/tables/${table}/fields`)
  },

  copy: (id: string) => {
    return apiClient.post<ApiResponse<DataSource>>(`/api/v1/datasources/copy/${id}`, {})
  },

  move: (id: string, target: string) => {
    return apiClient.post<ApiResponse<null>>(`/api/v1/datasources/move`, { id, target })
  },

  rename: (id: string, name: string) => {
    return apiClient.put<ApiResponse<DataSource>>(`/api/v1/datasources/${id}/rename`, { name })
  },

  search: (keyword: string, page?: number, pageSize?: number) => {
    return apiClient.get<ApiResponse<{datasources: DataSource[], total: number, page: number, pageSize: number}>>(`/api/v1/datasources/search`, {
      params: {
        keyword,
        page: page || 1,
        pageSize: pageSize || 10
      }
    })
  },

  listProfiles: () => {
    return apiClient.get<ApiResponse<any[]>>('/api/v1/datasources/profiles')
  }
}
