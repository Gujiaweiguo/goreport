<template>
  <div ref="previewRoot" class="dashboard-preview">
    <div class="preview-toolbar">
      <div class="toolbar-title">{{ title || '大屏预览' }}</div>
      <div class="toolbar-actions">
        <el-button :icon="Refresh" text :loading="isRefreshing" class="toolbar-btn" @click="handleRefresh">
          刷新
        </el-button>
        <el-button :icon="FullScreen" text class="toolbar-btn" @click="toggleFullscreen">
          {{ isFullscreen ? '退出全屏' : '全屏' }}
        </el-button>
        <el-button :icon="Close" text type="danger" class="toolbar-btn" @click="handleClose">
          关闭
        </el-button>
      </div>
    </div>

    <div class="preview-stage">
      <div
        v-for="component in visibleComponents"
        :key="component.id"
        class="preview-component"
        :class="`type-${component.type}`"
        :style="getComponentStyle(component)"
      >
        <div class="component-content" :style="component.style">
          <template v-if="component.type === 'text'">
            <p class="text-title">{{ component.title }}</p>
          </template>
          <template v-else-if="component.type === 'chart'">
            <div class="chart-container">
              <p class="component-label">图表</p>
              <p class="component-title">{{ component.title }}</p>
              <p class="component-hint">数据: {{ component.data?.dataSource || '未配置' }}</p>
            </div>
          </template>
          <template v-else-if="component.type === 'table'">
            <div class="table-container">
              <p class="component-label">表格</p>
              <p class="component-title">{{ component.title }}</p>
              <p class="component-hint">{{ component.data?.field || '未配置字段' }}</p>
            </div>
          </template>
          <template v-else-if="component.type === 'image'">
            <div class="image-container">
              <p class="component-label">图片</p>
              <p class="component-title">{{ component.title }}</p>
            </div>
          </template>
          <template v-else>
            <p class="component-label">{{ component.type }}</p>
            <p class="component-title">{{ component.title }}</p>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Close, FullScreen, Refresh } from '@element-plus/icons-vue'
import type { DashboardComponent } from '@/components/dashboard/PropertyPanel.vue'

const props = defineProps<{
  components: DashboardComponent[]
  title: string
}>()

const emit = defineEmits<{
  close: []
}>()

const previewRoot = ref<HTMLElement | null>(null)
const isRefreshing = ref(false)
const isFullscreen = ref(false)

const visibleComponents = computed(() => props.components.filter(component => component.visible))

function getComponentStyle(component: DashboardComponent) {
  return {
    left: `${component.x}px`,
    top: `${component.y}px`,
    width: `${component.width}px`,
    height: `${component.height}px`
  }
}

async function handleRefresh() {
  if (isRefreshing.value) return
  isRefreshing.value = true
  await new Promise(resolve => setTimeout(resolve, 700))
  isRefreshing.value = false
  ElMessage.success('预览已刷新')
}

async function toggleFullscreen() {
  if (!previewRoot.value) return
  try {
    if (!document.fullscreenElement) {
      await previewRoot.value.requestFullscreen()
      isFullscreen.value = true
      return
    }
    await document.exitFullscreen()
    isFullscreen.value = false
  } catch {
    ElMessage.error('全屏切换失败')
  }
}

async function exitFullscreen() {
  if (!document.fullscreenElement) return
  try {
    await document.exitFullscreen()
    isFullscreen.value = false
  } catch {
    ElMessage.error('退出全屏失败')
  }
}

async function handleClose() {
  await exitFullscreen()
  emit('close')
}

function syncFullscreenState() {
  isFullscreen.value = Boolean(document.fullscreenElement)
}

onMounted(() => {
  document.addEventListener('fullscreenchange', syncFullscreenState)
})

onBeforeUnmount(() => {
  document.removeEventListener('fullscreenchange', syncFullscreenState)
})

defineExpose({
  exitFullscreen
})
</script>

<style scoped>
.dashboard-preview {
  position: relative;
  width: 100%;
  height: 100%;
  background: #0a0e27;
  overflow: hidden;
}

.preview-toolbar {
  position: absolute;
  top: 16px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20px;
  min-width: 560px;
  padding: 10px 14px;
  border-radius: 14px;
  background: rgba(0, 0, 0, 0.48);
  border: 1px solid rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(6px);
}

.toolbar-title {
  font-size: 14px;
  font-weight: 600;
  color: #f5f7fa;
  letter-spacing: 0.02em;
  white-space: nowrap;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-btn {
  color: #ffffff;
}

.toolbar-btn:deep(.el-icon) {
  margin-right: 4px;
}

.preview-stage {
  position: relative;
  width: 100%;
  height: 100%;
  padding-top: 72px;
  box-sizing: border-box;
}

.preview-component {
  position: absolute;
  border-radius: 8px;
  box-shadow: 0 16px 36px rgba(5, 8, 24, 0.35);
  box-sizing: border-box;
  overflow: hidden;
}

.component-content {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 12px;
}

.text-title {
  margin: 0;
  font-size: 24px;
  font-weight: 700;
  color: #141b33;
  letter-spacing: 0.03em;
}

.component-label {
  margin: 0;
  color: #909399;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.component-title {
  margin: 8px 0 0;
  color: #303133;
  font-size: 16px;
  font-weight: 600;
}

.component-hint {
  margin: 4px 0 0;
  color: #c0c4cc;
  font-size: 12px;
}

.chart-container,
.table-container,
.image-container {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>
