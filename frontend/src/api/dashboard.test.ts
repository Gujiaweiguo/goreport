import { beforeEach, describe, expect, it, vi } from 'vitest'
import apiClient from './client'
import { dashboardApi } from './dashboard'

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

describe('dashboardApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('list', () => {
    it('calls list endpoint', () => {
      dashboardApi.list()
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/dashboard/list')
    })
  })

  describe('get', () => {
    it('calls get endpoint with id', () => {
      dashboardApi.get('dashboard-1')
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/dashboard/dashboard-1')
    })
  })

  describe('create', () => {
    it('calls create endpoint with data', () => {
      const data = {
        name: 'Test Dashboard',
        description: 'A test dashboard',
        config: {
          width: 1920,
          height: 1080,
          backgroundColor: '#ffffff'
        },
        components: []
      }
      dashboardApi.create(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/dashboard/create', data)
    })

    it('supports dashboard with components', () => {
      const data = {
        name: 'Dashboard with Charts',
        description: 'Contains multiple charts',
        config: {
          width: 1920,
          height: 1080,
          backgroundColor: '#f0f0f0'
        },
        components: [
          { id: 'chart-1', type: 'chart', x: 0, y: 0, width: 400, height: 300 }
        ]
      }
      dashboardApi.create(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/dashboard/create', data)
    })
  })

  describe('update', () => {
    it('calls update endpoint with id and data', () => {
      const data = {
        name: 'Updated Dashboard',
        description: 'Updated description'
      }
      dashboardApi.update('dashboard-1', data)
      expect(apiClient.put).toHaveBeenCalledWith('/api/v1/dashboard/dashboard-1', data)
    })

    it('supports partial config update', () => {
      const data = {
        config: {
          backgroundColor: '#000000'
        }
      }
      dashboardApi.update('dashboard-1', data)
      expect(apiClient.put).toHaveBeenCalledWith('/api/v1/dashboard/dashboard-1', data)
    })

    it('supports updating components', () => {
      const data = {
        components: [
          { id: 'chart-1', type: 'chart', x: 100, y: 100, width: 500, height: 400 }
        ]
      }
      dashboardApi.update('dashboard-1', data)
      expect(apiClient.put).toHaveBeenCalledWith('/api/v1/dashboard/dashboard-1', data)
    })
  })

  describe('delete', () => {
    it('calls delete endpoint with id', () => {
      dashboardApi.delete('dashboard-1')
      expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/dashboard/dashboard-1')
    })
  })
})
