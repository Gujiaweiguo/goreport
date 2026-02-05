import axios from 'axios'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    id: string
    username: string
    role: string
    tenantId: string
  }
}

export interface ApiResponse<T = any> {
  success: boolean
  result: T
  message: string
}

export const authApi = {
  login: (data: LoginRequest) => {
    return axios.post<ApiResponse<LoginResponse>>('/api/v1/auth/login', data)
  },

  logout: () => {
    return axios.post<ApiResponse<null>>('/api/v1/auth/logout')
  }
}

export const auth = {
  login: async (username: string, password: string) => {
    const response = await authApi.login({ username, password })
    if (response.data.success) {
      localStorage.setItem('token', response.data.result.token)
      localStorage.setItem('user', JSON.stringify(response.data.result.user))
      return { success: true, user: response.data.result.user }
    }
    return { success: false, message: response.data.message }
  },

  logout: async () => {
    try {
      await authApi.logout()
    } catch (error) {
      console.error('Logout failed:', error)
    } finally {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  },

  getToken: () => {
    return localStorage.getItem('token') || ''
  },

  getUser: () => {
    const userStr = localStorage.getItem('user')
    if (userStr) {
      try {
        return JSON.parse(userStr)
      } catch {
        return null
      }
    }
    return null
  },

  isAuthenticated: () => {
    return !!localStorage.getItem('token')
  }
}
