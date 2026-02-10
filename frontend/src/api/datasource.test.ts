import { beforeEach, describe, expect, it, vi } from 'vitest'
import { datasourceApi } from './datasource'
import apiClient from './client'

vi.mock('./client', () => {
  return {
    default: {
      get: vi.fn(),
      post: vi.fn(),
      put: vi.fn(),
      delete: vi.fn()
    }
  }
})

describe('datasourceApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('calls list endpoint', () => {
    datasourceApi.list()
    expect(apiClient.get).toHaveBeenCalledWith('/api/v1/datasource/list')
  })

  it('calls create endpoint', () => {
    const payload = {
      name: 'demo',
      type: 'mysql',
      host: 'localhost',
      port: 3306,
      database: 'goreport',
      username: 'root',
      password: 'root'
    }
    datasourceApi.create(payload)
    expect(apiClient.post).toHaveBeenCalledWith('/api/v1/datasource/create', payload)
  })

  it('calls update endpoint', () => {
    datasourceApi.update('ds-1', { name: 'new-name' })
    expect(apiClient.put).toHaveBeenCalledWith('/api/v1/datasource/ds-1', { name: 'new-name' })
  })

  it('calls delete endpoint', () => {
    datasourceApi.delete('ds-1')
    expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/datasource/ds-1')
  })

  it('calls test and metadata endpoints', () => {
    const payload = {
      name: 'demo',
      type: 'mysql',
      host: 'localhost',
      port: 3306,
      database: 'goreport',
      username: 'root',
      password: 'root'
    }

    datasourceApi.test(payload)
    datasourceApi.getTables('ds-1')
    datasourceApi.getFields('ds-1', 'users')

    expect(apiClient.post).toHaveBeenCalledWith('/api/v1/datasource/test', payload)
    expect(apiClient.get).toHaveBeenCalledWith('/api/v1/datasource/ds-1/tables')
    expect(apiClient.get).toHaveBeenCalledWith('/api/v1/datasource/ds-1/tables/users/fields')
  })
})
