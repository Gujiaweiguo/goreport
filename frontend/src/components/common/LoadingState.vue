<template>
  <div v-if="loading" class="loading-overlay">
    <div class="loading-content">
      <el-icon :size="size" class="is-loading">
        <Loading />
      </el-icon>
      <p v-if="text" class="loading-text">{{ text }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Loading } from '@element-plus/icons-vue'

interface Props {
  loading: boolean
  text?: string
  size?: number
  fullscreen?: boolean
}

withDefaults(defineProps<Props>(), {
  size: 48,
  fullscreen: false
})
</script>

<style scoped>
.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  backdrop-filter: blur(2px);
}

.loading-overlay.fullscreen {
  position: fixed;
}

.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.loading-text {
  margin: 0;
  font-size: 14px;
  color: #64748b;
}

.is-loading {
  animation: rotating 2s linear infinite;
}

@keyframes rotating {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
