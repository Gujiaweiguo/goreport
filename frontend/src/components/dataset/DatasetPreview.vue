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
import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'

interface Props {
  datasetId: string
  data: any[]
}

const props = defineProps<Props>()
const emit = defineEmits<{
  close: []
}>()

const chartRef = ref<HTMLElement | null>(null)
const chartInstance = ref<echarts.ECharts | null>(null)
const loading = ref(false)

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
    updateChart(newData)
  }
}, { deep: true })

function initChart() {
  if (chartRef.value) {
    chartInstance.value = echarts.init(chartRef.value)
  }
  updateChart(props.data)
}

function updateChart(data: any[]) {
  if (!chartInstance.value || !data || data.length === 0) return

  const regions = data.map((item: any) => item.region || '')
  const sales = data.map((item: any) => item.sales || 0)

  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: regions
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value: number) => (value / 1000).toFixed(2) + '万'
      }
    },
    series: [
      {
        name: '销售额',
        type: 'bar',
        data: sales,
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
