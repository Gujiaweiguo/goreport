<template>
  <div class="report-designer">
    <el-card class="designer-toolbar" shadow="never">
      <div class="toolbar-left">
        <span class="toolbar-title">报表设计器</span>
        <el-input
          v-model="reportName"
          size="small"
          class="name-input"
          placeholder="输入报表名称"
        />
      </div>
      <div class="toolbar-actions">
        <el-button @click="handleNew">新建</el-button>
        <el-button plain type="danger" :disabled="!selectedCell" @click="handleDeleteCell">
          删除
        </el-button>
        <el-button :disabled="!selectedCell" @click="handleMergeCell">合并</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </div>
    </el-card>

    <div class="designer-body">
      <div class="canvas-panel">
        <div ref="canvasWrapperRef" class="canvas-shell">
          <canvas
            ref="canvasRef"
            class="designer-canvas"
            @click="handleCanvasClick"
            @dblclick="handleCanvasDblClick"
          ></canvas>
          <div v-if="editing.visible" class="cell-editor" :style="editing.style">
            <el-input
              ref="editorRef"
              v-model="editing.value"
              size="small"
              @blur="commitEdit"
              @keyup.enter="commitEdit"
            />
          </div>
          <div v-if="selectedCell" class="cell-hint">
            选中 {{ selectedCellLabel }}
          </div>
        </div>
      </div>

      <div class="panel-panel">
        <PropertyPanel v-if="selectedCell" :cell="selectedCell" @update="handleCellUpdate" />
        <div v-else class="panel-empty">
          选择一个单元格以编辑属性
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { InputInstance } from 'element-plus'
import { reportApi } from '@/api/report'
import PropertyPanel from '@/components/report/PropertyPanel.vue'

interface CellStyle {
  fontSize: number
  fontWeight: 'normal' | 'bold'
  fontStyle: 'normal' | 'italic'
  align: 'left' | 'center' | 'right'
  color: string
  background: string
  borderColor: string
}

interface CellBinding {
  datasourceId?: string
  tableName?: string
  fieldName?: string
}

interface DesignerCell {
  row: number
  col: number
  text: string
  colSpan?: number
  style: CellStyle
  binding: CellBinding
}

const reportName = ref('')
const saving = ref(false)
const currentReportId = ref('')
const canvasRef = ref<HTMLCanvasElement | null>(null)
const canvasWrapperRef = ref<HTMLDivElement | null>(null)
const editorRef = ref<InputInstance>()
const selectedCell = ref<DesignerCell | null>(null)
const dpr = window.devicePixelRatio || 1

const gridConfig = reactive({
  rows: 18,
  cols: 12,
  cellWidth: 120,
  cellHeight: 44
})

const cells = reactive(new Map<string, DesignerCell>())

const editing = reactive({
  visible: false,
  row: 0,
  col: 0,
  value: '',
  style: {
    left: '0px',
    top: '0px',
    width: '120px',
    height: '40px'
  }
})

const selectedCellLabel = computed(() => {
  if (!selectedCell.value) return ''
  const rowIndex = selectedCell.value.row + 1
  const colLabel = String.fromCharCode(65 + selectedCell.value.col)
  return `${colLabel}${rowIndex}`
})

function createDefaultCell(row: number, col: number): DesignerCell {
  return {
    row,
    col,
    text: '',
    style: {
      fontSize: 14,
      fontWeight: 'normal',
      fontStyle: 'normal',
      align: 'left',
      color: '#0f172a',
      background: '#ffffff',
      borderColor: '#cbd5e1'
    },
    binding: {}
  }
}

function getCellKey(row: number, col: number) {
  return `${row}:${col}`
}

function getOrCreateCell(row: number, col: number) {
  const key = getCellKey(row, col)
  const existing = cells.get(key)
  if (existing) return existing
  const cell = createDefaultCell(row, col)
  cells.set(key, cell)
  return cell
}

function resolveMergedCell(row: number, col: number) {
  for (const cell of cells.values()) {
    const span = cell.colSpan ?? 1
    if (cell.row === row && col >= cell.col && col < cell.col + span) {
      return cell
    }
  }
  return getOrCreateCell(row, col)
}

function updateCanvasSize() {
  if (!canvasRef.value || !canvasWrapperRef.value) return
  const { clientWidth, clientHeight } = canvasWrapperRef.value
  canvasRef.value.width = clientWidth * dpr
  canvasRef.value.height = clientHeight * dpr
  canvasRef.value.style.width = `${clientWidth}px`
  canvasRef.value.style.height = `${clientHeight}px`
  renderGrid()
}

function renderGrid() {
  const canvas = canvasRef.value
  if (!canvas) return
  const ctx = canvas.getContext('2d')
  if (!ctx) return

  const width = canvas.width / dpr
  const height = canvas.height / dpr

  ctx.save()
  ctx.scale(dpr, dpr)
  ctx.clearRect(0, 0, width, height)

  ctx.fillStyle = '#f8fafc'
  ctx.fillRect(0, 0, width, height)

  ctx.lineWidth = 1
  ctx.strokeStyle = '#d0d5dd'
  for (let i = 0; i <= gridConfig.cols; i += 1) {
    const x = i * gridConfig.cellWidth
    ctx.beginPath()
    ctx.moveTo(x, 0)
    ctx.lineTo(x, gridConfig.rows * gridConfig.cellHeight)
    ctx.stroke()
  }
  for (let j = 0; j <= gridConfig.rows; j += 1) {
    const y = j * gridConfig.cellHeight
    ctx.beginPath()
    ctx.moveTo(0, y)
    ctx.lineTo(gridConfig.cols * gridConfig.cellWidth, y)
    ctx.stroke()
  }

  const coveredCells = new Set<string>()
  for (const cell of cells.values()) {
    const span = cell.colSpan ?? 1
    if (span > 1) {
      for (let offset = 1; offset < span; offset += 1) {
        coveredCells.add(getCellKey(cell.row, cell.col + offset))
      }
    }
  }

  for (const cell of cells.values()) {
    if (coveredCells.has(getCellKey(cell.row, cell.col))) continue
    const span = cell.colSpan ?? 1
    const x = cell.col * gridConfig.cellWidth
    const y = cell.row * gridConfig.cellHeight
    const width = gridConfig.cellWidth * span
    const height = gridConfig.cellHeight

    if (cell.style.background) {
      ctx.fillStyle = cell.style.background
      ctx.fillRect(x + 1, y + 1, width - 2, height - 2)
    }

    ctx.strokeStyle = cell.style.borderColor || '#cbd5e1'
    ctx.strokeRect(x, y, width, height)

    if (cell.text) {
      ctx.save()
      ctx.font = `${cell.style.fontStyle} ${cell.style.fontWeight} ${cell.style.fontSize}px "IBM Plex Sans", "Noto Sans", sans-serif`
      ctx.fillStyle = cell.style.color
      ctx.textBaseline = 'middle'
      const padding = 8
      let textX = x + padding
      if (cell.style.align === 'center') {
        ctx.textAlign = 'center'
        textX = x + width / 2
      } else if (cell.style.align === 'right') {
        ctx.textAlign = 'right'
        textX = x + width - padding
      } else {
        ctx.textAlign = 'left'
      }
      ctx.beginPath()
      ctx.rect(x + 2, y + 2, width - 4, height - 4)
      ctx.clip()
      ctx.fillText(cell.text, textX, y + height / 2)
      ctx.restore()
    }
  }

  if (selectedCell.value) {
    const span = selectedCell.value.colSpan ?? 1
    const x = selectedCell.value.col * gridConfig.cellWidth
    const y = selectedCell.value.row * gridConfig.cellHeight
    const width = gridConfig.cellWidth * span
    const height = gridConfig.cellHeight
    ctx.strokeStyle = '#2563eb'
    ctx.lineWidth = 2
    ctx.strokeRect(x + 1, y + 1, width - 2, height - 2)
  }

  ctx.restore()
}

function setSelectedCell(cell: DesignerCell) {
  selectedCell.value = cell
  renderGrid()
}

function getCellFromEvent(event: MouseEvent) {
  if (!canvasRef.value) return null
  const rect = canvasRef.value.getBoundingClientRect()
  const x = event.clientX - rect.left
  const y = event.clientY - rect.top
  const col = Math.floor(x / gridConfig.cellWidth)
  const row = Math.floor(y / gridConfig.cellHeight)
  if (row < 0 || col < 0 || row >= gridConfig.rows || col >= gridConfig.cols) {
    return null
  }
  return resolveMergedCell(row, col)
}

function handleCanvasClick(event: MouseEvent) {
  const cell = getCellFromEvent(event)
  if (!cell) return
  if (selectedCell.value && selectedCell.value.row === cell.row && selectedCell.value.col === cell.col) {
    startEditing(cell)
    return
  }
  editing.visible = false
  setSelectedCell(cell)
}

function handleCanvasDblClick(event: MouseEvent) {
  const cell = getCellFromEvent(event)
  if (!cell) return
  startEditing(cell)
}

function startEditing(cell: DesignerCell) {
  if (!canvasRef.value) return
  setSelectedCell(cell)
  const x = cell.col * gridConfig.cellWidth
  const y = cell.row * gridConfig.cellHeight
  editing.visible = true
  editing.value = cell.text
  editing.row = cell.row
  editing.col = cell.col
  editing.style = {
    left: `${x + 4}px`,
    top: `${y + 6}px`,
    width: `${gridConfig.cellWidth - 8}px`,
    height: `${gridConfig.cellHeight - 12}px`
  }
  nextTick(() => {
    editorRef.value?.focus()
  })
}

function commitEdit() {
  if (!editing.visible) return
  const cell = getOrCreateCell(editing.row, editing.col)
  cell.text = editing.value
  editing.visible = false
  renderGrid()
}

function handleCellUpdate(cell: DesignerCell) {
  const key = getCellKey(cell.row, cell.col)
  const updatedCell = {
    ...cell,
    style: {
      ...cell.style,
      borderColor: cell.style.borderColor || '#cbd5e1'
    }
  }
  cells.set(key, updatedCell)
  setSelectedCell(updatedCell)
  renderGrid()
}

function handleNew() {
  cells.clear()
  selectedCell.value = null
  reportName.value = ''
  currentReportId.value = ''
  renderGrid()
}

function handleDeleteCell() {
  if (!selectedCell.value) return
  const key = getCellKey(selectedCell.value.row, selectedCell.value.col)
  cells.delete(key)
  selectedCell.value = null
  renderGrid()
  ElMessage.success('单元格已清空')
}

function handleMergeCell() {
  if (!selectedCell.value) return
  const cell = selectedCell.value
  if (cell.col >= gridConfig.cols - 1) {
    ElMessage.warning('已到达行尾，无法合并')
    return
  }
  if (cell.colSpan && cell.colSpan > 1) {
    cell.colSpan = 1
    renderGrid()
    ElMessage.success('已取消合并')
    return
  }
  cell.colSpan = 2
  const coveredKey = getCellKey(cell.row, cell.col + 1)
  cells.delete(coveredKey)
  renderGrid()
  ElMessage.success('已合并到右侧单元格')
}

function serializeConfig() {
  return {
    grid: {
      rows: gridConfig.rows,
      cols: gridConfig.cols,
      cellWidth: gridConfig.cellWidth,
      cellHeight: gridConfig.cellHeight
    },
    cells: Array.from(cells.values())
  }
}

async function handleSave() {
  saving.value = true
  try {
    const payload = {
      name: reportName.value || '未命名报表',
      config: serializeConfig()
    }
    const response = currentReportId.value
      ? await reportApi.update({
          id: currentReportId.value,
          ...payload
        })
      : await reportApi.create(payload)

    if (response.data.success) {
      const savedId = response.data.result?.id || currentReportId.value
      if (savedId) {
        currentReportId.value = savedId
        localStorage.setItem('lastReportId', savedId)
      }
      ElMessage.success(savedId ? `报表已保存（ID: ${savedId}）` : '报表已保存')
    } else {
      ElMessage.error(response.data.message || '保存失败')
    }
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.message || '保存失败')
  } finally {
    saving.value = false
  }
}

let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  updateCanvasSize()
  resizeObserver = new ResizeObserver(() => {
    updateCanvasSize()
  })
  if (canvasWrapperRef.value) {
    resizeObserver.observe(canvasWrapperRef.value)
  }
})

onBeforeUnmount(() => {
  if (resizeObserver && canvasWrapperRef.value) {
    resizeObserver.unobserve(canvasWrapperRef.value)
  }
})
</script>

<style scoped>
.report-designer {
  min-height: 100vh;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  background: radial-gradient(circle at top left, #eef2f6, #f8fafc 35%, #ffffff);
  font-family: "IBM Plex Sans", "Noto Sans", sans-serif;
}

.designer-toolbar {
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  background: linear-gradient(145deg, #0f172a, #1f2937);
  color: #f8fafc;
}

.designer-toolbar :deep(.el-card__body) {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 16px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-title {
  font-size: 16px;
  font-weight: 600;
  letter-spacing: 0.6px;
}

.name-input {
  width: 220px;
}

.name-input :deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.18);
  box-shadow: none;
  color: #f8fafc;
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.designer-body {
  flex: 1;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 320px;
  gap: 16px;
}

.canvas-panel {
  background: #ffffff;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 12px 24px rgba(15, 23, 42, 0.08);
  padding: 12px;
  display: flex;
  flex-direction: column;
}

.canvas-shell {
  position: relative;
  flex: 1;
  min-height: 520px;
  border-radius: 10px;
  background: linear-gradient(145deg, #f8fafc, #f1f5f9);
  border: 1px solid #e2e8f0;
  overflow: hidden;
}

.designer-canvas {
  display: block;
  width: 100%;
  height: 100%;
  cursor: cell;
}

.cell-editor {
  position: absolute;
  z-index: 2;
}

.cell-editor :deep(.el-input__wrapper) {
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
}

.cell-hint {
  position: absolute;
  right: 12px;
  bottom: 12px;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.8);
  color: #f8fafc;
  font-size: 12px;
  letter-spacing: 0.4px;
}

.panel-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.panel-empty {
  padding: 24px;
  border-radius: 12px;
  border: 1px dashed #cbd5e1;
  background: #f8fafc;
  color: #64748b;
  text-align: center;
  font-size: 13px;
}

@media (max-width: 1200px) {
  .designer-body {
    grid-template-columns: 1fr;
  }

  .panel-panel {
    order: -1;
  }
}
</style>
