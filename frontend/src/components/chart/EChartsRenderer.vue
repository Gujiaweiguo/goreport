<template>
  <div ref="chartContainer" class="echarts-renderer"></div>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import * as echarts from 'echarts'
import type { ECharts, EChartsCoreOption } from 'echarts'

const props = withDefaults(
  defineProps<{
    option?: EChartsCoreOption
  }>(),
  {
    option: () => ({})
  }
)

const chartContainer = ref<HTMLElement | null>(null)
let chartInstance: ECharts | null = null
let resizeObserver: ResizeObserver | null = null

const chartPalette = ['#2d6cdf', '#22a699', '#f59e0b', '#ef4444', '#8b5cf6', '#06b6d4']

function ensureChartReady() {
  if (chartInstance || !chartContainer.value) {
    return
  }
  chartInstance = echarts.init(chartContainer.value)
}

function setOption(option: EChartsCoreOption, notMerge = false) {
  ensureChartReady()
  if (!chartInstance) {
    return
  }

  const mergedOption: EChartsCoreOption = {
    animation: false,
    color: chartPalette,
    ...option
  }
  chartInstance.setOption(mergedOption, { notMerge, lazyUpdate: true })
}

function resize() {
  chartInstance?.resize()
}

function getDom() {
  return chartInstance?.getDom() ?? null
}

onMounted(async () => {
  await nextTick()
  ensureChartReady()
  setOption(props.option)

  if (chartContainer.value) {
    resizeObserver = new ResizeObserver(() => {
      resize()
    })
    resizeObserver.observe(chartContainer.value)
  }
})

watch(
  () => props.option,
  option => {
    setOption(option)
  },
  { deep: true }
)

onBeforeUnmount(() => {
  resizeObserver?.disconnect()
  resizeObserver = null
  chartInstance?.dispose()
  chartInstance = null
})

defineExpose({
  setOption,
  resize,
  getDom
})
</script>

<style scoped>
.echarts-renderer {
  width: 100%;
  height: 100%;
  min-height: 260px;
}
</style>
