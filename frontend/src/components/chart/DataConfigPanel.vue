<template>
  <div class="data-config-panel">
    <div class="panel-header">
      <h3>数据配置</h3>
      <el-tag size="small" effect="plain" type="info">Data Binding</el-tag>
    </div>

    <el-form ref="formRef" :model="formState" :rules="rules" label-position="top" class="panel-form">
      <el-form-item label="绑定类型">
        <el-radio-group v-model="bindingType" @change="handleBindingTypeChange">
          <el-radio value="datasource">数据源</el-radio>
          <el-radio value="dataset">数据集</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-divider />

      <el-collapse v-model="activeGroups">
        <template v-if="bindingType === 'datasource'">
          <el-collapse-item title="数据源" name="source">
            <el-form-item label="数据源选择" prop="dataSourceId">
              <el-select
                v-model="formState.dataSourceId"
                placeholder="请选择数据源"
                filterable
                clearable
                :loading="loading.sources"
                @change="handleDataSourceChange"
              >
                <el-option v-for="item in dataSourceOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-collapse-item>

          <el-collapse-item title="数据表" name="table">
            <el-form-item label="数据表选择" prop="tableName">
              <el-select
                v-model="formState.tableName"
                placeholder="请选择数据表"
                filterable
                clearable
                :loading="loading.tables"
                :disabled="!formState.dataSourceId"
                @change="handleTableChange"
              >
                <el-option v-for="item in tableOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-collapse-item>

          <el-collapse-item title="字段" name="fields">
            <el-form-item label="字段选择（多选）" prop="fields">
              <el-select
                v-model="formState.fields"
                multiple
                collapse-tags
                collapse-tags-tooltip
                placeholder="请选择字段"
                :loading="loading.fields"
                :disabled="!formState.tableName"
                @change="emitChange"
              >
                <el-option v-for="item in fieldOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-collapse-item>
        </template>

        <el-collapse-item v-if="bindingType === 'dataset'" title="数据集配置" name="dataset">
          <el-form-item label="数据集" prop="datasetId">
            <el-select
              v-model="formState.datasetId"
              placeholder="请选择数据集"
              clearable
              filterable
              :loading="loading.datasets"
              @change="handleDatasetChange"
            >
              <el-option v-for="dataset in datasets" :key="dataset.id" :label="dataset.name" :value="dataset.id" />
            </el-select>
          </el-form-item>

          <el-form-item label="维度" prop="dimension">
            <el-select
              v-model="formState.dimension"
              placeholder="请选择维度"
              clearable
              filterable
              :disabled="!formState.datasetId"
              :loading="loading.datasetSchema"
              @change="handleDimensionChange"
            >
              <el-option
                v-for="dimension in dimensions"
                :key="dimension.id"
                :label="dimension.displayName || dimension.name"
                :value="dimension.name"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="指标" prop="measure">
            <el-select
              v-model="formState.measure"
              placeholder="请选择指标"
              clearable
              filterable
              :disabled="!formState.datasetId"
              :loading="loading.datasetSchema"
              @change="handleMeasureChange"
            >
              <el-option
                v-for="measure in measures"
                :key="measure.id"
                :label="measure.displayName || measure.name"
                :value="measure.name"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="聚合函数" prop="aggregation">
            <el-select
              v-model="formState.aggregation"
              placeholder="请选择聚合函数"
              :disabled="!formState.measure"
              @change="emitChange"
            >
              <el-option label="无" value="none" />
              <el-option label="求和 (SUM)" value="SUM" />
              <el-option label="平均值 (AVG)" value="AVG" />
              <el-option label="计数 (COUNT)" value="COUNT" />
              <el-option label="最大值 (MAX)" value="MAX" />
              <el-option label="最小值 (MIN)" value="MIN" />
            </el-select>
          </el-form-item>
        </el-collapse-item>

        <el-collapse-item title="参数" name="params">
          <div class="param-list">
            <div v-for="(param, index) in formState.params" :key="`param-${index}`" class="param-row">
              <el-input v-model="param.key" placeholder="参数名" @input="emitChange" />
              <el-input v-model="param.value" placeholder="参数值" @input="emitChange" />
              <el-button text type="danger" @click="removeParam(index)">删除</el-button>
            </div>
          </div>
          <el-button text type="primary" @click="addParam">+ 添加参数</el-button>
        </el-collapse-item>

        <el-collapse-item title="过滤" name="filter">
          <el-form-item label="数据过滤条件" prop="filter">
            <el-input
              v-model="formState.filter"
              type="textarea"
              :rows="3"
              placeholder="示例：amount > 1000 AND region = '华东'"
              @input="emitChange"
            />
          </el-form-item>
        </el-collapse-item>
      </el-collapse>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { datasourceApi } from '@/api/datasource'
import { datasetApi, type Dataset, type DatasetField } from '@/api/dataset'

type BindingType = 'datasource' | 'dataset'
type AggregationType = 'SUM' | 'AVG' | 'COUNT' | 'MAX' | 'MIN' | 'none'

interface OptionItem {
  label: string
  value: string
}

interface ParamItem {
  key: string
  value: string
}

export interface ChartDataConfig {
  dataSourceId: string
  tableName: string
  fields: string[]
  datasetId: string
  dimension: string
  measure: string
  aggregation: AggregationType
  params: ParamItem[]
  filter: string
}

const props = withDefaults(
  defineProps<{
    modelValue: ChartDataConfig
  }>(),
  {
    modelValue: () => ({
      dataSourceId: '',
      tableName: '',
      fields: [],
      datasetId: '',
      dimension: '',
      measure: '',
      aggregation: 'none',
      params: [],
      filter: ''
    })
  }
)

const emit = defineEmits<{
  'update:modelValue': [value: ChartDataConfig]
  change: [value: ChartDataConfig]
}>()

const formRef = ref<FormInstance>()
const activeGroups = ref(['source', 'table', 'fields', 'params', 'filter'])
const bindingType = ref<BindingType>('datasource')
const schemaDatasetId = ref('')

const loading = reactive({
  sources: false,
  tables: false,
  fields: false,
  datasets: false,
  datasetSchema: false
})

const dataSourceOptions = ref<OptionItem[]>([])
const tableOptions = ref<OptionItem[]>([])
const fieldOptions = ref<OptionItem[]>([])
const datasets = ref<Dataset[]>([])
const dimensions = ref<DatasetField[]>([])
const measures = ref<DatasetField[]>([])

const formState = reactive<ChartDataConfig>({
  dataSourceId: '',
  tableName: '',
  fields: [],
  datasetId: '',
  dimension: '',
  measure: '',
  aggregation: 'none',
  params: [],
  filter: ''
})

const rules = reactive<FormRules<ChartDataConfig>>({
  dataSourceId: [
    {
      validator: (_rule, value, callback) => {
        if (bindingType.value === 'datasource' && !value) {
          callback(new Error('请选择数据源'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ],
  tableName: [
    {
      validator: (_rule, value, callback) => {
        if (bindingType.value === 'datasource' && !value) {
          callback(new Error('请选择数据表'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ],
  fields: [
    {
      validator: (_rule, value, callback) => {
        if (bindingType.value === 'datasource' && (!Array.isArray(value) || value.length === 0)) {
          callback(new Error('请至少选择一个字段'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ],
  datasetId: [
    {
      validator: (_rule, value, callback) => {
        if (bindingType.value === 'dataset' && !value) {
          callback(new Error('请选择数据集'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ]
})

function setActiveGroupsByBindingType() {
  activeGroups.value = bindingType.value === 'dataset' ? ['dataset', 'params', 'filter'] : ['source', 'table', 'fields', 'params', 'filter']
}

function syncForm(value: ChartDataConfig) {
  formState.dataSourceId = value.dataSourceId || ''
  formState.tableName = value.tableName || ''
  formState.fields = [...(value.fields || [])]
  formState.datasetId = value.datasetId || ''
  formState.dimension = value.dimension || ''
  formState.measure = value.measure || ''
  formState.aggregation = value.aggregation || 'none'
  formState.params = [...(value.params || [])]
  formState.filter = value.filter || ''

  const hasDatasetBinding = !!(formState.datasetId || formState.dimension || formState.measure || formState.aggregation !== 'none')
  const hasDatasourceBinding = !!(formState.dataSourceId || formState.tableName || formState.fields.length)
  if (hasDatasetBinding) {
    bindingType.value = 'dataset'
  } else if (hasDatasourceBinding) {
    bindingType.value = 'datasource'
  }

  setActiveGroupsByBindingType()
}

function emitChange() {
  const nextValue: ChartDataConfig = {
    dataSourceId: formState.dataSourceId,
    tableName: formState.tableName,
    fields: [...formState.fields],
    datasetId: formState.datasetId,
    dimension: formState.dimension,
    measure: formState.measure,
    aggregation: formState.aggregation,
    params: formState.params.map(item => ({ ...item })),
    filter: formState.filter
  }
  emit('update:modelValue', nextValue)
  emit('change', nextValue)
  void formRef.value?.validate().catch(() => {
    return undefined
  })
}

function normalizeFieldOptions(result: any[]): OptionItem[] {
  return result
    .map(item => {
      if (typeof item === 'string') {
        return { label: item, value: item }
      }
      const name = item?.name || item?.fieldName || item?.columnName || ''
      if (!name) {
        return null
      }
      return { label: name, value: name }
    })
    .filter((item): item is OptionItem => item !== null)
}

async function loadDataSources() {
  loading.sources = true
  try {
    const response = await datasourceApi.list()
    if (response.data.success && Array.isArray(response.data.result)) {
      dataSourceOptions.value = response.data.result.map(item => ({
        label: item.name,
        value: item.id
      }))
      return
    }
  } catch {
    ElMessage.warning('数据源加载失败，已切换本地模拟数据')
  } finally {
    loading.sources = false
  }

  dataSourceOptions.value = [
    { label: 'SalesDB', value: 'sales-db' },
    { label: 'MarketingDB', value: 'marketing-db' },
    { label: 'DemoData', value: 'demo-data' }
  ]
}

async function handleDataSourceChange() {
  formState.tableName = ''
  formState.fields = []
  tableOptions.value = []
  fieldOptions.value = []

  if (!formState.dataSourceId) {
    emitChange()
    return
  }

  loading.tables = true
  try {
    const response = await datasourceApi.getTables(formState.dataSourceId)
    if (response.data.success && Array.isArray(response.data.result)) {
      tableOptions.value = response.data.result.map(name => ({ label: name, value: name }))
    } else {
      tableOptions.value = []
      ElMessage.error(response.data.message || '数据表加载失败')
    }
  } catch {
    tableOptions.value = [
      { label: 'sales_order', value: 'sales_order' },
      { label: 'sales_detail', value: 'sales_detail' },
      { label: 'sales_region', value: 'sales_region' }
    ]
    ElMessage.warning('数据表接口异常，使用模拟数据表')
  } finally {
    loading.tables = false
  }

  emitChange()
}

function handleBindingTypeChange() {
  if (bindingType.value === 'dataset') {
    formState.dataSourceId = ''
    formState.tableName = ''
    formState.fields = []
    tableOptions.value = []
    fieldOptions.value = []
    if (datasets.value.length === 0) {
      void loadDatasets()
    }
  } else {
    formState.datasetId = ''
    formState.dimension = ''
    formState.measure = ''
    formState.aggregation = 'none'
    dimensions.value = []
    measures.value = []
    schemaDatasetId.value = ''
  }

  setActiveGroupsByBindingType()
  emitChange()
}

async function loadDatasets() {
  loading.datasets = true
  try {
    const response = await datasetApi.list()
    if (response.data.success && Array.isArray(response.data.result)) {
      datasets.value = response.data.result
      return
    }
    datasets.value = []
    ElMessage.error(response.data.message || '数据集加载失败')
  } catch {
    datasets.value = []
    ElMessage.error('数据集加载失败')
  } finally {
    loading.datasets = false
  }
}

async function loadDatasetSchema(datasetId: string) {
  loading.datasetSchema = true
  try {
    const response = await datasetApi.getSchema(datasetId)
    if (response.data.success) {
      dimensions.value = response.data.result?.dimensions || []
      measures.value = response.data.result?.measures || []
      schemaDatasetId.value = datasetId
      return
    }
    dimensions.value = []
    measures.value = []
    ElMessage.error(response.data.message || '数据集字段加载失败')
  } catch {
    dimensions.value = []
    measures.value = []
    ElMessage.error('数据集字段加载失败')
  } finally {
    loading.datasetSchema = false
  }
}

async function handleDatasetChange() {
  formState.dimension = ''
  formState.measure = ''
  formState.aggregation = 'none'
  dimensions.value = []
  measures.value = []
  schemaDatasetId.value = ''

  if (!formState.datasetId) {
    emitChange()
    return
  }

  await loadDatasetSchema(formState.datasetId)
  emitChange()
}

function handleDimensionChange() {
  if (!formState.measure) {
    formState.aggregation = 'none'
  }
  emitChange()
}

function handleMeasureChange() {
  formState.aggregation = formState.measure ? 'SUM' : 'none'
  emitChange()
}

async function handleTableChange() {
  formState.fields = []
  fieldOptions.value = []

  if (!formState.dataSourceId || !formState.tableName) {
    emitChange()
    return
  }

  loading.fields = true
  try {
    const response = await datasourceApi.getFields(formState.dataSourceId, formState.tableName)
    if (response.data.success && Array.isArray(response.data.result)) {
      fieldOptions.value = normalizeFieldOptions(response.data.result)
    } else {
      fieldOptions.value = []
      ElMessage.error(response.data.message || '字段加载失败')
    }
  } catch {
    fieldOptions.value = [
      { label: 'month', value: 'month' },
      { label: 'region', value: 'region' },
      { label: 'amount', value: 'amount' },
      { label: 'count', value: 'count' }
    ]
    ElMessage.warning('字段接口异常，使用模拟字段')
  } finally {
    loading.fields = false
  }

  emitChange()
}

function addParam() {
  formState.params.push({ key: '', value: '' })
  emitChange()
}

function removeParam(index: number) {
  formState.params.splice(index, 1)
  emitChange()
}

watch(
  () => props.modelValue,
  async value => {
    syncForm(value)

    if (bindingType.value === 'dataset' && formState.datasetId) {
      if (schemaDatasetId.value !== formState.datasetId) {
        await loadDatasetSchema(formState.datasetId)
      }

      if (!dimensions.value.some(item => item.name === formState.dimension)) {
        formState.dimension = ''
      }
      if (!measures.value.some(item => item.name === formState.measure)) {
        formState.measure = ''
        formState.aggregation = 'none'
      }
    }
  },
  { immediate: true, deep: true }
)

onMounted(async () => {
  await Promise.all([loadDataSources(), loadDatasets()])
})
</script>

<style scoped>
.data-config-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 12px;
  border: 1px solid #dfe6f2;
  background: #ffffff;
  overflow: hidden;
}

.panel-header {
  height: 46px;
  padding: 0 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #e8edf5;
  background: linear-gradient(180deg, #f8faff 0%, #f2f6fc 100%);
}

.panel-header h3 {
  margin: 0;
  font-size: 14px;
  color: #1f2937;
}

.panel-form {
  flex: 1;
  overflow: auto;
  padding: 10px 12px;
}

.param-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 6px;
}

.param-row {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 8px;
}

:deep(.el-select),
:deep(.el-input),
:deep(.el-textarea) {
  width: 100%;
}

:deep(.el-form-item) {
  margin-bottom: 10px;
}

:deep(.el-collapse) {
  border-top: none;
  border-bottom: none;
}

:deep(.el-collapse-item__header) {
  color: #334155;
  font-weight: 600;
}
</style>
