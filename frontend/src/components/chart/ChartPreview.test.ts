import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import ChartPreview from './ChartPreview.vue'

vi.mock('echarts/core', () => ({
  init: vi.fn(() => ({
    setOption: vi.fn(),
    resize: vi.fn(),
    dispose: vi.fn(),
  })),
  use: vi.fn(),
}))

const globalStubs = {
  'el-tag': {
    template: '<span class="el-tag"><slot /></span>',
    props: ['size', 'effect', 'type'],
  },
  'el-skeleton': {
    template: '<div class="el-skeleton">Loading...</div>',
    props: ['animated', 'rows'],
  },
  EmptyState: {
    template: '<div class="empty-state"><slot /></div>',
    props: ['icon', 'iconSize', 'text', 'description'],
  },
  EChartsRenderer: {
    template: '<div class="echarts-renderer"></div>',
    props: ['option'],
    methods: {
      setOption: vi.fn(),
      resize: vi.fn(),
    },
  },
}

const defaultChartConfig = {
  title: 'Test Chart',
  width: 400,
  height: 300,
  margin: 10,
  theme: 'default' as const,
  color: '#2d6cdf',
  fontFamily: 'Arial',
  fontSize: 12,
  showLegend: true,
  hoverable: true,
  clickable: true,
  zoomable: false,
  mainTitle: 'Main Title',
  subTitle: 'Sub Title',
  titlePosition: 'center' as const,
}

const defaultChartData = {
  categories: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri'],
  series: [
    { type: 'bar', name: 'Sales', data: [120, 200, 150, 80, 70] },
  ],
}

describe('ChartPreview.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('Component Rendering', () => {
    it('renders the component', () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.chart-preview').exists()).toBe(true)
    })

    it('renders preview header', () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.preview-header').exists()).toBe(true)
      expect(wrapper.find('.preview-header span').text()).toBe('实时预览')
    })

    it('renders preview body', () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.preview-body').exists()).toBe(true)
    })
  })

  describe('Loading State', () => {
    it('shows loading skeleton when loading is true', () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
          loading: true,
        },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.loading-layer').exists()).toBe(true)
      expect(wrapper.find('.el-skeleton').exists()).toBe(true)
    })

    it('hides loading skeleton when loading is false', () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
          loading: false,
        },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.loading-layer').exists()).toBe(false)
    })
  })

  describe('Empty State', () => {
    it('shows empty state when no chart data', () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: { categories: [], series: [] },
        },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.empty-layer').exists()).toBe(true)
    })

    it('hides empty state when chart data exists', () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.empty-layer').exists()).toBe(false)
    })
  })

  describe('Dark Theme', () => {
    it('applies dark class when theme is dark', () => {
      const darkConfig = { ...defaultChartConfig, theme: 'dark' as const }
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: darkConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.chart-preview.dark').exists()).toBe(true)
    })

    it('does not apply dark class when theme is default', () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      expect(wrapper.find('.chart-preview.dark').exists()).toBe(false)
    })
  })

  describe('Preview Style', () => {
    it('computes preview style based on config', () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      const previewBody = wrapper.find('.preview-body')
      expect((previewBody.element as HTMLElement).style.minHeight).toBe('300px')
      expect((previewBody.element as HTMLElement).style.padding).toBe('10px')
    })

    it('updates preview style when config changes', async () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      
      const newConfig = { ...defaultChartConfig, height: 500, margin: 20 }
      await wrapper.setProps({ chartConfig: newConfig })
      
      const previewBody = wrapper.find('.preview-body')
      expect((previewBody.element as HTMLElement).style.minHeight).toBe('500px')
      expect((previewBody.element as HTMLElement).style.padding).toBe('20px')
    })
  })

  describe('Chart Option Building', () => {
    it('builds option with axis for bar chart', async () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.currentOption).toBeDefined()
      expect(vm.currentOption.xAxis).toBeDefined()
      expect(vm.currentOption.yAxis).toBeDefined()
    })

    it('builds option without axis for pie chart', async () => {
      const pieData = {
        categories: [],
        series: [{ type: 'pie', name: 'Share', data: [{ value: 100, name: 'A' }] }],
      }
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: pieData,
        },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.currentOption.series).toBeDefined()
    })

    it('includes dataZoom when zoomable is true', async () => {
      const zoomableConfig = { ...defaultChartConfig, zoomable: true }
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: zoomableConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.currentOption.dataZoom).toBeDefined()
    })
  })

  describe('Reactivity', () => {
    it('updates chart when chartConfig changes', async () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const newConfig = { ...defaultChartConfig, title: 'Updated Title' }
      await wrapper.setProps({ chartConfig: newConfig })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.currentOption).toBeDefined()
    })

    it('updates chart when chartData changes', async () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: defaultChartData,
        },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const newData = {
        categories: ['A', 'B', 'C'],
        series: [{ type: 'line', name: 'New', data: [1, 2, 3] }],
      }
      await wrapper.setProps({ chartData: newData })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.currentOption.series[0].type).toBe('line')
    })
  })

  describe('Edge Cases', () => {
    it('handles empty series array', async () => {
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: { categories: ['A', 'B'], series: [] },
        },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      expect(wrapper.find('.empty-layer').exists()).toBe(true)
    })

    it('handles gauge chart type', async () => {
      const gaugeData = {
        categories: [],
        series: [{ type: 'gauge', name: 'Progress', data: [{ value: 75 }] }],
      }
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: gaugeData,
        },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.currentOption.series).toBeDefined()
    })

    it('handles graph chart type', async () => {
      const graphData = {
        categories: [],
        series: [{ type: 'graph', name: 'Network', data: [] }],
      }
      const wrapper = mount(ChartPreview, {
        props: {
          chartConfig: defaultChartConfig,
          chartData: graphData,
        },
        global: { stubs: globalStubs },
      })
      await nextTick()
      
      const vm = wrapper.vm as any
      expect(vm.currentOption.xAxis).toBeUndefined()
    })
  })
})
