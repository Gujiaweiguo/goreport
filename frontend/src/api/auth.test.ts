import { beforeEach, describe, expect, it, vi } from 'vitest'
import axios from 'axios'
import apiClient from './client'
import { auth, authApi } from './auth'

vi.mock('axios', () => {
  return {
    default: {
      post: vi.fn()
    }
  }
})

vi.mock('./client', () => {
  return {
    default: {
      post: vi.fn()
    }
  }
})

describe('authApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('calls login endpoint', () => {
    authApi.login({ username: 'admin', password: 'admin123' })
    expect(axios.post).toHaveBeenCalledWith('/api/v1/auth/login', {
      username: 'admin',
      password: 'admin123'
    })
  })

  it('calls logout endpoint', () => {
    authApi.logout()
    expect(apiClient.post).toHaveBeenCalledWith('/api/v1/auth/logout')
  })
})

describe('auth session helper', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('persists token/user after successful login', async () => {
    vi.mocked(axios.post).mockResolvedValueOnce({
      data: {
        success: true,
        result: {
          token: 'jwt-token',
          user: {
            id: 'u-1',
            username: 'admin',
            role: 'admin',
            tenantId: 'default'
          }
        },
        message: 'ok'
      }
    })

    const result = await auth.login('admin', 'admin123')

    expect(result.success).toBe(true)
    expect(auth.getToken()).toBe('jwt-token')
    expect(auth.isAuthenticated()).toBe(true)
    expect(auth.getUser()).toMatchObject({ id: 'u-1', tenantId: 'default' })
  })

  it('clears local session after logout', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 'u-1' }))
    vi.mocked(apiClient.post).mockResolvedValueOnce({ data: { success: true } })

    await auth.logout()

    expect(auth.getToken()).toBe('')
    expect(auth.getUser()).toBeNull()
    expect(auth.isAuthenticated()).toBe(false)
  })

  it('supports explicit clearSession', () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 'u-1' }))

    auth.clearSession()

    expect(auth.getToken()).toBe('')
    expect(auth.getUser()).toBeNull()
  })
})
