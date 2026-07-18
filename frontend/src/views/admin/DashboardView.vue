<template>
  <div class="dashboard">
    <h2>仪表盘</h2>
    <el-row :gutter="20" v-loading="loading">
      <el-col :span="6" v-for="card in cards" :key="card.label">
        <el-card class="stat-card">
          <div class="stat-value">{{ card.value }}</div>
          <div class="stat-label">{{ card.label }}</div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { adminApi } from '@/api/admin'

interface DashboardStats {
  user_count: number
  sheet_count: number
  order_count: number
  today_pv: number
  pending_messages: number
}

const loading = ref(true)
const stats = ref<DashboardStats | null>(null)

const cards = computed(() => [
  { label: '用户总数', value: stats.value?.user_count ?? 0 },
  { label: '乐谱总数', value: stats.value?.sheet_count ?? 0 },
  { label: '订单总数', value: stats.value?.order_count ?? 0 },
  { label: '今日访问', value: stats.value?.today_pv ?? 0 },
  { label: '未读消息', value: stats.value?.pending_messages ?? 0 },
])

onMounted(async () => {
  try {
    const res = await adminApi.getStats()
    if (res.result) {
      stats.value = res.result
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}
.dashboard h2 {
  margin: 0 0 20px 0;
}
.stat-card {
  text-align: center;
  margin-bottom: 20px;
}
.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #409eff;
}
.stat-label {
  margin-top: 8px;
  color: #909399;
}
</style>
