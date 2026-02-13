<template>
  <div class="field-list">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="维度" name="dimensions">
        <el-table :data="dimensions" stripe>
          <el-table-column prop="name" label="字段名" min-width="150" />
          <el-table-column prop="displayName" label="显示名称" min-width="150" />
          <el-table-column prop="dataType" label="数据类型" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ row.dataType }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="isGroupable" label="可分组" width="100">
            <template #default="{ row }">
              <el-icon v-if="row.isGroupable" color="success"><Check /></el-icon>
            </template>
          </el-table-column>
          <el-table-column prop="isSortable" label="可排序" width="100">
            <template #default="{ row }">
              <el-icon v-if="row.isSortable" color="success"><Check /></el-icon>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="editField(row)">编辑</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="指标" name="measures">
        <el-table :data="measures" stripe>
          <el-table-column prop="name" label="字段名" min-width="150" />
          <el-table-column prop="displayName" label="显示名称" min-width="150" />
          <el-table-column prop="dataType" label="数据类型" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ row.dataType }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="editField(row)">编辑</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="计算字段" name="computed">
        <el-table :data="computedFields" stripe>
          <el-table-column prop="name" label="字段名" min-width="150" />
          <el-table-column prop="displayName" label="显示名称" min-width="150" />
          <el-table-column prop="type" label="类型" width="100">
            <template #default="{ row }">
              <el-tag :type="row.type === 'dimension' ? 'info' : 'warning'" size="small">
                {{ row.type === 'dimension' ? '维度' : '指标' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="expression" label="表达式" min-width="200" show-overflow-tooltip />
          <el-table-column label="操作" width="200" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="editField(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="deleteField(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>

        <el-button type="primary" @click="createComputedField">新建计算字段</el-button>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="editDialogVisible" :title="isEdit ? '编辑字段' : '编辑字段配置'" width="500px">
      <el-form :model="editForm" label-width="120px">
        <el-form-item label="显示名称">
          <el-input v-model="editForm.displayName" placeholder="请输入显示名称" />
        </el-form-item>

        <el-form-item label="数据类型">
          <el-select v-model="editForm.dataType">
            <el-option value="string">字符串</el-option>
            <el-option value="number">数字</el-option>
            <el-option value="date">日期</el-option>
            <el-option value="boolean">布尔</el-option>
          </el-select>
        </el-form-item>

        <el-form-item v-if="isEdit" label="默认排序">
          <el-select v-model="editForm.sortOrder">
            <el-option value="asc">升序</el-option>
            <el-option value="desc">降序</el-option>
            <el-option value="none">无</el-option>
          </el-select>
        </el-form-item>

        <el-form-item v-if="!isEdit && editForm.type === 'dimension'" label="可分组">
          <el-switch v-model="editForm.isGroupable" />
        </el-form-item>

        <el-form-item v-if="!isEdit" label="可排序">
          <el-switch v-model="editForm.isSortable" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveField" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import { datasetApi, type DatasetField, type UpdateFieldRequest } from '@/api/dataset'

interface Props {
  datasetId: string
  fields: DatasetField[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'refresh'): void
  (e: 'openComputedFieldEditor'): void
}>()

const activeTab = ref('dimensions')
const editDialogVisible = ref(false)
const isEdit = ref(false)
const saving = ref(false)

const editForm = ref<UpdateFieldRequest>({
  displayName: '',
  dataType: 'string',
  isGroupable: true,
  isSortable: true,
  sortOrder: 'none'
})

const editingField = ref<DatasetField | null>(null)

const dimensions = computed(() => props.fields.filter(f => f.type === 'dimension' && !f.isComputed))
const measures = computed(() => props.fields.filter(f => f.type === 'measure' && !f.isComputed))
const computedFields = computed(() => props.fields.filter(f => f.isComputed))

const editField = (field: DatasetField) => {
  editingField.value = field
  editForm.value = {
    displayName: field.displayName || '',
    dataType: field.dataType,
    isGroupable: field.isGroupable,
    isSortable: field.isSortable,
    sortOrder: field.defaultSortOrder
  }
  isEdit.value = field.isComputed
  editDialogVisible.value = true
}

const saveField = async () => {
  if (!editingField.value) return

  saving.value = true
  try {
    const response = await datasetApi.updateField(
      props.datasetId,
      editingField.value.id,
      editForm.value
    )

    if (response.data.success) {
      ElMessage.success('保存成功')
      editDialogVisible.value = false
      emit('refresh')
    } else {
      ElMessage.error(response.data.message || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const deleteField = async (field: DatasetField) => {
  try {
    const response = await datasetApi.deleteField(props.datasetId, field.id)
    if (response.data.success) {
      ElMessage.success('删除成功')
      emit('refresh')
    } else {
      ElMessage.error(response.data.message || '删除失败')
    }
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const createComputedField = () => {
  emit('openComputedFieldEditor')
}
</script>

<style scoped>
.field-list {
  padding: 20px;
}
</style>
