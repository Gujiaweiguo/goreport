import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import ChartPropertyPanel from './ChartPropertyPanel.vue'

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
  'el-collapse': {
    template: '<div class="el-collapse"><slot /></div>',
    props: ['modelValue'],
  },
  'el-collapse-item': {
    template: '<div class="el-collapse-item"><slot /></div>',
    props: ['title', 'name'],
  },
  'el-input': {
    template: '<input class="el-input" @input="$emit(\'input\', $event.target.value)" />',
    props: ['modelValue', 'placeholder'],
    emits: ['input', 'update:modelValue'],
  },
  'el-input-number': {
    template: '<input type="number" class="el-input-number" @change="$emit(\'change\', $event.target.value)" />',
    props: ['modelValue', 'min', 'max', 'step'],
    emits: ['change', 'update:modelValue'],
  },
  'el-select': {
    template: '<select class="el-select" @change="$emit(\'change\', $event.target.value)"><slot /></select>',
    props: ['modelValue'],
    emits: ['change', 'update:modelValue'],
  },
  'el-option': {
    template: '<option class="el-option"><slot /></option>',
    props: ['label', 'value'],
  },
  'el-switch': {
    template: '<input type="checkbox" class="el-switch" @change="$emit(\'change\', $event.target.checked)" />',
    props: ['modelValue'],
    emits: ['change', 'update:modelValue'],
  },
  'el-color-picker': {
    template: '<input type="color" class="el-color-picker" @change="$emit(\'change\', $event.target.value)" />',
    props: ['modelValue', 'showAlpha'],
    emits: ['change', 'update:modelValue'],
  },
  'el-tag': {
    template: '<span class="el-tag"><slot /></span>',
    props: ['size', 'effect', 'type'],
  },
}

const defaultConfig = {
  title: 'Test Chart',
  width: 760,
  height: 420,
  margin: 16,
  theme: 'default' as const,
  color: '#2d6cdf',
  fontFamily: "'Noto Serif SC', serif",
  fontSize: 14,
  showLegend: true,
  hoverable: true,
  clickable: true,
  zoomable: false,
  mainTitle: 'Main Title',
  subTitle: 'Sub Title',
  titlePosition: 'left' as const,
}

describe('ChartPropertyPanel.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Component Rendering', () => {
    it('renders the component', () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.chart-property-panel').exists()).toBe(true)
    })

    it('renders panel header', () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.panel-header').exists()).toBe(true)
      expect(wrapper.find('.panel-header h3').text()).toBe('图表属性')
    })

    it('renders property form', () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.property-form').exists()).toBe(true)
    })

    it('renders all collapse groups', () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      const items = wrapper.findAll('.el-collapse-item')
      expect(items.length).toBe(4)
    })
  })

  describe('Form State', () => {
    it('initializes formState from props', async () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.formState.title).toBe('Test Chart')
      expect(vm.formState.width).toBe(760)
      expect(vm.formState.height).toBe(420)
    })

    it('syncs formState when props change', async () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const newConfig = { ...defaultConfig, title: 'Updated Title' }
      await wrapper.setProps({ modelValue: newConfig })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.formState.title).toBe('Updated Title')
    })
  })

  describe('Emit Events', () => {
    it('emits update:modelValue on emitChange', async () => {
      const wrapper = mount(ChartPropertyPanel, {
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
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.emitChange()
      
      const emitted = wrapper.emitted('change')
      expect(emitted).toBeTruthy()
    })

    it('emitted value contains all config properties', async () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      vm.emitChange()
      
      const emitted = wrapper.emitted('update:modelValue')
      const value = emitted[0][0]
      
      expect(value).toHaveProperty('title')
      expect(value).toHaveProperty('width')
      expect(value).toHaveProperty('height')
      expect(value).toHaveProperty('margin')
      expect(value).toHaveProperty('theme')
      expect(value).toHaveProperty('color')
      expect(value).toHaveProperty('fontFamily')
      expect(value).toHaveProperty('fontSize')
      expect(value).toHaveProperty('showLegend')
      expect(value).toHaveProperty('hoverable')
      expect(value).toHaveProperty('clickable')
      expect(value).toHaveProperty('zoomable')
      expect(value).toHaveProperty('mainTitle')
      expect(value).toHaveProperty('subTitle')
      expect(value).toHaveProperty('titlePosition')
    })
  })

  describe('Form Validation Rules', () => {
    it('has validation rules for required fields', () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      const vm = wrapper.vm as any
      expect(vm.rules.title).toBeDefined()
      expect(vm.rules.width).toBeDefined()
      expect(vm.rules.height).toBeDefined()
      expect(vm.rules.margin).toBeDefined()
      expect(vm.rules.fontSize).toBeDefined()
    })
  })

  describe('Active Groups', () => {
    it('initializes with all groups expanded', () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      const vm = wrapper.vm as any
      expect(vm.activeGroups).toContain('basic')
      expect(vm.activeGroups).toContain('style')
      expect(vm.activeGroups).toContain('interaction')
      expect(vm.activeGroups).toContain('title-config')
    })
  })

  describe('Sync Form Function', () => {
    it('syncs all properties from value to formState', async () => {
      const wrapper = mount(ChartPropertyPanel, {
        props: { modelValue: defaultConfig },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      const newConfig = {
        title: 'New Title',
        width: 800,
        height: 500,
        margin: 20,
        theme: 'dark' as const,
        color: '#ff0000',
        fontFamily: 'Arial',
        fontSize: 16,
        showLegend: false,
        hoverable: false,
        clickable: false,
        zoomable: true,
        mainTitle: 'New Main',
        subTitle: 'New Sub',
        titlePosition: 'center' as const,
      }
      
      vm.syncForm(newConfig)
      
      expect(vm.formState.title).toBe('New Title')
      expect(vm.formState.width).toBe(800)
      expect(vm.formState.height).toBe(500)
      expect(vm.formState.margin).toBe(20)
      expect(vm.formState.theme).toBe('dark')
      expect(vm.formState.color).toBe('#ff0000')
      expect(vm.formState.showLegend).toBe(false)
      expect(vm.formState.zoomable).toBe(true)
      expect(vm.formState.titlePosition).toBe('center')
    })
  })

  describe('Default Props', () => {
    it('uses default modelValue when not provided', () => {
      const wrapper = mount(ChartPropertyPanel, {
        global: { stubs: globalStubs },
      })
      const vm = wrapper.vm as any
      expect(vm.formState.title).toBe('图表标题')
      expect(vm.formState.width).toBe(760)
      expect(vm.formState.height).toBe(420)
    })
  })
})
