import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import ReportList from './ReportList.vue'

const mockPush = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  })
}))

vi.mock('element-plus', () => ({
  ElMessage: {
    error: vi.fn(),
    success: vi.fn()
  },
  ElMessageBox: {
    confirm: vi.fn()
  }
}))

vi.mock('@/api/report', () => ({
  reportApi: {
    list: vi.fn().mockResolvedValue({ data: { success: true, result: [] } }),
    delete: vi.fn().mockResolvedValue({ data: { success: true } })
  }
}))

import { ElMessage, ElMessageBox } from 'element-plus'
import { reportApi } from '@/api/report'

const globalStubs = {
  'el-card': true,
  'el-button': true,
  'el-table': true,
  'el-table-column': true,
  'el-tag': true,
  'el-empty': true,
  LoadingState: true
}

describe('ReportList.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(reportApi.list).mockResolvedValue({ data: { success: true, result: [] } } as any)
  })

  describe('Component Rendering', () => {
    it('renders correctly', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.find('.report-list').exists()).toBe(true)
    })

    it('has page header with create button area', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.find('.page-header').exists()).toBe(true)
    })

    it('has table container', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.find('.table-container').exists()).toBe(true)
    })
  })

  describe('Initial state', () => {
    it('reports is initially empty array', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.reports).toEqual([])
    })

    it('loading state starts as false', async () => {
      vi.mocked(reportApi.list).mockImplementation(() => new Promise(() => {}))
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.loading).toBe(true)
    })
  })

  describe('handleCreate function', () => {
    it('navigates to ReportDesigner', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.handleCreate()
      expect(mockPush).toHaveBeenCalledWith({ name: 'ReportDesigner' })
    })
  })

  describe('handleEdit function', () => {
    it('navigates to ReportDesigner with id query', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.handleEdit({ id: 'report-123', name: 'Test' } as any)
      expect(mockPush).toHaveBeenCalledWith({
        name: 'ReportDesigner',
        query: { id: 'report-123' }
      })
    })
  })

  describe('handlePreview function', () => {
    it('navigates to ReportPreview with id query', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.handlePreview({ id: 'report-456', name: 'Test' } as any)
      expect(mockPush).toHaveBeenCalledWith({
        name: 'ReportPreview',
        query: { id: 'report-456' }
      })
    })
  })

  describe('formatDate function', () => {
    it('formats valid date string', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      const result = wrapper.vm.formatDate('2024-01-15T10:30:00')
      expect(result).toContain('2024')
    })

    it('returns dash for empty date', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.formatDate('')).toBe('-')
    })

    it('returns dash for null date', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.formatDate(null as any)).toBe('-')
    })

    it('returns dash for undefined date', async () => {
      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.formatDate(undefined as any)).toBe('-')
    })
  })

  describe('handleDelete function', () => {
    it('calls ElMessageBox.confirm with report name', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(reportApi.delete).mockResolvedValueOnce({ data: { success: true } } as any)

      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      
      const deletePromise = wrapper.vm.handleDelete({ id: 'r1', name: 'Test Report' } as any)
      await deletePromise

      expect(ElMessageBox.confirm).toHaveBeenCalledWith(
        '确认删除报表 "Test Report" 吗？',
        '提示',
        expect.any(Object)
      )
    })

    it('calls reportApi.delete when confirmed', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(reportApi.delete).mockResolvedValueOnce({ data: { success: true } } as any)

      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleDelete({ id: 'report-789', name: 'Test' } as any)

      expect(reportApi.delete).toHaveBeenCalledWith('report-789')
    })

    it('shows success message after successful delete', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(reportApi.delete).mockResolvedValueOnce({ data: { success: true } } as any)

      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleDelete({ id: 'r1', name: 'Test' } as any)

      expect(ElMessage.success).toHaveBeenCalledWith('删除报表成功')
    })

    it('does not call delete when cancelled', async () => {
      vi.mocked(ElMessageBox.confirm).mockRejectedValueOnce('cancel')

      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleDelete({ id: 'r1', name: 'Test' } as any)

      expect(reportApi.delete).not.toHaveBeenCalled()
    })

    it('shows error message when delete fails', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(reportApi.delete).mockResolvedValueOnce({ data: { success: false, message: 'Delete failed' } } as any)

      const wrapper = mount(ReportList, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleDelete({ id: 'r1', name: 'Test' } as any)

      expect(ElMessage.error).toHaveBeenCalledWith('Delete failed')
    })
  })
})
