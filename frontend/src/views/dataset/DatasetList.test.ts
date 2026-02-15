import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import DatasetList from './DatasetList.vue'

const mockPush = vi.fn()

vi.mock('vue-router', () => ({
  useRouter: () => ({
    push: mockPush
  })
}))

vi.mock('element-plus', () => ({
  ElMessage: {
    error: vi.fn(),
    success: vi.fn(),
    warning: vi.fn()
  },
  ElMessageBox: {
    confirm: vi.fn()
  }
}))

vi.mock('@/api/dataset', () => ({
  datasetApi: {
    list: vi.fn().mockResolvedValue({ data: { success: true, result: [], total: 0 } }),
    delete: vi.fn().mockResolvedValue({ data: { success: true } }),
    preview: vi.fn().mockResolvedValue({ data: { success: true, result: [] } })
  }
}))

import { ElMessage, ElMessageBox } from 'element-plus'
import { datasetApi } from '@/api/dataset'

const vLoading = {
  mounted: () => {},
  updated: () => {},
  unmounted: () => {},
}

const globalStubs = {
  'el-card': true,
  'el-button': true,
  'el-table': true,
  'el-table-column': true,
  'el-tag': true,
  'el-pagination': true,
  'el-dialog': true,
  DatasetPreview: true
}

const globalConfig = {
  stubs: globalStubs,
  directives: { loading: vLoading }
}

describe('DatasetList.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(datasetApi.list).mockResolvedValue({ data: { success: true, result: [], total: 0 } } as any)
  })

  describe('Component Rendering', () => {
    it('renders correctly', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.find('.dataset-list').exists()).toBe(true)
    })
  })

  describe('Initial state', () => {
    it('datasets is initially empty array', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.datasets).toEqual([])
    })

    it('currentPage starts at 1', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.currentPage).toBe(1)
    })

    it('pageSize starts at 10', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.pageSize).toBe(10)
    })

    it('total starts at 0', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.total).toBe(0)
    })

    it('previewVisible starts false', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.previewVisible).toBe(false)
    })
  })

  describe('createDataset function', () => {
    it('navigates to /dataset/create', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      wrapper.vm.createDataset()
      expect(mockPush).toHaveBeenCalledWith('/dataset/create')
    })
  })

  describe('editDataset function', () => {
    it('navigates to /dataset/edit/:id', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      wrapper.vm.editDataset({ id: 'dataset-123', name: 'Test' } as any)
      expect(mockPush).toHaveBeenCalledWith('/dataset/edit/dataset-123')
    })
  })

  describe('deleteDataset function', () => {
    it('calls ElMessageBox.confirm with dataset name', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(datasetApi.delete).mockResolvedValueOnce({ data: { success: true } } as any)

      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.deleteDataset({ id: 'd1', name: 'Test Dataset' } as any)

      expect(ElMessageBox.confirm).toHaveBeenCalledWith(
        '确定要删除数据集 "Test Dataset" 吗？',
        '确认删除',
        expect.any(Object)
      )
    })

    it('calls datasetApi.delete when confirmed', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(datasetApi.delete).mockResolvedValueOnce({ data: { success: true } } as any)

      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.deleteDataset({ id: 'dataset-789', name: 'Test' } as any)

      expect(datasetApi.delete).toHaveBeenCalledWith('dataset-789')
    })

    it('shows success message after successful delete', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(datasetApi.delete).mockResolvedValueOnce({ data: { success: true } } as any)

      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.deleteDataset({ id: 'd1', name: 'Test' } as any)

      expect(ElMessage.success).toHaveBeenCalledWith('删除成功')
    })

    it('does not call delete when cancelled', async () => {
      vi.mocked(ElMessageBox.confirm).mockRejectedValueOnce('cancel')

      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.deleteDataset({ id: 'd1', name: 'Test' } as any)

      expect(datasetApi.delete).not.toHaveBeenCalled()
    })

    it('shows error message when delete fails', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(datasetApi.delete).mockResolvedValueOnce({ data: { success: false, message: 'Delete failed' } } as any)

      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.deleteDataset({ id: 'd1', name: 'Test' } as any)

      expect(ElMessage.error).toHaveBeenCalledWith('Delete failed')
    })

    it('shows generic error when exception thrown', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      vi.mocked(datasetApi.delete).mockRejectedValueOnce(new Error('Network'))

      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.deleteDataset({ id: 'd1', name: 'Test' } as any)

      expect(ElMessage.error).toHaveBeenCalledWith('删除失败')
    })
  })

  describe('handlePageChange function', () => {
    it('updates currentPage and calls loadDatasets', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      vi.clearAllMocks()
      
      wrapper.vm.handlePageChange(2)
      
      expect(wrapper.vm.currentPage).toBe(2)
      expect(datasetApi.list).toHaveBeenCalled()
    })
  })

  describe('handleSizeChange function', () => {
    it('updates pageSize, resets currentPage to 1, and calls loadDatasets', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      vi.clearAllMocks()
      
      wrapper.vm.handleSizeChange(20)
      
      expect(wrapper.vm.pageSize).toBe(20)
      expect(wrapper.vm.currentPage).toBe(1)
      expect(datasetApi.list).toHaveBeenCalled()
    })
  })

  describe('getTypeLabel function', () => {
    it('returns SQL for sql type', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.getTypeLabel('sql')).toBe('SQL')
    })

    it('returns API for api type', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.getTypeLabel('api')).toBe('API')
    })

    it('returns 文件 for file type', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.getTypeLabel('file')).toBe('文件')
    })

    it('returns original value for unknown type', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.getTypeLabel('unknown')).toBe('unknown')
    })
  })

  describe('getTypeTagType function', () => {
    it('returns success for sql type', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.getTypeTagType('sql')).toBe('success')
    })

    it('returns warning for api type', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.getTypeTagType('api')).toBe('warning')
    })

    it('returns info for file type', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.getTypeTagType('file')).toBe('info')
    })

    it('returns empty string for unknown type', async () => {
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      expect(wrapper.vm.getTypeTagType('unknown')).toBe('')
    })
  })

  describe('previewDataset function', () => {
    it('sets previewDatasetId and opens dialog', async () => {
      vi.mocked(datasetApi.preview).mockResolvedValueOnce({ data: { success: true, result: [] } } as any)
      
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.previewDataset({ id: 'ds-1', name: 'Test' } as any)
      
      expect(wrapper.vm.previewDatasetId).toBe('ds-1')
      expect(wrapper.vm.previewVisible).toBe(true)
    })

    it('loads preview data from API', async () => {
      const mockData = [{ id: 1, name: 'row1' }]
      vi.mocked(datasetApi.preview).mockResolvedValueOnce({ data: { success: true, result: mockData } } as any)
      
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.previewDataset({ id: 'ds-1', name: 'Test' } as any)
      
      expect(wrapper.vm.previewData).toEqual(mockData)
    })

    it('shows mock data on API failure', async () => {
      vi.mocked(datasetApi.preview).mockResolvedValueOnce({ data: { success: false, message: 'Failed' } } as any)
      
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.previewDataset({ id: 'ds-1', name: 'Test' } as any)
      
      expect(wrapper.vm.previewData).toHaveLength(3)
      expect(ElMessage.warning).toHaveBeenCalled()
    })

    it('shows mock data on API error', async () => {
      vi.mocked(datasetApi.preview).mockRejectedValueOnce(new Error('Network'))
      
      const wrapper = mount(DatasetList, { global: globalConfig })
      await nextTick()
      
      await wrapper.vm.previewDataset({ id: 'ds-1', name: 'Test' } as any)
      
      expect(wrapper.vm.previewData).toHaveLength(3)
      expect(ElMessage.warning).toHaveBeenCalled()
    })
  })
})
