import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('axios', () => {
  const mockAxios = {
    create: vi.fn(() => mockAxiosInstance),
    interceptors: {
      request: {
        use: vi.fn()
      },
      response: {
        use: vi.fn()
      }
    }
  }
  
  const mockAxiosInstance = {
    interceptors: {
      request: {
        use: vi.fn((onFulfilled, onRejected) => {
          mockAxiosInstance._requestInterceptor = { onFulfilled, onRejected }
        })
      },
      response: {
        use: vi.fn((onFulfilled, onRejected) => {
          mockAxiosInstance._responseInterceptor = { onFulfilled, onRejected }
        })
      }
    },
    _requestInterceptor: null as any,
    _responseInterceptor: null as any
  }
  
  return {
    default: mockAxios
  }
})

describe('apiClient', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  describe('request interceptor', () => {
    it('adds token to headers when token exists', async () => {
      localStorage.setItem('token', 'test-token')
      
      const { default: apiClient } = await import('./client')
      
      const interceptor = (apiClient as any)._requestInterceptor
      expect(interceptor).toBeDefined()
      
      const config = { headers: {} }
      const result = interceptor.onFulfilled(config)
      
      expect(result.headers.Authorization).toBe('Bearer test-token')
    })

    it('does not add token when token does not exist', async () => {
      const { default: apiClient } = await import('./client')
      
      const interceptor = (apiClient as any)._requestInterceptor
      const config = { headers: {} }
      const result = interceptor.onFulfilled(config)
      
      expect(result.headers.Authorization).toBeUndefined()
    })
  })

  describe('response interceptor', () => {
    it('passes through successful responses', async () => {
      const { default: apiClient } = await import('./client')
      
      const interceptor = (apiClient as any)._responseInterceptor
      const response = { data: { success: true } }
      const result = interceptor.onFulfilled(response)
      
      expect(result).toBe(response)
    })

    it('handles 401 errors by clearing session and redirecting', async () => {
      const originalLocation = window.location
      const mockLocation = { href: '' }
      Object.defineProperty(window, 'location', {
        value: mockLocation,
        writable: true
      })

      localStorage.setItem('token', 'old-token')
      localStorage.setItem('user', JSON.stringify({ id: 'u-1' }))

      const { default: apiClient } = await import('./client')
      
      const interceptor = (apiClient as any)._responseInterceptor
      const error = {
        response: { status: 401 }
      }
      
      await expect(interceptor.onRejected(error)).rejects.toBeDefined()
      
      expect(localStorage.getItem('token')).toBeNull()
      expect(localStorage.getItem('user')).toBeNull()
      expect(mockLocation.href).toBe('/login')

      Object.defineProperty(window, 'location', {
        value: originalLocation,
        writable: true
      })
    })

    it('passes through non-401 errors', async () => {
      const { default: apiClient } = await import('./client')
      
      const interceptor = (apiClient as any)._responseInterceptor
      const error = {
        response: { status: 500 }
      }
      
      await expect(interceptor.onRejected(error)).rejects.toEqual(error)
    })

    it('handles errors without response', async () => {
      const { default: apiClient } = await import('./client')
      
      const interceptor = (apiClient as any)._responseInterceptor
      const error = new Error('Network error')
      
      await expect(interceptor.onRejected(error)).rejects.toEqual(error)
    })
  })
})
