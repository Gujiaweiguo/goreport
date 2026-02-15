import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import ChartTypeSelector from './ChartTypeSelector.vue'

const globalStubs = {
  'el-collapse': {
    template: '<div class="el-collapse"><slot /></div>',
    props: ['modelValue'],
  },
  'el-collapse-item': {
    template: '<div class="el-collapse-item"><slot name="title" /><slot /></div>',
    props: ['name', 'title'],
  },
  'el-card': {
    template: '<div class="el-card" @click="$emit(\'click\')"><slot /></div>',
    props: ['shadow', 'class'],
    emits: ['click'],
  },
  'el-tag': {
    template: '<span class="el-tag"><slot /></span>',
    props: ['size', 'effect', 'type'],
  },
  'el-icon': {
    template: '<i class="el-icon"><slot /></i>',
    props: ['size'],
  },
}

describe('ChartTypeSelector.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Component Rendering', () => {
    it('renders the component', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.chart-type-selector').exists()).toBe(true)
    })

    it('renders chart groups', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      const groups = wrapper.findAll('.el-collapse-item')
      expect(groups.length).toBe(5)
    })

    it('renders selector detail panel', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.selector-detail').exists()).toBe(true)
    })

    it('displays detail header for selected type', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.detail-header').exists()).toBe(true)
    })
  })

  describe('Type Selection', () => {
    it('emits update:modelValue when type is selected', async () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: '' },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const emitted = wrapper.emitted('update:modelValue')
      expect(emitted).toBeTruthy()
      expect(emitted![0][0]).toBe('bar')
    })

    it('emits change event with full item data', async () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: '' },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const emitted = wrapper.emitted('change')
      expect(emitted).toBeTruthy()
      expect(emitted![0][0]).toHaveProperty('key')
      expect(emitted![0][0]).toHaveProperty('name')
      expect(emitted![0][0]).toHaveProperty('category')
    })

    it('updates selection when modelValue prop changes', async () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      await wrapper.setProps({ modelValue: 'line' })
      await nextTick()
      
      const emitted = wrapper.emitted('update:modelValue')
      expect(emitted).toBeTruthy()
      const events = emitted!
      const lastEmit = events[events.length - 1]
      expect(lastEmit[0]).toBe('line')
    })
  })

  describe('Chart Groups', () => {
    it('contains basic chart group', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      const vm = wrapper.vm as any
      const basicGroup = vm.chartGroups.find((g: any) => g.key === 'basic')
      expect(basicGroup).toBeDefined()
      expect(basicGroup.items.length).toBe(4)
    })

    it('contains pie-extended chart group', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      const vm = wrapper.vm as any
      const pieGroup = vm.chartGroups.find((g: any) => g.key === 'pie-extended')
      expect(pieGroup).toBeDefined()
      expect(pieGroup.items.length).toBe(3)
    })

    it('contains relation chart group', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      const vm = wrapper.vm as any
      const relationGroup = vm.chartGroups.find((g: any) => g.key === 'relation')
      expect(relationGroup).toBeDefined()
      expect(relationGroup.items.length).toBe(3)
    })

    it('contains geo chart group', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      const vm = wrapper.vm as any
      const geoGroup = vm.chartGroups.find((g: any) => g.key === 'geo')
      expect(geoGroup).toBeDefined()
      expect(geoGroup.items.length).toBe(3)
    })

    it('contains combo chart group', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      const vm = wrapper.vm as any
      const comboGroup = vm.chartGroups.find((g: any) => g.key === 'combo')
      expect(comboGroup).toBeDefined()
      expect(comboGroup.items.length).toBe(2)
    })
  })

  describe('Computed Properties', () => {
    it('selectedType returns matched chart type', async () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'pie' },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.selectedType.key).toBe('pie')
      expect(vm.selectedType.name).toBe('饼图')
    })

    it('selectedType returns first item when no match', async () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'nonexistent' },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.selectedType.key).toBe('bar')
    })

    it('allChartTypes contains all chart items', () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'bar' },
        global: { stubs: globalStubs },
      })
      const vm = wrapper.vm as any
      expect(vm.allChartTypes.length).toBe(15)
    })
  })

  describe('Edge Cases', () => {
    it('handles empty modelValue', async () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: '' },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.currentTypeKey).toBe('bar')
    })

    it('handles invalid modelValue gracefully', async () => {
      const wrapper = mount(ChartTypeSelector, {
        props: { modelValue: 'invalid-type' },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.currentTypeKey).toBe('bar')
    })
  })
})
