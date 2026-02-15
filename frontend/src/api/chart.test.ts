import { beforeEach, describe, expect, it, vi } from 'vitest'
import apiClient from './client'
import { chartApi } from './chart'

vi.mock('./client', () => {
  return {
    default: {
      get: vi.fn(),
      post: vi.fn(),
      delete: vi.fn()
    }
  }
})

describe('chartApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('list', () => {
    it('calls list endpoint', () => {
      chartApi.list()
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/charts')
    })
  })

  describe('get', () => {
    it('calls get endpoint with id', () => {
      chartApi.get('chart-1')
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/charts/get', {
        params: { id: 'chart-1' }
      })
    })
  })

  describe('create', () => {
    it('calls create endpoint with data', () => {
      const data = {
        name: 'Test Chart',
        type: 'line',
        config: { x: 'date', y: 'value' }
      }
      chartApi.create(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/charts/create', data)
    })

    it('includes optional code field', () => {
      const data = {
        name: 'Test Chart',
        code: 'CHART_001',
        type: 'bar',
        config: {}
      }
      chartApi.create(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/charts/create', data)
    })
  })

  describe('update', () => {
    it('calls update endpoint with data', () => {
      const data = {
        id: 'chart-1',
        name: 'Updated Chart',
        config: { x: 'name', y: 'amount' }
      }
      chartApi.update(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/charts/update', data)
    })

    it('supports partial update', () => {
      const data = {
        id: 'chart-1',
        name: 'Only Name Changed'
      }
      chartApi.update(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/charts/update', data)
    })
  })

  describe('delete', () => {
    it('calls delete endpoint with id', () => {
      chartApi.delete('chart-1')
      expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/charts/delete', {
        params: { id: 'chart-1' }
      })
    })
  })

  describe('render', () => {
    it('calls render endpoint with id', () => {
      chartApi.render('chart-1')
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/charts/render', {
        params: { id: 'chart-1' }
      })
    })
  })
})
