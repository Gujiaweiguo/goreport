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

  it('calls login endpoint with different credentials', () => {
    authApi.login({ username: 'testuser', password: 'testpass' })
    expect(axios.post).toHaveBeenCalledWith('/api/v1/auth/login', {
      username: 'testuser',
      password: 'testpass'
    })
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

  it('returns failure when login response is not successful', async () => {
    vi.mocked(axios.post).mockResolvedValueOnce({
      data: {
        success: false,
        result: null,
        message: 'Invalid credentials'
      }
    })

    const result = await auth.login('wrong', 'wrong')

    expect(result.success).toBe(false)
    expect(result.message).toBe('Invalid credentials')
    expect(auth.isAuthenticated()).toBe(false)
  })

  it('clears session even when logout API fails', async () => {
    localStorage.setItem('token', 'jwt-token')
    localStorage.setItem('user', JSON.stringify({ id: 'u-1' }))
    vi.mocked(apiClient.post).mockRejectedValueOnce(new Error('Network error'))

    await auth.logout()

    expect(auth.getToken()).toBe('')
    expect(auth.getUser()).toBeNull()
  })

  it('getUser returns null when user data is invalid JSON', () => {
    localStorage.setItem('user', 'invalid-json')

    const user = auth.getUser()

    expect(user).toBeNull()
  })

  it('getUser returns null when no user data stored', () => {
    const user = auth.getUser()
    expect(user).toBeNull()
  })

  it('getUser returns parsed user object when valid', () => {
    const testUser = { id: 'u-2', username: 'test', role: 'user', tenantId: 't-1' }
    localStorage.setItem('user', JSON.stringify(testUser))

    const user = auth.getUser()

    expect(user).toMatchObject(testUser)
  })

  it('isAuthenticated returns false when no token', () => {
    expect(auth.isAuthenticated()).toBe(false)
  })

  it('isAuthenticated returns true when token exists', () => {
    localStorage.setItem('token', 'some-token')
    expect(auth.isAuthenticated()).toBe(true)
  })

  it('getToken returns empty string when no token', () => {
    expect(auth.getToken()).toBe('')
  })

  it('getToken returns token when exists', () => {
    localStorage.setItem('token', 'my-token')
    expect(auth.getToken()).toBe('my-token')
  })

  it('clearSession does not throw when storage is empty', () => {
    expect(() => auth.clearSession()).not.toThrow()
  })
})

describe('auth edge cases', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    localStorage.clear()
  })

  it('handles login with empty username and password', async () => {
    vi.mocked(axios.post).mockResolvedValueOnce({
      data: {
        success: false,
        result: null,
        message: 'Username and password required'
      }
    })

    const result = await auth.login('', '')

    expect(result.success).toBe(false)
  })

  it('handles login with special characters in credentials', async () => {
    vi.mocked(axios.post).mockResolvedValueOnce({
      data: {
        success: true,
        result: {
          token: 'jwt-token',
          user: {
            id: 'u-1',
            username: 'user@test.com',
            role: 'user',
            tenantId: 'default'
          }
        },
        message: 'ok'
      }
    })

    const result = await auth.login('user@test.com', 'p@ss!w0rd')

    expect(result.success).toBe(true)
    expect(axios.post).toHaveBeenCalledWith('/api/v1/auth/login', {
      username: 'user@test.com',
      password: 'p@ss!w0rd'
    })
  })

  it('clearSession can be called multiple times safely', () => {
    auth.clearSession()
    auth.clearSession()
    auth.clearSession()

    expect(auth.getToken()).toBe('')
    expect(auth.getUser()).toBeNull()
  })
})
