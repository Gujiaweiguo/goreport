<template>
  <div class="component-library">
    <el-collapse v-model="activeGroups" class="library-collapse">
      <el-collapse-item
        v-for="group in componentGroups"
        :key="group.key"
        :name="group.key"
        class="library-group"
      >
        <template #title>
          <div class="group-title">
            <span>{{ group.label }}</span>
            <el-tag size="small" effect="plain" type="info">{{ group.items.length }}</el-tag>
          </div>
        </template>

        <div class="group-list">
          <div
            v-for="item in group.items"
            :key="item.key"
            class="component-item"
            :class="{ dragging: draggingKey === item.key }"
            draggable="true"
            @dragstart="handleDragStart($event, item)"
            @dragend="handleDragEnd"
          >
            <div class="component-main">
              <el-icon class="item-icon" :size="16"><component :is="item.icon" /></el-icon>
              <span class="item-name">{{ item.name }}</span>
            </div>
            <el-tag size="small" effect="plain">{{ item.shortType }}</el-tag>
          </div>
        </div>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Collection, Edit, Grid, Picture, PieChart } from '@element-plus/icons-vue'

interface LibraryComponentItem {
  key: string
  name: string
  shortType: string
  icon: any
}

interface LibraryComponentGroup {
  key: string
  label: string
  items: LibraryComponentItem[]
}

const emit = defineEmits<{
  dragstart: [componentType: string]
  dragend: []
}>()

const activeGroups = ref(['text', 'chart', 'table', 'image', 'decorative'])
const draggingKey = ref('')

const componentGroups: LibraryComponentGroup[] = [
  {
    key: 'text',
    label: 'Text（文本组件）',
    items: [
      { key: 'text-title', name: '标题文本', shortType: 'Text', icon: Edit },
      { key: 'text-basic', name: '普通文本', shortType: 'Text', icon: Edit },
      { key: 'text-rich', name: '富文本', shortType: 'Text', icon: Edit }
    ]
  },
  {
    key: 'chart',
    label: 'Chart（图表组件）',
    items: [
      { key: 'chart-bar', name: '柱状图', shortType: 'Chart', icon: PieChart },
      { key: 'chart-line', name: '折线图', shortType: 'Chart', icon: PieChart },
      { key: 'chart-pie', name: '饼图', shortType: 'Chart', icon: PieChart },
      { key: 'chart-scatter', name: '散点图', shortType: 'Chart', icon: PieChart }
    ]
  },
  {
    key: 'table',
    label: 'Table（表格组件）',
    items: [
      { key: 'table-basic', name: '基础表格', shortType: 'Table', icon: Grid },
      { key: 'table-page', name: '分页表格', shortType: 'Table', icon: Grid }
    ]
  },
  {
    key: 'image',
    label: 'Image（图片组件）',
    items: [
      { key: 'image-basic', name: '图片', shortType: 'Image', icon: Picture },
      { key: 'image-border', name: '边框图片', shortType: 'Image', icon: Picture },
      { key: 'image-background', name: '背景图', shortType: 'Image', icon: Picture }
    ]
  },
  {
    key: 'decorative',
    label: 'Decorative（装饰组件）',
    items: [
      { key: 'decorative-line', name: '分割线', shortType: 'Decorative', icon: Collection },
      { key: 'decorative-frame', name: '装饰框', shortType: 'Decorative', icon: Collection },
      { key: 'decorative-corner', name: '角标', shortType: 'Decorative', icon: Collection }
    ]
  }
]

function handleDragStart(event: DragEvent, item: LibraryComponentItem) {
  draggingKey.value = item.key
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData('text/plain', item.key)
  }
  emit('dragstart', item.key)
}

function handleDragEnd() {
  draggingKey.value = ''
  emit('dragend')
}
</script>

<style scoped>
.component-library {
  height: 100%;
  overflow: auto;
}

.library-collapse {
  border: none;
}

.group-title {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-right: 12px;
  font-weight: 600;
}

.group-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 4px 0 8px;
}

.component-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 42px;
  padding: 0 10px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #ffffff;
  cursor: grab;
  transition: all 0.2s ease;
}

.component-item:hover {
  border-color: #7cb7ff;
  background: #f5faff;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.14);
}

.component-item:active {
  cursor: grabbing;
}

.component-item.dragging {
  opacity: 0.5;
  box-shadow: 0 10px 22px rgba(15, 23, 42, 0.2);
}

.component-main {
  min-width: 0;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.item-icon {
  color: #409eff;
}

.item-name {
  font-size: 13px;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

:deep(.el-collapse-item__header) {
  font-weight: 600;
  color: #303133;
}

:deep(.el-collapse-item__wrap) {
  border-bottom: none;
}
</style>
