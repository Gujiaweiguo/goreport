import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import ReportDesigner from './ReportDesigner.vue'

vi.mock('vue-router', () => ({
  useRoute: () => ({
    query: {}
  })
}))

vi.mock('element-plus', () => ({
  ElMessage: {
    error: vi.fn(),
    success: vi.fn(),
    warning: vi.fn()
  }
}))

vi.mock('@/api/report', () => ({
  reportApi: {
    get: vi.fn().mockResolvedValue({ data: { success: true, result: { id: '1', name: 'Test', config: '{}' } } }),
    create: vi.fn().mockResolvedValue({ data: { success: true } }),
    update: vi.fn().mockResolvedValue({ data: { success: true } })
  }
}))

vi.mock('@/api/dataset', () => ({
  datasetApi: {
    query: vi.fn().mockResolvedValue({ data: { success: true, result: { data: [] } } })
  }
}))

import { ElMessage } from 'element-plus'
import { reportApi } from '@/api/report'

const globalStubs = {
  'el-card': true,
  'el-button': true,
  'el-input': true,
  PropertyPanel: true
}

describe('ReportDesigner.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Component Rendering', () => {
    it('renders correctly', async () => {
      const wrapper = mount(ReportDesigner, { 
        global: { stubs: globalStubs },
        attachTo: document.body
      })
      await nextTick()
      expect(wrapper.find('.report-designer').exists()).toBe(true)
    })
  })

  describe('Initial state', () => {
    it('reportName starts empty', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.reportName).toBe('')
    })

    it('saving starts false', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.saving).toBe(false)
    })

    it('currentReportId starts empty', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.currentReportId).toBe('')
    })

    it('selectedCell starts null', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.selectedCell).toBeNull()
    })

    it('editing.visible starts false', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.editing.visible).toBe(false)
    })

    it('previewData starts empty', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.previewData).toEqual([])
    })

    it('loadingPreviewData starts false', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.loadingPreviewData).toBe(false)
    })
  })

  describe('gridConfig', () => {
    it('has correct default rows', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.gridConfig.rows).toBe(18)
    })

    it('has correct default cols', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.gridConfig.cols).toBe(12)
    })

    it('has correct default cellWidth', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.gridConfig.cellWidth).toBe(120)
    })

    it('has correct default cellHeight', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.gridConfig.cellHeight).toBe(44)
    })
  })

  describe('selectedCellLabel computed', () => {
    it('returns empty string when no cell selected', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.selectedCellLabel).toBe('')
    })

    it('returns correct label for cell at 0,0', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.selectedCell = { row: 0, col: 0, text: '', style: {} as any, binding: {} }
      expect(wrapper.vm.selectedCellLabel).toBe('A1')
    })

    it('returns correct label for cell at 2,3', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.selectedCell = { row: 2, col: 3, text: '', style: {} as any, binding: {} }
      expect(wrapper.vm.selectedCellLabel).toBe('D3')
    })

    it('returns correct label for cell at 9,11', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.selectedCell = { row: 9, col: 11, text: '', style: {} as any, binding: {} }
      expect(wrapper.vm.selectedCellLabel).toBe('L10')
    })
  })

  describe('getCellKey function', () => {
    it('returns correct key format', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.getCellKey(0, 0)).toBe('0:0')
      expect(wrapper.vm.getCellKey(5, 3)).toBe('5:3')
    })
  })

  describe('createDefaultCell function', () => {
    it('creates cell with correct row and col', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.createDefaultCell(3, 5)
      expect(cell.row).toBe(3)
      expect(cell.col).toBe(5)
    })

    it('creates cell with empty text', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.createDefaultCell(0, 0)
      expect(cell.text).toBe('')
    })

    it('creates cell with default style', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.createDefaultCell(0, 0)
      expect(cell.style.fontSize).toBe(14)
      expect(cell.style.fontWeight).toBe('normal')
      expect(cell.style.align).toBe('left')
    })

    it('creates cell with empty binding', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.createDefaultCell(0, 0)
      expect(cell.binding).toEqual({})
    })
  })

  describe('getOrCreateCell function', () => {
    it('creates new cell if not exists', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.getOrCreateCell(5, 5)
      expect(cell.row).toBe(5)
      expect(cell.col).toBe(5)
    })

    it('returns existing cell', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell1 = wrapper.vm.getOrCreateCell(2, 2)
      cell1.text = 'Test'
      const cell2 = wrapper.vm.getOrCreateCell(2, 2)
      expect(cell2.text).toBe('Test')
    })
  })

  describe('serializeConfig function', () => {
    it('serializes grid config', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const config = wrapper.vm.serializeConfig()
      expect(config.grid.rows).toBe(18)
      expect(config.grid.cols).toBe(12)
    })

    it('serializes cells as array', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const config = wrapper.vm.serializeConfig()
      expect(Array.isArray(config.cells)).toBe(true)
    })
  })

  describe('handleNew function', () => {
    it('clears all cells', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.getOrCreateCell(0, 0).text = 'Test'
      wrapper.vm.handleNew()
      expect(wrapper.vm.cells.size).toBe(0)
    })

    it('clears selectedCell', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.selectedCell = {} as any
      wrapper.vm.handleNew()
      expect(wrapper.vm.selectedCell).toBeNull()
    })

    it('clears reportName', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.reportName = 'Test Report'
      wrapper.vm.handleNew()
      expect(wrapper.vm.reportName).toBe('')
    })

    it('shows success message', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.handleNew()
      expect(ElMessage.success).toHaveBeenCalledWith('已新建空白报表')
    })
  })

  describe('handleDeleteCell function', () => {
    it('does nothing when no cell selected', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.getOrCreateCell(0, 0)
      const initialSize = wrapper.vm.cells.size
      wrapper.vm.selectedCell = null
      wrapper.vm.handleDeleteCell()
      expect(wrapper.vm.cells.size).toBe(initialSize)
    })

    it('deletes selected cell', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.getOrCreateCell(0, 0)
      wrapper.vm.selectedCell = cell
      wrapper.vm.handleDeleteCell()
      expect(wrapper.vm.cells.has('0:0')).toBe(false)
    })

    it('clears selectedCell after delete', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.getOrCreateCell(0, 0)
      wrapper.vm.selectedCell = cell
      wrapper.vm.handleDeleteCell()
      expect(wrapper.vm.selectedCell).toBeNull()
    })

    it('shows success message', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.getOrCreateCell(0, 0)
      wrapper.vm.selectedCell = cell
      wrapper.vm.handleDeleteCell()
      expect(ElMessage.success).toHaveBeenCalledWith('单元格已清空')
    })
  })

  describe('handleMergeCell function', () => {
    it('does nothing when no cell selected', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.selectedCell = null
      wrapper.vm.handleMergeCell()
      expect(ElMessage.warning).not.toHaveBeenCalled()
    })

    it('shows warning at row end', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.getOrCreateCell(0, 11)
      wrapper.vm.selectedCell = cell
      wrapper.vm.handleMergeCell()
      expect(ElMessage.warning).toHaveBeenCalledWith('已到达行尾，无法合并')
    })

    it('merges cell with right neighbor', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.getOrCreateCell(0, 0)
      wrapper.vm.selectedCell = cell
      wrapper.vm.handleMergeCell()
      expect(cell.colSpan).toBe(2)
      expect(ElMessage.success).toHaveBeenCalledWith('已合并到右侧单元格')
    })

    it('unmerges when already merged', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = wrapper.vm.getOrCreateCell(0, 0)
      cell.colSpan = 2
      wrapper.vm.selectedCell = cell
      wrapper.vm.handleMergeCell()
      expect(cell.colSpan).toBe(1)
      expect(ElMessage.success).toHaveBeenCalledWith('已取消合并')
    })
  })

  describe('commitEdit function', () => {
    it('does nothing when editing not visible', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.editing.visible = false
      wrapper.vm.commitEdit()
      expect(wrapper.vm.cells.size).toBe(0)
    })

    it('updates cell text', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.editing.visible = true
      wrapper.vm.editing.row = 0
      wrapper.vm.editing.col = 0
      wrapper.vm.editing.value = 'New Text'
      wrapper.vm.commitEdit()
      const cell = wrapper.vm.cells.get('0:0')
      expect(cell?.text).toBe('New Text')
    })

    it('hides editing overlay', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.editing.visible = true
      wrapper.vm.commitEdit()
      expect(wrapper.vm.editing.visible).toBe(false)
    })
  })

  describe('editing reactive state', () => {
    it('has correct initial style values', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.editing.style.left).toBe('0px')
      expect(wrapper.vm.editing.style.top).toBe('0px')
    })
  })

  describe('handleCellUpdate function', () => {
    it('updates cell in cells map', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const cell = {
        row: 0,
        col: 0,
        text: 'Updated',
        style: {
          fontSize: 16,
          fontWeight: 'bold' as const,
          fontStyle: 'normal' as const,
          align: 'center' as const,
          color: '#000',
          background: '#fff',
          borderColor: '#cbd5e1'
        },
        binding: {}
      }
      wrapper.vm.handleCellUpdate(cell)
      const updated = wrapper.vm.cells.get('0:0')
      expect(updated?.text).toBe('Updated')
      expect(updated?.style.fontSize).toBe(16)
    })
  })

  describe('handleSave function', () => {
    it('shows warning when report name is empty', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.reportName = '   '
      await wrapper.vm.handleSave()
      expect(ElMessage.warning).toHaveBeenCalledWith('请输入报表名称')
    })

    it('calls reportApi.create for new report', async () => {
      vi.mocked(reportApi.create).mockResolvedValueOnce({ 
        data: { success: true, result: { id: 'new-id' } } 
      } as any)
      
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.reportName = 'Test Report'
      await wrapper.vm.handleSave()
      
      expect(reportApi.create).toHaveBeenCalled()
      expect(ElMessage.success).toHaveBeenCalled()
    })

    it('calls reportApi.update for existing report', async () => {
      vi.mocked(reportApi.update).mockResolvedValueOnce({ 
        data: { success: true, result: { id: 'existing-id' } } 
      } as any)
      
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.reportName = 'Test Report'
      wrapper.vm.currentReportId = 'existing-id'
      await wrapper.vm.handleSave()
      
      expect(reportApi.update).toHaveBeenCalled()
    })

    it('sets saving flag during save', async () => {
      let resolveSave: () => void
      vi.mocked(reportApi.create).mockImplementation(() => new Promise(resolve => { resolveSave = resolve as any }))
      
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.reportName = 'Test'
      
      const savePromise = wrapper.vm.handleSave()
      expect(wrapper.vm.saving).toBe(true)
      
      resolveSave!()
      await savePromise
      expect(wrapper.vm.saving).toBe(false)
    })

    it('shows error on save failure', async () => {
      vi.mocked(reportApi.create).mockResolvedValueOnce({ 
        data: { success: false, message: 'Save failed' } 
      } as any)
      
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.reportName = 'Test'
      await wrapper.vm.handleSave()
      
      expect(ElMessage.error).toHaveBeenCalledWith('Save failed')
    })
  })

  describe('handlePreviewData function', () => {
    it('shows warning when no cell selected', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.selectedCell = null
      await wrapper.vm.handlePreviewData()
      expect(ElMessage.warning).toHaveBeenCalledWith('请先选择一个单元格并配置数据集绑定')
    })

    it('shows warning when no dataset binding', async () => {
      const wrapper = mount(ReportDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.selectedCell = { row: 0, col: 0, text: '', style: {} as any, binding: {} }
      await wrapper.vm.handlePreviewData()
      expect(ElMessage.warning).toHaveBeenCalledWith('请先选择一个单元格并配置数据集绑定')
    })
  })
})
