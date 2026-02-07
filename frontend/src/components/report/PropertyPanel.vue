<template>
  <div class="property-panel">
    <div class="panel-header">
      <span>单元格属性</span>
      <el-tag size="small" effect="dark" type="info">{{ cellLabel }}</el-tag>
    </div>
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="80px"
      class="panel-form"
    >
      <el-divider content-position="left">内容</el-divider>
      <el-form-item label="文本" prop="text">
        <el-input v-model="formData.text" placeholder="输入单元格文本" />
      </el-form-item>

      <el-divider content-position="left">样式</el-divider>
      <el-form-item label="字号" prop="fontSize">
        <el-input-number v-model="formData.fontSize" :min="10" :max="32" />
      </el-form-item>
      <el-form-item label="对齐" prop="align">
        <el-select v-model="formData.align" placeholder="请选择">
          <el-option label="左对齐" value="left" />
          <el-option label="居中" value="center" />
          <el-option label="右对齐" value="right" />
        </el-select>
      </el-form-item>
      <el-form-item label="加粗">
        <el-switch v-model="formData.bold" active-text="启用" />
      </el-form-item>
      <el-form-item label="斜体">
        <el-switch v-model="formData.italic" active-text="启用" />
      </el-form-item>
      <el-form-item label="文字色">
        <el-color-picker v-model="formData.color" show-alpha />
      </el-form-item>
      <el-form-item label="背景色">
        <el-color-picker v-model="formData.background" show-alpha />
      </el-form-item>

      <el-divider content-position="left">数据绑定</el-divider>
      <el-form-item label="数据源" prop="datasourceId">
        <el-select
          v-model="formData.datasourceId"
          placeholder="请选择数据源"
          clearable
          @change="handleDatasourceChange"
        >
          <el-option
            v-for="ds in datasources"
            :key="ds.id"
            :label="ds.name"
            :value="ds.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="数据表" prop="tableName">
        <el-select
          v-model="formData.tableName"
          placeholder="请选择数据表"
          clearable
          :disabled="!formData.datasourceId"
          :loading="tablesLoading"
          @change="handleTableChange"
        >
          <el-option v-for="table in tables" :key="table" :label="table" :value="table" />
        </el-select>
      </el-form-item>
      <el-form-item label="字段" prop="fieldName">
        <el-select
          v-model="formData.fieldName"
          placeholder="请选择字段"
          clearable
          :disabled="!formData.tableName"
          :loading="fieldsLoading"
        >
          <el-option
            v-for="field in fields"
            :key="field.name"
            :label="field.name"
            :value="field.name"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { datasourceApi, type DataSource } from '@/api/datasource'

interface CellStyle {
  fontSize: number
  fontWeight: 'normal' | 'bold'
  fontStyle: 'normal' | 'italic'
  align: 'left' | 'center' | 'right'
  color: string
  background: string
  borderColor?: string
}

interface CellBinding {
  datasourceId?: string
  tableName?: string
  fieldName?: string
}

interface CellData {
  row: number
  col: number
  text: string
  colSpan?: number
  style: CellStyle
  binding: CellBinding
}

interface PropertyForm {
  text: string
  fontSize: number
  align: 'left' | 'center' | 'right'
  bold: boolean
  italic: boolean
  color: string
  background: string
  datasourceId: string
  tableName: string
  fieldName: string
}

const props = defineProps<{ cell: CellData }>()
const emit = defineEmits<{ update: [cell: CellData] }>()

const formRef = ref<FormInstance>()
const datasources = ref<DataSource[]>([])
const tables = ref<string[]>([])
const fields = ref<Array<{ name: string }>>([])
const tablesLoading = ref(false)
const fieldsLoading = ref(false)
const saving = ref(false)

const formData = reactive<PropertyForm>({
  text: '',
  fontSize: 14,
  align: 'left',
  bold: false,
  italic: false,
  color: '#1f2933',
  background: '#ffffff',
  datasourceId: '',
  tableName: '',
  fieldName: ''
})

const formRules = reactive<FormRules<PropertyForm>>({
  tableName: [
    {
      validator: (_rule, value, callback) => {
        if (formData.datasourceId && !value) {
          callback(new Error('请选择数据表'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ],
  fieldName: [
    {
      validator: (_rule, value, callback) => {
        if (formData.datasourceId && formData.tableName && !value) {
          callback(new Error('请选择字段'))
          return
        }
        callback()
      },
      trigger: 'change'
    }
  ]
})

const cellLabel = computed(() => {
  const rowIndex = props.cell.row + 1
  const colLabel = String.fromCharCode(65 + props.cell.col)
  return `${colLabel}${rowIndex}`
})

function applyCellToForm(cell: CellData) {
  formData.text = cell.text || ''
  formData.fontSize = cell.style?.fontSize ?? 14
  formData.align = cell.style?.align ?? 'left'
  formData.bold = cell.style?.fontWeight === 'bold'
  formData.italic = cell.style?.fontStyle === 'italic'
  formData.color = cell.style?.color ?? '#1f2933'
  formData.background = cell.style?.background ?? '#ffffff'
  formData.datasourceId = cell.binding?.datasourceId ?? ''
  formData.tableName = cell.binding?.tableName ?? ''
  formData.fieldName = cell.binding?.fieldName ?? ''
}

async function loadDatasources() {
  try {
    const response = await datasourceApi.list()
    if (response.data.success) {
      datasources.value = response.data.result || []
    } else {
      ElMessage.error(response.data.message || '加载数据源失败')
    }
  } catch (error: any) {
    ElMessage.error('加载数据源失败')
  }
}

async function loadTables() {
  if (!formData.datasourceId) {
    tables.value = []
    return
  }

  tablesLoading.value = true
  try {
    const response = await datasourceApi.getTables(formData.datasourceId)
    if (response.data.success) {
      tables.value = response.data.result || []
    } else {
      ElMessage.error(response.data.message || '加载数据表失败')
    }
  } catch (error: any) {
    ElMessage.error('加载数据表失败')
  } finally {
    tablesLoading.value = false
  }
}

async function loadFields() {
  if (!formData.datasourceId || !formData.tableName) {
    fields.value = []
    return
  }

  fieldsLoading.value = true
  try {
    const response = await datasourceApi.getFields(formData.datasourceId, formData.tableName)
    if (response.data.success) {
      fields.value = response.data.result || []
    } else {
      ElMessage.error(response.data.message || '加载字段失败')
    }
  } catch (error: any) {
    ElMessage.error(error?.response?.data?.message || '加载字段失败')
  } finally {
    fieldsLoading.value = false
  }
}

function handleDatasourceChange() {
  formData.tableName = ''
  formData.fieldName = ''
  fields.value = []
  loadTables()
}

function handleTableChange() {
  formData.fieldName = ''
  loadFields()
}

async function handleSave() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  saving.value = true
  const updatedCell: CellData = {
    ...props.cell,
    text: formData.text,
    style: {
      fontSize: formData.fontSize,
      align: formData.align,
      fontWeight: formData.bold ? 'bold' : 'normal',
      fontStyle: formData.italic ? 'italic' : 'normal',
      color: formData.color,
      background: formData.background,
      borderColor: props.cell.style?.borderColor ?? '#cbd5e1'
    },
    binding: {
      datasourceId: formData.datasourceId || undefined,
      tableName: formData.tableName || undefined,
      fieldName: formData.fieldName || undefined
    }
  }

  emit('update', updatedCell)
  saving.value = false
  ElMessage.success('单元格属性已更新')
}

watch(
  () => props.cell,
  cell => {
    if (cell) {
      applyCellToForm(cell)
      if (cell.binding?.datasourceId) {
        loadTables().then(() => {
          if (cell.binding?.tableName) {
            loadFields()
          }
        })
      } else {
        tables.value = []
        fields.value = []
      }
    }
  },
  { immediate: true }
)

onMounted(() => {
  loadDatasources()
})
</script>

<style scoped>
.property-panel {
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

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  color: #0f172a;
  letter-spacing: 0.4px;
}

.panel-form {
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
