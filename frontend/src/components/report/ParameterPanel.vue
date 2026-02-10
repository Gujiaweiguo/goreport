<template>
  <div class="parameter-panel">
    <div class="panel-header">
      <span>报表参数</span>
      <el-button type="primary" size="small" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        添加参数
      </el-button>
    </div>

    <el-empty v-if="parameters.length === 0" description="暂无参数，点击上方按钮添加" />

    <div v-else class="parameter-list">
      <el-card v-for="(param, index) in parameters" :key="param.id" class="parameter-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span class="param-index">#{{ index + 1 }}</span>
            <el-input
              v-model="param.name"
              size="small"
              placeholder="参数名称"
              class="param-name-input"
            />
            <el-switch v-model="param.required" size="small" active-text="必填" />
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              circle
              @click="handleDelete(index)"
            />
          </div>
        </template>

        <el-form :model="param" label-width="80px" size="small">
          <el-form-item label="参数代码">
            <el-input v-model="param.code" placeholder="用于查询引用" />
          </el-form-item>

          <el-form-item label="参数类型">
            <el-select v-model="param.type" placeholder="请选择" @change="handleTypeChange(param)">
              <el-option label="字符串" value="string" />
              <el-option label="数字" value="number" />
              <el-option label="日期" value="date" />
              <el-option label="布尔值" value="boolean" />
              <el-option label="下拉选择" value="select" />
            </el-select>
          </el-form-item>

          <el-form-item label="默认值">
            <el-input
              v-if="param.type === 'string'"
              v-model="param.defaultValue"
              placeholder="默认值"
            />
            <el-input-number
              v-else-if="param.type === 'number'"
              v-model="param.defaultValue"
              placeholder="默认值"
            />
            <el-date-picker
              v-else-if="param.type === 'date'"
              v-model="param.defaultValue"
              type="date"
              placeholder="默认日期"
              value-format="YYYY-MM-DD"
            />
            <el-switch
              v-else-if="param.type === 'boolean'"
              v-model="param.defaultValue"
              active-text="是"
              inactive-text="否"
            />
            <div v-else-if="param.type === 'select'">
              <el-button size="small" @click="showOptionsEditor(param)">
                配置选项 ({{ param.options?.length || 0 }})
              </el-button>
            </div>
          </el-form-item>

          <el-form-item v-if="param.type === 'select'" label="可选值">
            <el-tag
              v-for="opt in param.options"
              :key="opt.value"
              closable
              @close="handleRemoveOption(param, opt)"
            >
              {{ opt.label }} ({{ opt.value }})
            </el-tag>
            <el-button
              size="small"
              @click="showAddOptionDialog(param)"
            >
              添加选项
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>

    <el-dialog
      v-model="addOptionDialog.visible"
      title="添加选项"
      width="500px"
      @close="addOptionDialog.visible = false"
    >
      <el-form :model="addOptionDialog" label-width="80px">
        <el-form-item label="显示标签">
          <el-input v-model="addOptionDialog.label" placeholder="如：销售部" />
        </el-form-item>
        <el-form-item label="选项值">
          <el-input v-model="addOptionDialog.value" placeholder="如：sales" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="addOptionDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="handleAddOptionConfirm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'

interface Option {
  label: string
  value: any
}

interface Parameter {
  id: string
  name: string
  code: string
  type: 'string' | 'number' | 'date' | 'boolean' | 'select'
  defaultValue: any
  required: boolean
  options: Option[]
}

const props = defineProps<{
  modelValue: Parameter[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: Parameter[]]
}>()

const parameters = ref<Parameter[]>([])
const addOptionDialog = reactive({
  visible: false,
  label: '',
  value: '',
  paramIndex: -1
})

function generateId() {
  return `param-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
}

function handleAdd() {
  const newParam: Parameter = {
    id: generateId(),
    name: '',
    code: '',
    type: 'string',
    defaultValue: '',
    required: false,
    options: []
  }
  parameters.value.push(newParam)
  emitUpdate()
}

function handleDelete(index: number) {
  parameters.value.splice(index, 1)
  emitUpdate()
}

function handleTypeChange(param: Parameter) {
  if (param.type === 'boolean') {
    param.defaultValue = false
  } else if (param.type === 'select') {
    param.defaultValue = ''
  } else if (param.type === 'number') {
    param.defaultValue = 0
  } else {
    param.defaultValue = ''
  }
  emitUpdate()
}

function showAddOptionDialog(param: Parameter) {
  const index = parameters.value.findIndex(p => p.id === param.id)
  addOptionDialog.visible = true
  addOptionDialog.label = ''
  addOptionDialog.value = ''
  addOptionDialog.paramIndex = index
}

function handleAddOptionConfirm() {
  if (!addOptionDialog.label || !addOptionDialog.value) {
    ElMessage.warning('请填写选项标签和值')
    return
  }

  const param = parameters.value[addOptionDialog.paramIndex]
  if (!param.options) {
    param.options = []
  }

  param.options.push({
    label: addOptionDialog.label,
    value: addOptionDialog.value
  })

  addOptionDialog.visible = false
  addOptionDialog.label = ''
  addOptionDialog.value = ''
  addOptionDialog.paramIndex = -1

  emitUpdate()
  ElMessage.success('选项已添加')
}

function handleRemoveOption(param: Parameter, option: Option) {
  const index = param.options?.findIndex(o => o.value === option.value)
  if (index !== undefined && index >= 0) {
    param.options!.splice(index, 1)
    emitUpdate()
  }
}

function emitUpdate() {
  emit('update:modelValue', JSON.parse(JSON.stringify(parameters.value)))
}

watch(
  () => props.modelValue,
  (newVal) => {
    if (newVal) {
      parameters.value = JSON.parse(JSON.stringify(newVal))
    }
  },
  { immediate: true, deep: true }
)
</script>

<style scoped>
.parameter-panel {
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
  font-weight: 600;
  color: #0f172a;
  letter-spacing: 0.4px;
}

.parameter-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 600px;
  overflow-y: auto;
}

.parameter-card {
  border: 1px solid #e2e8f0;
}

.parameter-card :deep(.el-card__header) {
  padding: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.param-index {
  font-weight: 600;
  color: #6366f1;
  min-width: 30px;
}

.param-name-input {
  flex: 1;
}

.param-name-input :deep(.el-input__wrapper) {
  border-radius: 6px;
}

:deep(.el-form-item) {
  margin-bottom: 12px;
}

:deep(.el-tag) {
  margin-right: 8px;
  margin-bottom: 8px;
}
</style>
