<template>
  <div class="chart-panel">
    <div class="panel-header">
      <span>图表库</span>
      <el-button type="primary" size="small" @click="handleCreateChart">
        <el-icon><Plus /></el-icon>
        新建图表
      </el-button>
    </div>

    <div class="chart-list">
      <div
        v-for="chart in charts"
        :key="chart.id"
        class="chart-item"
        draggable="true"
        @dragstart="handleDragStart(chart, $event)"
      >
        <div class="chart-icon">
          <el-icon><TrendCharts /></el-icon>
        </div>
        <div class="chart-info">
          <div class="chart-name">{{ chart.name }}</div>
          <div class="chart-type">{{ chartTypeLabel(chart.type) }}</div>
        </div>
        <el-button
          type="primary"
          size="small"
          @click.stop="handleEditChart(chart)"
        >
          编辑
        </el-button>
      </div>
    </div>

    <el-dialog
      v-model="chartEditor.visible"
      title="新建/编辑图表"
      width="600px"
      @close="chartEditor.visible = false"
    >
      <ChartEditor v-model="chartEditor.data" @update:model-value="handleChartUpdate" />
      <template #footer>
        <el-button @click="chartEditor.visible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveChart">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, TrendCharts } from '@element-plus/icons-vue'
import { chartApi, type Chart, type CreateChartRequest, type UpdateChartRequest } from '@/api/chart'
import ChartEditor from '@/components/chart/ChartEditor.vue'

const emit = defineEmits<{
  'chart-selected': [chart: Chart]
}>()

const charts = ref<Chart[]>([])
const chartEditor = reactive({
  visible: false,
  data: null as Chart | null,
  mode: 'create' as 'create' | 'update'
})

const chartTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    bar: '柱状图',
    line: '折线图',
    pie: '饼图'
  }
  return labels[type] || type
}

async function loadCharts() {
  try {
    const response = await chartApi.list()
    if (response.success) {
      charts.value = response.result || []
    }
  } catch (error: any) {
    ElMessage.error('加载图表列表失败')
  }
}

function handleCreateChart() {
  chartEditor.data = {
    id: '',
    tenantId: '',
    name: '',
    code: '',
    type: 'bar',
    config: {
      title: '新图表',
      series: []
    },
    status: 1,
    createdAt: '',
    updatedAt: ''
  }
  chartEditor.mode = 'create'
  chartEditor.visible = true
}

function handleEditChart(chart: Chart) {
  chartEditor.data = chart
  chartEditor.mode = 'update'
  chartEditor.visible = true
}

function handleChartUpdate(data: Chart) {
  chartEditor.data = data
}

async function handleSaveChart() {
  try {
    let response
    if (chartEditor.mode === 'create') {
      response = await chartApi.create({
        name: chartEditor.data!.name,
        code: chartEditor.data!.code,
        type: chartEditor.data!.type,
        config: chartEditor.data!.config
      })
    } else {
      response = await chartApi.update({
        id: chartEditor.data!.id,
        name: chartEditor.data!.name,
        code: chartEditor.data!.code,
        type: chartEditor.data!.type,
        config: chartEditor.data!.config
      })
    }

    if (response.success) {
      ElMessage.success(chartEditor.mode === 'create' ? '图表已创建' : '图表已更新')
      chartEditor.visible = false
      await loadCharts()
    }
  } catch (error: any) {
    ElMessage.error('保存图表失败')
  }
}

function handleDragStart(chart: Chart, event: DragEvent) {
  event.dataTransfer?.setData('application/json', JSON.stringify({
    type: 'chart',
    data: chart
  }))
  event.dataTransfer?.effectAllowed = 'copy'
}

onMounted(async () => {
  await loadCharts()
})
</script>

<style scoped>
.chart-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  background: var(--panel-bg, #f8fafc);
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
  font-family: "IBM Plex Sans", "Noto Sans", sans-serif;
  min-height: 400px;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  font-weight: 600;
  color: #0f172a;
  letter-spacing: 0.4px;
}

.chart-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
  overflow-y: auto;
}

.chart-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  cursor: move;
  transition: all 0.2s ease;
}

.chart-item:hover {
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.1);
}

.chart-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  color: #ffffff;
}

.chart-info {
  flex: 1;
}

.chart-name {
  font-weight: 600;
  color: #1f2937;
  font-size: 13px;
}

.chart-type {
  font-size: 11px;
  color: #6b7280;
}
</style>
