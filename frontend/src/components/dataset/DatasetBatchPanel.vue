<template>
  <el-empty
    v-if="!currentDatasetId"
    description="请先保存数据集后使用批量管理"
  />

  <template v-else>
    <el-card class="batch-panel" shadow="never">
      <template #header>
        <div class="batch-header">
          <span>批量字段操作</span>
          <div class="batch-header-actions">
            <el-tag type="info">已选 {{ selectedFields.length }} 个字段</el-tag>
            <el-button type="primary" plain @click="$emit('open-grouping-dialog')">新建分组字段</el-button>
          </div>
        </div>
      </template>

      <div class="batch-actions">
        <el-select v-model="batchConfig.type" clearable placeholder="批量设置字段类型">
          <el-option label="维度" value="dimension" />
          <el-option label="指标" value="measure" />
        </el-select>
        <el-select v-model="batchConfig.sortOrder" clearable placeholder="批量设置排序方式">
          <el-option label="不排序" value="none" />
          <el-option label="升序" value="asc" />
          <el-option label="降序" value="desc" />
        </el-select>
        <el-input v-model="batchConfig.displayNamePrefix" placeholder="显示名前缀（可选）" />
        <el-button
          type="primary"
          :loading="isLoading"
          @click="$emit('submit-batch')"
        >
          提交批量更新
        </el-button>
      </div>

      <el-table
        :data="allFields"
        row-key="id"
        @selection-change="(fields: DatasetField[]) => $emit('selection-change', fields)"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="name" label="字段名" min-width="160" />
        <el-table-column label="显示名" min-width="160">
          <template #default="scope">
            {{ scope.row.displayName || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="类型" width="120">
          <template #default="scope">
            <el-tag :type="scope.row.type === 'dimension' ? 'success' : 'warning'">
              {{ scope.row.type === 'dimension' ? '维度' : '指标' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="dataType" label="数据类型" width="120" />
      </el-table>
    </el-card>

    <el-dialog
      :model-value="groupingFieldDialogVisible"
      title="新建分组字段"
      width="520px"
      @update:model-value="(val: boolean) => $emit('update:groupingFieldDialogVisible', val)"
      @closed="$emit('reset-grouping-form')"
    >
      <el-form
        ref="groupingFieldFormRef"
        :model="groupingFieldForm"
        :rules="groupingFieldRules"
        label-width="100px"
        status-icon
      >
        <el-form-item label="字段名" prop="name">
          <el-input v-model="groupingFieldForm.name" placeholder="请输入字段名" />
        </el-form-item>
        <el-form-item label="显示名" prop="displayName">
          <el-input v-model="groupingFieldForm.displayName" placeholder="请输入显示名（可选）" />
        </el-form-item>
        <el-form-item label="数据类型" prop="dataType">
          <el-select v-model="groupingFieldForm.dataType" placeholder="请选择数据类型">
            <el-option label="字符串 (string)" value="string" />
            <el-option label="数值 (number)" value="number" />
            <el-option label="日期 (date)" value="date" />
            <el-option label="布尔 (boolean)" value="boolean" />
          </el-select>
        </el-form-item>
        <el-form-item label="分组规则" prop="groupingRule">
          <el-input
            v-model="groupingFieldForm.groupingRule"
            placeholder="例如：region_group 或 category_bucket"
          />
        </el-form-item>
        <el-form-item label="启用分组" prop="groupingEnabled">
          <el-switch v-model="groupingFieldForm.groupingEnabled" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="grouping-field-dialog-footer">
          <el-button @click="$emit('update:groupingFieldDialogVisible', false)">取消</el-button>
          <el-button
            type="primary"
            :loading="isSubmittingGrouping"
            @click="$emit('submit-grouping')"
          >
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>
  </template>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import type { DatasetField } from '@/api/dataset'
import type { BatchConfig, GroupingFieldForm } from '@/views/dataset/composables/useDatasetEdit'

const groupingFieldFormRef = ref<FormInstance>()

defineProps<{
  currentDatasetId: string
  allFields: DatasetField[]
  selectedFields: DatasetField[]
  batchConfig: BatchConfig
  groupingFieldDialogVisible: boolean
  groupingFieldForm: GroupingFieldForm
  groupingFieldRules: FormRules
  isLoading: boolean
  isSubmittingGrouping: boolean
}>()

defineEmits<{
  'selection-change': [fields: DatasetField[]]
  'open-grouping-dialog': []
  'reset-grouping-form': []
  'submit-batch': []
  'submit-grouping': []
  'update:groupingFieldDialogVisible': [value: boolean]
}>()

defineExpose({
  validateGroupingForm: () => groupingFieldFormRef.value?.validate()
})
</script>

<style scoped>
.batch-panel {
  border-radius: 12px;
}

.batch-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.batch-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.batch-actions {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
  margin-bottom: 14px;
}

.grouping-field-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

@media (max-width: 1280px) {
  .batch-actions {
    grid-template-columns: 1fr;
  }
}
</style>
