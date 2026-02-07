<template>
  <div class="data-config-panel">
    <div class="panel-header">
      <h3>数据配置</h3>
      <el-tag size="small" effect="plain" type="info">Data Binding</el-tag>
    </div>

    <el-form ref="formRef" :model="formState" :rules="rules" label-position="top" class="panel-form">
      <el-collapse v-model="activeGroups">
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

const loading = reactive({
  sources: false,
  tables: false,
  fields: false
})

const dataSourceOptions = ref<OptionItem[]>([])
const tableOptions = ref<OptionItem[]>([])
const fieldOptions = ref<OptionItem[]>([])

const formState = reactive<ChartDataConfig>({
  dataSourceId: '',
  tableName: '',
  fields: [],
  params: [],
  filter: ''
})

const rules = reactive<FormRules<ChartDataConfig>>({
  dataSourceId: [{ required: true, message: '请选择数据源', trigger: 'change' }],
  tableName: [{ required: true, message: '请选择数据表', trigger: 'change' }],
  fields: [{ required: true, message: '请至少选择一个字段', trigger: 'change' }]
})

function syncForm(value: ChartDataConfig) {
  formState.dataSourceId = value.dataSourceId || ''
  formState.tableName = value.tableName || ''
  formState.fields = [...(value.fields || [])]
  formState.params = [...(value.params || [])]
  formState.filter = value.filter || ''
}

function emitChange() {
  const nextValue: ChartDataConfig = {
    dataSourceId: formState.dataSourceId,
    tableName: formState.tableName,
    fields: [...formState.fields],
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
  value => {
    syncForm(value)
  },
  { immediate: true, deep: true }
)

onMounted(async () => {
  await loadDataSources()
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
