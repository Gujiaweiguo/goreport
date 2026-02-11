<template>
  <div class="dataset-preview">
    <div v-if="loading" class="loading-container">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
      <span>加载中...</span>
    </div>
    <div v-else-if="!datasetId" class="empty-container">
      <el-empty description="请选择要预览的数据集" />
    </div>
    <div v-else class="preview-content">
      <div class="preview-header">
        <h3>数据集预览</h3>
        <el-button @click="emit('close')">关闭</el-button>
      </div>
      <div class="chart-container">
        <div ref="chartRef" class="chart"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as echarts from 'echarts'
import { Loading } from '@element-plus/icons-vue'
import type { DatasetSchema } from '@/api/dataset'
import { buildPreviewChartModel } from './previewMapping'

interface Props {
  datasetId: string
  data: any[]
  schema?: DatasetSchema
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

const chartRef = ref<HTMLElement | null>(null)
const chartInstance = ref<echarts.ECharts | null>(null)
const loading = ref(false)
const previewModel = computed(() => buildPreviewChartModel(props.schema, props.data || []))

onMounted(() => {
  if (props.datasetId && props.data?.length > 0) {
    initChart()
    window.addEventListener('resize', handleResize)
  }
})

onBeforeUnmount(() => {
  if (chartInstance.value) {
    chartInstance.value.dispose()
    chartInstance.value = null
  }
  window.removeEventListener('resize', handleResize)
})

watch(() => props.data, (newData) => {
  if (newData?.length > 0 && chartInstance.value) {
    updateChart()
  }
}, { deep: true })

watch(previewModel, () => {
  if (chartInstance.value) {
    updateChart()
  }
}, { deep: true })

function initChart() {
  if (chartRef.value) {
    chartInstance.value = echarts.init(chartRef.value)
  }
  updateChart()
}

function updateChart() {
  if (!chartInstance.value) return
  const model = previewModel.value
  if (!model.categories.length || !model.values.length) return

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    title: model.fallbackMessage
      ? {
          text: '预览回退模式',
          subtext: model.fallbackMessage,
          textStyle: {
            fontSize: 14
          },
          subtextStyle: {
            fontSize: 12
          }
        }
      : undefined,
    xAxis: {
      type: 'category',
      name: model.categoryLabel,
      data: model.categories
    },
    yAxis: {
      type: 'value',
      name: model.valueLabel
    },
    series: [
      {
        name: model.valueLabel,
        type: 'bar',
        data: model.values,
        itemStyle: {
          color: '#409EFF'
        }
      }
    ]
  }

  chartInstance.value.setOption(option)
  chartInstance.value.resize()
}

function handleResize() {
  if (chartInstance.value) {
    chartInstance.value.resize()
  }
}
</script>

<style scoped>
.dataset-preview {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
  border-radius: 8px;
}

.loading-container {
  flex: 1;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.empty-container {
  flex: 1;
  align-items: center;
  justify-content: center;
  height: 100%;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
}

.preview-header h3 {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.chart-container {
  flex: 1;
  min-height: 300px;
}

.chart {
  width: 100%;
  height: 100%;
}
</style>
