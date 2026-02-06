<template>
  <div class="chart-editor">
    <aside class="editor-left">
      <ChartTypeSelector v-model="selectedChartTypeKey" @change="handleChartTypeChange" />
    </aside>

    <main class="editor-main">
      <div class="editor-toolbar">
        <el-button-group>
          <el-button :icon="Check" @click="handleValidate">验证配置</el-button>
          <el-button :icon="Document" @click="handleApplyConfig">应用配置</el-button>
          <el-button :icon="Refresh" :loading="previewLoading" @click="handleRefreshData">刷新数据</el-button>
        </el-button-group>
      </div>

      <div class="editor-config">
        <DataConfigPanel v-model="dataConfig" @change="handleDataConfigChange" />
      </div>

      <div class="editor-preview">
        <ChartPreview
          :chart-config="chartConfig"
          :chart-data="chartData"
          :loading="previewLoading"
          :renderer-component="EChartsRenderer"
        />
      </div>
    </main>

    <aside class="editor-right">
      <ChartPropertyPanel v-model="chartConfig" @change="handlePropertyChange" />
    </aside>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Check, Document, Refresh } from '@element-plus/icons-vue'
import ChartTypeSelector from '@/components/chart/ChartTypeSelector.vue'
import EChartsRenderer from '@/components/chart/EChartsRenderer.vue'
import DataConfigPanel from '@/components/chart/DataConfigPanel.vue'
import ChartPropertyPanel from '@/components/chart/ChartPropertyPanel.vue'
import ChartPreview from '@/components/chart/ChartPreview.vue'
import type { ChartTypeItem } from '@/components/chart/ChartTypeSelector.vue'
import type { ChartDataConfig } from '@/components/chart/DataConfigPanel.vue'
import type { ChartPropertyConfig } from '@/components/chart/ChartPropertyPanel.vue'
import type { ChartDataPayload } from '@/components/chart/ChartPreview.vue'

const selectedChartTypeKey = ref('bar')
const selectedChartType = ref<ChartTypeItem | null>(null)
const previewLoading = ref(false)

const dataConfig = reactive<ChartDataConfig>({
  dataSourceId: '',
  tableName: '',
  fields: [],
  params: [],
  filter: ''
})

const chartConfig = reactive<ChartPropertyConfig>({
  title: '柱状图',
  width: 760,
  height: 420,
  margin: 16,
  theme: 'default',
  color: '#2d6cdf',
  fontFamily: "'Noto Serif SC', serif",
  fontSize: 14,
  showLegend: true,
  hoverable: true,
  clickable: true,
  zoomable: false,
  mainTitle: '柱状图',
  subTitle: '图表编辑器实时预览',
  titlePosition: 'left'
})

const chartData = reactive<ChartDataPayload>(createChartData(selectedChartTypeKey.value))

function createChartData(typeKey: string): ChartDataPayload {
  const categories = ['1月', '2月', '3月', '4月', '5月', '6月']

  switch (typeKey) {
    case 'line':
      return {
        categories,
        series: [{ type: 'line', name: '趋势', smooth: true, data: [120, 152, 201, 164, 190, 330] }]
      }
    case 'pie':
      return {
        categories: [],
        series: [
          {
            type: 'pie',
            radius: '62%',
            name: '占比',
            data: [
              { value: 40, name: '华东' },
              { value: 26, name: '华北' },
              { value: 18, name: '华南' },
              { value: 16, name: '西南' }
            ]
          }
        ]
      }
    case 'donut':
      return {
        categories: [],
        series: [
          {
            type: 'pie',
            radius: ['42%', '68%'],
            name: '渠道分布',
            data: [
              { value: 35, name: '线上' },
              { value: 28, name: '线下' },
              { value: 21, name: '经销商' },
              { value: 16, name: '合作伙伴' }
            ]
          }
        ]
      }
    case 'rose':
      return {
        categories: [],
        series: [
          {
            type: 'pie',
            roseType: 'radius',
            radius: ['25%', '70%'],
            data: [
              { value: 30, name: '品牌 A' },
              { value: 24, name: '品牌 B' },
              { value: 18, name: '品牌 C' },
              { value: 14, name: '品牌 D' },
              { value: 8, name: '品牌 E' }
            ]
          }
        ]
      }
    case 'gauge':
      return {
        categories: [],
        series: [
          {
            type: 'gauge',
            min: 0,
            max: 100,
            detail: { formatter: '{value}%' },
            data: [{ value: 72, name: '完成率' }]
          }
        ]
      }
    case 'scatter':
      return {
        categories: [],
        series: [
          {
            type: 'scatter',
            name: '样本点',
            data: [
              [8, 12],
              [12, 22],
              [16, 32],
              [20, 35],
              [26, 48],
              [31, 53],
              [36, 68]
            ]
          }
        ]
      }
    case 'graph':
      return {
        categories: [],
        series: [
          {
            type: 'graph',
            layout: 'force',
            roam: true,
            data: [
              { name: '订单系统', value: 28, symbolSize: 50 },
              { name: '支付系统', value: 20, symbolSize: 42 },
              { name: '库存系统', value: 16, symbolSize: 36 },
              { name: '消息系统', value: 14, symbolSize: 34 }
            ],
            links: [
              { source: '订单系统', target: '支付系统' },
              { source: '订单系统', target: '库存系统' },
              { source: '支付系统', target: '消息系统' }
            ]
          }
        ]
      }
    case 'sankey':
      return {
        categories: [],
        series: [
          {
            type: 'sankey',
            data: [{ name: '访问' }, { name: '注册' }, { name: '下单' }, { name: '支付' }],
            links: [
              { source: '访问', target: '注册', value: 1200 },
              { source: '注册', target: '下单', value: 780 },
              { source: '下单', target: '支付', value: 520 }
            ]
          }
        ]
      }
    case 'tree':
      return {
        categories: [],
        series: [
          {
            type: 'tree',
            data: [
              {
                name: '总部',
                children: [
                  { name: '华东大区', children: [{ name: '上海' }, { name: '杭州' }] },
                  { name: '华南大区', children: [{ name: '广州' }, { name: '深圳' }] }
                ]
              }
            ],
            top: '10%',
            left: '12%',
            bottom: '10%',
            right: '24%'
          }
        ]
      }
    case 'map':
      return {
        categories: [],
        series: [
          {
            type: 'treemap',
            data: [
              { name: '华东', value: 320 },
              { name: '华南', value: 260 },
              { name: '华北', value: 210 },
              { name: '西南', value: 180 }
            ]
          }
        ]
      }
    case 'heatmap':
      return {
        categories: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri'],
        series: [
          {
            type: 'heatmap',
            data: [
              [0, 0, 7],
              [0, 1, 5],
              [1, 0, 4],
              [1, 1, 8],
              [2, 2, 9],
              [3, 1, 6],
              [4, 0, 3]
            ]
          }
        ],
        extra: {
          yAxis: {
            type: 'category',
            data: ['华东', '华北', '华南']
          }
        }
      }
    case 'map3d':
      return {
        categories: [],
        series: [
          {
            type: 'bar',
            name: '地理强度',
            data: [210, 260, 180, 240, 300, 220],
            itemStyle: {
              borderRadius: [6, 6, 0, 0]
            }
          }
        ]
      }
    case 'bar-line':
      return {
        categories,
        series: [
          {
            type: 'bar',
            name: '销售额',
            data: [180, 220, 260, 210, 300, 340],
            itemStyle: { borderRadius: [6, 6, 0, 0] }
          },
          {
            type: 'line',
            name: '增长率',
            smooth: true,
            data: [4, 8, 11, 7, 13, 16]
          }
        ]
      }
    case 'multi-y':
      return {
        categories,
        series: [
          {
            type: 'bar',
            name: '订单量',
            data: [420, 460, 520, 490, 610, 700],
            itemStyle: { borderRadius: [4, 4, 0, 0] }
          },
          {
            type: 'line',
            name: '客单价',
            yAxisIndex: 1,
            data: [68, 72, 80, 75, 86, 91],
            smooth: true
          }
        ],
        extra: {
          yAxis: [
            { type: 'value', name: '订单量' },
            { type: 'value', name: '客单价', position: 'right' }
          ]
        }
      }
    case 'bar':
    default:
      return {
        categories,
        series: [
          {
            type: 'bar',
            name: '销售额',
            data: [160, 230, 290, 210, 310, 360],
            itemStyle: {
              borderRadius: [6, 6, 0, 0]
            }
          }
        ]
      }
  }
}

function patchChartData(next: ChartDataPayload) {
  chartData.categories = next.categories
  chartData.series = next.series
  chartData.extra = next.extra
}

function handleChartTypeChange(type: ChartTypeItem) {
  selectedChartType.value = type
  selectedChartTypeKey.value = type.key
  chartConfig.title = type.name
  chartConfig.mainTitle = type.name
  chartConfig.subTitle = `${type.category} · 实时预览`
  patchChartData(createChartData(type.key))
}

function handleDataConfigChange(config: ChartDataConfig) {
  if (config.fields.length > 0) {
    const seriesName = config.fields[0]
    chartData.series = chartData.series.map((item, index) => {
      if (index !== 0) {
        return item
      }
      return {
        ...item,
        name: seriesName
      }
    })
  }
}

function handlePropertyChange(config: ChartPropertyConfig) {
  chartConfig.title = config.title
  chartConfig.width = config.width
  chartConfig.height = config.height
  chartConfig.margin = config.margin
  chartConfig.theme = config.theme
  chartConfig.color = config.color
  chartConfig.fontFamily = config.fontFamily
  chartConfig.fontSize = config.fontSize
  chartConfig.showLegend = config.showLegend
  chartConfig.hoverable = config.hoverable
  chartConfig.clickable = config.clickable
  chartConfig.zoomable = config.zoomable
  chartConfig.mainTitle = config.mainTitle
  chartConfig.subTitle = config.subTitle
  chartConfig.titlePosition = config.titlePosition
}

function handleValidate() {
  if (!selectedChartTypeKey.value) {
    ElMessage.error('请先选择图表类型')
    return
  }
  if (!dataConfig.dataSourceId || !dataConfig.tableName || dataConfig.fields.length === 0) {
    ElMessage.warning('数据配置不完整：请至少选择数据源、数据表和字段')
    return
  }
  ElMessage.success('图表配置校验通过')
}

function handleApplyConfig() {
  ElMessage.success('当前配置已应用到预览')
}

async function handleRefreshData() {
  previewLoading.value = true
  await new Promise(resolve => setTimeout(resolve, 500))
  patchChartData(createChartData(selectedChartTypeKey.value))
  previewLoading.value = false
  ElMessage.success('图表数据已刷新')
}
</script>

<style scoped>
.chart-editor {
  height: 100%;
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr) 320px;
  gap: 12px;
  padding: 12px;
  background: #eef2f8;
  box-sizing: border-box;
  overflow: hidden;
}

.editor-left,
.editor-right,
.editor-main {
  min-height: 0;
}

.editor-main {
  display: grid;
  grid-template-rows: 46px 320px minmax(0, 1fr);
  gap: 12px;
}

.editor-toolbar {
  border-radius: 12px;
  border: 1px solid #dfe6f2;
  background: #ffffff;
  display: flex;
  align-items: center;
  padding: 0 12px;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
}

.editor-config,
.editor-preview,
.editor-left,
.editor-right {
  border-radius: 12px;
  overflow: hidden;
}

.editor-config,
.editor-preview {
  min-height: 0;
}
</style>
