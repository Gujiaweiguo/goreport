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

  it('calls list endpoint with default pagination', () => {
    datasourceApi.list()
    expect(apiClient.get).toHaveBeenCalledWith('/api/v1/datasources', {
      params: {
        page: 1,
        pageSize: 10
      }
    })
  })

  it('calls list endpoint with explicit pagination', () => {
    datasourceApi.list(2, 50)
    expect(apiClient.get).toHaveBeenCalledWith('/api/v1/datasources', {
      params: {
        page: 2,
        pageSize: 50
      }
    })
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
    expect(apiClient.post).toHaveBeenCalledWith('/api/v1/datasources', payload)
  })

  it('normalizes postgresql type to postgres for create and test', () => {
    const payload = {
      name: 'pg-demo',
      type: 'postgresql',
      host: 'localhost',
      port: 5432,
      database: 'demo',
      username: 'postgres',
      password: 'postgres'
    }

    datasourceApi.create(payload)
    datasourceApi.test(payload)

    expect(apiClient.post).toHaveBeenCalledWith('/api/v1/datasources', {
      ...payload,
      type: 'postgres'
    })
    expect(apiClient.post).toHaveBeenCalledWith('/api/v1/datasources/test', {
      ...payload,
      type: 'postgres'
    })
  })

  it('calls update and delete endpoints', () => {
    datasourceApi.update('ds-1', { name: 'new-name' })
    datasourceApi.delete('ds-1')

    expect(apiClient.put).toHaveBeenCalledWith('/api/v1/datasources/ds-1', { name: 'new-name' })
    expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/datasources/ds-1')
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
    datasourceApi.testById('ds-1')
    datasourceApi.getTables('ds-1')
    datasourceApi.getFields('ds-1', 'users')

    expect(apiClient.post).toHaveBeenCalledWith('/api/v1/datasources/test', payload)
    expect(apiClient.post).toHaveBeenCalledWith('/api/v1/datasources/ds-1/test', {})
    expect(apiClient.get).toHaveBeenCalledWith('/api/v1/datasources/ds-1/tables')
    expect(apiClient.get).toHaveBeenCalledWith('/api/v1/datasources/ds-1/tables/users/fields')
  })
})
