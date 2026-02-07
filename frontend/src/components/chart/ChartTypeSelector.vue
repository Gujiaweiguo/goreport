<template>
  <div class="chart-type-selector">
    <div class="selector-list">
      <el-collapse v-model="activeGroups" class="selector-collapse">
        <el-collapse-item
          v-for="group in chartGroups"
          :key="group.key"
          :name="group.key"
          class="selector-group"
        >
          <template #title>
            <div class="group-title">
              <span>{{ group.label }}</span>
              <el-tag size="small" effect="plain" type="info">{{ group.items.length }}</el-tag>
            </div>
          </template>

          <div class="type-grid">
            <el-card
              v-for="item in group.items"
              :key="item.key"
              class="type-card"
              :class="{ selected: currentTypeKey === item.key }"
              shadow="hover"
              @click="handleSelect(item)"
            >
              <div class="type-head">
                <el-icon :size="18" class="type-icon"><component :is="item.icon" /></el-icon>
                <div class="type-title">
                  <span>{{ item.name }}</span>
                  <el-tag size="small" effect="dark" type="primary">{{ item.category }}</el-tag>
                </div>
              </div>
              <p class="type-desc">{{ item.description }}</p>
            </el-card>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>

    <div class="selector-detail" v-if="selectedType">
      <div class="detail-header">
        <el-icon :size="22" class="detail-icon"><component :is="selectedType.icon" /></el-icon>
        <div>
          <h4>{{ selectedType.name }}</h4>
          <el-tag size="small" effect="plain" type="success">{{ selectedType.category }}</el-tag>
        </div>
      </div>
      <p class="detail-description">{{ selectedType.description }}</p>
      <div class="detail-note">
        <span>推荐场景</span>
        <p>{{ selectedType.scene }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import {
  Connection,
  DataAnalysis,
  DataLine,
  DataBoard,
  Grid,
  Histogram,
  PieChart,
  Share,
  TrendCharts
} from '@element-plus/icons-vue'

export interface ChartTypeItem {
  key: string
  category: string
  name: string
  description: string
  scene: string
  icon: any
}

interface ChartTypeGroup {
  key: string
  label: string
  items: ChartTypeItem[]
}

const props = withDefaults(
  defineProps<{
    modelValue?: string
  }>(),
  {
    modelValue: ''
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
  change: [value: ChartTypeItem]
}>()

const chartGroups: ChartTypeGroup[] = [
  {
    key: 'basic',
    label: '基础图表',
    items: [
      {
        key: 'bar',
        category: '基础图表',
        name: '柱状图',
        description: '用于横向对比不同维度的数据表现。\n适合展示销量、访问量和排行。',
        scene: '分类对比、指标排行、月度趋势对照',
        icon: Histogram
      },
      {
        key: 'line',
        category: '基础图表',
        name: '折线图',
        description: '强调连续时间维度上的变化趋势。\n可用于观察峰值和转折点。',
        scene: '趋势分析、实时监控、同比环比',
        icon: TrendCharts
      },
      {
        key: 'pie',
        category: '基础图表',
        name: '饼图',
        description: '突出整体与部分之间的占比关系。\n数据类别不宜过多以保证可读性。',
        scene: '占比结构、渠道份额、成本构成',
        icon: PieChart
      },
      {
        key: 'scatter',
        category: '基础图表',
        name: '散点图',
        description: '展示两个变量间的相关性与分布。\n便于发现异常点与聚类。',
        scene: '相关性分析、离群点检测、样本分布',
        icon: Grid
      }
    ]
  },
  {
    key: 'pie-extended',
    label: '饼图扩展',
    items: [
      {
        key: 'donut',
        category: '饼图',
        name: '环形图',
        description: '在饼图基础上保留中心留白区域。\n可放置关键 KPI 增强信息密度。',
        scene: '核心指标卡、占比+总量联合展示',
        icon: PieChart
      },
      {
        key: 'rose',
        category: '饼图',
        name: '南丁玫瑰图',
        description: '使用半径或面积表达维度差异。\n视觉冲击更强，适合展示层级差。',
        scene: '品牌偏好、地域贡献、层级比较',
        icon: DataAnalysis
      },
      {
        key: 'gauge',
        category: '饼图',
        name: '仪表盘',
        description: '适合单值指标的状态展示。\n可配阈值区间体现风险等级。',
        scene: '完成率、健康度、风险预警',
        icon: DataBoard
      }
    ]
  },
  {
    key: 'relation',
    label: '关系图',
    items: [
      {
        key: 'graph',
        category: '关系图',
        name: '关系图',
        description: '以节点和连线刻画实体关系网络。\n适合展示连接强度和传播链路。',
        scene: '组织关系、社交网络、调用链路',
        icon: Share
      },
      {
        key: 'sankey',
        category: '关系图',
        name: '桑基图',
        description: '展示流量在多个阶段间的流向。\n节点宽度体现体量变化。',
        scene: '转化漏斗、能耗流向、资金去向',
        icon: Connection
      },
      {
        key: 'tree',
        category: '关系图',
        name: '树图',
        description: '用于层级结构和父子关系表达。\n支持自上而下的结构浏览。',
        scene: '组织架构、目录结构、权限树',
        icon: DataLine
      }
    ]
  },
  {
    key: 'geo',
    label: '地理图',
    items: [
      {
        key: 'map',
        category: '地理图',
        name: '地图',
        description: '以地域维度呈现指标分布状态。\n当前编辑器提供基础地理预览。',
        scene: '区域分布、门店布局、地理覆盖',
        icon: DataLine
      },
      {
        key: 'heatmap',
        category: '地理图',
        name: '热力图',
        description: '通过颜色深浅表达密度差异。\n可快速定位热点与冷区。',
        scene: '访问热点、人流密度、风险热区',
        icon: DataAnalysis
      },
      {
        key: 'map3d',
        category: '地理图',
        name: '3D 地图',
        description: '强调空间层次和视觉纵深。\n编辑器使用轻量模拟配置进行预览。',
        scene: '大屏展示、城市态势、空间分层',
        icon: DataBoard
      }
    ]
  },
  {
    key: 'combo',
    label: '组合图',
    items: [
      {
        key: 'bar-line',
        category: '组合图',
        name: '柱线组合图',
        description: '柱图展示规模，折线展示趋势。\n兼顾比较与变化两个视角。',
        scene: '销量与增速、收入与成本、量价关系',
        icon: TrendCharts
      },
      {
        key: 'multi-y',
        category: '组合图',
        name: '多 Y 轴图',
        description: '在同一图中展示不同量纲指标。\n通过双轴提高信息承载能力。',
        scene: '金额与比率、数量与单价、温度与湿度',
        icon: DataBoard
      }
    ]
  }
]

const activeGroups = ref(chartGroups.map(group => group.key))
const currentTypeKey = ref('')

const allChartTypes = chartGroups.flatMap(group => group.items)

const selectedType = computed(() => {
  const matched = allChartTypes.find(item => item.key === currentTypeKey.value)
  return matched || allChartTypes[0]
})

function setCurrentType(typeKey: string) {
  const matched = allChartTypes.find(item => item.key === typeKey) || allChartTypes[0]
  currentTypeKey.value = matched.key
  emit('update:modelValue', matched.key)
  emit('change', matched)
}

function handleSelect(item: ChartTypeItem) {
  setCurrentType(item.key)
}

watch(
  () => props.modelValue,
  value => {
    if (!value) {
      if (!currentTypeKey.value) {
        setCurrentType(allChartTypes[0].key)
      }
      return
    }
    if (value !== currentTypeKey.value) {
      const exists = allChartTypes.some(item => item.key === value)
      setCurrentType(exists ? value : allChartTypes[0].key)
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.chart-type-selector {
  height: 100%;
  display: grid;
  grid-template-columns: minmax(0, 1fr) 220px;
  gap: 12px;
  padding: 12px;
  box-sizing: border-box;
  background: linear-gradient(135deg, #f3f6fb 0%, #eef2f9 100%);
}

.selector-list {
  border-radius: 12px;
  border: 1px solid #e3e8f2;
  background: rgba(255, 255, 255, 0.88);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.7), 0 8px 24px rgba(15, 23, 42, 0.06);
  overflow: auto;
}

.selector-collapse {
  border: none;
}

.group-title {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-right: 10px;
  font-weight: 600;
  color: #1f2937;
}

.type-grid {
  display: grid;
  gap: 10px;
  padding: 8px 0 14px;
}

.type-card {
  border-radius: 10px;
  border-color: #e5eaf3;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, border-color 0.2s ease;
}

.type-card:hover {
  transform: translateY(-1px);
  border-color: #9ab5e0;
  box-shadow: 0 8px 18px rgba(45, 77, 128, 0.12);
}

.type-card.selected {
  border-color: #2d6cdf;
  background: linear-gradient(180deg, rgba(45, 108, 223, 0.06), rgba(45, 108, 223, 0.02));
  box-shadow: 0 10px 22px rgba(45, 108, 223, 0.2);
}

.type-head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.type-icon {
  color: #2d6cdf;
}

.type-title {
  min-width: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  flex: 1;
}

.type-title span {
  font-size: 13px;
  font-weight: 600;
  color: #273142;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.type-desc {
  margin: 0;
  white-space: pre-line;
  color: #6b7280;
  font-size: 12px;
  line-height: 1.55;
  min-height: 38px;
}

.selector-detail {
  border-radius: 12px;
  border: 1px solid #dde5f0;
  background: linear-gradient(160deg, #ffffff 0%, #f6f9ff 100%);
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.08);
  padding: 14px;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.detail-icon {
  color: #2d6cdf;
}

.detail-header h4 {
  margin: 0 0 6px;
  font-size: 15px;
  color: #1f2937;
}

.detail-description {
  margin: 0;
  white-space: pre-line;
  font-size: 12px;
  color: #5f6b7a;
  line-height: 1.6;
}

.detail-note {
  margin-top: 14px;
  padding-top: 10px;
  border-top: 1px dashed #d8e1f0;
}

.detail-note span {
  display: inline-block;
  margin-bottom: 6px;
  color: #475569;
  font-size: 12px;
  font-weight: 600;
}

.detail-note p {
  margin: 0;
  color: #64748b;
  font-size: 12px;
  line-height: 1.5;
}

:deep(.el-collapse-item__header) {
  border-bottom: 1px solid #edf2fb;
  color: #334155;
  font-weight: 600;
}

:deep(.el-collapse-item__wrap) {
  border-bottom: none;
}

:deep(.el-card__body) {
  padding: 10px;
}
</style>
