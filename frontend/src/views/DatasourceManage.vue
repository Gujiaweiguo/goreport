<template>
  <div class="datasource-manage">
    <h2>数据源管理</h2>
    <div class="page-actions">
      <el-button type="primary" @click="showCreateDialog">
        <el-icon><Plus /></el-icon>
        创建数据源
      </el-button>
      <el-tooltip content="批量导入功能即将上线，敬请期待" placement="top">
        <el-button type="primary" plain disabled>
          <span>批量导入</span>
        </el-button>
      </el-tooltip>
    </div>

    <el-input
      v-model="searchKeyword"
      placeholder="搜索数据源名称、类型或连接地址"
      class="search-input"
      @input="handleSearch"
      clearable
    >
      <template #prefix>
        <el-icon><Search /></el-icon>
      </template>
    </el-input>

    <el-table
      :data="datasources"
      @selection-change="handleSelectionChange"
      v-loading="loading"
      stripe
      border
      @row-click="handleRowClick"
      style="width: 100%"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="name" label="名称" width="180" />
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.type === 'mysql'" type="success">MySQL</el-tag>
          <el-tag v-else-if="row.type === 'postgresql' || row.type === 'postgres'" type="warning">PostgreSQL</el-tag>
          <el-tag v-else-if="row.type === 'sqlserver'" type="info">SQL Server</el-tag>
          <el-tag v-else-if="row.type === 'excel'" type="info">Excel</el-tag>
          <el-tag v-else-if="row.type === 'csv'" type="info">CSV</el-tag>
          <el-tag v-else-if="row.type === 'api'" type="warning">API 数据源</el-tag>
          <el-tag v-else>Unknown</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="连接地址" width="250">
        <template #default="{ row }">
          {{ row.host }}:{{ row.port }}
        </template>
      </el-table-column>
      <el-table-column prop="database" label="数据库" width="150" />
      <el-table-column prop="username" label="用户名" width="120" />
      <el-table-column prop="updatedAt" label="更新时间" width="130">
        <template #default="{ row }">
          <span>{{ formatDate(row.updatedAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="operations" label="操作" width="350">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click.stop="showEditDialog(row)">编辑</el-button>
          <el-button link type="success" size="small" @click.stop="handleCopy(row)">复制</el-button>
          <el-button link type="warning" size="small" @click.stop="showMoveDialog(row)">移动</el-button>
          <el-button link type="info" size="small" @click.stop="showRenameDialog(row)">重命名</el-button>
          <el-button link type="primary" size="small" @click.stop="openTestDialog(row)">测试连接</el-button>
          <el-button link type="danger" size="small" @click.stop="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      class="pagination-bg"
      layout="total, prev, pager, next"
      :total="total"
      @current-change="handlePageChange"
      @size-change="handleSizeChange"
    />
  </div>

  <el-dialog
    v-model="dialogVisible"
    :title="dialogTitle"
    width="880px"
    @close="handleDialogClose"
  >
    <el-steps v-if="!isEdit" :active="createStep" simple class="create-steps">
      <el-step title="选择数据源" />
      <el-step title="配置信息" />
    </el-steps>

    <div v-if="!isEdit && createStep === 1" class="source-selector">
      <div class="source-categories">
        <button
          v-for="category in sourceCategories"
          :key="category.key"
          type="button"
          class="category-item"
          :class="{ active: selectedCategoryKey === category.key }"
          @click="selectedCategoryKey = category.key"
        >
          {{ category.label }}
        </button>
      </div>
      <div class="source-cards">
        <div class="source-cards-title">{{ activeCategoryLabel }}</div>
        <div class="source-card-grid">
          <button
            v-for="item in activeTemplates"
            :key="item.key"
            type="button"
            class="source-card"
            :class="{
              active: selectedTemplateKey === item.key,
              unsupported: !item.supported
            }"
            @click="selectTemplate(item)"
          >
            <span>{{ item.label }}</span>
            <el-tag size="small" :type="item.supported ? 'success' : 'info'">
              {{ item.supported ? '已支持' : '即将支持' }}
            </el-tag>
          </button>
        </div>
      </div>
    </div>

    <el-form
      v-else
      ref="formRef"
      :model="form"
      :rules="formRules"
      label-width="100px"
    >
      <el-form-item label="名称" prop="name">
        <el-input v-model="form.name" placeholder="请输入数据源名称" />
      </el-form-item>

      <el-form-item label="类型" prop="type">
        <el-tag v-if="!isEdit">{{ selectedTypeLabel }}</el-tag>
        <el-select v-else v-model="form.type" placeholder="请选择数据源类型" @change="handleTypeChange">
          <el-option label="MySQL" value="mysql" />
          <el-option label="PostgreSQL" value="postgres" />
          <el-option label="SQL Server" value="sqlserver" />
          <el-option label="MongoDB" value="mongodb" />
          <el-option label="Excel" value="excel" />
          <el-option label="CSV" value="csv" />
          <el-option label="API" value="api" />
        </el-select>
      </el-form-item>

      <el-form-item label="主机" prop="host">
        <el-input v-model="form.host" placeholder="请输入主机地址" />
      </el-form-item>

      <el-form-item label="端口" prop="port">
        <el-input-number v-model="form.port" :min="1" :max="65535" />
      </el-form-item>

      <el-form-item label="数据库" prop="database" v-if="needsDatabase">
        <el-input v-model="form.database" placeholder="请输入数据库名称" />
      </el-form-item>

      <el-form-item label="用户名" prop="username" v-if="needsAuth">
        <el-input v-model="form.username" placeholder="请输入用户名" />
      </el-form-item>

      <el-form-item label="密码" prop="password" v-if="needsAuth">
        <el-input
          v-model="form.password"
          type="password"
          :placeholder="isEdit ? '留空则保持原密码不变' : '请输入密码'"
          @focus="handlePasswordFocus"
          show-password
        />
      </el-form-item>

      <el-divider>高级配置</el-divider>

      <el-form-item label="SSH 隧道" v-if="supportsSSH">
        <el-checkbox v-model="enableSSH">启用 SSH 隧道</el-checkbox>
      </el-form-item>

      <template v-if="enableSSH">
        <el-form-item label="SSH 主机">
          <el-input v-model="form.advanced.sshHost" placeholder="请输入 SSH 主机地址" />
        </el-form-item>
        <el-form-item label="SSH 端口">
          <el-input-number v-model="form.advanced.sshPort" :min="1" :max="65535" />
        </el-form-item>
        <el-form-item label="SSH 用户名">
          <el-input v-model="form.advanced.sshUsername" placeholder="请输入 SSH 用户名" />
        </el-form-item>
        <el-form-item label="SSH 认证方式">
          <el-radio-group v-model="sshAuthType">
            <el-radio value="password">密码认证</el-radio>
            <el-radio value="key">密钥认证</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="SSH 密码" v-if="sshAuthType === 'password'">
          <el-input
            v-model="form.advanced.sshPassword"
            type="password"
            placeholder="请输入 SSH 密码"
            show-password
          />
        </el-form-item>
        <el-form-item label="SSH 密钥" v-if="sshAuthType === 'key'">
          <el-input
            v-model="form.advanced.sshKey"
            type="textarea"
            :rows="4"
            placeholder="请输入 SSH 私钥"
          />
        </el-form-item>
        <el-form-item label="SSH 密钥密码" v-if="sshAuthType === 'key'">
          <el-input
            v-model="form.advanced.sshKeyPhrase"
            type="password"
            placeholder="请输入密钥密码（可选）"
            show-password
          />
        </el-form-item>
      </template>

      <el-form-item label="最大连接数">
        <el-input-number v-model="form.advanced.maxConnections" :min="1" :max="100" placeholder="1-100" />
      </el-form-item>

      <el-form-item label="查询超时（秒）">
        <el-input-number v-model="form.advanced.queryTimeoutSeconds" :min="5" :max="300" placeholder="5-300" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="handleDialogClose">取消</el-button>
      <template v-if="!isEdit && createStep === 1">
        <el-button type="primary" :disabled="!canGoNext" @click="goNextStep">下一步</el-button>
      </template>
      <template v-else>
        <el-button v-if="!isEdit" @click="createStep = 1">上一步</el-button>
        <el-button @click="testConnectionInDialog" :loading="testingInDialog">测试连接</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </template>
  </el-dialog>

  <el-dialog
    v-model="moveDialogVisible"
    title="移动数据源"
    width="400px"
  >
    <el-form :model="moveForm" label-width="80px">
      <el-form-item label="目标位置">
        <el-input v-model="moveForm.target" placeholder="请输入目标位置" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="moveDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleMoveSubmit">移动</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="renameDialogVisible"
    title="重命名数据源"
    width="400px"
  >
    <el-form :model="renameForm" label-width="80px">
      <el-form-item label="新名称">
        <el-input v-model="renameForm.name" placeholder="请输入新名称" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="renameDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="handleRenameSubmit">确定</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="testDialogVisible"
    title="测试连接"
    width="400px"
  >
    <div v-loading="testing">
      <p v-if="!testResult">点击下方按钮测试连接...</p>
      <div v-else>
        <el-result
          :icon="testResult.success ? 'success' : 'error'"
          :title="testResult.success ? '连接成功' : '连接失败'"
          :sub-title="testResult.message"
        />
      </div>
    </div>
    <template #footer>
      <el-button @click="testDialogVisible = false">关闭</el-button>
      <el-button type="primary" :loading="testing" @click="runTest">
        测试
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { datasourceApi, type DataSource } from '@/api/datasource'
import { getErrorMessage } from '@/utils/errorHandling'

const dialogVisible = ref(false)
const testDialogVisible = ref(false)
const moveDialogVisible = ref(false)
const renameDialogVisible = ref(false)
const dialogTitle = ref('创建数据源')
const loading = ref(false)
const submitting = ref(false)
const testing = ref(false)
const isEdit = ref(false)
const currentEditId = ref('')
const formRef = ref<FormInstance>()

const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')
const datasources = ref<DataSource[]>([])
const testResult = ref<{ success: boolean; message: string } | null>(null)
const currentTestDatasourceId = ref('')
const testingInDialog = ref(false)
const maskedPassword = '********'
const createStep = ref(1)
const selectedCategoryKey = ref('oltp')
const selectedTemplateKey = ref('mysql')

const sshAuthType = ref<'password' | 'key'>('password')
const enableSSH = ref(false)

const moveForm = reactive({
  id: '',
  target: ''
})

const renameForm = reactive({
  id: '',
  name: ''
})

interface SourceCategory {
  key: string
  label: string
}

interface SourceTemplate {
  key: string
  label: string
  supported: boolean
  backendType?: string
  defaultPort?: number
}

const sourceCategories: SourceCategory[] = [
  { key: 'oltp', label: 'OLTP' },
  { key: 'olap', label: 'OLAP' },
  { key: 'lake', label: '数据湖' },
  { key: 'api', label: 'API数据' },
  { key: 'file', label: '文件' }
]

const sourceTemplateMap: Record<string, SourceTemplate[]> = {
  oltp: [
    { key: 'db2', label: 'Db2', supported: false },
    { key: 'mysql', label: 'MySQL', supported: true, backendType: 'mysql', defaultPort: 3306 },
    { key: 'mariadb', label: 'MariaDB', supported: false },
    { key: 'mongodb-bi', label: 'Mongodb-BI', supported: true, backendType: 'mongodb', defaultPort: 27017 },
    { key: 'oracle', label: 'Oracle', supported: false },
    { key: 'postgres', label: 'PostgreSQL', supported: true, backendType: 'postgres', defaultPort: 5432 },
    { key: 'sqlserver', label: 'SQL Server', supported: true, backendType: 'sqlserver', defaultPort: 1433 },
    { key: 'tidb', label: 'TiDB', supported: false }
  ],
  olap: [
    { key: 'impala', label: 'Apache Impala', supported: false },
    { key: 'doris', label: 'Apache Doris', supported: false },
    { key: 'clickhouse', label: 'ClickHouse', supported: false },
    { key: 'elasticsearch', label: 'Elasticsearch', supported: false },
    { key: 'starrocks', label: 'StarRocks', supported: false }
  ],
  lake: [
    { key: 'redshift', label: 'AWS Redshift', supported: false }
  ],
  api: [
    { key: 'api', label: 'API', supported: true, backendType: 'api' }
  ],
  file: [
    { key: 'excel-local', label: '本地 Excel/CSV', supported: true, backendType: 'excel' },
    { key: 'excel-remote', label: '远程 Excel/CSV', supported: true, backendType: 'csv' }
  ]
}

const datasourceTypeLabelMap: Record<string, string> = {
  mysql: 'MySQL',
  postgres: 'PostgreSQL',
  sqlserver: 'SQL Server',
  mongodb: 'MongoDB',
  excel: 'Excel',
  csv: 'CSV',
  api: 'API'
}

interface DataSourceForm {
  name: string
  type: string
  host: string
  port: number
  database: string
  username: string
  password: string
  advanced: {
    sshHost: string
    sshPort: number
    sshUsername: string
    sshPassword: string
    sshKey: string
    sshKeyPhrase: string
    maxConnections: number
    queryTimeoutSeconds: number
  }
}

const form = reactive<DataSourceForm>({
  name: '',
  type: 'mysql',
  host: 'localhost',
  port: 3306,
  database: '',
  username: '',
  password: '',
  advanced: {
    sshHost: '',
    sshPort: 22,
    sshUsername: '',
    sshPassword: '',
    sshKey: '',
    sshKeyPhrase: '',
    maxConnections: 10,
    queryTimeoutSeconds: 30
  }
})

const formRules = computed<FormRules<DataSourceForm>>(() => ({
  name: [{ required: true, message: '请输入数据源名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择数据源类型', trigger: 'change' }],
  host: [{ required: true, message: '请输入主机地址', trigger: 'blur' }],
  port: [{ required: true, message: '请输入端口号', trigger: 'blur' }],
  database: [{ required: true, message: '请输入数据库名称', trigger: 'blur' }],
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: !isEdit.value, message: '请输入密码', trigger: 'blur' }]
}))

const activeTemplates = computed(() => {
  return sourceTemplateMap[selectedCategoryKey.value] || []
})

const activeCategoryLabel = computed(() => {
  return sourceCategories.find(category => category.key === selectedCategoryKey.value)?.label || ''
})

const selectedTemplate = computed(() => {
  return activeTemplates.value.find(item => item.key === selectedTemplateKey.value) || null
})

const canGoNext = computed(() => {
  return Boolean(selectedTemplate.value?.supported && selectedTemplate.value?.backendType)
})

const selectedTypeLabel = computed(() => {
  return datasourceTypeLabelMap[form.type] || form.type
})

const needsDatabase = computed(() => {
  return ['mysql', 'postgres', 'sqlserver', 'mongodb'].includes(form.type)
})

const needsAuth = computed(() => {
  return ['mysql', 'postgres', 'sqlserver', 'mongodb'].includes(form.type)
})

const supportsSSH = computed(() => {
  return ['mysql', 'postgres', 'mongodb'].includes(form.type)
})

const formatDate = (date?: string) => {
  if (!date) return '-'
  const d = new Date(date)
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function loadDatasources() {
  loading.value = true
  try {
    const response = await datasourceApi.list(currentPage.value, pageSize.value)
    if (response.data.success) {
      datasources.value = response.data.result.datasources || []
      total.value = response.data.result.total || 0
    } else {
      ElMessage.error(response.data.message || '加载数据源失败')
    }
  } catch (error: unknown) {
    ElMessage.error('加载数据源失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  currentPage.value = 1
  loadDatasources()
}

function handlePageChange(page: number) {
  currentPage.value = page
  loadDatasources()
}

function handleSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  loadDatasources()
}

function handleSelectionChange(selection: DataSource[]) {
  // 处理选择变化
}

function handleRowClick(row: DataSource) {
  // 处理行点击
}

function applyTemplate(template: SourceTemplate) {
  if (!template.backendType) return
  form.type = template.backendType
  if (template.defaultPort) {
    form.port = template.defaultPort
  }
  handleTypeChange()
}

function selectTemplate(template: SourceTemplate) {
  selectedTemplateKey.value = template.key
  if (template.supported) {
    applyTemplate(template)
  }
}

function goNextStep() {
  if (!selectedTemplate.value || !selectedTemplate.value.supported || !selectedTemplate.value.backendType) {
    ElMessage.warning('该数据源即将支持，请选择已支持的数据源')
    return
  }
  applyTemplate(selectedTemplate.value)
  createStep.value = 2
}

function showCreateDialog() {
  isEdit.value = false
  dialogTitle.value = '创建数据源'
  resetForm()
  createStep.value = 1
  selectedCategoryKey.value = 'oltp'
  selectedTemplateKey.value = 'mysql'
  applyTemplate(sourceTemplateMap.oltp[1])
  dialogVisible.value = true
}

function showEditDialog(row: DataSource) {
  isEdit.value = true
  dialogTitle.value = '编辑数据源'
  form.name = row.name
  form.type = row.type
  form.host = row.host
  form.port = row.port
  form.database = row.database || ''
  form.username = row.username || ''
  form.password = maskedPassword
  currentEditId.value = row.id
  createStep.value = 2
  dialogVisible.value = true
}

async function showImportDialog() {
  ElMessage.info('批量导入功能开发中...')
}

async function handleCopy(row: DataSource) {
  try {
    const response = await datasourceApi.copy(row.id)
    if (response.data.success) {
      ElMessage.success('复制数据源成功')
      await loadDatasources()
    } else {
      ElMessage.error(response.data.message || '复制数据源失败')
    }
  } catch (error: unknown) {
    ElMessage.error('复制数据源失败')
  }
}

function showMoveDialog(row: DataSource) {
  moveForm.id = row.id
  moveForm.target = ''
  moveDialogVisible.value = true
}

async function handleMoveSubmit() {
  if (!moveForm.target) {
    ElMessage.warning('请输入目标位置')
    return
  }
  try {
    const response = await datasourceApi.move(moveForm.id, moveForm.target)
    if (response.data.success) {
      ElMessage.success('移动数据源成功')
      moveDialogVisible.value = false
      await loadDatasources()
    } else {
      ElMessage.error(response.data.message || '移动数据源失败')
    }
  } catch (error: unknown) {
    ElMessage.error('移动数据源失败')
  }
}

function showRenameDialog(row: DataSource) {
  renameForm.id = row.id
  renameForm.name = row.name
  renameDialogVisible.value = true
}

async function handleRenameSubmit() {
  if (!renameForm.name.trim()) {
    ElMessage.warning('请输入新名称')
    return
  }
  try {
    const response = await datasourceApi.rename(renameForm.id, renameForm.name)
    if (response.data.success) {
      ElMessage.success('重命名成功')
      renameDialogVisible.value = false
      await loadDatasources()
    } else {
      ElMessage.error(response.data.message || '重命名失败')
    }
  } catch (error: unknown) {
    ElMessage.error('重命名失败')
  }
}

async function handleDelete(row: DataSource) {
  try {
    await ElMessageBox.confirm(`确认删除数据源 "${row.name}" 吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    const response = await datasourceApi.delete(row.id)
    if (response.data.success) {
      ElMessage.success('删除数据源成功')
      await loadDatasources()
    } else {
      ElMessage.error(response.data.message || '删除数据源失败')
    }
  } catch (error: unknown) {
    if (error !== 'cancel') {
      ElMessage.error('删除数据源失败')
    }
  }
}

function openTestDialog(row: DataSource) {
  currentTestDatasourceId.value = row.id
  testDialogVisible.value = true
  testResult.value = null
}

function runTest() {
  testing.value = true
  const testData = {
    name: form.name,
    type: form.type,
    host: form.host,
    port: form.port,
    database: form.database,
    username: form.username,
    password: form.password,
    advanced: enableSSH.value ? {
      sshHost: form.advanced.sshHost,
      sshPort: form.advanced.sshPort,
      sshUsername: form.advanced.sshUsername,
      sshPassword: form.advanced.sshPassword,
      sshKey: form.advanced.sshKey,
      sshKeyPhrase: form.advanced.sshKeyPhrase,
      maxConnections: form.advanced.maxConnections,
      queryTimeoutSeconds: form.advanced.queryTimeoutSeconds
    } : undefined
  }

  const testRequest = currentTestDatasourceId.value
    ? datasourceApi.testById(currentTestDatasourceId.value)
    : datasourceApi.test(testData)

  testRequest
    .then(response => {
      testResult.value = {
        success: response.data.success,
        message: response.data.message || (response.data.success ? '连接成功' : '连接失败')
      }
    })
    .catch(error => {
      testResult.value = {
        success: false,
        message: getErrorMessage(error, '连接测试失败')
      }
    })
    .finally(() => {
      testing.value = false
    })
}

function resetForm() {
  form.name = ''
  form.type = 'mysql'
  form.host = 'localhost'
  form.port = 3306
  form.database = ''
  form.username = ''
  form.password = ''
  form.advanced.sshHost = ''
  form.advanced.sshPort = 22
  form.advanced.sshUsername = ''
  form.advanced.sshPassword = ''
  form.advanced.sshKey = ''
  form.advanced.sshKeyPhrase = ''
  form.advanced.maxConnections = 10
  form.advanced.queryTimeoutSeconds = 30
  enableSSH.value = false
  sshAuthType.value = 'password'
}

function handleDialogClose() {
  dialogVisible.value = false
  resetForm()
  createStep.value = 1
  selectedCategoryKey.value = 'oltp'
  selectedTemplateKey.value = 'mysql'
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

function handleTypeChange() {
  if (!needsDatabase.value) {
    form.database = ''
  }
  if (!needsAuth.value) {
    form.username = ''
    form.password = ''
  }
}

function handlePasswordFocus() {
  if (isEdit.value && form.password === maskedPassword) {
    form.password = ''
  }
}

async function testConnectionInDialog() {
  testingInDialog.value = true

  try {
    let response
    const hasRealPassword = form.password.trim() !== '' && form.password !== maskedPassword

    if (isEdit.value && currentEditId.value && !hasRealPassword) {
      response = await datasourceApi.testById(currentEditId.value)
    } else {
      const testData = {
        name: form.name,
        type: form.type,
        host: form.host,
        port: form.port,
        database: form.database,
        username: form.username,
        password: hasRealPassword ? form.password : '',
        advanced: enableSSH.value ? {
          sshHost: form.advanced.sshHost,
          sshPort: form.advanced.sshPort,
          sshUsername: form.advanced.sshUsername,
          sshPassword: form.advanced.sshPassword,
          sshKey: form.advanced.sshKey,
          sshKeyPhrase: form.advanced.sshKeyPhrase,
          maxConnections: form.advanced.maxConnections,
          queryTimeoutSeconds: form.advanced.queryTimeoutSeconds
        } : undefined
      }
      response = await datasourceApi.test(testData)
    }

    if (response.data.success) {
      ElMessage.success('连接测试成功')
    } else {
      ElMessage.error(response.data.message || '连接测试失败')
    }
  } catch (error: unknown) {
    ElMessage.error(getErrorMessage(error, '连接测试失败'))
  } finally {
    testingInDialog.value = false
  }
}

async function handleSubmit() {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  submitting.value = true

  try {
    if (isEdit.value) {
      const updateData: import('@/api/datasource').UpdateDataSourceRequest = {
        name: form.name,
        type: form.type,
        host: form.host,
        port: form.port,
        database: form.database,
        username: form.username
      }
      const hasRealPassword = form.password.trim() !== '' && form.password !== maskedPassword
      if (hasRealPassword) {
        updateData.password = form.password
      }
      const response = await datasourceApi.update(currentEditId.value, updateData)
      if (response.data.success) {
        ElMessage.success('更新数据源成功')
        await loadDatasources()
        handleDialogClose()
      } else {
        ElMessage.error(response.data.message || '更新数据源失败')
      }
    } else {
      const createData: import('@/api/datasource').CreateDataSourceRequest = {
        name: form.name,
        type: form.type,
        host: form.host,
        port: form.port,
        database: form.database,
        username: form.username,
        password: form.password
      }
      const response = await datasourceApi.create(createData)
      if (response.data.success) {
        ElMessage.success('创建数据源成功')
        await loadDatasources()
        handleDialogClose()
      } else {
        ElMessage.error(response.data.message || '创建数据源失败')
      }
    }
  } catch (error: unknown) {
    ElMessage.error('操作失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadDatasources()
})
</script>

<style scoped>
.datasource-manage {
  padding: 20px;
}

.page-actions {
  margin-bottom: 20px;
  display: flex;
  gap: 16px;
}

.search-input {
  margin-bottom: 20px;
}

.table-container {
  min-height: 400px;
  background: #f5f7fa;
  border: 1px solid #e0e0e0e;
  border-radius: 8px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
}

.el-button {
  margin-right: 8px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.el-table {
  font-size: 13px;
}

.el-form {
  max-width: 400px;
}

.create-steps {
  margin-bottom: 16px;
}

.source-selector {
  display: flex;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
  min-height: 360px;
}

.source-categories {
  width: 160px;
  background: #f8fafc;
  border-right: 1px solid #e4e7ed;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.category-item {
  width: 100%;
  border: none;
  background: transparent;
  border-radius: 6px;
  text-align: left;
  padding: 10px 12px;
  cursor: pointer;
  color: #606266;
}

.category-item.active {
  background: #ecf5ff;
  color: #409eff;
  font-weight: 600;
}

.source-cards {
  flex: 1;
  padding: 16px;
}

.source-cards-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
}

.source-card-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.source-card {
  border: 1px solid #dcdfe6;
  border-radius: 8px;
  background: #fff;
  padding: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
  transition: all 0.2s;
}

.source-card:hover {
  border-color: #409eff;
}

.source-card.active {
  border-color: #409eff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.12);
}

.source-card.unsupported {
  opacity: 0.7;
}

.source-card span {
  font-size: 14px;
}

@media (max-width: 900px) {
  .source-selector {
    flex-direction: column;
  }

  .source-categories {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #e4e7ed;
    flex-direction: row;
    overflow-x: auto;
  }

  .category-item {
    flex: 0 0 auto;
    white-space: nowrap;
  }

  .source-card-grid {
    grid-template-columns: 1fr;
  }
}

.el-pagination {
  margin-top: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-controls {
  display: flex;
  gap: 12px;
  align-items: center;
}

.el-empty-text {
  color: #909399;
}

.pagination-bg {
  background-color: #f5f7fa;
}
</style>
