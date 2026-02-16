<template>
  <aside class="workflow-sidebar">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px" status-icon>
      <el-form-item label="数据集名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入数据集名称" />
      </el-form-item>

      <el-form-item label="数据集类型" prop="type">
        <el-radio-group v-model="formData.type" @change="$emit('type-change')">
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
          @change="(val: string) => $emit('datasource-change', val)"
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
            @click="$emit('update:activeSidebarTableKey', item.key)"
            @dblclick="$emit('sidebar-item-dblclick', item)"
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
          :on-change="(file: UploadUserFile) => $emit('file-change', file)"
          :file-list="fileList"
        >
          <el-button>选择文件</el-button>
          <template #tip>支持 Excel、CSV 格式</template>
        </el-upload>
      </el-form-item>
    </el-form>
  </aside>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { InfoFilled } from '@element-plus/icons-vue'
import type { FormInstance, FormRules, UploadUserFile } from 'element-plus'
import type { DataSource } from '@/api/datasource'
import type { CreateDatasetRequest } from '@/api/dataset'
import type { SidebarTableItem } from '@/views/dataset/composables/useDatasetEdit'

const formRef = ref<FormInstance>()

defineProps<{
  formData: CreateDatasetRequest
  formRules: FormRules
  datasources: DataSource[]
  datasourcesLoading: boolean
  tablesLoading: boolean
  sidebarTableItems: SidebarTableItem[]
  activeSidebarTableKey: string
  fileList: UploadUserFile[]
}>()

defineEmits<{
  'update:activeSidebarTableKey': [value: string]
  'type-change': []
  'datasource-change': [value: string]
  'file-change': [file: UploadUserFile]
  'sidebar-item-dblclick': [item: SidebarTableItem]
}>()

defineExpose({
  validate: () => formRef.value?.validate(),
  clearValidate: () => formRef.value?.clearValidate()
})
</script>

<style scoped>
.workflow-sidebar {
  padding: 16px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: linear-gradient(180deg, #ffffff 0%, #f9fafb 100%);
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

.no-datasource-tip {
  margin-top: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #909399;
}
</style>
