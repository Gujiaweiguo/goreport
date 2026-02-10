<template>
  <el-dialog v-model="visible" :title="isEdit ? '编辑计算字段' : '新建计算字段'" width="900px" @close="handleClose">
    <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="字段名称" prop="name">
            <el-input v-model="formData.name" placeholder="请输入字段名称" />
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="显示名称" prop="displayName">
            <el-input v-model="formData.displayName" placeholder="请输入显示名称" />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="字段类型" prop="type">
            <el-radio-group v-model="formData.type">
              <el-radio value="dimension">维度</el-radio>
              <el-radio value="measure">指标</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item label="数据类型" prop="dataType">
            <el-select v-model="formData.dataType">
              <el-option value="string">字符串</el-option>
              <el-option value="number">数字</el-option>
              <el-option value="date">日期</el-option>
              <el-option value="boolean">布尔</el-option>
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <el-form-item label="表达式" prop="expression">
        <el-input
          v-model="formData.expression"
          type="textarea"
          :rows="4"
          placeholder="请输入计算字段表达式，如: [price] * [quantity]"
          @input="validateExpression"
        />
        <div v-if="expressionError" class="error-message">{{ expressionError }}</div>
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="12">
          <div class="field-list-container">
            <h4>可用字段</h4>
            <el-scrollbar max-height="200px">
              <div class="field-list">
                <el-tag
                  v-for="field in fields"
                  :key="field.name"
                  class="field-tag"
                  @click="insertField(field)"
                >
                  {{ field.displayName || field.name }}
                </el-tag>
              </div>
            </div>
          </el-scrollbar>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="function-list-container">
            <h4>函数库</h4>
            <el-scrollbar max-height="200px">
              <div class="function-list">
                <div
                  v-for="category in functionCategories"
                  :key="category.name"
                  class="function-category"
                >
                  <div class="category-name">{{ category.name }}</div>
                  <el-tag
                    v-for="fn in category.functions"
                    :key="fn.name"
                    size="small"
                    class="function-tag"
                    @click="insertFunction(fn)"
                  >
                    {{ fn.name }}
                  </el-tag>
                </div>
              </div>
            </div>
          </el-scrollbar>
          </div>
        </el-col>
      </el-row>
    </el-form>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" @click="save" :loading="saving">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { datasetApi, type CreateFieldRequest, type DatasetField } from '@/api/dataset'

interface Props {
  visible: boolean
  datasetId: string
  fields: DatasetField[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'refresh'): void
}>()

const formRef = ref()
const saving = ref(false)
const expressionError = ref('')

const formData = ref<CreateFieldRequest>({
  name: '',
  displayName: '',
  type: 'dimension',
  dataType: 'string',
  expression: ''
})

const formRules = {
  name: [{ required: true, message: '请输入字段名称', trigger: 'blur' }],
  expression: [{ required: true, message: '请输入表达式', trigger: 'blur' }],
  type: [{ required: true, message: '请选择字段类型', trigger: 'change' }],
  dataType: [{ required: true, message: '请选择数据类型', trigger: 'change' }]
}

const isEdit = computed(() => !!props.visible && props.datasetId)

const functionCategories = [
  {
    name: '聚合函数',
    functions: [
      { name: 'SUM', signature: 'SUM([field])' },
      { name: 'AVG', signature: 'AVG([field])' },
      { name: 'COUNT', signature: 'COUNT([field])' },
      { name: 'MAX', signature: 'MAX([field])' },
      { name: 'MIN', signature: 'MIN([field])' }
    ]
  },
  {
    name: '字符串函数',
    functions: [
      { name: 'CONCAT', signature: 'CONCAT([field1], [field2])' },
      { name: 'SUBSTRING', signature: 'SUBSTRING([field], start, length)' },
      { name: 'LENGTH', signature: 'LENGTH([field])' },
      { name: 'UPPER', signature: 'UPPER([field])' },
      { name: 'LOWER', signature: 'LOWER([field])' },
      { name: 'TRIM', signature: 'TRIM([field])' }
    ]
  },
  {
    name: '日期函数',
    functions: [
      { name: 'DATE_FORMAT', signature: 'DATE_FORMAT([field], format)' },
      { name: 'DATE_ADD', signature: 'DATE_ADD([field], interval)' },
      { name: 'DATEDIFF', signature: 'DATEDIFF([field1], [field2])' },
      { name: 'NOW', signature: 'NOW()' },
      { name: 'YEAR', signature: 'YEAR([field])' },
      { name: 'MONTH', signature: 'MONTH([field])' }
    ]
  },
  {
    name: '数学函数',
    functions: [
      { name: 'ROUND', signature: 'ROUND([field], decimals)' },
      { name: 'CEIL', signature: 'CEIL([field])' },
      { name: 'FLOOR', signature: 'FLOOR([field])' },
      { name: 'ABS', signature: 'ABS([field])' }
    ]
  }
]

const insertField = (field: DatasetField) => {
  const textarea = document.querySelector('textarea[placeholder*="计算字段表达式"]') as HTMLTextAreaElement
  if (textarea) {
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = formData.value.expression
    const fieldRef = `[${field.name}]`
    formData.value.expression = value.substring(0, start) + fieldRef + value.substring(end)
    textarea.focus()
  }
}

const insertFunction = (fn: any) => {
  const textarea = document.querySelector('textarea[placeholder*="计算字段表达式"]') as HTMLTextAreaElement
  if (textarea) {
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = formData.value.expression
    formData.value.expression = value.substring(0, start) + fn.name + '()' + value.substring(end)
    textarea.focus()
    const cursorPos = start + fn.name.length + 1
    textarea.setSelectionRange(cursorPos, cursorPos)
  }
}

const validateExpression = () => {
  const expression = formData.value.expression

  if (!expression || expression.trim() === '') {
    expressionError.value = '表达式不能为空'
    return
  }

  const fieldPattern = /\[([^\]]+)\]/g
  const matches = expression.match(fieldPattern)

  if (matches) {
    const fieldNames = matches.map(m => m.slice(1, -1))
    const validFieldNames = props.fields.map(f => f.name)
    for (const fieldName of fieldNames) {
      if (!validFieldNames.includes(fieldName)) {
        expressionError.value = `字段 [${fieldName}] 不存在`
        return
      }
    }
  }

  expressionError.value = ''
}

const save = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    validateExpression()

    if (expressionError.value) {
      return
    }

    saving.value = true
    const response = await datasetApi.createField(props.datasetId, formData.value)

    if (response.success) {
      ElMessage.success('保存成功')
      handleClose()
      emit('refresh')
    } else {
      ElMessage.error(response.message || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleClose = () => {
  formData.value = {
    name: '',
    displayName: '',
    type: 'dimension',
    dataType: 'string',
    expression: ''
  }
  expressionError.value = ''
  emit('update:visible', false)
}
</script>

<style scoped>
.field-list-container,
.function-list-container {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 10px;
  height: 250px;
}

.field-list-container h4,
.function-list-container h4 {
  margin: 0 0 10px 0;
  font-size: 14px;
  font-weight: bold;
}

.field-list,
.function-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.function-category {
  margin-bottom: 12px;
}

.category-name {
  font-size: 12px;
  color: #909399;
  margin-bottom: 6px;
}

.field-tag,
.function-tag {
  cursor: pointer;
  transition: all 0.3s;
}

.field-tag:hover,
.function-tag:hover {
  background-color: #409eff;
  border-color: #409eff;
  color: #fff;
}

.error-message {
  color: #f56c6c;
  font-size: 12px;
  margin-top: 4px;
}
</style>
