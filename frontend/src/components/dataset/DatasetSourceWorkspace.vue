<template>
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
        @click="$emit('update:activeSourceKey', source.key)"
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
</template>

<script setup lang="ts">
import type { CreateDatasetRequest } from '@/api/dataset'
import type { SelectedSource, JoinConfigState, SourceType } from '@/views/dataset/composables/useDatasetEdit'

defineProps<{
  formData: CreateDatasetRequest
  selectedSources: SelectedSource[]
  activeSourceKey: string
  joinConfig: JoinConfigState
  joinFieldLoading: { left: boolean; right: boolean }
  leftJoinFieldOptions: string[]
  rightJoinFieldOptions: string[]
  showSqlEditor: boolean
}>()

defineEmits<{
  'update:activeSourceKey': [value: string]
}>()
</script>

<style scoped>
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

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: #909399;
}

@media (max-width: 1280px) {
  .join-editor-grid {
    grid-template-columns: 1fr;
  }
}
</style>
