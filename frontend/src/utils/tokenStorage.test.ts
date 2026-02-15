import { describe, it, expect, beforeEach, vi } from 'vitest'
import { SecureTokenStorage, type TokenInfo } from './tokenStorage'

describe('SecureTokenStorage', () => {
  let storage: SecureTokenStorage
  const mockLocalStorage = new Map<string, string>()

  const validToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c'

  const mockUserInfo: TokenInfo = {
    token: validToken,
    userId: 'user-123',
    username: 'testuser',
    tenantId: 'tenant-123',
    expiresAt: Date.now() + 3600000
  }

  beforeEach(() => {
    mockLocalStorage.clear()
    
    vi.stubGlobal('localStorage', {
      getItem: (key: string) => mockLocalStorage.get(key) || null,
      setItem: (key: string, value: string) => mockLocalStorage.set(key, value),
      removeItem: (key: string) => mockLocalStorage.delete(key),
      clear: () => mockLocalStorage.clear(),
      length: mockLocalStorage.size,
      key: (index: number) => Array.from(mockLocalStorage.keys())[index] || null
    })
    
    storage = new SecureTokenStorage()
  })

  describe('setToken', () => {
    it('should store valid token', () => {
      storage.setToken(validToken, mockUserInfo, 3600)
      expect(storage.getToken()).toBe(validToken)
    })

    it('should throw error for invalid token format', () => {
      expect(() => storage.setToken('invalid-token', mockUserInfo)).toThrow()
    })

    it('should throw error for empty token', () => {
      expect(() => storage.setToken('', mockUserInfo)).toThrow()
    })

    it('should store user info', () => {
      storage.setToken(validToken, mockUserInfo, 3600)
      const userInfo = storage.getUserInfo()
      expect(userInfo).not.toBeNull()
      expect(userInfo?.userId).toBe('user-123')
      expect(userInfo?.username).toBe('testuser')
    })
  })

  describe('getToken', () => {
    it('should return null when no token stored', () => {
      expect(storage.getToken()).toBeNull()
    })

    it('should return token when stored', () => {
      storage.setToken(validToken, mockUserInfo, 3600)
      expect(storage.getToken()).toBe(validToken)
    })
  })

  describe('clearToken', () => {
    it('should remove token from storage', () => {
      storage.setToken(validToken, mockUserInfo, 3600)
      storage.clearToken()
      expect(storage.getToken()).toBeNull()
    })

    it('should remove user info from storage', () => {
      storage.setToken(validToken, mockUserInfo, 3600)
      storage.clearToken()
      expect(storage.getUserInfo()).toBeNull()
    })
  })

  describe('isTokenValid', () => {
    it('should return false for empty token', () => {
      expect(storage.isTokenValid('')).toBe(false)
    })

    it('should return false for undefined token', () => {
      expect(storage.isTokenValid(undefined)).toBe(false)
    })

    it('should return true for valid non-expired token', () => {
      storage.setToken(validToken, mockUserInfo, 3600)
      expect(storage.isTokenValid(validToken)).toBe(true)
    })
  })

  describe('getUserInfo', () => {
    it('should return null when no user info stored', () => {
      expect(storage.getUserInfo()).toBeNull()
    })

    it('should return user info when stored', () => {
      storage.setToken(validToken, mockUserInfo, 3600)
      const info = storage.getUserInfo()
      expect(info).not.toBeNull()
      expect(info?.userId).toBe('user-123')
      expect(info?.tenantId).toBe('tenant-123')
    })
  })

  describe('cleanExpiredTokens', () => {
    it('should do nothing when no token stored', () => {
      storage.cleanExpiredTokens()
      expect(storage.getToken()).toBeNull()
    })

    it('should keep valid token', () => {
      storage.setToken(validToken, mockUserInfo, 3600)
      storage.cleanExpiredTokens()
      expect(storage.getToken()).toBe(validToken)
    })
  })

  describe('validateTokenFormat', () => {
    it('should validate JWT format with three parts', () => {
      storage.setToken(validToken, mockUserInfo, 3600)
      expect(storage.getToken()).toBe(validToken)
    })

    it('should reject token with wrong number of parts', () => {
      expect(() => storage.setToken('part1.part2', mockUserInfo)).toThrow()
    })

    it('should reject token with empty parts', () => {
      expect(() => storage.setToken('part1..part3', mockUserInfo)).toThrow()
    })
  })
})
