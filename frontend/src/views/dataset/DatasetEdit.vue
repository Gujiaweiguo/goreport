<template>
  <div class="dataset-edit">
    <el-card class="workflow-shell">
      <template #header>
        <DatasetEditHeader
          :current-dataset-id="currentDatasetId"
          :is-loading="workflowState.status === 'loading'"
          :operation="workflowState.operation"
          @save="handleSaveAction('save')"
          @save-and-return="handleSaveAction('save_and_return')"
          @refresh="handleRefreshData"
          @back="goBack"
        />
      </template>

      <el-alert
        v-if="workflowState.status === 'success' && workflowState.message"
        type="success"
        :title="workflowState.message"
        show-icon
        class="workflow-alert"
      />
      <el-alert
        v-if="workflowState.status === 'error' && workflowState.message"
        type="error"
        :title="workflowState.message"
        show-icon
        class="workflow-alert"
      />

      <div class="workflow-layout">
        <DatasetEditSidebar
          ref="sidebarRef"
          :form-data="formData"
          :form-rules="formRules"
          :datasources="datasources"
          :datasources-loading="datasourcesLoading"
          :tables-loading="tablesLoading"
          :sidebar-table-items="sidebarTableItems"
          :active-sidebar-table-key="activeSidebarTableKey"
          :file-list="fileList"
          @update:activeSidebarTableKey="activeSidebarTableKey = $event"
          @type-change="handleTypeChange"
          @datasource-change="handleDatasourceChange"
          @file-change="handleFileChange"
          @sidebar-item-dblclick="handleSidebarItemDblClick"
        />

        <section class="workflow-main" v-loading="workflowState.status === 'loading' && workflowState.operation === 'tab_switch'">
          <DatasetSourceWorkspace
            :form-data="formData"
            :selected-sources="selectedSources"
            :active-source-key="activeSourceKey"
            :join-config="joinConfig"
            :join-field-loading="joinFieldLoading"
            :left-join-field-options="leftJoinFieldOptions"
            :right-join-field-options="rightJoinFieldOptions"
            :show-sql-editor="showSqlEditor"
            @update:activeSourceKey="activeSourceKey = $event"
          />

          <el-tabs v-model="activeTab" @tab-change="handleTabChange">
            <el-tab-pane label="数据预览" name="preview">
              <div class="tab-content">
                <div class="tab-actions">
                  <el-button
                    type="primary"
                    :loading="workflowState.status === 'loading' && workflowState.operation === 'preview'"
                    @click="previewDataset"
                  >
                    预览数据
                  </el-button>
                  <el-button @click="closePreview">清空预览</el-button>
                </div>

                <el-empty
                  v-if="!previewData.length"
                  description="暂无预览数据，请先执行预览或刷新数据"
                />
                <DatasetPreview
                  v-else
                  :dataset-id="currentDatasetId"
                  :data="previewData"
                  :schema="schema"
                  @close="closePreview"
                />
              </div>
            </el-tab-pane>

            <el-tab-pane label="批量管理" name="batch">
              <div class="tab-content">
                <DatasetBatchPanel
                  ref="batchPanelRef"
                  :current-dataset-id="currentDatasetId"
                  :all-fields="allFields"
                  :selected-fields="selectedFields"
                  :batch-config="batchConfig"
                  :grouping-field-dialog-visible="groupingFieldDialogVisible"
                  :grouping-field-form="groupingFieldForm"
                  :grouping-field-rules="groupingFieldRules"
                  :is-loading="workflowState.status === 'loading' && workflowState.operation === 'batch_update'"
                  :is-submitting-grouping="workflowState.status === 'loading' && workflowState.operation === 'create_grouping_field'"
                  @selection-change="handleSelectionChange"
                  @open-grouping-dialog="openGroupingFieldDialog"
                  @reset-grouping-form="resetGroupingFieldForm"
                  @submit-batch="submitBatchUpdate"
                  @submit-grouping="submitGroupingField"
                  @update:groupingFieldDialogVisible="groupingFieldDialogVisible = $event"
                />
              </div>
            </el-tab-pane>
          </el-tabs>
        </section>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import DatasetEditHeader from '@/components/dataset/DatasetEditHeader.vue'
import DatasetEditSidebar from '@/components/dataset/DatasetEditSidebar.vue'
import DatasetSourceWorkspace from '@/components/dataset/DatasetSourceWorkspace.vue'
import DatasetBatchPanel from '@/components/dataset/DatasetBatchPanel.vue'
import DatasetPreview from '@/components/dataset/DatasetPreview.vue'
import { useDatasetEdit } from './composables/useDatasetEdit'

const route = useRoute()
const datasetId = (route.params.id as string) || ''

const {
  formRef,
  groupingFieldFormRef,
  currentDatasetId,
  activeTab,
  formData,
  formRules,
  workflowState,
  datasources,
  datasourcesLoading,
  tables,
  tablesLoading,
  activeSidebarTableKey,
  activeSourceKey,
  selectedSources,
  leftJoinFieldOptions,
  rightJoinFieldOptions,
  joinFieldLoading,
  joinConfig,
  fileList,
  previewData,
  schema,
  loadedSchemaDatasetId,
  loadedPreviewDatasetId,
  selectedFields,
  groupingFieldDialogVisible,
  groupingFieldForm,
  groupingFieldRules,
  batchConfig,
  allFields,
  hasCustomSqlSource,
  showSqlEditor,
  sidebarTableItems,
  handleSidebarItemDblClick,
  handleTypeChange,
  handleDatasourceChange,
  handleFileChange,
  previewDataset,
  handleSaveAction,
  handleRefreshData,
  handleTabChange,
  handleSelectionChange,
  submitBatchUpdate,
  resetGroupingFieldForm,
  openGroupingFieldDialog,
  submitGroupingField,
  closePreview,
  goBack,
  initialize
} = useDatasetEdit(datasetId)

// Component refs
const sidebarRef = ref<InstanceType<typeof DatasetEditSidebar>>()
const batchPanelRef = ref<InstanceType<typeof DatasetBatchPanel>>()

// Forward form validation from sidebar component
Object.defineProperty(formRef, 'value', {
  get: () => ({
    validate: () => sidebarRef.value?.validate() ?? Promise.resolve(),
    clearValidate: () => sidebarRef.value?.clearValidate()
  }),
  set: () => {}
})

// Forward grouping form validation from batch panel
Object.defineProperty(groupingFieldFormRef, 'value', {
  get: () => ({
    validate: () => batchPanelRef.value?.validateGroupingForm() ?? Promise.resolve(),
    clearValidate: () => {}
  }),
  set: () => {}
})

onMounted(() => {
  initialize()
})
</script>

<style scoped>
.dataset-edit {
  padding: 20px;
}

.workflow-shell {
  border-radius: 16px;
  border: 1px solid #d7dbe2;
  box-shadow: 0 16px 36px rgba(15, 23, 42, 0.08);
}

.workflow-alert {
  margin-bottom: 14px;
}

.workflow-layout {
  display: grid;
  grid-template-columns: 360px minmax(0, 1fr);
  gap: 16px;
  min-height: 640px;
}

.workflow-main {
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #ffffff;
}

.tab-content {
  min-height: 520px;
}

.tab-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 14px;
}

@media (max-width: 1280px) {
  .workflow-layout {
    grid-template-columns: 1fr;
  }
}
</style>
