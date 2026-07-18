<template>
  <div class="users-view">
    <h2>用户管理</h2>
    <el-card>
      <el-table :data="users" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="points" label="积分" width="90" />
        <el-table-column label="管理员" width="90">
          <template #default="{ row }">
            <el-tag v-if="row.is_admin" type="danger" size="small">是</el-tag>
            <span v-else>否</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '正常' : '封禁' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 1"
              type="danger"
              size="small"
              @click="handleBan(row)"
            >封禁</el-button>
            <el-button
              v-else
              type="success"
              size="small"
              @click="handleUnban(row)"
            >解封</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminApi } from '@/api/admin'

interface AdminUser {
  id: number
  username: string
  nickname: string
  email: string
  points: number
  is_admin: boolean
  status: number
}

const users = ref<AdminUser[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await adminApi.getUsers(page.value, pageSize.value)
    if (res.result) {
      users.value = res.result.list || []
      total.value = res.result.total
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

const handlePageChange = (val: number) => {
  page.value = val
  fetchUsers()
}

const handleBan = async (row: AdminUser) => {
  try {
    await ElMessageBox.confirm(`确定封禁用户「${row.username}」吗？`, '提示', { type: 'warning' })
    await adminApi.banUser(row.id)
    ElMessage.success('已封禁')
    fetchUsers()
  } catch {
    // 取消或失败
  }
}

const handleUnban = async (row: AdminUser) => {
  try {
    await adminApi.unbanUser(row.id)
    ElMessage.success('已解封')
    fetchUsers()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.users-view h2 {
  margin: 0 0 20px 0;
}
.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
