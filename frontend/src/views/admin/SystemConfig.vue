<template>
  <div class="system-config">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统配置管理</span>
          <el-button type="primary" :icon="Plus" @click="handleAdd">添加配置</el-button>
        </div>
      </template>

      <el-table :data="configs" v-loading="loading" style="width: 100%">
        <el-table-column prop="config_key" label="配置键" width="200" />
        <el-table-column label="配置值" min-width="300">
          <template #default="{ row }">
            <el-input
              v-if="row.editing"
              v-model="row.tempValue"
              type="textarea"
              :rows="4"
              @blur="handleSave(row)"
              @keyup.enter="handleSave(row)"
            />
            <div v-else class="editable-cell" @click="handleEdit(row)">
              <div
                v-if="isRichTextConfig(row.config_key)"
                v-html="row.config_val || '空'"
                class="rich-text-preview"
              ></div>
              <div v-else>{{ row.config_val || '空' }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" width="200" />
        <el-table-column label="更新时间" width="180">
          <template #default="{ row }">{{ formatTime(row.updated_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button v-if="!row.editing" type="primary" size="small" :icon="Edit" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button v-else type="success" size="small" :icon="Check" @click="handleSave(row)">
              保存
            </el-button>
            <el-button type="danger" size="small" :icon="Delete" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加配置对话框 -->
    <el-dialog v-model="dialogVisible" title="添加系统配置" width="500px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="配置键" prop="key">
          <el-input v-model="form.key" placeholder="请输入配置键" />
        </el-form-item>
        <el-form-item label="配置值" prop="value">
          <el-input v-model="form.value" type="textarea" :rows="3" placeholder="请输入配置值" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="form.remark" placeholder="请输入备注说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Edit, Check, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getSystemConfigs,
  createSystemConfig,
  updateSystemConfig,
  deleteSystemConfig,
  type SystemConfig,
} from '@/api/system-config'

interface EditableConfig extends SystemConfig {
  editing: boolean
  tempValue: string
}

const RICH_TEXT_KEYS = ['about_us', 'privacy_policy', 'terms_of_service', 'site_description']

const configs = ref<EditableConfig[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const submitting = ref(false)
const formRef = ref<FormInstance>()

const form = ref({ key: '', value: '', remark: '' })

const rules: FormRules = {
  key: [
    { required: true, message: '请输入配置键', trigger: 'blur' },
    { min: 1, max: 100, message: '配置键长度在 1 到 100 个字符', trigger: 'blur' },
  ],
  value: [{ required: true, message: '请输入配置值', trigger: 'blur' }],
}

const formatTime = (val: string) => (val ? new Date(val).toLocaleString('zh-CN') : '')

const isRichTextConfig = (key: string) => RICH_TEXT_KEYS.includes(key)

const fetchConfigs = async () => {
  loading.value = true
  try {
    const res = await getSystemConfigs()
    configs.value = (res.result || []).map((c) => ({ ...c, editing: false, tempValue: c.config_val }))
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

const handleEdit = (row: EditableConfig) => {
  row.editing = true
  row.tempValue = row.config_val
}

const handleSave = async (row: EditableConfig) => {
  if (!row.editing) return
  try {
    await updateSystemConfig(row.config_key, row.tempValue, row.remark)
    row.config_val = row.tempValue
    row.editing = false
    ElMessage.success('配置更新成功')
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const handleDelete = async (row: EditableConfig) => {
  try {
    await ElMessageBox.confirm(`确定要删除配置 "${row.config_key}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })
  } catch {
    return
  }
  try {
    await deleteSystemConfig(row.config_key)
    ElMessage.success('配置删除成功')
    fetchConfigs()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const handleAdd = () => {
  form.value = { key: '', value: '', remark: '' }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    return
  }
  try {
    submitting.value = true
    await createSystemConfig(form.value.key, form.value.value, form.value.remark)
    ElMessage.success('配置添加成功')
    dialogVisible.value = false
    fetchConfigs()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchConfigs()
})
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.editable-cell {
  cursor: pointer;
  padding: 5px;
  border: 1px solid transparent;
  border-radius: 4px;
  min-height: 20px;
}

.editable-cell:hover {
  border-color: #dcdfe6;
  background-color: #f5f7fa;
}

.rich-text-preview {
  max-height: 100px;
  overflow: hidden;
  line-height: 1.5;
}

.rich-text-preview :deep(p) {
  margin: 0 0 8px 0;
}

.rich-text-preview :deep(h1),
.rich-text-preview :deep(h2),
.rich-text-preview :deep(h3),
.rich-text-preview :deep(h4),
.rich-text-preview :deep(h5),
.rich-text-preview :deep(h6) {
  margin: 8px 0;
  font-weight: bold;
}

.rich-text-preview :deep(ul),
.rich-text-preview :deep(ol) {
  margin: 8px 0;
  padding-left: 20px;
}

.rich-text-preview :deep(li) {
  margin: 4px 0;
}
</style>
