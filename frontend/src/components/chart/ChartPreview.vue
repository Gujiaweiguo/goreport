<template>
  <div ref="previewRoot" class="chart-preview" :class="{ dark: chartConfig.theme === 'dark' }">
    <div class="preview-header">
      <span>实时预览</span>
      <el-tag size="small" effect="plain" type="success">Auto Sync</el-tag>
    </div>

    <div class="preview-body" :style="previewStyle">
      <div v-if="loading" class="loading-layer">
        <el-skeleton animated :rows="5" />
      </div>

      <div v-else-if="!hasChartData" class="empty-layer">
        <EmptyState
          :icon="DataAnalysis"
          :icon-size="36"
          text="暂无可预览数据"
          description="请先在数据配置中选择字段或刷新图表数据"
        />
      </div>

      <component
        :is="rendererComponent || EChartsRenderer"
        v-else
        ref="rendererRef"
        class="renderer"
        :option="currentOption"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { Component, CSSProperties } from 'vue'
import type { EChartsCoreOption } from 'echarts'
import { DataAnalysis } from '@element-plus/icons-vue'
import EmptyState from '@/components/common/EmptyState.vue'
import EChartsRenderer from '@/components/chart/EChartsRenderer.vue'

interface ChartSeriesItem {
  type: string
  name?: string
  data?: any[]
  [key: string]: any
}

export interface ChartDataPayload {
  categories: string[]
  series: ChartSeriesItem[]
  extra?: {
    yAxis?: any
  }
}

export interface ChartPreviewConfig {
  title: string
  width: number
  height: number
  margin: number
  theme: 'default' | 'dark'
  color: string
  fontFamily: string
  fontSize: number
  showLegend: boolean
  hoverable: boolean
  clickable: boolean
  zoomable: boolean
  mainTitle: string
  subTitle: string
  titlePosition: 'left' | 'center' | 'right'
}

const props = defineProps<{
  chartConfig: ChartPreviewConfig
  chartData: ChartDataPayload
  loading?: boolean
  rendererComponent?: Component
}>()

const rendererRef = ref<InstanceType<typeof EChartsRenderer> | null>(null)
const previewRoot = ref<HTMLElement | null>(null)
const currentOption = ref<EChartsCoreOption>({})
let resizeObserver: ResizeObserver | null = null

const hasChartData = computed(() => props.chartData.series.length > 0)

const previewStyle = computed<CSSProperties>(() => ({
  minHeight: `${props.chartConfig.height}px`,
  padding: `${props.chartConfig.margin}px`
}))

function buildBaseOption(config: ChartPreviewConfig, chartData: ChartDataPayload): EChartsCoreOption {
  const tooltipTrigger = chartData.series.some(item => item.type === 'pie' || item.type === 'gauge') ? 'item' : 'axis'

  const option: EChartsCoreOption = {
    animation: false,
    title: {
      text: config.mainTitle || config.title,
      subtext: config.subTitle,
      left: config.titlePosition,
      textStyle: {
        color: config.theme === 'dark' ? '#f3f4f6' : '#1f2937',
        fontSize: config.fontSize + 2,
        fontFamily: config.fontFamily
      },
      subtextStyle: {
        color: config.theme === 'dark' ? '#cbd5e1' : '#6b7280',
        fontFamily: config.fontFamily
      }
    },
    legend: {
      show: config.showLegend,
      bottom: 0,
      textStyle: {
        color: config.theme === 'dark' ? '#e5e7eb' : '#334155',
        fontFamily: config.fontFamily,
        fontSize: config.fontSize - 1
      }
    },
    tooltip: {
      show: config.hoverable,
      trigger: tooltipTrigger
    },
    backgroundColor: config.theme === 'dark' ? '#0f172a' : '#ffffff',
    textStyle: {
      color: config.theme === 'dark' ? '#e2e8f0' : '#334155',
      fontFamily: config.fontFamily,
      fontSize: config.fontSize
    },
    color: [config.color],
    grid: {
      left: config.margin,
      right: config.margin,
      top: config.margin + 42,
      bottom: config.margin + 28,
      containLabel: true
    }
  }

  const firstSeriesType = chartData.series[0]?.type || ''
  const nonAxisChart = firstSeriesType === 'pie' || firstSeriesType === 'gauge' || firstSeriesType === 'graph' || firstSeriesType === 'sankey' || firstSeriesType === 'tree' || firstSeriesType === 'treemap'

  if (!nonAxisChart) {
    option.xAxis = {
      type: 'category',
      data: chartData.categories,
      axisLabel: {
        color: config.theme === 'dark' ? '#cbd5e1' : '#475569'
      }
    }
    option.yAxis = chartData.extra?.yAxis || {
      type: 'value',
      axisLabel: {
        color: config.theme === 'dark' ? '#cbd5e1' : '#475569'
      },
      splitLine: {
        lineStyle: {
          color: config.theme === 'dark' ? 'rgba(148,163,184,0.2)' : 'rgba(148,163,184,0.25)'
        }
      }
    }
  }

  if (config.zoomable && !nonAxisChart) {
    option.dataZoom = [
      { type: 'inside' },
      { type: 'slider', bottom: 22 }
    ]
  }

  return option
}

function updateChart(config: ChartPreviewConfig) {
  const option = {
    ...buildBaseOption(config, props.chartData),
    series: props.chartData.series
  }
  currentOption.value = option
  rendererRef.value?.setOption(option)
}

function updateData(data: ChartDataPayload) {
  const base = buildBaseOption(props.chartConfig, data)
  const option = {
    ...base,
    series: data.series
  }
  currentOption.value = option
  rendererRef.value?.setOption({
    ...base,
    series: data.series
  })
}

watch(
  () => props.chartConfig,
  newConfig => {
    updateChart(newConfig)
  },
  { deep: true, immediate: true }
)

watch(
  () => props.chartData,
  newData => {
    updateData(newData)
  },
  { deep: true, immediate: true }
)

onMounted(() => {
  if (previewRoot.value) {
    resizeObserver = new ResizeObserver(() => {
      rendererRef.value?.resize()
    })
    resizeObserver.observe(previewRoot.value)
  }
})

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
  resizeObserver = null
})
</script>

<style scoped>
.chart-preview {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  border: 1px solid #dfe6f2;
  background: #f8fafc;
  overflow: hidden;
}

.chart-preview.dark {
  background: #0b1220;
  border-color: #1f2937;
}

.preview-header {
  height: 44px;
  padding: 0 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e8edf5;
  font-size: 13px;
  color: #1e293b;
  font-weight: 600;
}

.chart-preview.dark .preview-header {
  border-color: #1f2937;
  color: #f1f5f9;
}

.preview-body {
  flex: 1;
  position: relative;
  background: #ffffff;
}

.chart-preview.dark .preview-body {
  background: #0f172a;
}

.renderer,
.loading-layer,
.empty-layer {
  width: 100%;
  height: 100%;
}

.loading-layer {
  padding: 14px;
  box-sizing: border-box;
}

.empty-layer {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
