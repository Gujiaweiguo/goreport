import { beforeEach, describe, expect, it, vi } from 'vitest'
import apiClient from './client'
import { reportApi } from './report'

vi.mock('./client', () => {
  return {
    default: {
      get: vi.fn(),
      post: vi.fn(),
      delete: vi.fn()
    }
  }
})

describe('reportApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('list', () => {
    it('calls list endpoint', () => {
      reportApi.list()
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/jmreport/list')
    })
  })

  describe('get', () => {
    it('calls get endpoint with id', () => {
      reportApi.get('report-1')
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/jmreport/get', {
        params: { id: 'report-1' }
      })
    })
  })

  describe('create', () => {
    it('calls create endpoint with data', () => {
      const data = {
        name: 'Test Report',
        config: { template: 'default' }
      }
      reportApi.create(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/jmreport/create', data)
    })

    it('includes optional fields', () => {
      const data = {
        name: 'Full Report',
        code: 'RPT_001',
        type: 'table',
        config: { columns: ['id', 'name', 'amount'] }
      }
      reportApi.create(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/jmreport/create', data)
    })
  })

  describe('update', () => {
    it('calls update endpoint with data', () => {
      const data = {
        id: 'report-1',
        name: 'Updated Report',
        config: { template: 'custom' }
      }
      reportApi.update(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/jmreport/update', data)
    })

    it('supports partial update', () => {
      const data = {
        id: 'report-1',
        name: 'Only Name Changed'
      }
      reportApi.update(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/jmreport/update', data)
    })
  })

  describe('delete', () => {
    it('calls delete endpoint with id', () => {
      reportApi.delete('report-1')
      expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/jmreport/delete', {
        params: { id: 'report-1' }
      })
    })
  })

  describe('preview', () => {
    it('calls preview endpoint with id', () => {
      const data = { id: 'report-1' }
      reportApi.preview(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/jmreport/preview', data)
    })

    it('supports preview with params', () => {
      const data = {
        id: 'report-1',
        params: { startDate: '2024-01-01', endDate: '2024-12-31' }
      }
      reportApi.preview(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/jmreport/preview', data)
    })
  })

  describe('export', () => {
    it('calls export endpoint with format', () => {
      const data = {
        id: 'report-1',
        format: 'pdf'
      }
      reportApi.export(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/jmreport/export', data)
    })

    it('supports different formats', () => {
      const formats = ['pdf', 'excel', 'csv', 'word']
      formats.forEach(format => {
        const data = { id: 'report-1', format }
        reportApi.export(data)
        expect(apiClient.post).toHaveBeenCalledWith('/api/v1/jmreport/export', data)
      })
    })

    it('supports export with params', () => {
      const data = {
        id: 'report-1',
        format: 'excel',
        params: { department: 'sales' }
      }
      reportApi.export(data)
      expect(apiClient.post).toHaveBeenCalledWith('/api/v1/jmreport/export', data)
    })
  })
})
