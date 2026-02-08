<template>
  <div class="dataset-edit">
    <el-card>
      <template #header>
        <div class="card-header">
          <h2>{{ isEdit ? '编辑数据集' : '新建数据集' }}</h2>
          <el-button @click="goBack">返回</el-button>
        </div>
      </template>
 
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
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
          <el-select v-model="formData.datasourceId" placeholder="请选择数据源" :loading="datasourcesLoading" :no-data-text="'暂无数据源，请先创建数据源'">
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
 
        <el-form-item v-if="formData.type === 'sql'" label="SQL 查询" prop="config.query">
          <el-input
            v-model="formData.config.query"
            type="textarea"
            :rows="6"
            placeholder="请输入 SQL 查询语句"
          />
          <div class="form-tip">示例: SELECT * FROM users WHERE status = 1</div>
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
            <template #tip>
              支持 Excel、CSV 格式
            </template>
          </el-upload>
        </el-form-item>
 
        <el-form-item>
          <el-button type="primary" @click="previewDataset" :loading="previewing">预览数据</el-button>
          <el-button @click="saveDataset" :loading="saving">保存</el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
    
    <!-- 数据预览图表 -->
    <el-card v-if="datasetId" class="preview-card">
      <template #header>
        <div class="card-header">
          <h2>数据预览</h2>
          <el-button @click="closePreview">关闭</el-button>
        </div>
      </template>
      <DatasetPreview :dataset-id="datasetId" :data="previewData" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { datasourceApi, type DataSource } from '@/api/datasource'
import { datasetApi, type CreateDatasetRequest, type Dataset } from '@/api/dataset'

const route = useRoute()
const router = useRouter()

const isEdit = computed(() => !!route.params.id)
const datasetId = computed(() => route.params.id as string)

const formRef = ref()
const formData = ref<CreateDatasetRequest>({
  name: '',
  type: 'sql',
  datasourceId: '',
  config: {}
})

const formRules = {
  name: [{ required: true, message: '请输入数据集名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择数据集类型', trigger: 'change' }],
  datasourceId: [{ required: true, message: '请选择数据源', trigger: 'change' }]
}

const datasources = ref<DataSource[]>([])
const fileList = ref<any[]>([])
const previewData = ref<Record<string, any>[]>([])
const previewColumns = ref<string[]>([])
const previewing = ref(false)
const saving = ref(false)

const loadDatasources = async () => {
  try {
    const response = await datasourceApi.list()
    if (response.success) {
      datasources.value = response.result || []
    }
  } catch (error) {
    ElMessage.error('加载数据源失败')
  }
}

const loadDataset = async () => {
  try {
    const response = await datasetApi.get(datasetId.value)
    if (response.success) {
      const dataset = response.result
      formData.value = {
        name: dataset.name,
        type: dataset.type,
        datasourceId: dataset.datasourceId,
        config: typeof dataset.config === 'string' ? JSON.parse(dataset.config) : dataset.config
      }
    }
  } catch (error) {
    ElMessage.error('加载数据集失败')
  }
}

const handleTypeChange = () => {
  formData.value.config = {}
  fileList.value = []
  previewData.value = []
  previewColumns.value = []
}

const handleFileChange = (file: any) => {
  fileList.value = [file]
  formData.value.config.file = file.raw
}

const previewDataset = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    previewing.value = true
    let response

    if (isEdit.value) {
      const updateData = {
        name: formData.value.name,
        config: formData.value.config
      }
      await datasetApi.update(datasetId.value, updateData)
      response = await datasetApi.preview(datasetId.value)
    } else {
      const createData = {
        name: formData.value.name,
        type: formData.value.type,
        datasourceId: formData.value.datasourceId,
        config: formData.value.config
      }
      const createResponse = await datasetApi.create(createData)
      if (createResponse.success) {
        response = await datasetApi.preview(createResponse.result.id)
      }
    }

    if (response.success) {
      previewData.value = response.result
      if (response.result.length > 0) {
        previewColumns.value = Object.keys(response.result[0])
      }
      ElMessage.success('预览成功')
    }
  } catch (error) {
    ElMessage.error('预览失败')
  } finally {
    previewing.value = false
  }
}

const saveDataset = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    saving.value = true
    let response

    if (isEdit.value) {
      const updateData = {
        name: formData.value.name,
        config: formData.value.config
      }
      response = await datasetApi.update(datasetId.value, updateData)
    } else {
      response = await datasetApi.create(formData.value)
    }

    if (response.success) {
      ElMessage.success('保存成功')
      goBack()
    } else {
      ElMessage.error(response.message || '保存失败')
    }
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const goBack = () => {
  router.push('/dataset')
}

onMounted(async () => {
  await loadDatasources()
  if (isEdit.value) {
    await loadDataset()
  }
})
</script>

<style scoped>
.dataset-edit {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h2 {
  margin: 0;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.preview-section {
  margin-top: 30px;
}

.preview-section h3 {
  margin-bottom: 15px;
}
</style>
