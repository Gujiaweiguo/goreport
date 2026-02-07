<template>
  <div class="layer-panel">
    <div class="panel-header">
      <h3>图层管理</h3>
      <el-tag size="small" type="info" effect="plain">{{ layers.length }}</el-tag>
    </div>

    <div v-if="!layers.length" class="panel-empty">
      <el-icon :size="22"><Collection /></el-icon>
      <span>暂无图层</span>
    </div>

    <div v-else class="layer-list">
      <div
        v-for="(layer, index) in layers"
        :key="layer.id"
        class="layer-item"
        :class="{
          active: layer.id === selectedId,
          locked: layer.locked,
          'drop-before': dragOverId === layer.id && dragPosition === 'before',
          'drop-after': dragOverId === layer.id && dragPosition === 'after'
        }"
        :draggable="!layer.locked"
        @click="handleSelect(layer.id)"
        @dragstart="handleDragStart($event, layer.id)"
        @dragover.prevent="handleDragOver($event, layer.id)"
        @dragenter.prevent
        @dragleave="handleDragLeave(layer.id)"
        @drop.prevent="handleDrop(layer.id)"
        @dragend="handleDragEnd"
      >
        <div class="layer-main">
          <el-icon class="type-icon"><component :is="resolveIcon(layer.type)" /></el-icon>
          <span class="layer-name" :title="layer.name">{{ layer.name }}</span>
          <span class="layer-level">{{ layers.length - index }}</span>
        </div>

        <div class="layer-actions" @click.stop>
          <el-button
            link
            :type="layer.visible ? 'primary' : 'info'"
            :title="layer.visible ? '隐藏图层' : '显示图层'"
            @click="emit('toggle-visibility', layer.id)"
          >
            <el-icon><component :is="layer.visible ? View : Hide" /></el-icon>
          </el-button>

          <el-button
            link
            :type="layer.locked ? 'warning' : 'info'"
            :title="layer.locked ? '解锁图层' : '锁定图层'"
            @click="emit('toggle-lock', layer.id)"
          >
            <el-icon><component :is="layer.locked ? Lock : Unlock" /></el-icon>
          </el-button>

          <el-button link :disabled="index === 0 || layer.locked" title="上移" @click="moveUp(index)">
            <el-icon><ArrowUpBold /></el-icon>
          </el-button>
          <el-button
            link
            :disabled="index === layers.length - 1 || layer.locked"
            title="下移"
            @click="moveDown(index)"
          >
            <el-icon><ArrowDownBold /></el-icon>
          </el-button>

          <el-button link type="danger" title="删除图层" @click="emit('delete', layer.id)">
            <el-icon><Delete /></el-icon>
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import {
  ArrowDownBold,
  ArrowUpBold,
  Collection,
  DataLine,
  Delete,
  Grid,
  Hide,
  Histogram,
  Lock,
  Picture,
  Unlock,
  View
} from '@element-plus/icons-vue'

export interface LayerItem {
  id: string
  name: string
  type: string
  visible: boolean
  locked: boolean
}

const props = defineProps<{
  layers: LayerItem[]
  selectedId?: string | null
}>()

const emit = defineEmits<{
  select: [id: string]
  'toggle-visibility': [id: string]
  'toggle-lock': [id: string]
  delete: [id: string]
  reorder: [layers: LayerItem[]]
}>()

const draggingId = ref<string | null>(null)
const dragOverId = ref<string | null>(null)
const dragPosition = ref<'before' | 'after' | null>(null)

const typeIconMap = computed(() => ({
  text: Grid,
  chart: Histogram,
  image: Picture,
  table: DataLine
}))

function resolveIcon(type: string) {
  return typeIconMap.value[type as keyof typeof typeIconMap.value] || Grid
}

function handleSelect(id: string) {
  emit('select', id)
}

function moveUp(index: number) {
  if (index <= 0) return
  const newLayers = [...props.layers]
  const current = newLayers[index]
  if (current.locked) return
  newLayers.splice(index, 1)
  newLayers.splice(index - 1, 0, current)
  emit('reorder', newLayers)
}

function moveDown(index: number) {
  if (index >= props.layers.length - 1) return
  const newLayers = [...props.layers]
  const current = newLayers[index]
  if (current.locked) return
  newLayers.splice(index, 1)
  newLayers.splice(index + 1, 0, current)
  emit('reorder', newLayers)
}

function handleDragStart(event: DragEvent, id: string) {
  draggingId.value = id
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/plain', id)
  }
}

function handleDragOver(event: DragEvent, overId: string) {
  if (!draggingId.value || draggingId.value === overId) return
  dragOverId.value = overId
  const target = event.currentTarget as HTMLElement
  const rect = target.getBoundingClientRect()
  const offsetY = event.clientY - rect.top
  dragPosition.value = offsetY < rect.height / 2 ? 'before' : 'after'
}

function handleDragLeave(overId: string) {
  if (dragOverId.value === overId) {
    dragOverId.value = null
    dragPosition.value = null
  }
}

function handleDrop(overId: string) {
  if (!draggingId.value || !dragPosition.value || draggingId.value === overId) {
    handleDragEnd()
    return
  }
  const draggedIndex = props.layers.findIndex(layer => layer.id === draggingId.value)
  const overIndex = props.layers.findIndex(layer => layer.id === overId)
  if (draggedIndex === -1 || overIndex === -1) {
    handleDragEnd()
    return
  }

  const targetLayer = props.layers[overIndex]
  if (targetLayer.locked) {
    handleDragEnd()
    return
  }

  const newLayers = [...props.layers]
  const [draggedLayer] = newLayers.splice(draggedIndex, 1)
  const baseIndex = newLayers.findIndex(layer => layer.id === overId)
  const insertIndex = dragPosition.value === 'before' ? baseIndex : baseIndex + 1
  newLayers.splice(insertIndex, 0, draggedLayer)
  emit('reorder', newLayers)
  handleDragEnd()
}

function handleDragEnd() {
  draggingId.value = null
  dragOverId.value = null
  dragPosition.value = null
}
</script>

<style scoped>
.layer-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  background: #ffffff;
  box-shadow: 0 6px 16px rgba(15, 23, 42, 0.08);
  border: 1px solid #ebeef5;
  overflow: hidden;
}

.panel-header {
  padding: 12px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #ebeef5;
}

.panel-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.panel-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 8px;
  color: #909399;
  font-size: 13px;
}

.layer-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.layer-item {
  position: relative;
  border-radius: 8px;
  border: 1px solid #ebeef5;
  background: #ffffff;
  padding: 8px 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.layer-item:hover {
  border-color: #c6e2ff;
  box-shadow: 0 4px 10px rgba(64, 158, 255, 0.12);
}

.layer-item.active {
  border-color: #409eff;
  background: #ecf5ff;
}

.layer-item.locked {
  opacity: 0.86;
}

.layer-item.drop-before::before,
.layer-item.drop-after::after {
  content: '';
  position: absolute;
  left: 6px;
  right: 6px;
  height: 2px;
  background: #67c23a;
}

.layer-item.drop-before::before {
  top: -5px;
}

.layer-item.drop-after::after {
  bottom: -5px;
}

.layer-main {
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.type-icon {
  color: #606266;
}

.layer-name {
  max-width: 118px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
  color: #303133;
}

.layer-level {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 22px;
  height: 22px;
  border-radius: 11px;
  background: #f4f4f5;
  color: #606266;
  font-size: 12px;
}

.layer-actions {
  display: inline-flex;
  align-items: center;
  gap: 2px;
}
</style>
