<template>
  <div class="banned-ip">
    <h2>IP封禁管理</h2>
    <el-card>
      <div class="ban-form">
        <el-input v-model="ipAddress" placeholder="输入IP地址" class="ban-input" />
        <el-input v-model="remark" placeholder="备注（可选）" class="ban-input" />
        <el-button type="primary" @click="handleBanIp">封禁IP</el-button>
      </div>

      <el-table :data="ips" v-loading="loading" style="width: 100%">
        <el-table-column prop="ip" label="IP地址" width="180" />
        <el-table-column label="封禁时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column prop="reason" label="备注" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleUnbanIp(row)">解封</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > 0"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="fetchIps"
        @current-change="fetchIps"
        style="margin-top: 16px; justify-content: flex-end"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { adminApi } from '@/api/admin'

interface BannedIp {
  id: number
  ip: string
  reason: string
  created_at: string
}

const ips = ref<BannedIp[]>([])
const loading = ref(false)
const ipAddress = ref('')
const remark = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const formatTime = (val: string) => (val ? new Date(val).toLocaleString('zh-CN') : '')

const fetchIps = async () => {
  loading.value = true
  try {
    const res = await adminApi.getBannedIPs(page.value, pageSize.value)
    if (res.result) {
      ips.value = res.result.list || []
      total.value = res.result.total || 0
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

const handleBanIp = async () => {
  if (!ipAddress.value) {
    ElMessage.warning('请输入IP地址')
    return
  }
  try {
    await adminApi.createBannedIP({ ip: ipAddress.value, reason: remark.value })
    ElMessage.success('封禁成功')
    ipAddress.value = ''
    remark.value = ''
    page.value = 1
    fetchIps()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const handleUnbanIp = async (row: BannedIp) => {
  try {
    await adminApi.deleteBannedIP(row.id)
    ElMessage.success('解封成功')
    fetchIps()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

onMounted(() => {
  fetchIps()
})
</script>

<style scoped>
.banned-ip h2 {
  margin: 0 0 20px 0;
}

.ban-form {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.ban-input {
  width: 200px;
}
</style>
