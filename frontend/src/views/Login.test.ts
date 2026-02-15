import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import Login from './Login.vue'

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  })
}))

vi.mock('element-plus', () => ({
  ElMessage: {
    error: vi.fn(),
    success: vi.fn()
  }
}))

vi.mock('@/api/auth', () => ({
  auth: {
    login: vi.fn()
  }
}))

import { ElMessage } from 'element-plus'
import { auth } from '@/api/auth'

const mockPush = vi.fn()

const globalStubs = {
  'el-card': true,
  'el-form': true,
  'el-form-item': true,
  'el-input': true,
  'el-button': true
}

describe('Login.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('validateUsername', () => {
    const validateUsername = (_rule: unknown, value: string | null | undefined, callback: (error?: Error) => void) => {
      if (!value) {
        callback(new Error('请输入用户名'))
      } else if (value.length < 3) {
        callback(new Error('用户名长度不能少于 3 位'))
      } else {
        callback()
      }
    }

    it('should fail when username is empty', () => {
      let error: Error | undefined
      validateUsername({}, '', (e) => { error = e })
      expect(error).toBeInstanceOf(Error)
      expect(error?.message).toBe('请输入用户名')
    })

    it('should fail when username is null', () => {
      let error: Error | undefined
      validateUsername({}, null, (e) => { error = e })
      expect(error).toBeInstanceOf(Error)
      expect(error?.message).toBe('请输入用户名')
    })

    it('should fail when username is undefined', () => {
      let error: Error | undefined
      validateUsername({}, undefined, (e) => { error = e })
      expect(error).toBeInstanceOf(Error)
      expect(error?.message).toBe('请输入用户名')
    })

    it('should fail when username length is less than 3', () => {
      let error: Error | undefined
      validateUsername({}, 'ab', (e) => { error = e })
      expect(error).toBeInstanceOf(Error)
      expect(error?.message).toBe('用户名长度不能少于 3 位')
    })

    it('should pass when username length is exactly 3', () => {
      let called = false
      let error: Error | undefined
      validateUsername({}, 'abc', (e) => { called = true; error = e })
      expect(called).toBe(true)
      expect(error).toBeUndefined()
    })

    it('should pass when username is valid', () => {
      let called = false
      let error: Error | undefined
      validateUsername({}, 'admin', (e) => { called = true; error = e })
      expect(called).toBe(true)
      expect(error).toBeUndefined()
    })

    it('should pass with long username', () => {
      let called = false
      let error: Error | undefined
      validateUsername({}, 'verylongusername', (e) => { called = true; error = e })
      expect(called).toBe(true)
      expect(error).toBeUndefined()
    })
  })

  describe('validatePassword', () => {
    const validatePassword = (_rule: unknown, value: string | null | undefined, callback: (error?: Error) => void) => {
      if (!value) {
        callback(new Error('请输入密码'))
      } else if (value.length < 6) {
        callback(new Error('密码长度不能少于 6 位'))
      } else {
        callback()
      }
    }

    it('should fail when password is empty', () => {
      let error: Error | undefined
      validatePassword({}, '', (e) => { error = e })
      expect(error).toBeInstanceOf(Error)
      expect(error?.message).toBe('请输入密码')
    })

    it('should fail when password is null', () => {
      let error: Error | undefined
      validatePassword({}, null, (e) => { error = e })
      expect(error).toBeInstanceOf(Error)
      expect(error?.message).toBe('请输入密码')
    })

    it('should fail when password length is less than 6', () => {
      let error: Error | undefined
      validatePassword({}, '12345', (e) => { error = e })
      expect(error).toBeInstanceOf(Error)
      expect(error?.message).toBe('密码长度不能少于 6 位')
    })

    it('should pass when password length is exactly 6', () => {
      let called = false
      let error: Error | undefined
      validatePassword({}, '123456', (e) => { called = true; error = e })
      expect(called).toBe(true)
      expect(error).toBeUndefined()
    })

    it('should pass when password is valid', () => {
      let called = false
      let error: Error | undefined
      validatePassword({}, 'admin123', (e) => { called = true; error = e })
      expect(called).toBe(true)
      expect(error).toBeUndefined()
    })

    it('should pass with complex password', () => {
      let called = false
      let error: Error | undefined
      validatePassword({}, 'P@ssw0rd!123', (e) => { called = true; error = e })
      expect(called).toBe(true)
      expect(error).toBeUndefined()
    })
  })

  describe('LoginForm Interface', () => {
    it('should have username and password fields', () => {
      interface LoginForm {
        username: string
        password: string
      }

      const form: LoginForm = {
        username: 'test',
        password: 'test123'
      }

      expect(form.username).toBe('test')
      expect(form.password).toBe('test123')
    })
  })

  describe('Component Rendering', () => {
    it('renders correctly', async () => {
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.find('.login-container').exists()).toBe(true)
    })

    it('has login card', async () => {
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.find('.login-card').exists()).toBe(true)
    })
  })

  describe('Initial state', () => {
    it('loading is initially false', async () => {
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.loading).toBe(false)
    })

    it('loginForm has empty username', async () => {
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.loginForm.username).toBe('')
    })

    it('loginForm has empty password', async () => {
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.loginForm.password).toBe('')
    })
  })

  describe('handleLogin function', () => {
    it('does nothing when form ref is null', async () => {
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.loginFormRef = null
      await wrapper.vm.handleLogin()
      expect(auth.login).not.toHaveBeenCalled()
    })

    it('calls auth.login with username and password', async () => {
      vi.mocked(auth.login).mockResolvedValueOnce({ success: true } as any)
      
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.loginForm.username = 'admin'
      wrapper.vm.loginForm.password = 'admin123'
      
      wrapper.vm.loginFormRef = {
        validate: vi.fn().mockResolvedValue(true)
      } as any
      
      await wrapper.vm.handleLogin()
      
      expect(auth.login).toHaveBeenCalledWith('admin', 'admin123')
    })

    it('shows success message on successful login', async () => {
      vi.mocked(auth.login).mockResolvedValueOnce({ success: true } as any)
      
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.loginForm.username = 'admin'
      wrapper.vm.loginForm.password = 'admin123'
      wrapper.vm.loginFormRef = {
        validate: vi.fn().mockResolvedValue(true)
      } as any
      
      await wrapper.vm.handleLogin()
      
      expect(ElMessage.success).toHaveBeenCalledWith('登录成功')
    })

    it('navigates to home on successful login', async () => {
      vi.mocked(auth.login).mockResolvedValueOnce({ success: true } as any)
      
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.loginForm.username = 'admin'
      wrapper.vm.loginForm.password = 'admin123'
      wrapper.vm.loginFormRef = {
        validate: vi.fn().mockResolvedValue(true)
      } as any
      
      await wrapper.vm.handleLogin()
      
      expect(mockPush).toHaveBeenCalledWith('/')
    })

    it('shows error message on failed login', async () => {
      vi.mocked(auth.login).mockResolvedValueOnce({ success: false, message: 'Invalid credentials' } as any)
      
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.loginForm.username = 'admin'
      wrapper.vm.loginForm.password = 'wrong'
      wrapper.vm.loginFormRef = {
        validate: vi.fn().mockResolvedValue(true)
      } as any
      
      await wrapper.vm.handleLogin()
      
      expect(ElMessage.error).toHaveBeenCalledWith('Invalid credentials')
    })

    it('shows error message on exception', async () => {
      vi.mocked(auth.login).mockRejectedValueOnce(new Error('Network Error'))
      
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.loginForm.username = 'admin'
      wrapper.vm.loginForm.password = 'admin123'
      wrapper.vm.loginFormRef = {
        validate: vi.fn().mockResolvedValue(true)
      } as any
      
      await wrapper.vm.handleLogin()
      
      expect(ElMessage.error).toHaveBeenCalled()
    })

    it('sets loading to false after login attempt', async () => {
      vi.mocked(auth.login).mockResolvedValueOnce({ success: true } as any)
      
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.loginForm.username = 'admin'
      wrapper.vm.loginForm.password = 'admin123'
      wrapper.vm.loginFormRef = {
        validate: vi.fn().mockResolvedValue(true)
      } as any
      
      await wrapper.vm.handleLogin()
      
      expect(wrapper.vm.loading).toBe(false)
    })

    it('does not call login when validation fails', async () => {
      const wrapper = mount(Login, { global: { stubs: globalStubs } })
      await nextTick()
      
      wrapper.vm.loginForm.username = 'admin'
      wrapper.vm.loginForm.password = 'admin123'
      wrapper.vm.loginFormRef = {
        validate: vi.fn().mockRejectedValue(new Error('Validation failed'))
      } as any
      
      await wrapper.vm.handleLogin()
      
      expect(auth.login).not.toHaveBeenCalled()
    })
  })
})
