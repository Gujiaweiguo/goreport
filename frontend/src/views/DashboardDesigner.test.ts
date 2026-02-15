import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'
import DashboardDesigner from './DashboardDesigner.vue'

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

vi.mock('@/api/datasource', () => ({
  datasourceApi: {
    list: vi.fn().mockResolvedValue({ data: { success: true, result: { datasources: [] } } })
  }
}))

vi.mock('@/api/dashboard', () => ({
  dashboardApi: {
    create: vi.fn().mockResolvedValue({ data: { success: true } }),
    list: vi.fn().mockResolvedValue({ data: { success: true, result: [] } })
  }
}))

import { ElMessage, ElMessageBox } from 'element-plus'
import { datasourceApi } from '@/api/datasource'
import { dashboardApi } from '@/api/dashboard'

const globalStubs = {
  'el-card': true,
  'el-button': true,
  'el-button-group': true,
  'el-tabs': true,
  'el-tab-pane': true,
  'el-icon': true,
  'el-tag': true,
  ComponentLibrary: true,
  PropertyPanel: true,
  LayerPanel: true,
  DashboardPreview: true
}

describe('DashboardDesigner.vue', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(datasourceApi.list).mockResolvedValue({ 
      data: { success: true, result: { datasources: [] } } 
    } as any)
  })

  describe('Component Rendering', () => {
    it('renders correctly', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.find('.dashboard-designer').exists()).toBe(true)
    })
  })

  describe('Initial state', () => {
    it('activeTab starts as properties', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.activeTab).toBe('properties')
    })

    it('isPreviewMode starts false', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.isPreviewMode).toBe(false)
    })

    it('isDragging starts false', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.isDragging).toBe(false)
    })

    it('isDropInvalid starts false', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.isDropInvalid).toBe(false)
    })

    it('savingDashboard starts false', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.savingDashboard).toBe(false)
    })

    it('loadingDashboard starts false', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.loadingDashboard).toBe(false)
    })

    it('components starts with preset items', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.components.length).toBeGreaterThan(0)
    })

    it('selectedComponentId starts with first component id', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.selectedComponentId).toBe(wrapper.vm.components[0]?.id)
    })
  })

  describe('selectedComponent computed', () => {
    it('returns null when no component selected', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.selectedComponentId = null
      expect(wrapper.vm.selectedComponent).toBeNull()
    })

    it('returns component when selected', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const firstComponent = wrapper.vm.components[0]
      wrapper.vm.selectedComponentId = firstComponent.id
      expect(wrapper.vm.selectedComponent).toEqual(firstComponent)
    })
  })

  describe('buildLayers function', () => {
    it('creates layers from components', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const layers = wrapper.vm.buildLayers(wrapper.vm.components)
      expect(layers.length).toBe(wrapper.vm.components.length)
      expect(layers[0]).toHaveProperty('id')
      expect(layers[0]).toHaveProperty('name')
      expect(layers[0]).toHaveProperty('type')
      expect(layers[0]).toHaveProperty('visible')
      expect(layers[0]).toHaveProperty('locked')
    })
  })

  describe('syncLayersFromComponents function', () => {
    it('syncs layers with components', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const initialLayersCount = wrapper.vm.layers.length
      wrapper.vm.syncLayersFromComponents()
      expect(wrapper.vm.layers.length).toBe(initialLayersCount)
    })
  })

  describe('handleComponentSelect function', () => {
    it('selects component when visible and unlocked', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const componentId = wrapper.vm.components[0].id
      wrapper.vm.handleComponentSelect(componentId)
      expect(wrapper.vm.selectedComponentId).toBe(componentId)
      expect(wrapper.vm.activeTab).toBe('properties')
    })

    it('does not select locked component', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const component = wrapper.vm.components[0]
      component.locked = true
      wrapper.vm.selectedComponentId = null
      wrapper.vm.handleComponentSelect(component.id)
      expect(wrapper.vm.selectedComponentId).toBeNull()
    })

    it('does not select hidden component', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const component = wrapper.vm.components[0]
      component.visible = false
      wrapper.vm.selectedComponentId = null
      wrapper.vm.handleComponentSelect(component.id)
      expect(wrapper.vm.selectedComponentId).toBeNull()
    })
  })

  describe('handlePropertyUpdate function', () => {
    it('updates component properties', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const component = wrapper.vm.components[0]
      wrapper.vm.selectedComponentId = component.id
      
      const updated = { ...component, title: 'Updated Title' }
      wrapper.vm.handlePropertyUpdate(updated)
      
      expect(wrapper.vm.components[0].title).toBe('Updated Title')
    })
  })

  describe('handleLayerSelect function', () => {
    it('calls handleComponentSelect with layer id', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const componentId = wrapper.vm.components[0].id
      wrapper.vm.handleLayerSelect(componentId)
      expect(wrapper.vm.selectedComponentId).toBe(componentId)
    })
  })

  describe('handleToggleVisibility function', () => {
    it('toggles component visibility', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const component = wrapper.vm.components[0]
      const initialVisibility = component.visible
      wrapper.vm.handleToggleVisibility(component.id)
      expect(component.visible).toBe(!initialVisibility)
    })

    it('deselects component when hidden', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const component = wrapper.vm.components[0]
      component.visible = true
      wrapper.vm.selectedComponentId = component.id
      wrapper.vm.handleToggleVisibility(component.id)
      expect(wrapper.vm.selectedComponentId).toBeNull()
    })
  })

  describe('handleToggleLock function', () => {
    it('toggles component lock state', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const component = wrapper.vm.components[0]
      const initialLock = component.locked
      wrapper.vm.handleToggleLock(component.id)
      expect(component.locked).toBe(!initialLock)
    })

    it('deselects component when locked', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const component = wrapper.vm.components[0]
      component.locked = false
      wrapper.vm.selectedComponentId = component.id
      wrapper.vm.handleToggleLock(component.id)
      expect(wrapper.vm.selectedComponentId).toBeNull()
    })
  })

  describe('handleDeleteLayer function', () => {
    it('deletes component', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const initialCount = wrapper.vm.components.length
      const componentId = wrapper.vm.components[0].id
      wrapper.vm.handleDeleteLayer(componentId)
      expect(wrapper.vm.components.length).toBe(initialCount - 1)
    })

    it('selects first component after deletion of selected', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const component = wrapper.vm.components[0]
      wrapper.vm.selectedComponentId = component.id
      wrapper.vm.handleDeleteLayer(component.id)
      expect(wrapper.vm.selectedComponentId).toBe(wrapper.vm.components[0]?.id ?? null)
    })
  })

  describe('handleReorder function', () => {
    it('reorders layers and components', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const reversedLayers = [...wrapper.vm.layers].reverse()
      wrapper.vm.handleReorder(reversedLayers)
      expect(wrapper.vm.layers[0].id).toBe(reversedLayers[0].id)
    })
  })

  describe('handleLibraryDragStart function', () => {
    it('sets drag state', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.handleLibraryDragStart('text-title')
      expect(wrapper.vm.draggedComponentType).toBe('text-title')
      expect(wrapper.vm.isDragging).toBe(true)
      expect(wrapper.vm.isDropInvalid).toBe(false)
    })
  })

  describe('handleLibraryDragEnd function', () => {
    it('resets drag state', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.draggedComponentType = 'text-title'
      wrapper.vm.isDragging = true
      wrapper.vm.handleLibraryDragEnd()
      expect(wrapper.vm.draggedComponentType).toBeNull()
      expect(wrapper.vm.isDragging).toBe(false)
    })
  })

  describe('resetDraggingState function', () => {
    it('resets all drag related state', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.draggedComponentType = 'test'
      wrapper.vm.isDragging = true
      wrapper.vm.isDropInvalid = true
      wrapper.vm.resetDraggingState()
      expect(wrapper.vm.draggedComponentType).toBeNull()
      expect(wrapper.vm.isDragging).toBe(false)
      expect(wrapper.vm.isDropInvalid).toBe(false)
    })
  })

  describe('enterPreviewMode function', () => {
    it('sets isPreviewMode to true', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.enterPreviewMode()
      expect(wrapper.vm.isPreviewMode).toBe(true)
    })
  })

  describe('exitPreviewMode function', () => {
    it('sets isPreviewMode to false', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.isPreviewMode = true
      await wrapper.vm.exitPreviewMode()
      expect(wrapper.vm.isPreviewMode).toBe(false)
    })
  })

  describe('handleSaveDashboard function', () => {
    it('shows warning when no components', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      wrapper.vm.components = []
      await wrapper.vm.handleSaveDashboard()
      expect(ElMessage.warning).toHaveBeenCalledWith('请先添加组件再保存')
    })

    it('calls dashboardApi.create when components exist', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      await wrapper.vm.handleSaveDashboard()
      expect(dashboardApi.create).toHaveBeenCalled()
      expect(ElMessage.success).toHaveBeenCalledWith('大屏保存成功')
    })

    it('shows error on save failure', async () => {
      vi.mocked(dashboardApi.create).mockResolvedValueOnce({ 
        data: { success: false, message: 'Save failed' } 
      } as any)
      
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      await wrapper.vm.handleSaveDashboard()
      expect(ElMessage.error).toHaveBeenCalledWith('Save failed')
    })
  })

  describe('handleLoadDashboard function', () => {
    it('shows warning when no saved dashboards', async () => {
      vi.mocked(dashboardApi.list).mockResolvedValueOnce({ 
        data: { success: true, result: [] } 
      } as any)
      
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      vi.clearAllMocks()
      
      await wrapper.vm.handleLoadDashboard()
      expect(ElMessage.warning).toHaveBeenCalledWith('暂无已保存的大屏')
    })

    it('loads dashboard components', async () => {
      const mockDashboard = {
        name: 'Test Dashboard',
        components: [
          { id: 'c1', title: 'Component 1', type: 'text', x: 0, y: 0, width: 100, height: 50, visible: true, locked: false }
        ]
      }
      vi.mocked(dashboardApi.list).mockResolvedValueOnce({ 
        data: { success: true, result: [mockDashboard] } 
      } as any)
      
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      vi.clearAllMocks()
      
      await wrapper.vm.handleLoadDashboard()
      expect(ElMessage.success).toHaveBeenCalled()
    })
  })

  describe('handleClearDashboard function', () => {
    it('clears all components after confirmation', async () => {
      vi.mocked(ElMessageBox.confirm).mockResolvedValueOnce('confirm' as any)
      
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      
      await wrapper.vm.handleClearDashboard()
      
      expect(wrapper.vm.components.length).toBe(0)
      expect(wrapper.vm.selectedComponentId).toBeNull()
      expect(ElMessage.success).toHaveBeenCalledWith('大屏已清空')
    })

    it('does not clear when cancelled', async () => {
      vi.mocked(ElMessageBox.confirm).mockRejectedValueOnce('cancel')
      
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      const initialCount = wrapper.vm.components.length
      
      await wrapper.vm.handleClearDashboard()
      
      expect(wrapper.vm.components.length).toBe(initialCount)
    })
  })

  describe('libraryComponentPresets', () => {
    it('has text-title preset', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.libraryComponentPresets['text-title']).toBeDefined()
      expect(wrapper.vm.libraryComponentPresets['text-title'].type).toBe('text')
    })

    it('has chart-bar preset', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.libraryComponentPresets['chart-bar']).toBeDefined()
      expect(wrapper.vm.libraryComponentPresets['chart-bar'].type).toBe('chart')
    })

    it('has table-basic preset', async () => {
      const wrapper = mount(DashboardDesigner, { global: { stubs: globalStubs } })
      await nextTick()
      expect(wrapper.vm.libraryComponentPresets['table-basic']).toBeDefined()
      expect(wrapper.vm.libraryComponentPresets['table-basic'].type).toBe('table')
    })
  })
})
