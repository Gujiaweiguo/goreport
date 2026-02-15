import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import ChartEditor from './ChartEditor.vue'

vi.mock('@/api/dataset', () => ({
  datasetApi: {
    list: vi.fn(),
    getSchema: vi.fn(),
  },
}))

vi.mock('element-plus', () => ({
  ElMessage: {
    success: vi.fn(),
    error: vi.fn(),
  },
}))

import { datasetApi } from '@/api/dataset'

const mockDatasets = [
  { id: 'ds-1', name: 'Sales Dataset', tenantId: 'tenant-1' },
  { id: 'ds-2', name: 'Revenue Dataset', tenantId: 'tenant-1' },
]

const mockSchema = {
  dimensions: [
    { id: 'f-1', name: 'category', displayName: 'Category', fieldType: 'dimension' },
    { id: 'f-2', name: 'region', displayName: 'Region', fieldType: 'dimension' },
  ],
  measures: [
    { id: 'f-3', name: 'amount', displayName: 'Amount', fieldType: 'measure' },
    { id: 'f-4', name: 'count', displayName: 'Count', fieldType: 'measure' },
  ],
}

const defaultChartData = {
  title: 'Test Chart',
  type: 'bar' as const,
  datasetId: '',
  dimension: '',
  measure: '',
  aggregation: 'none' as const,
}

const globalStubs = {
  'el-form': {
    template: '<form class="el-form"><slot /></form>',
    props: ['model', 'labelWidth'],
  },
  'el-form-item': {
    template: '<div class="el-form-item"><slot /></div>',
    props: ['label', 'prop'],
  },
  'el-input': {
    template: '<input class="el-input" />',
    props: ['modelValue', 'placeholder'],
  },
  'el-select': {
    template: '<select class="el-select" @change="$emit(\'change\')"><slot /></select>',
    props: ['modelValue', 'placeholder', 'clearable', 'disabled'],
    emits: ['change', 'update:modelValue'],
  },
  'el-option': {
    template: '<option class="el-option"><slot /></option>',
    props: ['label', 'value'],
  },
  'el-tag': {
    template: '<span class="el-tag"><slot /></span>',
    props: ['size', 'effect', 'type'],
  },
  'el-divider': {
    template: '<hr class="el-divider" />',
    props: ['contentPosition'],
  },
}

describe('ChartEditor.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(datasetApi.list).mockResolvedValue({
      data: { success: true, result: mockDatasets },
    } as any)
    vi.mocked(datasetApi.getSchema).mockResolvedValue({
      data: { success: true, result: mockSchema },
    } as any)
  })

  describe('Component Rendering', () => {
    it('renders the component', () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.chart-editor').exists()).toBe(true)
    })

    it('renders editor header', () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.editor-header').exists()).toBe(true)
    })

    it('renders editor form', () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.editor-form').exists()).toBe(true)
    })
  })

  describe('Chart Type Label', () => {
    it('returns correct label for bar chart', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: { ...defaultChartData, type: 'bar' } },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.chartTypeLabel).toBe('柱状图')
    })

    it('returns correct label for line chart', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: { ...defaultChartData, type: 'line' } },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.chartTypeLabel).toBe('折线图')
    })

    it('returns correct label for pie chart', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: { ...defaultChartData, type: 'pie' } },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.chartTypeLabel).toBe('饼图')
    })
  })

  describe('Loading Datasets', () => {
    it('loads datasets on mount', async () => {
      mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      await nextTick()
      await nextTick()
      
      expect(datasetApi.list).toHaveBeenCalled()
    })

    it('handles load error gracefully', async () => {
      vi.mocked(datasetApi.list).mockRejectedValueOnce(new Error('Network error'))
      
      mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      await nextTick()
      await nextTick()
      
      expect(datasetApi.list).toHaveBeenCalled()
    })
  })

  describe('Form Data Initialization', () => {
    it('initializes formData from props', async () => {
      const chartData = {
        ...defaultChartData,
        title: 'Custom Title',
        type: 'line' as const,
        datasetId: 'ds-1',
      }
      const wrapper = mount(ChartEditor, {
        props: { modelValue: chartData },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.formData.title).toBe('Custom Title')
      expect(vm.formData.type).toBe('line')
      expect(vm.formData.datasetId).toBe('ds-1')
    })
  })

  describe('Handle Type Change', () => {
    it('emits update when type changes', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.handleTypeChange()
      
      const emitted = wrapper.emitted('update:modelValue')
      expect(emitted).toBeTruthy()
    })
  })

  describe('Handle Dataset Change', () => {
    it('clears dimensions and measures when dataset changes', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: { ...defaultChartData, datasetId: 'ds-1', dimension: 'cat', measure: 'amt' } },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.formData.datasetId = 'ds-2'
      vm.handleDatasetChange()
      
      expect(vm.dimensions).toEqual([])
      expect(vm.measures).toEqual([])
      expect(vm.formData.dimension).toBe('')
      expect(vm.formData.measure).toBe('')
    })

    it('loads schema when dataset is set', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.formData.datasetId = 'ds-1'
      vm.handleDatasetChange()
      await nextTick()
      await new Promise(resolve => setTimeout(resolve, 100))
      
      expect(datasetApi.getSchema).toHaveBeenCalledWith('ds-1')
    })
  })

  describe('Handle Measure Change', () => {
    it('emits update when measure changes', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.handleMeasureChange()
      
      const emitted = wrapper.emitted('update:modelValue')
      expect(emitted).toBeTruthy()
    })
  })

  describe('Emit Update', () => {
    it('emits complete chart data', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.formData.title = 'New Title'
      vm.formData.type = 'pie'
      vm.formData.datasetId = 'ds-1'
      vm.emitUpdate()
      
      const emitted = wrapper.emitted('update:modelValue')
      const lastEmit = emitted[emitted.length - 1][0]
      
      expect(lastEmit.title).toBe('New Title')
      expect(lastEmit.type).toBe('pie')
      expect(lastEmit.datasetId).toBe('ds-1')
    })
  })

  describe('Props Watch', () => {
    it('updates formData when props change', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const newData = { ...defaultChartData, title: 'Updated Title', type: 'line' as const }
      await wrapper.setProps({ modelValue: newData })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.formData.title).toBe('Updated Title')
      expect(vm.formData.type).toBe('line')
    })
  })

  describe('Load Dataset Schema', () => {
    it('does nothing if no datasetId', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: defaultChartData },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      vi.clearAllMocks()
      
      const vm = wrapper.vm as any
      await vm.loadDatasetSchema()
      
      expect(datasetApi.getSchema).not.toHaveBeenCalled()
    })

    it('loads schema for valid datasetId', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: { ...defaultChartData, datasetId: 'ds-1' } },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      await vm.loadDatasetSchema()
      
      expect(datasetApi.getSchema).toHaveBeenCalledWith('ds-1')
    })

    it('sets dimensions and measures from response', async () => {
      const wrapper = mount(ChartEditor, {
        props: { modelValue: { ...defaultChartData, datasetId: 'ds-1' } },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      await vm.loadDatasetSchema()
      await nextTick()
      
      expect(vm.dimensions.length).toBe(2)
      expect(vm.measures.length).toBe(2)
    })
  })
})
