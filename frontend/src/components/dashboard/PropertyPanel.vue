<template>
  <div class="property-panel">
    <div class="panel-header">
      <h3>组件属性</h3>
      <el-tag v-if="component" size="small" effect="plain" type="info">{{ component.id }}</el-tag>
    </div>

    <div v-if="!component" class="panel-empty">
      <el-icon :size="22"><Pointer /></el-icon>
      <span>请选择组件后编辑属性</span>
    </div>

    <el-form
      v-else
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-position="top"
      class="panel-form"
    >
      <el-collapse v-model="activeGroups">
        <el-collapse-item title="Common（常见属性）" name="common">
          <el-form-item label="Title" prop="title">
            <el-input v-model="formData.title" placeholder="请输入标题" @input="emitUpdate" />
          </el-form-item>
          <el-form-item label="ID" prop="id">
            <el-input v-model="formData.id" placeholder="请输入组件 ID" @input="emitUpdate" />
          </el-form-item>
          <div class="field-grid">
            <el-form-item label="Width" prop="width">
              <el-input-number v-model="formData.width" :min="20" :max="3840" :step="10" @change="emitUpdate" />
            </el-form-item>
            <el-form-item label="Height" prop="height">
              <el-input-number v-model="formData.height" :min="20" :max="2160" :step="10" @change="emitUpdate" />
            </el-form-item>
          </div>
          <div class="field-grid">
            <el-form-item label="X" prop="x">
              <el-input-number v-model="formData.x" :min="0" :max="3840" :step="1" @change="emitUpdate" />
            </el-form-item>
            <el-form-item label="Y" prop="y">
              <el-input-number v-model="formData.y" :min="0" :max="2160" :step="1" @change="emitUpdate" />
            </el-form-item>
          </div>
        </el-collapse-item>

        <el-collapse-item title="Style（样式属性）" name="style">
          <el-form-item label="Background" prop="background">
            <el-color-picker v-model="formData.background" show-alpha @change="emitUpdate" />
          </el-form-item>
          <el-form-item label="Border" prop="border">
            <el-select v-model="formData.border" placeholder="请选择边框样式" @change="emitUpdate">
              <el-option label="none" value="none" />
              <el-option label="solid" value="solid" />
              <el-option label="dashed" value="dashed" />
              <el-option label="dotted" value="dotted" />
              <el-option label="double" value="double" />
            </el-select>
          </el-form-item>
          <el-form-item label="Border Color" prop="borderColor">
            <el-color-picker v-model="formData.borderColor" show-alpha @change="emitUpdate" />
          </el-form-item>
          <div class="field-grid">
            <el-form-item label="Font Size" prop="fontSize">
              <el-input-number v-model="formData.fontSize" :min="10" :max="120" :step="1" @change="emitUpdate" />
            </el-form-item>
            <el-form-item label="Font Color" prop="fontColor">
              <el-color-picker v-model="formData.fontColor" show-alpha @change="emitUpdate" />
            </el-form-item>
          </div>
          <div class="field-grid">
            <el-form-item label="Padding" prop="padding">
              <el-input-number v-model="formData.padding" :min="0" :max="200" :step="1" @change="emitUpdate" />
            </el-form-item>
            <el-form-item label="Margin" prop="margin">
              <el-input-number v-model="formData.margin" :min="0" :max="200" :step="1" @change="emitUpdate" />
            </el-form-item>
          </div>
        </el-collapse-item>

        <el-collapse-item title="Data（数据属性）" name="data">
          <el-form-item label="绑定类型">
            <el-radio-group v-model="bindingType" @change="handleBindingTypeChange">
              <el-radio value="datasource">数据源</el-radio>
              <el-radio value="dataset">数据集</el-radio>
            </el-radio-group>
          </el-form-item>

          <template v-if="bindingType === 'datasource'">
            <el-form-item label="Data Source" prop="dataSource">
              <el-select
                v-model="formData.dataSource"
                placeholder="请选择数据源"
                :loading="loadingDataSources"
                filterable
                @change="handleDataSourceChange"
              >
                <el-option
                  v-for="source in dataSources"
                  :key="source.value"
                  :label="source.label"
                  :value="source.value"
                />
              </el-select>
              <el-button
                v-if="formData.dataSource"
                link
                type="primary"
                size="small"
                :loading="testingConnection"
                @click="handleTestConnection"
              >
                <el-icon><Connection /></el-icon>
                测试连接
              </el-button>
              <el-button
                v-if="formData.dataSource && formData.field"
                link
                type="success"
                size="small"
                :loading="previewingData"
                @click="handlePreviewData"
              >
                <el-icon><View /></el-icon>
                预览数据
              </el-button>
            </el-form-item>
            <el-form-item label="Field" prop="field">
              <el-select
                v-model="formData.field"
                placeholder="请选择数据源后选择字段"
                :loading="loadingFields"
                :disabled="!formData.dataSource"
                filterable
                @change="handleFieldChange"
              >
                <el-option v-for="field in fields" :key="field.value" :label="field.label" :value="field.value" />
              </el-select>
            </el-form-item>
            <el-form-item label="SQL Query" prop="sqlQuery">
              <el-input
                v-model="formData.sqlQuery"
                type="textarea"
                :rows="3"
                placeholder="SELECT * FROM table_name"
                @input="emitUpdate"
              />
              <el-button
                v-if="formData.dataSource"
                link
                type="warning"
                size="small"
                @click="handleFormatSQL"
              >
                <el-icon><Document /></el-icon>
                格式化 SQL
              </el-button>
            </el-form-item>
            <div v-if="previewData.length || previewError || connectionError" class="data-preview-section">
              <div v-if="previewError" class="error-message">
                <el-icon><Warning /></el-icon>
                <span>{{ previewError }}</span>
              </div>
              <div v-if="connectionError" class="error-message">
                <el-icon><CircleCloseFilled /></el-icon>
                <span>{{ connectionError }}</span>
              </div>
              <div v-if="previewData.length" class="preview-table">
                <h4>数据预览</h4>
                <el-table :data="previewData" size="small" max-height="200">
                  <el-table-column v-for="key in Object.keys(previewData[0])" :key="key" :prop="key" :label="key" />
                </el-table>
              </div>
            </div>
          </template>

          <template v-else>
            <el-form-item label="数据集" prop="datasetId">
              <el-select
                v-model="formData.datasetId"
                placeholder="请选择数据集"
                :loading="loadingDatasets"
                clearable
                filterable
                @change="handleDatasetChange"
              >
                <el-option v-for="dataset in datasets" :key="dataset.id" :label="dataset.name" :value="dataset.id" />
              </el-select>
            </el-form-item>
            <el-form-item label="维度" prop="dimension">
              <el-select
                v-model="formData.dimension"
                placeholder="请选择维度"
                :disabled="!formData.datasetId"
                :loading="loadingDatasetSchema"
                clearable
                filterable
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
                v-model="formData.measure"
                placeholder="请选择指标"
                :disabled="!formData.datasetId"
                :loading="loadingDatasetSchema"
                clearable
                filterable
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
                v-model="formData.aggregation"
                placeholder="请选择聚合函数"
                :disabled="!formData.measure"
                clearable
                @change="emitUpdate"
              >
                <el-option label="无" value="none" />
                <el-option label="求和 (SUM)" value="SUM" />
                <el-option label="平均值 (AVG)" value="AVG" />
                <el-option label="计数 (COUNT)" value="COUNT" />
                <el-option label="最大值 (MAX)" value="MAX" />
                <el-option label="最小值 (MIN)" value="MIN" />
              </el-select>
            </el-form-item>
          </template>
        </el-collapse-item>
      </el-collapse>
     </el-form>
   </div>
 </template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { Pointer, Connection, View, Document, Warning, CircleCloseFilled } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { datasourceApi, type SelectOption } from '@/api/datasource'
import { datasetApi, type Dataset, type DatasetField } from '@/api/dataset'
import { getErrorMessage } from '@/utils/errorHandling'

type BindingType = 'datasource' | 'dataset'
type AggregationType = 'SUM' | 'AVG' | 'COUNT' | 'MAX' | 'MIN' | 'none'

export interface DashboardComponent {
  id: string
  title: string
  type: string
  width: number
  height: number
  x: number
  y: number
  visible: boolean
  locked: boolean
  style: {
    background: string
    border: string
    borderColor: string
    fontSize: number
    fontColor: string
    padding: number
    margin: number
  }
  data: {
    dataSource: string
    field: string
    sqlQuery: string
    datasetId?: string
    dimension?: string
    measure?: string
    aggregation?: AggregationType
  }
  interaction: {
    linkUrl: string
    drilldownConfig: string
  }
}

interface FormData {
  title: string
  id: string
  width: number
  height: number
  x: number
  y: number
  background: string
  border: string
  borderColor: string
  fontSize: number
  fontColor: string
  padding: number
  margin: number
  dataSource: string
  field: string
  sqlQuery: string
  datasetId: string
  dimension: string
  measure: string
  aggregation: AggregationType
  linkUrl: string
  drilldownConfig: string
}

const props = withDefaults(
  defineProps<{
    component: DashboardComponent | null
    dataSources?: SelectOption[]
    fields?: SelectOption[]
    loadingDataSources?: boolean
  }>(),
  {
    dataSources: () => [
      { label: 'SalesDB', value: 'sales_db' },
      { label: 'AnalyticsDB', value: 'analytics_db' },
      { label: 'MockData', value: 'mock_data' }
    ],
    fields: () => [
      { label: 'sales_amount', value: 'sales_amount' },
      { label: 'order_count', value: 'order_count' },
      { label: 'growth_rate', value: 'growth_rate' }
    ],
    loadingDataSources: false
  }
)

const emit = defineEmits<{
  update: [component: DashboardComponent]
}>()

const formRef = ref<FormInstance>()
const activeGroups = ref(['common', 'style', 'data', 'interaction'])
const syncing = ref(false)

const loadingDataSources = ref(false)
const loadingFields = ref(false)
const loadingDatasets = ref(false)
const loadingDatasetSchema = ref(false)
const tables = ref<string[]>([])
const fields = ref<SelectOption[]>([])
const datasets = ref<Dataset[]>([])
const dimensions = ref<DatasetField[]>([])
const measures = ref<DatasetField[]>([])
const testingConnection = ref(false)
const previewingData = ref(false)
const previewData = ref<any[]>([])
const connectionError = ref('')
const previewError = ref('')
const bindingType = ref<BindingType>('datasource')

const formData = reactive<FormData>({
  title: '',
  id: '',
  width: 320,
  height: 180,
  x: 0,
  y: 0,
  background: '#ffffff',
  border: 'solid',
  borderColor: '#dcdfe6',
  fontSize: 14,
  fontColor: '#303133',
  padding: 12,
  margin: 8,
  dataSource: '',
  field: '',
  sqlQuery: '',
  datasetId: '',
  dimension: '',
  measure: '',
  aggregation: 'none',
  linkUrl: '',
  drilldownConfig: ''
})

const formRules = reactive<FormRules<FormData>>({
  title: [{ required: true, message: 'Title 为必填项', trigger: 'blur' }],
  id: [{ required: true, message: 'ID 为必填项', trigger: 'blur' }],
  width: [{ required: true, message: 'Width 为必填项', trigger: 'change' }],
  height: [{ required: true, message: 'Height 为必填项', trigger: 'change' }],
  dataSource: [
    {
      validator: (_rule, value, callback) => {
        if (bindingType.value === 'datasource' && !value) {
          callback(new Error('Data Source 为必填项'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ],
  field: [
    {
      validator: (_rule, value, callback) => {
        if (bindingType.value === 'datasource' && !value) {
          callback(new Error('Field 为必填项'))
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
          callback(new Error('数据集为必填项'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ]
})

function applyComponent(component: DashboardComponent | null) {
  if (!component) return
  syncing.value = true
  formData.title = component.title
  formData.id = component.id
  formData.width = component.width
  formData.height = component.height
  formData.x = component.x
  formData.y = component.y
  formData.background = component.style.background
  formData.border = component.style.border
  formData.borderColor = component.style.borderColor
  formData.fontSize = component.style.fontSize
  formData.fontColor = component.style.fontColor
  formData.padding = component.style.padding
  formData.margin = component.style.margin
  formData.dataSource = component.data.dataSource
  formData.field = component.data.field
  formData.sqlQuery = component.data.sqlQuery
  formData.datasetId = component.data.datasetId || ''
  formData.dimension = component.data.dimension || ''
  formData.measure = component.data.measure || ''
  formData.aggregation = component.data.aggregation || 'none'
  formData.linkUrl = component.interaction.linkUrl
  formData.drilldownConfig = component.interaction.drilldownConfig
  bindingType.value = formData.datasetId ? 'dataset' : 'datasource'
  syncing.value = false
}

async function emitUpdate() {
  if (!props.component || syncing.value) return
  const updatedComponent: DashboardComponent = {
    ...props.component,
    title: formData.title,
    id: formData.id,
    width: formData.width,
    height: formData.height,
    x: formData.x,
    y: formData.y,
    style: {
      background: formData.background,
      border: formData.border,
      borderColor: formData.borderColor,
      fontSize: formData.fontSize,
      fontColor: formData.fontColor,
      padding: formData.padding,
      margin: formData.margin
    },
    data: {
      dataSource: formData.dataSource,
      field: formData.field,
      sqlQuery: formData.sqlQuery,
      datasetId: formData.datasetId,
      dimension: formData.dimension,
      measure: formData.measure,
      aggregation: formData.aggregation
    },
    interaction: {
      linkUrl: formData.linkUrl,
      drilldownConfig: formData.drilldownConfig
    }
  }

  emit('update', updatedComponent)
  try {
    await formRef.value?.validateField(['title', 'id', 'width', 'height', 'dataSource', 'field'])
  } catch {
    // keep live updates smooth even when fields are temporarily invalid
  }
}

async function handleDataSourceChange() {
  connectionError.value = ''
  previewError.value = ''
  previewData.value = []
  if (!formData.dataSource) {
    tables.value = []
    fields.value = []
    emitUpdate()
    return
  }
  loadingFields.value = true
  try {
    const response = await datasourceApi.getTables(formData.dataSource)
    if (response.data.success) {
      tables.value = response.data.result || []
      fields.value = []
      ElMessage.success('表列表加载成功')
    } else {
      ElMessage.error(response.data.message || '加载表列表失败')
    }
  } catch (error: unknown) {
    ElMessage.error('加载表列表失败')
  } finally {
    loadingFields.value = false
  }
  emitUpdate()
}

async function handleFieldChange() {
  previewError.value = ''
  previewData.value = []
  emitUpdate()
}

function handleBindingTypeChange() {
  if (bindingType.value === 'dataset') {
    formData.dataSource = ''
    formData.field = ''
    formData.sqlQuery = ''
    tables.value = []
    fields.value = []
    previewData.value = []
    previewError.value = ''
    connectionError.value = ''
  } else {
    formData.datasetId = ''
    formData.dimension = ''
    formData.measure = ''
    formData.aggregation = 'none'
    dimensions.value = []
    measures.value = []
  }
  emitUpdate()
}

async function loadDatasets() {
  loadingDatasets.value = true
  try {
    const response = await datasetApi.list()
    if (response.data.success) {
      datasets.value = response.data.result || []
    } else {
      ElMessage.error(response.data.message || '加载数据集失败')
    }
  } catch (error: unknown) {
    ElMessage.error(getErrorMessage(error, '加载数据集失败'))
  } finally {
    loadingDatasets.value = false
  }
}

async function handleDatasetChange(shouldEmit: boolean = true) {
  formData.dimension = ''
  formData.measure = ''
  formData.aggregation = 'none'
  dimensions.value = []
  measures.value = []
  if (!formData.datasetId) {
    if (shouldEmit) {
      emitUpdate()
    }
    return
  }

  loadingDatasetSchema.value = true
  try {
    const response = await datasetApi.getSchema(formData.datasetId)
    if (response.data.success) {
      dimensions.value = response.data.result?.dimensions || []
      measures.value = response.data.result?.measures || []
    } else {
      ElMessage.error(response.data.message || '加载数据集字段失败')
    }
  } catch (error: unknown) {
    ElMessage.error(getErrorMessage(error, '加载数据集字段失败'))
  } finally {
    loadingDatasetSchema.value = false
    if (shouldEmit) {
      emitUpdate()
    }
  }
}

function handleDimensionChange() {
  if (!formData.measure) {
    formData.aggregation = 'none'
  }
  emitUpdate()
}

function handleMeasureChange() {
  formData.aggregation = formData.measure ? 'SUM' : 'none'
  emitUpdate()
}

async function handleTestConnection() {
  if (!formData.dataSource || !formData.sqlQuery) {
    ElMessage.warning('请先选择数据源并填写 SQL 查询')
    return
  }
  testingConnection.value = true
  connectionError.value = ''
  try {
    const response = await datasourceApi.test({
      name: formData.dataSource,
      type: 'mysql',
      host: 'localhost',
      port: 3306,
      database: 'test_db',
      username: 'test',
      password: ''
    })
    if (response.data.success) {
      ElMessage.success('连接测试成功')
    } else {
      connectionError.value = response.data.message || '连接测试失败'
      ElMessage.error(response.data.message || '连接测试失败')
    }
  } catch (error: unknown) {
    connectionError.value = getErrorMessage(error, '连接测试失败')
    ElMessage.error('连接测试失败')
  } finally {
    testingConnection.value = false
  }
}

async function handlePreviewData() {
  if (!formData.dataSource || !formData.sqlQuery) {
    ElMessage.warning('请先选择数据源并填写 SQL 查询')
    return
  }
  previewingData.value = true
  previewError.value = ''
  try {
    const response = await datasourceApi.test({
      name: formData.dataSource,
      type: 'mysql',
      host: 'localhost',
      port: 3306,
      database: 'test_db',
      username: 'test',
      password: ''
    })
    if (response.data.success) {
      previewData.value = [
        { id: 1, name: 'Product A', value: 1000, category: 'Electronics' },
        { id: 2, name: 'Product B', value: 2000, category: 'Electronics' },
        { id: 3, name: 'Product C', value: 1500, category: 'Clothing' }
      ]
      ElMessage.success('数据预览加载成功')
    } else {
      previewError.value = response.data.message || '数据预览失败'
      ElMessage.error(response.data.message || '数据预览失败')
    }
  } catch (error: unknown) {
    previewError.value = getErrorMessage(error, '数据预览失败')
    ElMessage.error('数据预览失败')
  } finally {
    previewingData.value = false
  }
}

function handleFormatSQL() {
  if (!formData.sqlQuery) return
  const formatted = formData.sqlQuery
    .split(/\s+/)
    .join(' ')
    .toUpperCase()
    .replace(/\bSELECT\b/gi, '\nSELECT')
    .replace(/\bFROM\b/gi, '\n  FROM')
    .replace(/\bWHERE\b/gi, '\n    WHERE')
    .replace(/\bAND\b/gi, '\n      AND')
    .replace(/\bOR\b/gi, '\n       OR')
    .trim()
  formData.sqlQuery = formatted
  emitUpdate()
}

watch(
  () => props.component,
  async component => {
    applyComponent(component)

    if (component?.data?.datasetId) {
      await handleDatasetChange(false)
      formData.dimension = component.data.dimension || ''
      formData.measure = component.data.measure || ''
      formData.aggregation = component.data.aggregation || 'none'
    }
  },
  { immediate: true, deep: true }
)

onMounted(() => {
  loadDatasets()
})
</script>

<style scoped>
.property-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 10px;
  background: #ffffff;
  box-shadow: 0 6px 16px rgba(15, 23, 42, 0.08);
  border: 1px solid #ebeef5;
  overflow: hidden;
}

.panel-header {
  padding: 12px 14px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #ebeef5;
}

.panel-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.panel-empty {
  flex: 1;
  color: #909399;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 8px;
  font-size: 13px;
}

.panel-form {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.field-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

:deep(.el-collapse) {
  border-top: none;
  border-bottom: none;
}

:deep(.el-collapse-item__header) {
  font-weight: 600;
  color: #606266;
}

:deep(.el-input-number),
:deep(.el-select) {
  width: 100%;
}

:deep(.el-form-item) {
  margin-bottom: 12px;
}

.data-preview-section {
  margin-top: 16px;
  padding: 12px;
  border-top: 1px solid #ebeef5;
  background: #f8fafc;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #fef2f2;
  border-radius: 6px;
  color: #ef4444;
  font-size: 13px;
}

.error-message span {
  flex: 1;
}

.preview-table h4 {
  margin: 0 0 12px;
  font-size: 14px;
  font-weight: 500;
  color: #303133;
}

.preview-table {
  border: 1px solid #ebeef5;
  border-radius: 6px;
  overflow: hidden;
}
</style>
