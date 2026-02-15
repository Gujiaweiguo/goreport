import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import Home from './Home.vue'

vi.mock('axios', () => ({
  default: {
    get: vi.fn()
  }
}))

import axios from 'axios'

const globalStubs = {
  'el-button': {
    template: '<button><slot /></button>',
    props: ['type']
  }
}

describe('Home.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Component Rendering', () => {
    it('renders correctly', () => {
      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      expect(wrapper.find('.home').exists()).toBe(true)
    })

    it('displays the title', () => {
      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      expect(wrapper.find('h1').text()).toBe('goReport')
    })

    it('displays the welcome message', () => {
      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      expect(wrapper.find('p').text()).toContain('欢迎使用 goReport 报表系统')
    })

    it('has a test connection button', () => {
      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      const button = wrapper.find('button')
      expect(button.exists()).toBe(true)
      expect(button.text()).toBe('测试连接')
    })

    it('does not show health status initially', () => {
      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      const paragraphs = wrapper.findAll('p')
      const statusParagraph = paragraphs.find(p => p.text().includes('状态:'))
      expect(statusParagraph).toBeUndefined()
    })
  })

  describe('testConnection function', () => {
    it('calls axios.get with /health endpoint', async () => {
      const mockGet = vi.mocked(axios.get)
      mockGet.mockResolvedValueOnce({ data: { status: 'ok' } })

      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      await wrapper.find('button').trigger('click')

      expect(mockGet).toHaveBeenCalledWith('/health')
    })

    it('displays health status on successful connection', async () => {
      const mockGet = vi.mocked(axios.get)
      mockGet.mockResolvedValueOnce({ data: { status: 'ok', message: 'healthy' } })

      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      await wrapper.find('button').trigger('click')
      await nextTick()

      const paragraphs = wrapper.findAll('p')
      const statusParagraph = paragraphs.find(p => p.text().includes('状态:'))
      expect(statusParagraph).toBeDefined()
      expect(statusParagraph?.text()).toContain('status')
    })

    it('displays connection failed message on error', async () => {
      const mockGet = vi.mocked(axios.get)
      mockGet.mockRejectedValueOnce(new Error('Network Error'))

      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      await wrapper.find('button').trigger('click')
      await nextTick()

      const paragraphs = wrapper.findAll('p')
      const statusParagraph = paragraphs.find(p => p.text().includes('状态:'))
      expect(statusParagraph).toBeDefined()
      expect(statusParagraph?.text()).toContain('连接失败')
    })

    it('stringifies the response data for display', async () => {
      const mockGet = vi.mocked(axios.get)
      const responseData = { status: 'ok', version: '1.0.0' }
      mockGet.mockResolvedValueOnce({ data: responseData })

      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      await wrapper.find('button').trigger('click')
      await nextTick()

      const paragraphs = wrapper.findAll('p')
      const statusParagraph = paragraphs.find(p => p.text().includes('状态:'))
      expect(statusParagraph?.text()).toContain(JSON.stringify(responseData))
    })
  })

  describe('healthStatus state', () => {
    it('healthStatus is initially empty', () => {
      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      expect(wrapper.vm.healthStatus).toBe('')
    })

    it('updates healthStatus on successful connection', async () => {
      const mockGet = vi.mocked(axios.get)
      mockGet.mockResolvedValueOnce({ data: { status: 'ok' } })

      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      await wrapper.find('button').trigger('click')
      await nextTick()

      expect(wrapper.vm.healthStatus).toBe('{"status":"ok"}')
    })

    it('updates healthStatus on connection failure', async () => {
      const mockGet = vi.mocked(axios.get)
      mockGet.mockRejectedValueOnce(new Error('Network Error'))

      const wrapper = mount(Home, { global: { stubs: globalStubs } })
      await wrapper.find('button').trigger('click')
      await nextTick()

      expect(wrapper.vm.healthStatus).toBe('连接失败')
    })
  })
})
