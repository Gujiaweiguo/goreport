import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import DataConfigPanel from './DataConfigPanel.vue'

vi.mock('@/api/datasource', () => ({
  datasourceApi: {
    list: vi.fn(),
    getTables: vi.fn(),
    getFields: vi.fn(),
  },
}))

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
    warning: vi.fn(),
  },
}))

import { datasourceApi } from '@/api/datasource'
import { datasetApi } from '@/api/dataset'

const defaultConfig = {
  dataSourceId: '',
  tableName: '',
  fields: [] as string[],
  datasetId: '',
  dimension: '',
  measure: '',
  aggregation: 'none' as const,
  params: [] as { key: string; value: string }[],
  filter: '',
}

const globalStubs = {
  'el-form': {
    template: '<form class="el-form"><slot /></form>',
    props: ['model', 'rules', 'labelPosition'],
    methods: {
      validate: vi.fn().mockResolvedValue(true),
    },
  },
  'el-form-item': {
    template: '<div class="el-form-item"><slot /></div>',
    props: ['label', 'prop'],
  },
  'el-radio-group': {
    template: '<div class="el-radio-group"><slot /></div>',
    props: ['modelValue'],
    emits: ['change', 'update:modelValue'],
  },
  'el-radio': {
    template: '<label class="el-radio"><slot /></label>',
    props: ['value'],
  },
  'el-select': {
    template: '<select class="el-select"><slot /></select>',
    props: ['modelValue', 'placeholder', 'filterable', 'clearable', 'loading', 'disabled', 'multiple'],
    emits: ['change', 'update:modelValue'],
  },
  'el-option': {
    template: '<option class="el-option"><slot /></option>',
    props: ['label', 'value'],
  },
  'el-input': {
    template: '<input class="el-input" />',
    props: ['modelValue', 'placeholder'],
  },
  'el-button': {
    template: '<button class="el-button"><slot /></button>',
    props: ['type', 'text'],
  },
  'el-tag': {
    template: '<span class="el-tag"><slot /></span>',
    props: ['size', 'effect', 'type'],
  },
  'el-collapse': {
    template: '<div class="el-collapse"><slot /></div>',
    props: ['modelValue'],
  },
  'el-collapse-item': {
    template: '<div class="el-collapse-item"><slot /></div>',
    props: ['title', 'name'],
  },
  'el-divider': {
    template: '<hr class="el-divider" />',
  },
}

describe('DataConfigPanel.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(datasourceApi.list).mockResolvedValue({
      data: { success: true, result: [{ id: 'ds-1', name: 'Sales DB' }] },
    } as any)
    vi.mocked(datasourceApi.getTables).mockResolvedValue({
      data: { success: true, result: ['users', 'orders'] },
    } as any)
    vi.mocked(datasourceApi.getFields).mockResolvedValue({
      data: { success: true, result: ['id', 'name', 'amount'] },
    } as any)
    vi.mocked(datasetApi.list).mockResolvedValue({
      data: { success: true, result: [{ id: 'dataset-1', name: 'Sales Dataset' }] },
    } as any)
    vi.mocked(datasetApi.getSchema).mockResolvedValue({
      data: {
        success: true,
        result: {
          dimensions: [{ id: 'd1', name: 'category', displayName: 'Category' }],
          measures: [{ id: 'm1', name: 'amount', displayName: 'Amount' }],
        },
      },
    } as any)
  })

  describe('Component Rendering', () => {
    it('renders the component', () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.data-config-panel').exists()).toBe(true)
    })

    it('renders panel header', () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.panel-header').exists()).toBe(true)
    })

    it('renders binding type selector', () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.el-radio-group').exists()).toBe(true)
    })
  })

  describe('Form State Initialization', () => {
    it('initializes formState from props', async () => {
      const config = {
        ...defaultConfig,
        dataSourceId: 'ds-1',
        tableName: 'users',
        fields: ['id', 'name'],
      }
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: config },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.formState.dataSourceId).toBe('ds-1')
      expect(vm.formState.tableName).toBe('users')
      expect(vm.formState.fields).toEqual(['id', 'name'])
    })

    it('syncs formState when props change', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const newConfig = { ...defaultConfig, datasetId: 'dataset-1' }
      await wrapper.setProps({ modelValue: newConfig })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.formState.datasetId).toBe('dataset-1')
    })
  })

  describe('Binding Type Detection', () => {
    it('sets datasource binding type when dataSourceId present', async () => {
      const config = { ...defaultConfig, dataSourceId: 'ds-1' }
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: config },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.bindingType).toBe('datasource')
    })

    it('sets dataset binding type when datasetId present', async () => {
      const config = { ...defaultConfig, datasetId: 'dataset-1' }
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: config },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.bindingType).toBe('dataset')
    })
  })

  describe('Emit Events', () => {
    it('emits update:modelValue on emitChange', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.emitChange()
      
      const emitted = wrapper.emitted('update:modelValue')
      expect(emitted).toBeTruthy()
    })

    it('emits change event on emitChange', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.emitChange()
      
      const emitted = wrapper.emitted('change')
      expect(emitted).toBeTruthy()
    })
  })

  describe('Normalize Field Options', () => {
    it('normalizes string array to options', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      const result = vm.normalizeFieldOptions(['id', 'name'])
      
      expect(result).toEqual([
        { label: 'id', value: 'id' },
        { label: 'name', value: 'name' },
      ])
    })

    it('normalizes object array to options', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      const result = vm.normalizeFieldOptions([
        { name: 'field1' },
        { fieldName: 'field2' },
        { columnName: 'field3' },
      ])
      
      expect(result).toEqual([
        { label: 'field1', value: 'field1' },
        { label: 'field2', value: 'field2' },
        { label: 'field3', value: 'field3' },
      ])
    })

    it('filters out invalid items', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      const result = vm.normalizeFieldOptions([
        'valid',
        { other: 'data' },
        { name: 'also_valid' },
      ])
      
      expect(result.length).toBe(2)
    })
  })

  describe('Parameter Management', () => {
    it('adds a parameter', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: { ...defaultConfig, params: [] } },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.addParam()
      
      expect(vm.formState.params.length).toBe(1)
      expect(vm.formState.params[0]).toEqual({ key: '', value: '' })
    })

    it('removes a parameter', async () => {
      const config = {
        ...defaultConfig,
        params: [{ key: 'p1', value: 'v1' }, { key: 'p2', value: 'v2' }],
      }
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: config },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.removeParam(0)
      
      expect(vm.formState.params.length).toBe(1)
      expect(vm.formState.params[0]).toEqual({ key: 'p2', value: 'v2' })
    })
  })

  describe('Active Groups', () => {
    it('sets datasource groups by default', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.activeGroups).toContain('source')
      expect(vm.activeGroups).toContain('table')
      expect(vm.activeGroups).toContain('fields')
    })

    it('sets dataset groups when binding type is dataset', async () => {
      const config = { ...defaultConfig, datasetId: 'dataset-1' }
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: config },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.activeGroups).toContain('dataset')
      expect(vm.activeGroups).not.toContain('source')
    })
  })

  describe('Loading States', () => {
    it('has loading state structure defined', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.loading).toBeDefined()
      expect(vm.loading).toHaveProperty('sources')
      expect(vm.loading).toHaveProperty('tables')
      expect(vm.loading).toHaveProperty('fields')
      expect(vm.loading).toHaveProperty('datasets')
      expect(vm.loading).toHaveProperty('datasetSchema')
    })
  })

  describe('Sync Form Function', () => {
    it('syncs all properties correctly', async () => {
      const wrapper = mount(DataConfigPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      const newConfig = {
        dataSourceId: 'ds-2',
        tableName: 'orders',
        fields: ['order_id', 'total'],
        datasetId: 'dataset-2',
        dimension: 'category',
        measure: 'sales',
        aggregation: 'SUM' as const,
        params: [{ key: 'year', value: '2024' }],
        filter: 'status = "active"',
      }
      
      vm.syncForm(newConfig)
      
      expect(vm.formState.dataSourceId).toBe('ds-2')
      expect(vm.formState.tableName).toBe('orders')
      expect(vm.formState.fields).toEqual(['order_id', 'total'])
      expect(vm.formState.datasetId).toBe('dataset-2')
      expect(vm.formState.dimension).toBe('category')
      expect(vm.formState.measure).toBe('sales')
      expect(vm.formState.aggregation).toBe('SUM')
      expect(vm.formState.params).toEqual([{ key: 'year', value: '2024' }])
      expect(vm.formState.filter).toBe('status = "active"')
    })
  })
})
