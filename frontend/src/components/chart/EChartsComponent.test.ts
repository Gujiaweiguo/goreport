import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

vi.mock('echarts/core', () => ({
  init: vi.fn(() => ({
    setOption: vi.fn(),
    resize: vi.fn(),
    dispose: vi.fn(),
  })),
  use: vi.fn(),
}))

vi.mock('echarts/components', () => ({
  TitleComponent: {},
  TooltipComponent: {},
  LegendComponent: {},
  GridComponent: {},
}))

vi.mock('echarts/charts', () => ({
  BarChart: {},
  LineChart: {},
  PieChart: {},
}))

vi.mock('echarts/renderers', () => ({
  CanvasRenderer: {},
}))

import EChartsComponent from './EChartsComponent.vue'
import * as echarts from 'echarts/core'

const defaultOption = {
  title: { text: 'Test Chart' },
  xAxis: { type: 'category' },
  yAxis: { type: 'value' },
  series: [{ type: 'bar', data: [1, 2, 3] }],
} as any

describe('EChartsComponent.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  describe('Component Rendering', () => {
    it('renders a div element', () => {
      const wrapper = mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      expect(wrapper.find('div').exists()).toBe(true)
    })

    it('has default width of 100%', () => {
      const wrapper = mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      const div = wrapper.find('div')
      expect(div.element.style.width).toBe('100%')
    })

    it('has default height of 400px', () => {
      const wrapper = mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      const div = wrapper.find('div')
      expect(div.element.style.height).toBe('400px')
    })

    it('accepts custom width prop', () => {
      const wrapper = mount(EChartsComponent, {
        props: { option: defaultOption, width: '500px' },
      })
      const div = wrapper.find('div')
      expect(div.element.style.width).toBe('500px')
    })

    it('accepts custom height prop', () => {
      const wrapper = mount(EChartsComponent, {
        props: { option: defaultOption, height: '300px' },
      })
      const div = wrapper.find('div')
      expect(div.element.style.height).toBe('300px')
    })
  })

  describe('Chart Initialization', () => {
    it('initializes chart on mount', async () => {
      mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      await nextTick()
      expect(echarts.init).toHaveBeenCalled()
    })

    it('calls setOption with initial option on mount', async () => {
      const mockInstance = {
        setOption: vi.fn(),
        resize: vi.fn(),
        dispose: vi.fn(),
      }
      vi.mocked(echarts.init).mockReturnValueOnce(mockInstance as any)
      
      mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      await nextTick()
      expect(mockInstance.setOption).toHaveBeenCalledWith(defaultOption, true)
    })
  })

  describe('Option Updates', () => {
    it('updates chart when option prop changes', async () => {
      const mockInstance = {
        setOption: vi.fn(),
        resize: vi.fn(),
        dispose: vi.fn(),
      }
      vi.mocked(echarts.init).mockReturnValueOnce(mockInstance as any)
      
      const wrapper = mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      await nextTick()
      mockInstance.setOption.mockClear()

      const newOption = {
        ...defaultOption,
        title: { text: 'Updated Chart' },
      }
      await wrapper.setProps({ option: newOption })
      await nextTick()

      expect(mockInstance.setOption).toHaveBeenCalledWith(newOption, true)
    })
  })

  describe('Resize Handling', () => {
    it('adds resize event listener on mount', async () => {
      const addSpy = vi.spyOn(window, 'addEventListener')
      mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      await nextTick()
      expect(addSpy).toHaveBeenCalledWith('resize', expect.any(Function))
      addSpy.mockRestore()
    })

    it('removes resize event listener on unmount', async () => {
      const removeSpy = vi.spyOn(window, 'removeEventListener')
      const wrapper = mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      await nextTick()
      wrapper.unmount()
      expect(removeSpy).toHaveBeenCalledWith('resize', expect.any(Function))
      removeSpy.mockRestore()
    })

    it('calls resize on window resize event', async () => {
      const mockInstance = {
        setOption: vi.fn(),
        resize: vi.fn(),
        dispose: vi.fn(),
      }
      vi.mocked(echarts.init).mockReturnValueOnce(mockInstance as any)
      
      mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      await nextTick()
      mockInstance.resize.mockClear()

      window.dispatchEvent(new Event('resize'))
      expect(mockInstance.resize).toHaveBeenCalled()
    })
  })

  describe('Cleanup', () => {
    it('disposes chart on unmount', async () => {
      const mockInstance = {
        setOption: vi.fn(),
        resize: vi.fn(),
        dispose: vi.fn(),
      }
      vi.mocked(echarts.init).mockReturnValueOnce(mockInstance as any)
      
      const wrapper = mount(EChartsComponent, {
        props: { option: defaultOption },
      })
      await nextTick()
      wrapper.unmount()
      expect(mockInstance.dispose).toHaveBeenCalled()
    })
  })

  describe('Edge Cases', () => {
    it('handles empty option', async () => {
      const mockInstance = {
        setOption: vi.fn(),
        resize: vi.fn(),
        dispose: vi.fn(),
      }
      vi.mocked(echarts.init).mockReturnValueOnce(mockInstance as any)
      
      mount(EChartsComponent, {
        props: { option: {} },
      })
      await nextTick()
      expect(mockInstance.setOption).toHaveBeenCalledWith({}, true)
    })

    it('handles complex nested option', async () => {
      const mockInstance = {
        setOption: vi.fn(),
        resize: vi.fn(),
        dispose: vi.fn(),
      }
      vi.mocked(echarts.init).mockReturnValueOnce(mockInstance as any)
      
      const complexOption = {
        title: { text: 'Complex', subtext: 'Subtitle' },
        tooltip: { trigger: 'axis' },
        legend: { data: ['Sales', 'Expenses'] },
        xAxis: [
          { type: 'category', data: ['Mon', 'Tue', 'Wed'] },
        ],
        yAxis: [
          { type: 'value', name: 'Sales' },
          { type: 'value', name: 'Expenses' },
        ],
        series: [
          { name: 'Sales', type: 'bar', data: [2, 4, 6] },
          { name: 'Expenses', type: 'line', data: [1, 3, 5] },
        ],
      } as any
      mount(EChartsComponent, {
        props: { option: complexOption },
      })
      await nextTick()
      expect(mockInstance.setOption).toHaveBeenCalledWith(complexOption, true)
    })
  })
})
