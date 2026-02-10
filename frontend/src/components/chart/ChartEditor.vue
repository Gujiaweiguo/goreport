<template>
  <div class="chart-editor">
    <div class="editor-header">
      <span>图表属性</span>
      <el-tag size="small" effect="dark" type="info">{{ chartTypeLabel }}</el-tag>
    </div>

    <el-form ref="formRef" :model="formData" label-width="90px" class="editor-form">
      <el-divider content-position="left">基础配置</el-divider>
      
      <el-form-item label="图表标题" prop="title">
        <el-input v-model="formData.title" placeholder="输入图表标题" />
      </el-form-item>

      <el-form-item label="图表类型">
        <el-select v-model="formData.type" placeholder="请选择" @change="handleTypeChange">
          <el-option label="柱状图" value="bar" />
          <el-option label="折线图" value="line" />
          <el-option label="饼图" value="pie" />
        </el-select>
      </el-form-item>

      <el-divider content-position="left">数据绑定</el-divider>
      
      <el-form-item label="数据集">
        <el-select
          v-model="formData.datasetId"
          placeholder="请选择数据集"
          clearable
          @change="handleDatasetChange"
        >
          <el-option
            v-for="ds in datasets"
            :key="ds.id"
            :label="ds.name"
            :value="ds.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item v-if="formData.datasetId" label="维度字段">
        <el-select
          v-model="formData.dimension"
          placeholder="请选择维度"
          clearable
        :disabled="!formData.datasetId"
        >
          <el-option
            v-for="dim in dimensions"
            :key="dim.id"
            :label="dim.displayName || dim.name"
            :value="dim.name"
          />
        </el-select>
      </el-form-item>

      <el-form-item v-if="formData.datasetId" label="指标字段">
        <el-select
          v-model="formData.measure"
          placeholder="请选择指标"
          clearable
          :disabled="!formData.datasetId"
          @change="handleMeasureChange"
        >
          <el-option
            v-for="m in measures"
            :key="m.id"
            :label="m.displayName || m.name"
            :value="m.name"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="聚合方式">
        <el-select v-model="formData.aggregation" clearable>
          <el-option label="无" value="none" />
          <el-option label="求和 (SUM)" value="SUM" />
          <el-option label="平均值 (AVG)" value="AVG" />
          <el-option label="计数 (COUNT)" value="COUNT" />
          <el-option label="最大值 (MAX)" value="MAX" />
          <el-option label="最小值 (MIN)" value="MIN" />
        </el-select>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { datasetApi, type Dataset, type DatasetField } from '@/api/dataset'

interface ChartData {
  title: string
  type: 'bar' | 'line' | 'pie'
  datasetId: string
  dimension: string
  measure: string
  aggregation: 'SUM' | 'AVG' | 'COUNT' | 'MAX' | 'MIN' | 'none'
}

const props = defineProps<{
  modelValue: ChartData
}>()

const emit = defineEmits<{
  'update:modelValue': [value: ChartData]
}>()

const formRef = ref<FormInstance>()
const datasets = ref<Dataset[]>([])
const dimensions = ref<DatasetField[]>([])
const measures = ref<DatasetField[]>([])

const formData = reactive<ChartData>({
  title: '',
  type: 'bar',
  datasetId: '',
  dimension: '',
  measure: '',
  aggregation: 'none'
})

const chartTypeLabel = computed(() => {
  const labels: Record<string, string> = {
    bar: '柱状图',
    line: '折线图',
    pie: '饼图'
  }
  return labels[formData.type] || formData.type
})

function handleTypeChange() {
  emitUpdate()
}

function handleDatasetChange() {
  dimensions.value = []
  measures.value = []
  formData.dimension = ''
  formData.measure = ''
  
  if (formData.datasetId) {
    loadDatasetSchema()
  }
  
  emitUpdate()
}

function handleMeasureChange() {
  emitUpdate()
}

async function loadDatasets() {
  try {
    const response = await datasetApi.list()
    if (response.success) {
      datasets.value = response.result || []
    } else {
      ElMessage.error(response.message || '加载数据集失败')
    }
  } catch (error: any) {
    ElMessage.error('加载数据集失败')
  }
}

async function loadDatasetSchema() {
  if (!formData.datasetId) return
  
  try {
    const response = await datasetApi.getSchema(formData.datasetId)
    if (response.success) {
      dimensions.value = response.result.dimensions || []
      measures.value = response.result.measures || []
    } else {
      ElMessage.error(response.message || '加载数据集 Schema 失败')
    }
  } catch (error: any) {
    ElMessage.error('加载数据集 Schema 失败')
  }
}

function emitUpdate() {
  emit('update:modelValue', {
    title: formData.title,
    type: formData.type,
    datasetId: formData.datasetId,
    dimension: formData.dimension,
    measure: formData.measure,
    aggregation: formData.aggregation
  })
}

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal) {
      formData.title = newVal.title || ''
      formData.type = newVal.type || 'bar'
      formData.datasetId = newVal.datasetId || ''
      formData.dimension = newVal.dimension || ''
      formData.measure = newVal.measure || ''
      formData.aggregation = newVal.aggregation || 'none'
      
      if (newVal.datasetId) {
        loadDatasetSchema()
      }
    }
  },
  { immediate: true, deep: true }
)

onMounted(async () => {
  await loadDatasets()
})
</script>

<style scoped>
.chart-editor {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  background: var(--panel-bg, #f8fafc);
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 8px 20px rgba(15, 23, 42, 0.06);
  font-family: "IBM Plex Sans", "Noto Sans", sans-serif;
}

.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  color: #0f172a;
  letter-spacing: 0.4px;
}

.editor-form {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

:deep(.el-divider__text) {
  color: #1f2933;
  font-weight: 600;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
}

:deep(.el-select) {
  width: 100%;
}
</style>
