import { beforeEach, describe, expect, it, vi } from 'vitest'
import apiClient from './client'
import { userApi } from './user'

vi.mock('./client', () => {
  return {
    default: {
      get: vi.fn()
    }
  }
})

describe('userApi', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('getMe', () => {
    it('calls me endpoint', () => {
      userApi.getMe()
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/users/me')
    })
  })
})
