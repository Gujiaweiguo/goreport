import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import ChartPanel from './ChartPanel.vue'

vi.mock('@/api/chart', () => ({
  chartApi: {
    list: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
  },
}))

vi.mock('element-plus', () => ({
  ElMessage: {
    success: vi.fn(),
    error: vi.fn(),
  },
}))

import { chartApi } from '@/api/chart'

const mockCharts = [
  { id: '1', tenantId: 'tenant-1', name: 'Sales Chart', code: 'sales', type: 'bar', config: '{}', status: 1, createdAt: '', updatedAt: '' },
  { id: '2', tenantId: 'tenant-1', name: 'Trend Chart', code: 'trend', type: 'line', config: '{}', status: 1, createdAt: '', updatedAt: '' },
  { id: '3', tenantId: 'tenant-1', name: 'Distribution', code: 'dist', type: 'pie', config: '{}', status: 1, createdAt: '', updatedAt: '' },
]

const globalStubs = {
  'el-button': {
    template: '<button class="el-button" @click="$emit(\'click\')"><slot /></button>',
    props: ['type', 'size'],
    emits: ['click'],
  },
  'el-icon': {
    template: '<i class="el-icon"><slot /></i>',
    props: ['size'],
  },
  'el-dialog': {
    template: '<div class="el-dialog" v-if="modelValue"><slot /><slot name="footer" /></div>',
    props: ['modelValue', 'title', 'width'],
    emits: ['update:modelValue', 'close'],
  },
  ChartEditor: {
    template: '<div class="chart-editor"></div>',
    props: ['modelValue'],
    emits: ['update:modelValue'],
  },
}

describe('ChartPanel.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(chartApi.list).mockResolvedValue({
      data: { success: true, result: mockCharts },
    } as any)
  })

  describe('Component Rendering', () => {
    it('renders the component', () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.chart-panel').exists()).toBe(true)
    })

    it('renders panel header', () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.panel-header').exists()).toBe(true)
      expect(wrapper.find('.panel-header span').text()).toBe('图表库')
    })

    it('renders create chart button', () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      const buttons = wrapper.findAll('.el-button')
      expect(buttons.length).toBeGreaterThan(0)
    })
  })

  describe('Loading Charts', () => {
    it('loads charts on mount', async () => {
      mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      await nextTick()
      expect(chartApi.list).toHaveBeenCalled()
    })

    it('displays loaded charts', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      await nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
      
      const items = wrapper.findAll('.chart-item')
      expect(items.length).toBe(3)
    })

    it('handles load error gracefully', async () => {
      vi.mocked(chartApi.list).mockRejectedValueOnce(new Error('Network error'))
      
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
      
      const items = wrapper.findAll('.chart-item')
      expect(items.length).toBe(0)
    })
  })

  describe('Chart Type Labels', () => {
    it('returns correct label for bar chart', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.chartTypeLabel('bar')).toBe('柱状图')
    })

    it('returns correct label for line chart', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.chartTypeLabel('line')).toBe('折线图')
    })

    it('returns correct label for pie chart', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.chartTypeLabel('pie')).toBe('饼图')
    })

    it('returns original type for unknown chart', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.chartTypeLabel('scatter')).toBe('scatter')
    })
  })

  describe('Create Chart', () => {
    it('opens dialog with create mode', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.handleCreateChart()
      
      expect(vm.chartEditor.mode).toBe('create')
      expect(vm.chartEditor.visible).toBe(true)
      expect(vm.chartEditor.data).toBeDefined()
    })

    it('initializes new chart with default values', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.handleCreateChart()
      
      expect(vm.chartEditor.data.type).toBe('bar')
      expect(vm.chartEditor.data.name).toBe('')
    })
  })

  describe('Edit Chart', () => {
    it('opens dialog with update mode', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.handleEditChart(mockCharts[0])
      
      expect(vm.chartEditor.mode).toBe('update')
      expect(vm.chartEditor.visible).toBe(true)
    })

    it('sets chart data for editing', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.handleEditChart(mockCharts[1])
      
      expect(vm.chartEditor.data).toEqual(mockCharts[1])
    })
  })

  describe('Chart Update Handler', () => {
    it('updates chartEditor data', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      const newData = { ...mockCharts[0], name: 'Updated Name' }
      vm.handleChartUpdate(newData)
      
      expect(vm.chartEditor.data).toEqual(newData)
    })
  })

  describe('Drag and Drop', () => {
    it('sets drag data on dragstart', async () => {
      const wrapper = mount(ChartPanel, {
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      const mockEvent = {
        dataTransfer: {
          setData: vi.fn(),
          effectAllowed: null,
        },
      }
      
      vm.handleDragStart(mockCharts[0], mockEvent as any)
      
      expect(mockEvent.dataTransfer.setData).toHaveBeenCalledWith(
        'application/json',
        expect.stringContaining('"type":"chart"')
      )
      expect(mockEvent.dataTransfer.effectAllowed).toBe('copy')
    })
  })
})
