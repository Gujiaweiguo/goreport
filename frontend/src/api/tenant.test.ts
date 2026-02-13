import { beforeEach, describe, expect, it, vi } from 'vitest'
import apiClient from './client'
import { tenantApi } from './tenant'

vi.mock('./client', () => {
  return {
    default: {
      get: vi.fn()
    }
  }
})

describe('tenantApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('list', () => {
    it('calls list endpoint', () => {
      tenantApi.list()
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/tenants')
    })
  })

  describe('getCurrent', () => {
    it('calls current endpoint', () => {
      tenantApi.getCurrent()
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/tenants/current')
    })
  })
})
