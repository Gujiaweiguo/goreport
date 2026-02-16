<template>
  <div class="top-action-bar">
    <div class="title-wrap">
      <h2>{{ currentDatasetId ? '编辑数据集' : '新建数据集' }}</h2>
      <span class="title-sub">数据集编辑工作流</span>
    </div>
    <div class="action-buttons">
      <el-button
        type="primary"
        :loading="isLoading && operation === 'save'"
        @click="$emit('save')"
      >
        保存
      </el-button>
      <el-button
        type="success"
        :loading="isLoading && operation === 'save_and_return'"
        @click="$emit('save-and-return')"
      >
        保存并返回
      </el-button>
      <el-button
        :loading="isLoading && operation === 'refresh'"
        @click="$emit('refresh')"
      >
        刷新数据
      </el-button>
      <el-button @click="$emit('back')">返回</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  currentDatasetId: string
  isLoading: boolean
  operation: string
}>()

defineEmits<{
  save: []
  'save-and-return': []
  refresh: []
  back: []
}>()
</script>

<style scoped>
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

@media (max-width: 1280px) {
  .top-action-bar {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
