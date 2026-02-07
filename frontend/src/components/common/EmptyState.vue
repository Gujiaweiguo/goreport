<template>
  <div class="empty-state">
    <div class="empty-icon">
      <el-icon :size="iconSize">
        <component :is="iconComponent" />
      </el-icon>
    </div>
    <p class="empty-text">{{ text }}</p>
    <p v-if="description" class="empty-description">{{ description }}</p>
    <div v-if="action" class="empty-action">
      <el-button :type="actionType" :icon="actionIcon" @click="handleAction">
        {{ actionText }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Loading, Warning, Lock, InfoFilled, CircleClose } from '@element-plus/icons-vue'

interface Props {
  icon?: any
  iconSize?: number
  text?: string
  description?: string
  action?: () => void
  actionText?: string
  actionType?: 'primary' | 'default' | 'success' | 'warning' | 'danger'
  actionIcon?: any
}

const props = withDefaults(defineProps<Props>(), {
  iconSize: 64,
  text: '暂无数据',
  actionType: 'primary'
})

const iconComponent = computed(() => props.icon || InfoFilled)

function handleAction() {
  if (props.action) {
    props.action()
  }
}
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #909399;
  gap: 16px;
}

.empty-icon {
  color: #d1d5db;
  opacity: 0.6;
}

.empty-text {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  color: #64748b;
}

.empty-description {
  margin: 0;
  font-size: 14px;
  color: #909399;
  max-width: 400px;
  text-align: center;
  line-height: 1.6;
}

.empty-action {
  margin-top: 12px;
}
</style>
