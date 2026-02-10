<template>
  <div class="dataset-list">
    <el-card>
      <template #header>
        <div class="card-header">
          <h2>数据集管理</h2>
          <el-button type="primary" @click="createDataset">新建数据集</el-button>
        </div>
      </template>

      <el-table :data="datasets" v-loading="loading" stripe>
        <el-table-column prop="name" label="名称" min-width="200" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTagType(row.type)">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="datasourceName" label="数据源" width="150" />
        <el-table-column prop="fieldCount" label="字段数" width="100" />
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="editDataset(row)">编辑</el-button>
            <el-button size="small" type="success" @click="previewDataset(row)">预览</el-button>
            <el-button size="small" type="danger" @click="deleteDataset(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
    
    <!-- 数据预览弹窗 -->
    <el-dialog
      v-model="previewVisible"
      title="数据集预览"
      width="80%"
      :close-on-click-modal="false"
    >
      <DatasetPreview v-if="previewDatasetId" :dataset-id="previewDatasetId" :data="previewData" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { datasetApi, type Dataset } from '@/api/dataset'
import DatasetPreview from '@/components/dataset/DatasetPreview.vue'

const router = useRouter()

const datasets = ref<Dataset[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const previewVisible = ref(false)
const previewDatasetId = ref('')
const previewData = ref<any[]>([])

const loadDatasets = async () => {
  loading.value = true
  try {
    const response = await datasetApi.list(currentPage.value, pageSize.value)
    const { data } = response
    if (data.success) {
      datasets.value = data.result || []
      total.value = data.total || 0
    } else {
      ElMessage.error(data.message || '加载数据集失败')
    }
  } catch (error) {
    ElMessage.error('加载数据集失败')
  } finally {
    loading.value = false
  }
}

const createDataset = () => {
  router.push('/dataset/create')
}

const editDataset = (dataset: Dataset) => {
  router.push(`/dataset/edit/${dataset.id}`)
}

const deleteDataset = async (dataset: Dataset) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除数据集 "${dataset.name}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const response = await datasetApi.delete(dataset.id)
    const { data } = response
    if (data.success) {
      ElMessage.success('删除成功')
      await loadDatasets()
    } else {
      ElMessage.error(data.message || '删除失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  loadDatasets()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadDatasets()
}

const previewDataset = async (dataset: Dataset) => {
  previewDatasetId.value = dataset.id
  previewData.value = []
  previewVisible.value = true
  
  // 模拟预览数据加载 - 实际应该从后端 API 获取
  try {
    const response = await datasetApi.preview(dataset.id)
    const { data } = response
    if (data.success) {
      previewData.value = data.result || []
    } else {
      ElMessage.warning(data.message || '预览数据加载失败，显示模拟数据')
      previewData.value = [
        { id: 1, name: '测试数据 1', value: 1000 },
        { id: 2, name: '测试数据 2', value: 2000 },
        { id: 3, name: '测试数据 3', value: 1500 }
      ]
    }
  } catch (error) {
    ElMessage.warning('预览失败，显示模拟数据')
    previewData.value = [
      { id: 1, name: '测试数据 1', value: 1000 },
      { id: 2, name: '测试数据 2', value: 2000 },
      { id: 3, name: '测试数据 3', value: 1500 }
    ]
  }
}

const getTypeLabel = (type: string) => {
  const typeMap: Record<string, string> = {
    sql: 'SQL',
    api: 'API',
    file: '文件'
  }
  return typeMap[type] || type
}

const getTypeTagType = (type: string) => {
  const typeMap: Record<string, any> = {
    sql: 'success',
    api: 'warning',
    file: 'info'
  }
  return typeMap[type] || ''
}

onMounted(() => {
  loadDatasets()
})
</script>

<style scoped>
.dataset-list {
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

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
