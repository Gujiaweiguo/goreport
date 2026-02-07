<template>
  <div class="home">
    <h1>goReport</h1>
    <p>欢迎使用 goReport 报表系统</p>
    <el-button type="primary" @click="testConnection">测试连接</el-button>
    <p v-if="healthStatus">
      状态: {{ healthStatus }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'

const healthStatus = ref<string>('')

async function testConnection() {
  try {
    const response = await axios.get('/health')
    healthStatus.value = JSON.stringify(response.data)
  } catch (error) {
    healthStatus.value = '连接失败'
  }
}
</script>

<style scoped>
.home {
  padding: 40px;
  text-align: center;
}

h1 {
  font-size: 32px;
  margin-bottom: 20px;
}

p {
  margin: 10px 0;
}
</style>
