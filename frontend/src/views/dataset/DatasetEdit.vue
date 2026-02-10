<template>
  <div class="dataset-edit">
    <el-card class="workflow-shell">
      <template #header>
        <div class="top-action-bar">
          <div class="title-wrap">
            <h2>{{ currentDatasetId ? '编辑数据集' : '新建数据集' }}</h2>
            <span class="title-sub">数据集编辑工作流</span>
          </div>
          <div class="action-buttons">
            <el-button
              type="primary"
              :loading="workflowState.status === 'loading' && workflowState.operation === 'save'"
              @click="handleSaveAction('save')"
            >
              保存
            </el-button>
            <el-button
              type="success"
              :loading="workflowState.status === 'loading' && workflowState.operation === 'save_and_return'"
              @click="handleSaveAction('save_and_return')"
            >
              保存并返回
            </el-button>
            <el-button
              :loading="workflowState.status === 'loading' && workflowState.operation === 'refresh'"
              @click="handleRefreshData"
            >
              刷新数据
            </el-button>
            <el-button @click="goBack">返回</el-button>
          </div>
        </div>
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
        <aside class="workflow-sidebar">
          <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" status-icon>
            <el-form-item label="数据集名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入数据集名称" />
            </el-form-item>

            <el-form-item label="数据集类型" prop="type">
              <el-radio-group v-model="formData.type" @change="handleTypeChange">
                <el-radio value="sql">SQL 查询</el-radio>
                <el-radio value="api">API 接口</el-radio>
                <el-radio value="file">文件导入</el-radio>
              </el-radio-group>
            </el-form-item>

            <el-form-item v-if="formData.type === 'sql'" label="数据源" prop="datasourceId">
              <el-select
                v-model="formData.datasourceId"
                placeholder="请选择数据源"
                :loading="datasourcesLoading"
                :no-data-text="'暂无数据源，请先创建数据源'"
                filterable
                @change="handleDatasourceChange"
              >
                <el-option
                  v-for="ds in datasources"
                  :key="ds.id"
                  :label="ds.name"
                  :value="ds.id"
                />
              </el-select>
              <div v-if="!datasources.length && !datasourcesLoading" class="no-datasource-tip">
                <el-icon><InfoFilled /></el-icon>
                <span>提示：请先在数据源管理中创建数据源</span>
              </div>
            </el-form-item>

            <el-form-item v-if="formData.type === 'sql'" label="数据表列表">
              <div class="table-list-shell" v-loading="tablesLoading">
                <div class="table-list-hint">双击添加到右侧来源面板</div>
                <div
                  v-for="item in sidebarTableItems"
                  :key="item.key"
                  class="table-list-item"
                  :class="{ 'is-active': activeSidebarTableKey === item.key }"
                  @click="activeSidebarTableKey = item.key"
                  @dblclick="handleSidebarItemDblClick(item)"
                >
                  <span>{{ item.label }}</span>
                  <el-tag v-if="item.type === 'custom_sql'" size="small" type="warning">SQL</el-tag>
                </div>
                <el-empty
                  v-if="!formData.datasourceId"
                  :image-size="46"
                  description="请先选择数据源加载数据表"
                />
              </div>
            </el-form-item>

            <el-form-item v-if="formData.type === 'api'" label="API 端点" prop="config.url">
              <el-input v-model="formData.config.url" placeholder="请输入 API 端点 URL" />
            </el-form-item>

            <el-form-item v-if="formData.type === 'api'" label="请求方法">
              <el-select v-model="formData.config.method">
                <el-option value="GET">GET</el-option>
                <el-option value="POST">POST</el-option>
              </el-select>
            </el-form-item>

            <el-form-item v-if="formData.type === 'file'" label="文件上传" prop="config.file">
              <el-upload
                :auto-upload="false"
                :on-change="handleFileChange"
                :file-list="fileList"
              >
                <el-button>选择文件</el-button>
                <template #tip>支持 Excel、CSV 格式</template>
              </el-upload>
            </el-form-item>
          </el-form>
        </aside>

        <section class="workflow-main" v-loading="workflowState.status === 'loading' && workflowState.operation === 'tab_switch'">
          <div v-if="formData.type === 'sql'" class="source-workspace">
            <div class="source-workspace-header">
              <span>来源面板</span>
              <span class="source-workspace-tip">双击左侧表格可添加或更新来源</span>
            </div>

            <div class="source-card-grid">
              <el-empty v-if="!selectedSources.length" :image-size="52" description="暂无来源，双击左侧列表添加" />
              <div
                v-for="source in selectedSources"
                :key="source.key"
                class="source-card"
                :class="{ 'is-active': activeSourceKey === source.key }"
                @click="activeSourceKey = source.key"
              >
                <div class="source-card-title">{{ source.label }}</div>
                <div class="source-card-desc">{{ source.type === 'custom_sql' ? '自定义SQL来源' : '数据表来源' }}</div>
              </div>
            </div>

            <div v-if="selectedSources.length >= 2" class="join-editor">
              <div class="join-editor-title">关联配置</div>
              <div class="join-editor-grid">
                <el-select v-model="joinConfig.leftSourceKey" placeholder="左来源">
                  <el-option
                    v-for="source in selectedSources"
                    :key="`left-${source.key}`"
                    :label="source.label"
                    :value="source.key"
                  />
                </el-select>
                <el-select v-model="joinConfig.relationType" placeholder="关联类型">
                  <el-option label="内连接（inner）" value="inner" />
                  <el-option label="左连接（left）" value="left" />
                  <el-option label="右连接（right）" value="right" />
                  <el-option label="全连接（full）" value="full" />
                </el-select>
                <el-select v-model="joinConfig.rightSourceKey" placeholder="右来源">
                  <el-option
                    v-for="source in selectedSources"
                    :key="`right-${source.key}`"
                    :label="source.label"
                    :value="source.key"
                  />
                </el-select>
                <el-select
                  v-model="joinConfig.leftField"
                  placeholder="左字段"
                  :loading="joinFieldLoading.left"
                  :disabled="!leftJoinFieldOptions.length"
                >
                  <el-option
                    v-for="field in leftJoinFieldOptions"
                    :key="`left-field-${field}`"
                    :label="field"
                    :value="field"
                  />
                </el-select>
                <el-select
                  v-model="joinConfig.rightField"
                  placeholder="右字段"
                  :loading="joinFieldLoading.right"
                  :disabled="!rightJoinFieldOptions.length"
                >
                  <el-option
                    v-for="field in rightJoinFieldOptions"
                    :key="`right-field-${field}`"
                    :label="field"
                    :value="field"
                  />
                </el-select>
              </div>
            </div>

            <div v-if="showSqlEditor" class="custom-sql-editor">
              <div class="join-editor-title">自定义SQL</div>
              <el-input
                v-model="formData.config.query"
                type="textarea"
                :rows="7"
                placeholder="请输入 SQL 查询语句"
              />
              <div class="form-tip">示例: SELECT * FROM users WHERE status = 1</div>
            </div>
          </div>

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
                  @close="closePreview"
                />
              </div>
            </el-tab-pane>

            <el-tab-pane label="批量管理" name="batch">
              <div class="tab-content">
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
                          <el-button type="primary" plain @click="openGroupingFieldDialog">新建分组字段</el-button>
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
                        :loading="workflowState.status === 'loading' && workflowState.operation === 'batch_update'"
                        @click="submitBatchUpdate"
                      >
                        提交批量更新
                      </el-button>
                    </div>

                    <el-table
                      :data="allFields"
                      row-key="id"
                      @selection-change="handleSelectionChange"
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
                    v-model="groupingFieldDialogVisible"
                    title="新建分组字段"
                    width="520px"
                    @closed="resetGroupingFieldForm"
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
                        <el-button @click="groupingFieldDialogVisible = false">取消</el-button>
                        <el-button
                          type="primary"
                          :loading="workflowState.status === 'loading' && workflowState.operation === 'create_grouping_field'"
                          @click="submitGroupingField"
                        >
                          确定
                        </el-button>
                      </div>
                    </template>
                  </el-dialog>
                </template>
              </div>
            </el-tab-pane>
          </el-tabs>
        </section>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { InfoFilled } from '@element-plus/icons-vue'
import type { FormInstance, FormRules, UploadUserFile } from 'element-plus'
import { datasourceApi, type DataSource } from '@/api/datasource'
import {
  datasetApi,
  type CreateFieldRequest,
  type CreateDatasetRequest,
  type DatasetField,
  type DatasetSchema,
  type UpdateDatasetRequest,
  type UpdateFieldRequest
} from '@/api/dataset'
import DatasetPreview from '@/components/dataset/DatasetPreview.vue'

type WorkflowStatus = 'idle' | 'loading' | 'success' | 'error'

interface WorkflowState {
  status: WorkflowStatus
  operation: 'init' | 'save' | 'save_and_return' | 'preview' | 'refresh' | 'tab_switch' | 'batch_update' | 'create_grouping_field' | ''
  message: string
}

interface BatchConfig {
  type?: 'dimension' | 'measure'
  sortOrder?: 'asc' | 'desc' | 'none'
  displayNamePrefix: string
}

interface GroupingFieldForm {
  name: string
  displayName: string
  dataType: 'string' | 'number' | 'date' | 'boolean'
  groupingRule: string
  groupingEnabled: boolean
}

type SourceType = 'custom_sql' | 'table'
type JoinType = 'inner' | 'left' | 'right' | 'full'

interface SelectedSource {
  key: string
  label: string
  type: SourceType
  table?: string
}

interface SidebarTableItem {
  key: string
  label: string
  type: SourceType
  table?: string
}

interface JoinConfigState {
  leftSourceKey: string
  rightSourceKey: string
  relationType: JoinType
  leftField: string
  rightField: string
}

type CreateGroupingFieldRequest = CreateFieldRequest & {
  isGroupingField: boolean
  groupingRule: string
  groupingEnabled: boolean
}

const route = useRoute()
const router = useRouter()

const formRef = ref<FormInstance>()
const activeTab = ref('preview')
const currentDatasetId = ref<string>((route.params.id as string) || '')

const formData = ref<CreateDatasetRequest>({
  name: '',
  type: 'sql',
  datasourceId: '',
  config: {}
})

const formRules = reactive<FormRules>({
  name: [{ required: true, message: '请输入数据集名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择数据集类型', trigger: 'change' }],
  datasourceId: [{ required: true, message: '请选择数据源', trigger: 'change' }]
})

const workflowState = ref<WorkflowState>({ status: 'idle', operation: '', message: '' })
const datasources = ref<DataSource[]>([])
const datasourcesLoading = ref(false)
const tables = ref<string[]>([])
const tablesLoading = ref(false)
const activeSidebarTableKey = ref<string>('custom_sql')
const activeSourceKey = ref<string>('')
const selectedSources = ref<SelectedSource[]>([])
const leftJoinFieldOptions = ref<string[]>([])
const rightJoinFieldOptions = ref<string[]>([])
const joinFieldLoading = reactive({ left: false, right: false })
const joinConfig = ref<JoinConfigState>({
  leftSourceKey: '',
  rightSourceKey: '',
  relationType: 'inner',
  leftField: '',
  rightField: ''
})
const fileList = ref<UploadUserFile[]>([])
const previewData = ref<Record<string, any>[]>([])
const schema = ref<DatasetSchema>({ dimensions: [], measures: [], computed: [] })
const selectedFields = ref<DatasetField[]>([])
const groupingFieldDialogVisible = ref(false)
const groupingFieldFormRef = ref<FormInstance>()
const groupingFieldForm = ref<GroupingFieldForm>({
  name: '',
  displayName: '',
  dataType: 'string',
  groupingRule: '',
  groupingEnabled: false
})
const groupingFieldRules = reactive<FormRules>({
  name: [{ required: true, message: '请输入字段名', trigger: 'blur' }],
  groupingRule: [{ required: true, message: '请输入分组规则', trigger: 'blur' }],
  dataType: [{ required: true, message: '请选择数据类型', trigger: 'change' }]
})

const batchConfig = ref<BatchConfig>({
  displayNamePrefix: ''
})

const CUSTOM_SQL_KEY = 'custom_sql'
const CUSTOM_SQL_LABEL = '自定义SQL'
const MAX_SELECTED_SOURCES = 2

const allFields = computed(() => [...schema.value.dimensions, ...schema.value.measures, ...schema.value.computed])
const hasCustomSqlSource = computed(() => selectedSources.value.some((source) => source.type === 'custom_sql'))
const showSqlEditor = computed(() => formData.value.type === 'sql' && hasCustomSqlSource.value)
const sidebarTableItems = computed<SidebarTableItem[]>(() => {
  const baseItems: SidebarTableItem[] = [{ key: CUSTOM_SQL_KEY, label: CUSTOM_SQL_LABEL, type: 'custom_sql' }]
  if (!formData.value.datasourceId) {
    return baseItems
  }

  const tableItems = tables.value.map((table) => ({
    key: `table:${table}`,
    label: table,
    type: 'table' as const,
    table
  }))

  return [...baseItems, ...tableItems]
})

const getErrorMessage = (error: unknown, fallback: string): string => {
  if (typeof error === 'object' && error !== null) {
    const maybeError = error as { response?: { data?: { message?: string } }; message?: string }
    return maybeError.response?.data?.message || maybeError.message || fallback
  }
  return fallback
}

const runWithState = async <T>(
  operation: WorkflowState['operation'],
  runner: () => Promise<T>,
  options?: { successMessage?: string; errorMessage?: string; notifySuccess?: boolean }
): Promise<T> => {
  workflowState.value = {
    status: 'loading',
    operation,
    message: ''
  }

  try {
    const result = await runner()
    workflowState.value = {
      status: 'success',
      operation,
      message: options?.successMessage || ''
    }
    if (options?.successMessage && options.notifySuccess !== false) {
      ElMessage.success(options.successMessage)
    }
    return result
  } catch (error) {
    const message = getErrorMessage(error, options?.errorMessage || '操作失败')
    workflowState.value = {
      status: 'error',
      operation,
      message
    }
    ElMessage.error(message)
    throw error
  }
}

const loadDatasources = async () => {
  datasourcesLoading.value = true
  try {
    const response = await datasourceApi.list()
    if (response.data.success) {
      datasources.value = response.data.result || []
      return
    }
    throw new Error(response.data.message || '加载数据源失败')
  } catch (error) {
    ElMessage.error(getErrorMessage(error, '加载数据源失败'))
  } finally {
    datasourcesLoading.value = false
  }
}

const loadTables = async (datasourceId: string) => {
  if (!datasourceId) {
    tables.value = []
    return
  }

  try {
    tablesLoading.value = true
    const response = await datasourceApi.getTables(datasourceId)
    if (response.data.success) {
      tables.value = response.data.result || []
      return
    }
    tables.value = []
  } catch {
    ElMessage.error('加载表格列表失败')
    tables.value = []
  } finally {
    tablesLoading.value = false
  }
}

const ensureSqlConfig = () => {
  if (!formData.value.config || typeof formData.value.config !== 'object') {
    formData.value.config = {}
  }

  if (!Array.isArray(formData.value.config.sources)) {
    formData.value.config.sources = []
  }

  if (!formData.value.config.join || typeof formData.value.config.join !== 'object') {
    formData.value.config.join = {
      leftSourceKey: '',
      rightSourceKey: '',
      relationType: 'inner',
      leftField: '',
      rightField: ''
    }
  }

  if (typeof formData.value.config.query !== 'string') {
    formData.value.config.query = ''
  }
}

const syncSourcesToConfig = () => {
  ensureSqlConfig()
  formData.value.config.sources = selectedSources.value.map((source, index) => ({
    key: source.key,
    label: source.label,
    type: source.type,
    table: source.table,
    sortIndex: index
  }))
}

const syncJoinToConfig = () => {
  ensureSqlConfig()
  formData.value.config.join = { ...joinConfig.value }
}

const resetJoinFieldOptions = () => {
  leftJoinFieldOptions.value = []
  rightJoinFieldOptions.value = []
  joinFieldLoading.left = false
  joinFieldLoading.right = false
}

const resetSqlWorkspace = () => {
  selectedSources.value = []
  activeSourceKey.value = ''
  joinConfig.value = {
    leftSourceKey: '',
    rightSourceKey: '',
    relationType: 'inner',
    leftField: '',
    rightField: ''
  }
  resetJoinFieldOptions()

  if (formData.value.type === 'sql') {
    ensureSqlConfig()
    formData.value.config.sources = []
    formData.value.config.join = { ...joinConfig.value }
    formData.value.config.query = ''
  }
}

const normalizeFieldOptions = (fields: any[]): string[] => {
  return fields
    .map((field) => {
      if (typeof field === 'string') {
        return field
      }
      if (field && typeof field.name === 'string') {
        return field.name
      }
      if (field && typeof field.columnName === 'string') {
        return field.columnName
      }
      return ''
    })
    .filter((field): field is string => !!field)
}

const loadJoinFieldsBySource = async (sourceKey: string, side: 'left' | 'right') => {
  const source = selectedSources.value.find((item) => item.key === sourceKey)
  if (!source || source.type !== 'table' || !source.table || !formData.value.datasourceId) {
    if (side === 'left') {
      leftJoinFieldOptions.value = []
    } else {
      rightJoinFieldOptions.value = []
    }
    return
  }

  joinFieldLoading[side] = true
  try {
    const response = await datasourceApi.getFields(formData.value.datasourceId, source.table)
    const options = response.data.success ? normalizeFieldOptions(response.data.result || []) : []
    if (side === 'left') {
      leftJoinFieldOptions.value = options
    } else {
      rightJoinFieldOptions.value = options
    }
  } catch {
    if (side === 'left') {
      leftJoinFieldOptions.value = []
    } else {
      rightJoinFieldOptions.value = []
    }
    ElMessage.error('加载关联字段失败')
  } finally {
    joinFieldLoading[side] = false
  }
}

const joinKeywordMap: Record<JoinType, string> = {
  inner: 'INNER JOIN',
  left: 'LEFT JOIN',
  right: 'RIGHT JOIN',
  full: 'FULL JOIN'
}

const syncAutoQueryFromSources = () => {
  if (formData.value.type !== 'sql') {
    return
  }
  ensureSqlConfig()
  if (hasCustomSqlSource.value) {
    return
  }

  const tableSources = selectedSources.value.filter((source) => source.type === 'table' && source.table)
  if (!tableSources.length) {
    formData.value.config.query = ''
    return
  }

  if (tableSources.length === 1) {
    formData.value.config.query = `SELECT * FROM ${tableSources[0].table}`
    return
  }

  const leftSource = selectedSources.value.find((source) => source.key === joinConfig.value.leftSourceKey)
  const rightSource = selectedSources.value.find((source) => source.key === joinConfig.value.rightSourceKey)
  const fallbackLeft = tableSources[0]
  const fallbackRight = tableSources[1]
  const actualLeft = leftSource?.type === 'table' ? leftSource : fallbackLeft
  const actualRight = rightSource?.type === 'table' ? rightSource : fallbackRight

  if (!actualLeft?.table || !actualRight?.table) {
    formData.value.config.query = `SELECT * FROM ${fallbackLeft.table}`
    return
  }

  if (!joinConfig.value.leftField || !joinConfig.value.rightField) {
    formData.value.config.query = `SELECT * FROM ${actualLeft.table}`
    return
  }

  formData.value.config.query = `SELECT * FROM ${actualLeft.table} t1 ${joinKeywordMap[joinConfig.value.relationType]} ${actualRight.table} t2 ON t1.${joinConfig.value.leftField} = t2.${joinConfig.value.rightField}`
}

const applyDefaultJoinSources = () => {
  if (selectedSources.value.length < 2) {
    joinConfig.value.leftSourceKey = ''
    joinConfig.value.rightSourceKey = ''
    joinConfig.value.leftField = ''
    joinConfig.value.rightField = ''
    return
  }

  if (!selectedSources.value.some((source) => source.key === joinConfig.value.leftSourceKey)) {
    joinConfig.value.leftSourceKey = selectedSources.value[0].key
    joinConfig.value.leftField = ''
  }

  if (!selectedSources.value.some((source) => source.key === joinConfig.value.rightSourceKey)) {
    const fallbackRight = selectedSources.value[1] || selectedSources.value[0]
    joinConfig.value.rightSourceKey = fallbackRight.key
    joinConfig.value.rightField = ''
  }
}

const handleSidebarItemDblClick = async (item: SidebarTableItem) => {
  if (formData.value.type !== 'sql') {
    return
  }

  ensureSqlConfig()

  const existingIndex = selectedSources.value.findIndex((source) => source.key === item.key)
  if (existingIndex >= 0) {
    const existing = selectedSources.value[existingIndex]
    selectedSources.value.splice(existingIndex, 1)
    selectedSources.value.push(existing)
  } else {
    if (selectedSources.value.length >= MAX_SELECTED_SOURCES) {
      selectedSources.value.shift()
      ElMessage.info('最多配置两个来源，已替换最早添加的来源')
    }

    selectedSources.value.push({
      key: item.key,
      label: item.label,
      type: item.type,
      table: item.table
    })
  }

  activeSourceKey.value = item.key

  if (item.type === 'custom_sql' && !formData.value.config.query) {
    formData.value.config.query = ''
  }

  applyDefaultJoinSources()
  syncSourcesToConfig()
  syncJoinToConfig()
  syncAutoQueryFromSources()
  await Promise.all([
    loadJoinFieldsBySource(joinConfig.value.leftSourceKey, 'left'),
    loadJoinFieldsBySource(joinConfig.value.rightSourceKey, 'right')
  ])
}

const parseSelectedTable = (query: unknown): string => {
  if (typeof query !== 'string') {
    return ''
  }

  const match = query.match(/^\s*SELECT\s+\*\s+FROM\s+([`"\[]?[\w.]+[`"\]]?)\s*;?\s*$/i)
  return match?.[1]?.replace(/[`"\[\]]/g, '') || ''
}

const restoreSqlWorkspaceFromConfig = async () => {
  if (formData.value.type !== 'sql') {
    return
  }

  ensureSqlConfig()
  selectedSources.value = []

  if (Array.isArray(formData.value.config.sources) && formData.value.config.sources.length) {
    selectedSources.value = formData.value.config.sources
      .slice(0, MAX_SELECTED_SOURCES)
      .map((source: any) => ({
        key: source.key || (source.type === 'custom_sql' ? CUSTOM_SQL_KEY : `table:${source.table}`),
        label: source.label || (source.type === 'custom_sql' ? CUSTOM_SQL_LABEL : source.table || ''),
        type: source.type === 'custom_sql' ? 'custom_sql' : 'table',
        table: source.table
      }))
      .filter((source: SelectedSource) => !!source.label)
  }

  if (!selectedSources.value.length) {
    const parsedTable = parseSelectedTable(formData.value.config.query)
    if (parsedTable) {
      selectedSources.value = [{ key: `table:${parsedTable}`, label: parsedTable, type: 'table', table: parsedTable }]
    } else if (formData.value.config.query) {
      selectedSources.value = [{ key: CUSTOM_SQL_KEY, label: CUSTOM_SQL_LABEL, type: 'custom_sql' }]
    }
  }

  const configJoin = formData.value.config.join || {}
  joinConfig.value = {
    leftSourceKey: typeof configJoin.leftSourceKey === 'string' ? configJoin.leftSourceKey : '',
    rightSourceKey: typeof configJoin.rightSourceKey === 'string' ? configJoin.rightSourceKey : '',
    relationType: ['inner', 'left', 'right', 'full'].includes(configJoin.relationType) ? configJoin.relationType : 'inner',
    leftField: typeof configJoin.leftField === 'string' ? configJoin.leftField : '',
    rightField: typeof configJoin.rightField === 'string' ? configJoin.rightField : ''
  }

  activeSourceKey.value = selectedSources.value[0]?.key || ''
  applyDefaultJoinSources()
  syncSourcesToConfig()
  syncJoinToConfig()
  syncAutoQueryFromSources()

  await Promise.all([
    loadJoinFieldsBySource(joinConfig.value.leftSourceKey, 'left'),
    loadJoinFieldsBySource(joinConfig.value.rightSourceKey, 'right')
  ])
}

const loadDataset = async (id: string) => {
  const response = await datasetApi.get(id)
  if (!response.data.success) {
    throw new Error(response.data.message || '加载数据集失败')
  }

  const dataset = response.data.result
  formData.value = {
    name: dataset.name,
    type: dataset.type,
    datasourceId: dataset.datasourceId,
    config: typeof dataset.config === 'string' ? JSON.parse(dataset.config) : dataset.config || {}
  }
}

const loadSchema = async (id: string) => {
  const response = await datasetApi.getSchema(id)
  if (!response.data.success) {
    throw new Error(response.data.message || '加载字段结构失败')
  }
  schema.value = response.data.result || { dimensions: [], measures: [], computed: [] }
}

const ensureDatasetPersisted = async () => {
  if (currentDatasetId.value) {
    return currentDatasetId.value
  }

  const response = await datasetApi.create({
    name: formData.value.name,
    type: formData.value.type,
    datasourceId: formData.value.datasourceId,
    config: formData.value.config
  })

  if (!response.data.success || !response.data.result?.id) {
    throw new Error(response.data.message || '创建数据集失败')
  }

  const newId = response.data.result.id
  currentDatasetId.value = newId
  await router.replace(`/dataset/edit/${newId}`)
  return newId
}

const fetchPreview = async (id: string) => {
  const response = await datasetApi.preview(id)
  if (!response.data.success) {
    throw new Error(response.data.message || '预览失败')
  }
  previewData.value = response.data.result || []
}

const handleTypeChange = () => {
  formData.value.config = {}
  fileList.value = []
  previewData.value = []
  activeSidebarTableKey.value = CUSTOM_SQL_KEY
  resetSqlWorkspace()

  if (formData.value.type === 'sql' && formData.value.datasourceId) {
    ensureSqlConfig()
    loadTables(formData.value.datasourceId)
    return
  }

  tables.value = []
}

const handleDatasourceChange = async (datasourceId: string) => {
  if (formData.value.type === 'sql') {
    resetSqlWorkspace()
  }

  if (formData.value.type === 'sql' && datasourceId) {
    ensureSqlConfig()
    await loadTables(datasourceId)
    return
  }

  tables.value = []
}

const handleFileChange = (file: UploadUserFile) => {
  fileList.value = [file]
  formData.value.config.file = (file as UploadUserFile & { raw?: File }).raw
}

const previewDataset = async () => {
  if (!formRef.value) {
    return
  }

  await formRef.value.validate()
  await runWithState(
    'preview',
    async () => {
      const id = await ensureDatasetPersisted()
      await fetchPreview(id)
    },
    { successMessage: '预览成功', errorMessage: '预览失败' }
  )
}

const handleSaveAction = async (action: 'save' | 'save_and_return') => {
  if (!formRef.value) {
    return
  }

  await formRef.value.validate()

  await runWithState(
    action,
    async () => {
      if (currentDatasetId.value) {
        const updatePayload: UpdateDatasetRequest = {
          name: formData.value.name,
          config: formData.value.config,
          action
        }
        const updateResponse = await datasetApi.update(currentDatasetId.value, updatePayload)
        if (!updateResponse.data.success) {
          throw new Error(updateResponse.data.message || '保存失败')
        }
      } else {
        const createResponse = await datasetApi.create({
          name: formData.value.name,
          type: formData.value.type,
          datasourceId: formData.value.datasourceId,
          config: formData.value.config
        })
        if (!createResponse.data.success || !createResponse.data.result?.id) {
          throw new Error(createResponse.data.message || '保存失败')
        }

        currentDatasetId.value = createResponse.data.result.id
        if (action === 'save') {
          await router.replace(`/dataset/edit/${currentDatasetId.value}`)
        }
      }

      await loadSchema(currentDatasetId.value)

      if (action === 'save_and_return') {
        goBack()
      }
    },
    {
      successMessage: action === 'save' ? '保存成功' : '保存成功，已返回列表',
      errorMessage: '保存失败'
    }
  )
}

const handleRefreshData = async () => {
  await runWithState(
    'refresh',
    async () => {
      if (!currentDatasetId.value) {
        throw new Error('请先保存数据集后再刷新数据')
      }
      await Promise.all([fetchPreview(currentDatasetId.value), loadSchema(currentDatasetId.value)])
    },
    { successMessage: '数据刷新成功', errorMessage: '数据刷新失败' }
  )
}

const handleTabChange = async (name: string | number) => {
  const targetTab = String(name)
  await runWithState(
    'tab_switch',
    async () => {
      activeTab.value = targetTab

      if (targetTab === 'batch' && currentDatasetId.value && !allFields.value.length) {
        await loadSchema(currentDatasetId.value)
      }

      if (targetTab === 'preview' && currentDatasetId.value && !previewData.value.length) {
        await fetchPreview(currentDatasetId.value)
      }
    },
    {
      errorMessage: '切换标签失败',
      notifySuccess: false
    }
  )
}

const handleSelectionChange = (fields: DatasetField[]) => {
  selectedFields.value = fields
}

const submitBatchUpdate = async () => {
  if (!currentDatasetId.value) {
    ElMessage.warning('请先保存数据集')
    return
  }
  if (!selectedFields.value.length) {
    ElMessage.warning('请先选择要批量更新的字段')
    return
  }

  const hasPatch = !!(batchConfig.value.type || batchConfig.value.sortOrder || batchConfig.value.displayNamePrefix.trim())
  if (!hasPatch) {
    ElMessage.warning('请至少配置一个批量更新项')
    return
  }

  await runWithState(
    'batch_update',
    async () => {
      const failures: string[] = []

      for (const field of selectedFields.value) {
        const payload: UpdateFieldRequest = {}
        if (batchConfig.value.type) {
          payload.type = batchConfig.value.type
        }
        if (batchConfig.value.sortOrder) {
          payload.sortOrder = batchConfig.value.sortOrder
        }
        if (batchConfig.value.displayNamePrefix.trim()) {
          payload.displayName = `${batchConfig.value.displayNamePrefix.trim()}${field.displayName || field.name}`
        }

        const response = await datasetApi.updateField(currentDatasetId.value, field.id, payload)
        if (!response.data.success) {
          failures.push(field.name)
        }
      }

      if (failures.length) {
        throw new Error(`以下字段更新失败：${failures.join('、')}`)
      }

      await loadSchema(currentDatasetId.value)
      selectedFields.value = []
    },
    {
      successMessage: '批量更新成功',
      errorMessage: '批量更新失败'
    }
  )
}

const resetGroupingFieldForm = () => {
  groupingFieldForm.value = {
    name: '',
    displayName: '',
    dataType: 'string',
    groupingRule: '',
    groupingEnabled: false
  }
  groupingFieldFormRef.value?.clearValidate()
}

const openGroupingFieldDialog = () => {
  groupingFieldDialogVisible.value = true
}

const submitGroupingField = async () => {
  if (!currentDatasetId.value) {
    ElMessage.warning('请先保存数据集')
    return
  }

  if (!groupingFieldFormRef.value) {
    return
  }

  await groupingFieldFormRef.value.validate()

  await runWithState(
    'create_grouping_field',
    async () => {
      const payload: CreateGroupingFieldRequest = {
        name: groupingFieldForm.value.name.trim(),
        displayName: groupingFieldForm.value.displayName.trim() || undefined,
        type: 'dimension',
        dataType: groupingFieldForm.value.dataType,
        isGroupingField: true,
        groupingRule: groupingFieldForm.value.groupingRule.trim(),
        groupingEnabled: groupingFieldForm.value.groupingEnabled
      }

      const response = await datasetApi.createField(currentDatasetId.value, payload)
      if (!response.data.success) {
        throw new Error(response.data.message || '创建分组字段失败')
      }

      await loadSchema(currentDatasetId.value)
      groupingFieldDialogVisible.value = false
      resetGroupingFieldForm()
    },
    {
      successMessage: '分组字段创建成功',
      errorMessage: '创建分组字段失败'
    }
  )
}

const closePreview = () => {
  previewData.value = []
}

const goBack = () => {
  router.push('/dataset')
}

watch(
  () => [joinConfig.value.leftSourceKey, joinConfig.value.rightSourceKey],
  async () => {
    if (formData.value.type !== 'sql') {
      return
    }
    await Promise.all([
      loadJoinFieldsBySource(joinConfig.value.leftSourceKey, 'left'),
      loadJoinFieldsBySource(joinConfig.value.rightSourceKey, 'right')
    ])
  }
)

watch(
  () => joinConfig.value,
  () => {
    if (formData.value.type !== 'sql') {
      return
    }
    syncJoinToConfig()
    syncAutoQueryFromSources()
  },
  { deep: true }
)

watch(
  selectedSources,
  () => {
    if (formData.value.type !== 'sql') {
      return
    }
    syncSourcesToConfig()
    applyDefaultJoinSources()
    syncJoinToConfig()
    syncAutoQueryFromSources()
  },
  { deep: true }
)

onMounted(async () => {
  await loadDatasources()

  if (!currentDatasetId.value) {
    return
  }

  try {
    await runWithState(
      'init',
      async () => {
        await Promise.all([loadDataset(currentDatasetId.value), loadSchema(currentDatasetId.value)])

        if (formData.value.type === 'sql' && formData.value.datasourceId) {
          await loadTables(formData.value.datasourceId)
          await restoreSqlWorkspaceFromConfig()
        }
      },
      {
        notifySuccess: false,
        errorMessage: '初始化页面失败'
      }
    )
  } catch {
    // 初始化失败后保留现有表单上下文，避免用户输入丢失。
  }
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

.top-action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.title-wrap h2 {
  margin: 0;
  font-size: 22px;
  line-height: 1.2;
  color: #1f2937;
}

.title-sub {
  font-size: 12px;
  color: #6b7280;
}

.action-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
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

.workflow-sidebar {
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: linear-gradient(180deg, #ffffff 0%, #f9fafb 100%);
}

.workflow-main {
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #ffffff;
}

.table-list-shell {
  width: 100%;
  border: 1px solid #e4e7ed;
  border-radius: 10px;
  background: #fff;
  max-height: 360px;
  overflow-y: auto;
  padding: 8px;
}

.table-list-hint {
  font-size: 12px;
  color: #909399;
  margin: 0 0 8px;
}

.table-list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 8px;
  cursor: pointer;
  color: #303133;
  transition: all 0.2s ease;
}

.table-list-item:hover {
  background: #f5f7fa;
}

.table-list-item.is-active {
  background: #ecf5ff;
  color: #1d4ed8;
}

.source-workspace {
  margin-bottom: 16px;
  padding: 14px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: linear-gradient(160deg, #ffffff 0%, #f8fafc 100%);
}

.source-workspace-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
  color: #111827;
  font-weight: 600;
}

.source-workspace-tip {
  font-size: 12px;
  color: #6b7280;
  font-weight: 400;
}

.source-card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 10px;
  margin-bottom: 12px;
}

.source-card {
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 10px 12px;
  cursor: pointer;
  background: #ffffff;
  transition: all 0.2s ease;
}

.source-card:hover {
  border-color: #93c5fd;
  box-shadow: 0 6px 18px rgba(59, 130, 246, 0.12);
}

.source-card.is-active {
  border-color: #3b82f6;
  background: #eff6ff;
}

.source-card-title {
  font-weight: 600;
  color: #1f2937;
}

.source-card-desc {
  margin-top: 2px;
  font-size: 12px;
  color: #6b7280;
}

.join-editor,
.custom-sql-editor {
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 12px;
  background: #fff;
}

.custom-sql-editor {
  margin-top: 10px;
}

.join-editor-title {
  margin-bottom: 10px;
  font-size: 13px;
  font-weight: 600;
  color: #1f2937;
}

.join-editor-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.tab-content {
  min-height: 520px;
}

.tab-actions {
  display: flex;
  gap: 10px;
  margin-bottom: 14px;
}

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

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
}

.no-datasource-tip {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #909399;
}

.grouping-field-dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

@media (max-width: 1280px) {
  .workflow-layout {
    grid-template-columns: 1fr;
  }

  .batch-actions {
    grid-template-columns: 1fr;
  }

  .top-action-bar {
    flex-direction: column;
    align-items: flex-start;
  }

  .join-editor-grid {
    grid-template-columns: 1fr;
  }
}
</style>
