<template>
  <div class="dashboard-designer" :class="{ 'preview-layout': isPreviewMode }">
    <DashboardPreview
      v-if="isPreviewMode"
      ref="previewRef"
      :components="components"
      title="业务大屏预览"
      @close="exitPreviewMode"
    />

    <template v-else>
      <aside class="designer-left">
      <div class="panel-header">
        <h3>组件面板</h3>
      </div>
      <div class="panel-content library-panel">
        <ComponentLibrary @dragstart="handleLibraryDragStart" @dragend="handleLibraryDragEnd" />
      </div>
      </aside>

      <div class="designer-canvas">
        <div class="canvas-toolbar">
          <el-button-group>
            <el-button :icon="ZoomIn">放大</el-button>
            <el-button :icon="ZoomOut">缩小</el-button>
            <el-button :icon="Refresh">重置</el-button>
          </el-button-group>
          <div class="toolbar-divider"></div>
          <el-button-group>
            <el-button :icon="FolderOpened" @click="handleLoadDashboard">加载</el-button>
            <el-button :icon="Document" :loading="savingDashboard" @click="handleSaveDashboard">保存</el-button>
            <el-button :icon="Delete" @click="handleClearDashboard" type="danger">清空</el-button>
          </el-button-group>
          <div class="mock-component-strip">
            <span class="strip-label">画布组件（模拟点击）</span>
            <el-tag
              v-for="component in components"
              :key="component.id"
              size="small"
              :type="selectedComponent?.id === component.id ? 'primary' : 'info'"
              effect="plain"
              class="component-tag"
              @click="handleComponentSelect(component.id)"
            >
              {{ component.title }}
            </el-tag>
          </div>
          <el-button :icon="View" type="primary" @click="enterPreviewMode">预览</el-button>
        </div>
        <div
          ref="canvasAreaRef"
          class="canvas-area"
          :class="{
            'drag-active': isDragging,
            'drag-invalid': isDragging && isDropInvalid
          }"
          @dragover.prevent="handleCanvasDragOver"
          @dragleave="handleCanvasDragLeave"
          @drop.prevent="handleCanvasDrop"
        >
          <div v-if="isDragging" class="drop-indicator">
            <span>{{ isDropInvalid ? '无效放置区域' : '释放以添加组件' }}</span>
          </div>
          <div class="empty-state">
            <el-icon :size="64"><Monitor /></el-icon>
            <p>{{ isDragging ? '释放以添加组件' : '拖拽左侧组件到画布开始设计' }}</p>
          </div>
        </div>
      </div>

      <aside class="designer-right">
        <el-tabs v-model="activeTab" type="border-card">
          <el-tab-pane label="属性" name="properties">
            <div class="panel-content">
              <PropertyPanel
                :component="selectedComponent"
                :dataSources="dataSources"
                :loadingDataSources="loadingDataSources"
                @update="handlePropertyUpdate"
              />
            </div>
          </el-tab-pane>
          <el-tab-pane label="图层" name="layers">
            <div class="panel-content">
              <LayerPanel
                :layers="layers"
                :selected-id="selectedComponent?.id"
                @select="handleLayerSelect"
                @toggle-visibility="handleToggleVisibility"
                @toggle-lock="handleToggleLock"
                @delete="handleDeleteLayer"
                @reorder="handleReorder"
              />
            </div>
          </el-tab-pane>
        </el-tabs>
      </aside>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { Monitor, ZoomIn, ZoomOut, Refresh, View, FolderOpened, Document, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ComponentLibrary from '@/components/dashboard/ComponentLibrary.vue'
import PropertyPanel from '@/components/dashboard/PropertyPanel.vue'
import LayerPanel from '@/components/dashboard/LayerPanel.vue'
import DashboardPreview from '@/components/dashboard/DashboardPreview.vue'
import type { DashboardComponent } from '@/components/dashboard/PropertyPanel.vue'
import type { LayerItem } from '@/components/dashboard/LayerPanel.vue'
import { datasourceApi, type SelectOption } from '@/api/datasource'
import { dashboardApi, type Dashboard, type CreateDashboardRequest } from '@/api/dashboard'

const activeTab = ref('properties')
const selectedComponentId = ref<string | null>(null)
const isPreviewMode = ref(false)
const previewRef = ref<{ exitFullscreen: () => Promise<void> } | null>(null)
const canvasAreaRef = ref<HTMLElement | null>(null)

const draggedComponentType = ref<string | null>(null)
const isDragging = ref(false)
const isDropInvalid = ref(false)

const loadingDataSources = ref(false)
const dataSources = ref<SelectOption[]>([])
const savingDashboard = ref(false)
const loadingDashboard = ref(false)

const components = ref<DashboardComponent[]>([
  {
    id: 'cmp_title_001',
    title: '销售总览',
    type: 'text',
    width: 480,
    height: 120,
    x: 80,
    y: 48,
    visible: true,
    locked: false,
    style: {
      background: '#ffffff',
      border: 'solid',
      borderColor: '#dcdfe6',
      fontSize: 24,
      fontColor: '#303133',
      padding: 16,
      margin: 0
    },
    data: {
      dataSource: 'sales_db',
      field: 'sales_amount',
      sqlQuery: 'SELECT sales_amount FROM overview LIMIT 1'
    },
    interaction: {
      linkUrl: '',
      drilldownConfig: ''
    }
  },
  {
    id: 'cmp_chart_001',
    title: '地区销售趋势',
    type: 'chart',
    width: 620,
    height: 320,
    x: 80,
    y: 190,
    visible: true,
    locked: false,
    style: {
      background: '#ffffff',
      border: 'solid',
      borderColor: '#dcdfe6',
      fontSize: 14,
      fontColor: '#303133',
      padding: 12,
      margin: 0
    },
    data: {
      dataSource: 'analytics_db',
      field: 'growth_rate',
      sqlQuery: 'SELECT region, growth_rate FROM trend_data'
    },
    interaction: {
      linkUrl: '/dashboard/detail',
      drilldownConfig: '{"target":"region-detail"}'
    }
  },
  {
    id: 'cmp_table_001',
    title: '订单明细',
    type: 'table',
    width: 620,
    height: 280,
    x: 740,
    y: 190,
    visible: true,
    locked: false,
    style: {
      background: '#ffffff',
      border: 'solid',
      borderColor: '#dcdfe6',
      fontSize: 13,
      fontColor: '#303133',
      padding: 12,
      margin: 0
    },
    data: {
      dataSource: 'sales_db',
      field: 'order_count',
      sqlQuery: 'SELECT * FROM orders LIMIT 100'
    },
    interaction: {
      linkUrl: '',
      drilldownConfig: ''
    }
  }
])

const layers = ref<LayerItem[]>(buildLayers(components.value))
selectedComponentId.value = components.value[0]?.id ?? null

const selectedComponent = computed(() => {
  if (!selectedComponentId.value) return null
  return components.value.find(component => component.id === selectedComponentId.value) ?? null
})

function buildLayers(componentList: DashboardComponent[]): LayerItem[] {
  return componentList.map(component => ({
    id: component.id,
    name: component.title,
    type: component.type,
    visible: component.visible,
    locked: component.locked
  }))
}

function syncLayersFromComponents() {
  layers.value = buildLayers(components.value)
}

async function loadDataSources() {
  loadingDataSources.value = true
  try {
    const response = await datasourceApi.list()
    if (response.data.success) {
      dataSources.value = (response.data.result?.datasources || []).map(ds => ({
        label: `${ds.name} (${ds.type})`,
        value: ds.id
      }))
      ElMessage.success('数据源列表加载成功')
    } else {
      ElMessage.error(response.data.message || '加载数据源列表失败')
    }
  } catch (error: any) {
    ElMessage.error('加载数据源列表失败')
  } finally {
    loadingDataSources.value = false
  }
}

function handleComponentSelect(componentId: string) {
  const target = components.value.find(component => component.id === componentId)
  if (!target || target.locked || !target.visible) return
  selectedComponentId.value = componentId
  activeTab.value = 'properties'
  if (target.data.dataSource && !loadingDataSources.value) {
    loadDataSources()
  }
}

function handlePropertyUpdate(updatedComponent: DashboardComponent) {
  const targetIndex = components.value.findIndex(component => component.id === selectedComponentId.value)
  if (targetIndex === -1) return
  components.value.splice(targetIndex, 1, {
    ...components.value[targetIndex],
    ...updatedComponent
  })
  selectedComponentId.value = updatedComponent.id
  syncLayersFromComponents()
}

function handleLayerSelect(layerId: string) {
  handleComponentSelect(layerId)
}

function handleToggleVisibility(layerId: string) {
  const target = components.value.find(component => component.id === layerId)
  if (!target) return
  target.visible = !target.visible
  if (!target.visible && selectedComponentId.value === layerId) {
    selectedComponentId.value = null
  }
  syncLayersFromComponents()
}

function handleToggleLock(layerId: string) {
  const target = components.value.find(component => component.id === layerId)
  if (!target) return
  target.locked = !target.locked
  if (target.locked && selectedComponentId.value === layerId) {
    selectedComponentId.value = null
  }
  syncLayersFromComponents()
}

function handleDeleteLayer(layerId: string) {
  const index = components.value.findIndex(component => component.id === layerId)
  if (index === -1) return
  components.value.splice(index, 1)
  if (selectedComponentId.value === layerId) {
    selectedComponentId.value = components.value[0]?.id ?? null
  }
  syncLayersFromComponents()
}

function handleReorder(nextLayers: LayerItem[]) {
  layers.value = [...nextLayers]
  const orderMap = new Map(nextLayers.map((layer, index) => [layer.id, index]))
  components.value = [...components.value].sort((a, b) => {
    return (orderMap.get(a.id) ?? 0) - (orderMap.get(b.id) ?? 0)
  })
}

const libraryComponentPresets: Record<
  string,
  { title: string; type: DashboardComponent['type']; width: number; height: number }
> = {
  'text-title': { title: '标题文本', type: 'text', width: 520, height: 90 },
  'text-basic': { title: '普通文本', type: 'text', width: 420, height: 140 },
  'text-rich': { title: '富文本', type: 'text', width: 520, height: 200 },
  'chart-bar': { title: '柱状图', type: 'chart', width: 560, height: 320 },
  'chart-line': { title: '折线图', type: 'chart', width: 560, height: 320 },
  'chart-pie': { title: '饼图', type: 'chart', width: 480, height: 320 },
  'chart-scatter': { title: '散点图', type: 'chart', width: 560, height: 320 },
  'table-basic': { title: '基础表格', type: 'table', width: 620, height: 300 },
  'table-page': { title: '分页表格', type: 'table', width: 620, height: 340 },
  'image-basic': { title: '图片', type: 'image', width: 420, height: 240 },
  'image-border': { title: '边框图片', type: 'image', width: 420, height: 240 },
  'image-background': { title: '背景图', type: 'image', width: 640, height: 360 },
  'decorative-line': { title: '分割线', type: 'decorative', width: 400, height: 28 },
  'decorative-frame': { title: '装饰框', type: 'decorative', width: 520, height: 280 },
  'decorative-corner': { title: '角标', type: 'decorative', width: 120, height: 120 }
}

function handleLibraryDragStart(componentType: string) {
  draggedComponentType.value = componentType
  isDragging.value = true
  isDropInvalid.value = false
}

function handleLibraryDragEnd() {
  resetDraggingState()
}

function handleCanvasDragOver(event: DragEvent) {
  if (!draggedComponentType.value || !canvasAreaRef.value) return
  if (event.dataTransfer) {
    event.dataTransfer.dropEffect = 'copy'
  }
  const isInBounds = isDropPointInCanvas(event)
  isDragging.value = true
  isDropInvalid.value = !isInBounds
}

function handleCanvasDragLeave(event: DragEvent) {
  const relatedTarget = event.relatedTarget as Node | null
  if (!canvasAreaRef.value?.contains(relatedTarget)) {
    isDropInvalid.value = false
  }
}

function handleCanvasDrop(event: DragEvent) {
  if (!draggedComponentType.value || !canvasAreaRef.value) {
    ElMessage.error('未识别拖拽组件类型，放置失败')
    resetDraggingState()
    return
  }

  const preset = libraryComponentPresets[draggedComponentType.value]
  if (!preset) {
    ElMessage.error('组件类型不受支持，放置失败')
    resetDraggingState()
    return
  }

  const canvasRect = canvasAreaRef.value.getBoundingClientRect()
  const dropX = event.clientX - canvasRect.left + canvasAreaRef.value.scrollLeft
  const dropY = event.clientY - canvasRect.top + canvasAreaRef.value.scrollTop

  if (!isDropPointInCanvas(event)) {
    ElMessage.error('无效放置位置，请在画布区域内放置组件')
    resetDraggingState()
    return
  }

  const maxX = Math.max(0, canvasAreaRef.value.clientWidth - preset.width)
  const maxY = Math.max(0, canvasAreaRef.value.clientHeight - preset.height)
  const safeX = Math.min(Math.max(0, Math.round(dropX - preset.width / 2)), maxX)
  const safeY = Math.min(Math.max(0, Math.round(dropY - preset.height / 2)), maxY)

  const newComponent: DashboardComponent = {
    id: `cmp_${preset.type}_${Date.now()}`,
    title: preset.title,
    type: preset.type,
    width: preset.width,
    height: preset.height,
    x: safeX,
    y: safeY,
    visible: true,
    locked: false,
    style: {
      background: '#ffffff',
      border: 'solid',
      borderColor: '#dcdfe6',
      fontSize: preset.type === 'text' ? 24 : 14,
      fontColor: '#303133',
      padding: 12,
      margin: 0
    },
    data: {
      dataSource: '',
      field: '',
      sqlQuery: ''
    },
    interaction: {
      linkUrl: '',
      drilldownConfig: ''
    }
  }

  components.value.push(newComponent)
  selectedComponentId.value = newComponent.id
  activeTab.value = 'properties'
  syncLayersFromComponents()
  ElMessage.success(`已添加组件：${preset.title}`)
  resetDraggingState()
}

async function handleSaveDashboard() {
  if (!components.value.length) {
    ElMessage.warning('请先添加组件再保存')
    return
  }
  savingDashboard.value = true
  try {
    const response = await dashboardApi.create({
      name: `业务大屏_${Date.now()}`,
      description: '业务数据大屏',
      config: {
        width: 1920,
        height: 1080,
        backgroundColor: '#0a0e27'
      },
      components: components.value.map(component => ({
        ...component
      }))
    })
    if (response.data.success) {
      ElMessage.success('大屏保存成功')
    } else {
      ElMessage.error(response.data.message || '保存失败')
    }
  } catch (error: any) {
    const backendMessage = error?.response?.data?.message
    ElMessage.error(backendMessage || error.message || '保存失败')
  } finally {
    savingDashboard.value = false
  }
}

async function handleLoadDashboard() {
  loadingDashboard.value = true
  try {
    const response = await dashboardApi.list()
    if (response.data.success && response.data.result && response.data.result.length > 0) {
      const latestDashboard = response.data.result[0]
      components.value = latestDashboard.components.map((component: any) => ({
        ...component,
        style: component.style || {
          background: '#ffffff',
          border: 'solid',
          borderColor: '#dcdfe6',
          fontSize: 14,
          fontColor: '#303133',
          padding: 12,
          margin: 0
        },
        data: component.data || {
          dataSource: '',
          field: '',
          sqlQuery: ''
        },
        interaction: component.interaction || {
          linkUrl: '',
          drilldownConfig: ''
        }
      }))
      syncLayersFromComponents()
      ElMessage.success(`已加载大屏：${latestDashboard.name}`)
    } else {
      ElMessage.warning('暂无已保存的大屏')
    }
  } catch (error: any) {
    const backendMessage = error?.response?.data?.message
    ElMessage.error(backendMessage || error.message || '加载失败')
  } finally {
    loadingDashboard.value = false
  }
}

async function handleClearDashboard() {
  try {
    await ElMessageBox.confirm('确认清空当前大屏所有组件吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    components.value = []
    selectedComponentId.value = null
    syncLayersFromComponents()
    ElMessage.success('大屏已清空')
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

function isDropPointInCanvas(event: DragEvent) {
  if (!canvasAreaRef.value) return false
  const rect = canvasAreaRef.value.getBoundingClientRect()
  return (
    event.clientX >= rect.left &&
    event.clientX <= rect.right &&
    event.clientY >= rect.top &&
    event.clientY <= rect.bottom
  )
}

function resetDraggingState() {
  draggedComponentType.value = null
  isDragging.value = false
  isDropInvalid.value = false
}

function enterPreviewMode() {
  isPreviewMode.value = true
}

async function exitPreviewMode() {
  await previewRef.value?.exitFullscreen()
  isPreviewMode.value = false
}

async function handleEscKey(event: KeyboardEvent) {
  if (event.key !== 'Escape' || !isPreviewMode.value) return
  await exitPreviewMode()
}

onMounted(() => {
  window.addEventListener('keydown', handleEscKey)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleEscKey)
})

loadDataSources()
</script>

<style scoped>
.dashboard-designer {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.dashboard-designer.preview-layout {
  flex-direction: column;
}

.designer-left,
.designer-right {
  width: 280px;
  background-color: #fff;
  border-right:1px solid #dcdfe6;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.designer-right {
  border-right: none;
  border-left: 1px solid #dcdfe6;
}

.panel-header {
  padding: 12px 16px;
  border-bottom: 1px solid #dcdfe6;
  background-color: #f5f7fa;
}

.panel-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
}

.panel-content {
  flex: 1;
  overflow: auto;
  padding: 16px;
}

.library-panel {
  padding: 10px 12px;
}

.designer-canvas {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #f0f2f5;
  overflow: hidden;
}

.canvas-toolbar {
  padding: 12px;
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 8px;
  flex-wrap: wrap;
}

.toolbar-divider {
  width: 1px;
  height: 24px;
  background: #dcdfe6;
  margin: 0 8px;
}

.toolbar-divider {
  width: 1px;
  height: 24px;
  background: #dcdfe6;
  margin: 0 8px;
}

.mock-component-strip {
  min-width: 0;
  flex: 1;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.strip-label {
  font-size: 12px;
  color: #909399;
}

.component-tag {
  cursor: pointer;
}

.canvas-area {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: auto;
  transition: background-color 0.2s ease;
}

.canvas-area.drag-active {
  background: rgba(64, 158, 255, 0.06);
}

.canvas-area.drag-active::after {
  content: '';
  position: absolute;
  top: 18px;
  right: 18px;
  bottom: 18px;
  left: 18px;
  border: 2px dashed #409eff;
  border-radius: 12px;
  pointer-events: none;
}

.canvas-area.drag-invalid::after {
  border-color: #f56c6c;
}

.drop-indicator {
  position: absolute;
  top: 32px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 2;
  padding: 8px 14px;
  border-radius: 999px;
  font-size: 12px;
  color: #ffffff;
  background: rgba(15, 23, 42, 0.72);
  box-shadow: 0 8px 18px rgba(15, 23, 42, 0.2);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #909399;
  gap: 12px;
  padding: 40px 0;
}

.empty-state p {
  margin: 0;
  font-size: 14px;
}
</style>
