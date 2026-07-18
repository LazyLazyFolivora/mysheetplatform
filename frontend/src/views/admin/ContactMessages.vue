<template>
  <div class="contact-messages">
    <h2>联系消息</h2>
    <el-card>
      <el-table :data="messages" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.read ? 'success' : 'info'" size="small">
              {{ row.read ? '已读' : '未读' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="留言时间" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleView(row)">查看</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchMessages"
          @current-change="fetchMessages"
        />
      </div>
    </el-card>

    <!-- 留言详情 -->
    <el-dialog v-model="dialogVisible" title="留言详情" width="500px">
      <div v-if="currentMessage" class="message-detail">
        <div class="detail-item">
          <span class="label">姓名</span>
          <span>{{ currentMessage.name }}</span>
        </div>
        <div class="detail-item">
          <span class="label">邮箱</span>
          <span>{{ currentMessage.email || '未填写' }}</span>
        </div>
        <div class="detail-item">
          <span class="label">内容</span>
          <div class="message-content">{{ currentMessage.content }}</div>
        </div>
        <div class="detail-item">
          <span class="label">留言时间</span>
          <span>{{ formatTime(currentMessage.created_at) }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
        <el-button v-if="currentMessage && !currentMessage.read" type="primary" @click="handleMarkAsRead">
          标记为已读
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminApi } from '@/api/admin'

interface ContactMessage {
  id: number
  name: string
  email: string
  subject: string
  content: string
  read: boolean
  created_at: string
}

const messages = ref<ContactMessage[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const dialogVisible = ref(false)
const currentMessage = ref<ContactMessage | null>(null)

const formatTime = (val: string) => (val ? new Date(val).toLocaleString('zh-CN') : '')

const fetchMessages = async () => {
  loading.value = true
  try {
    const res = await adminApi.getMessages(page.value, pageSize.value)
    if (res.result) {
      messages.value = res.result.list || []
      total.value = res.result.total
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

const handleView = (row: ContactMessage) => {
  currentMessage.value = row
  dialogVisible.value = true
}

const handleMarkAsRead = async () => {
  if (!currentMessage.value) return
  try {
    await adminApi.markMessageRead(currentMessage.value.id)
    currentMessage.value.read = true
    ElMessage.success('已标记为已读')
    fetchMessages()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const handleDelete = async (row: ContactMessage) => {
  try {
    await ElMessageBox.confirm('确定删除这条留言吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await adminApi.deleteMessage(row.id)
    ElMessage.success('删除成功')
    fetchMessages()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

onMounted(() => {
  fetchMessages()
})
</script>

<style scoped>
.contact-messages h2 {
  margin: 0 0 20px 0;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.message-detail {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.detail-item .label {
  font-weight: bold;
  color: #606266;
}

.message-content {
  background: #f5f7fa;
  border-radius: 4px;
  padding: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
