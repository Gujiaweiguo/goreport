import apiClient from './client'

export interface Dataset {
  id: string
  tenantId: string
  datasourceId?: string
  name: string
  type: 'sql' | 'api' | 'file'
  config: any
  status: number
  createdBy: string
  createdAt: string
  updatedAt: string
  fields?: DatasetField[]
}

export interface DatasetField {
  id: string
  datasetId: string
  name: string
  displayName?: string
  type: 'dimension' | 'measure'
  dataType: 'string' | 'number' | 'date' | 'boolean'
  isComputed: boolean
  expression?: string
  isSortable: boolean
  isGroupable: boolean
  defaultSortOrder: 'asc' | 'desc' | 'none'
  sortIndex: number
  config?: any
  createdAt: string
  updatedAt: string
  isGroupingField?: boolean
  groupingRule?: string
  groupingEnabled?: boolean
}

export interface DatasetSource {
  id: string
  datasetId: string
  sourceType: 'datasource' | 'api' | 'file'
  sourceId?: string
  sourceConfig: any
  joinType: 'inner' | 'left' | 'right' | 'full'
  joinCondition?: string
  sortIndex: number
  createdAt: string
  updatedAt: string
}

export interface DatasetSchema {
  dimensions: DatasetField[]
  measures: DatasetField[]
  computed: DatasetField[]
}

export interface CreateDatasetRequest {
  name: string
  type: 'sql' | 'api' | 'file'
  datasourceId?: string
  config: any
}

export interface UpdateDatasetRequest {
  name?: string
  config?: any
  status?: number
  action?: 'save' | 'save_and_return'
}

export interface CreateFieldRequest {
  name: string
  displayName?: string
  type: 'dimension' | 'measure'
  dataType: 'string' | 'number' | 'date' | 'boolean'
  expression?: string
}

export interface UpdateFieldRequest {
  displayName?: string
  type?: 'dimension' | 'measure'
  dataType?: 'string' | 'number' | 'date' | 'boolean'
  isSortable?: boolean
  isGroupable?: boolean
  sortOrder?: 'asc' | 'desc' | 'none'
  isGroupingField?: boolean
  groupingRule?: string
  groupingEnabled?: boolean
}

export interface BatchUpdateFieldRequest extends UpdateFieldRequest {
  fieldId: string
}

export interface BatchUpdateFieldsRequest {
  fields: BatchUpdateFieldRequest[]
}

export interface BatchUpdateFieldsResponse {
  success: boolean
  updatedFields: string[]
  errors: FieldError[]
}

export interface FieldError {
  fieldId: string
  message: string
}

export const formatBatchFieldErrors = (errors: FieldError[]): string => {
  if (!errors.length) {
    return ''
  }

  return errors
    .map((item) => {
      const fieldId = item.fieldId || 'unknown'
      const message = item.message || '未知错误'
      return `${fieldId}: ${message}`
    })
    .join('；')
}

interface ErrorPayload {
  message?: string
}

interface HttpErrorLike {
  isAxiosError?: boolean
  message?: string
  response?: {
    data?: ErrorPayload
  }
  request?: unknown
}

export const getApiErrorMessage = (
  error: unknown,
  fallbackMessage: string,
  transportFallbackMessage: string = '网络连接失败，请检查网络后重试'
): string => {
  if (!error || typeof error !== 'object') {
    return fallbackMessage
  }

  const maybeError = error as HttpErrorLike
  if (maybeError.response) {
    return maybeError.response.data?.message || maybeError.message || fallbackMessage
  }

  if (maybeError.isAxiosError && maybeError.request) {
    return transportFallbackMessage
  }

  return maybeError.message || fallbackMessage
}

export interface QueryRequest {
  fields?: string[]
  filters?: Filter[]
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
  page?: number
  pageSize?: number
  groupBy?: string[]
  aggregations?: Record<string, Aggregation>
}

export interface Filter {
  field: string
  operator: 'eq' | 'neq' | 'gt' | 'gte' | 'lt' | 'lte' | 'like' | 'in'
  value: any
}

export interface Aggregation {
  function: 'SUM' | 'AVG' | 'COUNT' | 'MAX' | 'MIN'
  field: string
}

export interface QueryResponse {
  data: Record<string, any>[]
  total: number
  page: number
  pageSize: number
  executionTime: number
  aggregations?: Record<string, any>
}

export interface ApiResponse<T = any> {
  success: boolean
  result: T
  message: string
  total?: number
  page?: number
  pageSize?: number
}

export const datasetApi = {
  list: (page: number = 1, pageSize: number = 10) => {
    return apiClient.get<ApiResponse<Dataset[]>>(`/api/v1/datasets?page=${page}&pageSize=${pageSize}`)
  },

  get: (id: string) => {
    return apiClient.get<ApiResponse<Dataset>>(`/api/v1/datasets/${id}`)
  },

  create: (data: CreateDatasetRequest) => {
    return apiClient.post<ApiResponse<Dataset>>('/api/v1/datasets', data)
  },

  update: (id: string, data: UpdateDatasetRequest) => {
    return apiClient.put<ApiResponse<Dataset>>(`/api/v1/datasets/${id}`, data)
  },

  delete: (id: string) => {
    return apiClient.delete<ApiResponse<null>>(`/api/v1/datasets/${id}`)
  },

  preview: (id: string) => {
    return apiClient.get<ApiResponse<Record<string, any>[]>>(`/api/v1/datasets/${id}/preview`)
  },

  query: (id: string, query: QueryRequest) => {
    return apiClient.post<ApiResponse<QueryResponse>>(`/api/v1/datasets/${id}/data`, query)
  },

  getDimensions: (id: string) => {
    return apiClient.get<ApiResponse<DatasetField[]>>(`/api/v1/datasets/${id}/dimensions`)
  },

  getMeasures: (id: string) => {
    return apiClient.get<ApiResponse<DatasetField[]>>(`/api/v1/datasets/${id}/measures`)
  },

  getSchema: (id: string) => {
    return apiClient.get<ApiResponse<DatasetSchema>>(`/api/v1/datasets/${id}/schema`)
  },

  createField: (id: string, data: CreateFieldRequest) => {
    return apiClient.post<ApiResponse<DatasetField>>(`/api/v1/datasets/${id}/fields`, data)
  },

  updateField: (id: string, fieldId: string, data: UpdateFieldRequest) => {
    return apiClient.put<ApiResponse<DatasetField>>(`/api/v1/datasets/${id}/fields/${fieldId}`, data)
  },

  deleteField: (id: string, fieldId: string) => {
    return apiClient.delete<ApiResponse<null>>(`/api/v1/datasets/${id}/fields/${fieldId}`)
  },

  batchUpdateFields: (id: string, data: BatchUpdateFieldsRequest) => {
    return apiClient.patch<ApiResponse<BatchUpdateFieldsResponse>>(`/api/v1/datasets/${id}/fields`, data)
  }
}
