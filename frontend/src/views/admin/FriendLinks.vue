<template>
  <div class="friend-links-manage">
    <div class="toolbar">
      <h2>友情链接管理</h2>
      <el-button type="primary" @click="openCreate">添加链接</el-button>
    </div>

    <el-card>
      <el-table :data="links" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="名称" width="140" />
        <el-table-column prop="url" label="链接" show-overflow-tooltip />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column label="Logo" width="100">
          <template #default="{ row }">
            <img v-if="row.logo" :src="row.logo" class="logo-img" alt="logo" />
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.is_enabled ? 'success' : 'danger'" size="small">
              {{ row.is_enabled ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
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
          @size-change="fetchLinks"
          @current-change="fetchLinks"
        />
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑友情链接' : '添加友情链接'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="form.name" maxlength="100" placeholder="2-100 个字符" />
        </el-form-item>
        <el-form-item label="链接" required>
          <el-input v-model="form.url" placeholder="https://" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" maxlength="200" />
        </el-form-item>
        <el-form-item label="Logo">
          <div class="logo-upload">
            <img v-if="form.logo" :src="form.logo" class="logo-preview" alt="logo" />
            <input ref="logoInputRef" type="file" accept="image/*" hidden @change="handleLogoChange" />
            <el-button :loading="uploading" @click="logoInputRef?.click()">
              {{ form.logo ? '更换 Logo' : '上传 Logo' }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort_order" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="form.is_enabled" active-text="显示" inactive-text="隐藏" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminApi } from '@/api/admin'
import { uploadImage } from '@/api/file'

interface FriendLink {
  id: number
  name: string
  url: string
  description: string
  logo: string
  is_enabled: boolean
  sort_order: number
  created_at: string
}

const links = ref<FriendLink[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const dialogVisible = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)
const uploading = ref(false)
const logoInputRef = ref<HTMLInputElement>()

const emptyForm = () => ({
  name: '',
  url: '',
  description: '',
  logo: '',
  sort_order: 0,
  is_enabled: true,
})
const form = ref(emptyForm())

const formatTime = (val: string) => (val ? new Date(val).toLocaleString('zh-CN') : '')

const fetchLinks = async () => {
  loading.value = true
  try {
    const res = await adminApi.getFriendLinks(page.value, pageSize.value)
    if (res.result) {
      links.value = res.result.list || []
      total.value = res.result.total
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  editingId.value = null
  form.value = emptyForm()
  dialogVisible.value = true
}

const openEdit = (row: FriendLink) => {
  editingId.value = row.id
  form.value = {
    name: row.name,
    url: row.url,
    description: row.description || '',
    logo: row.logo || '',
    sort_order: row.sort_order,
    is_enabled: row.is_enabled,
  }
  dialogVisible.value = true
}

const handleLogoChange = async (e: Event) => {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  try {
    uploading.value = true
    const res = await uploadImage(file)
    if (res.result?.url) {
      form.value.logo = res.result.url
      ElMessage.success('Logo 上传成功')
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    uploading.value = false
  }
}

const handleSave = async () => {
  if (!form.value.name || form.value.name.length < 2) {
    ElMessage.warning('名称长度需在 2 个字符以上')
    return
  }
  if (!/^https?:\/\/.+/.test(form.value.url)) {
    ElMessage.warning('请输入有效的链接地址（http:// 或 https:// 开头）')
    return
  }
  try {
    saving.value = true
    if (editingId.value) {
      await adminApi.updateFriendLink(editingId.value, form.value)
      ElMessage.success('更新成功')
    } else {
      await adminApi.createFriendLink(form.value)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchLinks()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    saving.value = false
  }
}

const handleDelete = async (row: FriendLink) => {
  try {
    await ElMessageBox.confirm(`确定删除链接「${row.name}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await adminApi.deleteFriendLink(row.id)
    ElMessage.success('删除成功')
    fetchLinks()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

onMounted(() => {
  fetchLinks()
})
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.toolbar h2 {
  margin: 0;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.logo-img {
  height: 32px;
  max-width: 80px;
  object-fit: contain;
}

.logo-upload {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo-preview {
  height: 48px;
  max-width: 120px;
  object-fit: contain;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}
</style>
