<template>
  <div class="permission-manage">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card>
          <div class="tree-title">权限模块</div>
          <el-tree
            :data="moduleTree"
            :props="{ label: 'name', children: 'children' }"
            node-key="id"
            highlight-current
            default-expand-all
            @node-click="handleModuleClick"
          >
            <template #default="{ node, data }">
              <span>{{ node.label }}</span>
              <span class="module-actions">
                <el-button link type="primary" @click.stop="openEditModule(data)">编辑</el-button>
                <el-button link type="danger" @click.stop="handleDeleteModule(data)">删除</el-button>
                <el-button link type="success" @click.stop="openAssignDialog(data)">分配权限</el-button>
              </span>
            </template>
          </el-tree>
          <el-button type="primary" size="small" style="margin-top: 10px" @click="openCreateModule">
            新增模块
          </el-button>
        </el-card>
      </el-col>

      <el-col :span="18">
        <el-card>
          <div class="toolbar">
            <el-button type="primary" @click="openCreatePermission">新增权限</el-button>
            <el-input
              v-model="search"
              placeholder="搜索权限名称/标识"
              style="width: 200px; margin-left: 10px"
              clearable
            />
            <el-tag v-if="currentModule" closable style="margin-left: 10px" @close="clearModuleFilter">
              {{ currentModule.name }}
            </el-tag>
          </div>
          <el-table :data="filteredPermissions" v-loading="loading" style="width: 100%; margin-top: 10px">
            <el-table-column prop="name" label="权限名称" />
            <el-table-column prop="code" label="权限标识" />
            <el-table-column prop="url" label="接口路径" />
            <el-table-column prop="method" label="请求方式" width="100" />
            <el-table-column label="操作" width="180">
              <template #default="{ row }">
                <el-button size="small" @click="openEditPermission(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="handleDeletePermission(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 新增/编辑权限 -->
    <el-dialog v-model="permVisible" :title="permForm.id ? '编辑权限' : '新增权限'" width="520px">
      <el-form :model="permForm" label-width="90px">
        <el-form-item label="权限名称" required>
          <el-input v-model="permForm.name" maxlength="100" />
        </el-form-item>
        <el-form-item label="权限标识" required>
          <el-input v-model="permForm.code" maxlength="100" placeholder="如 sheet-music:create" />
        </el-form-item>
        <el-form-item label="接口路径">
          <el-input v-model="permForm.url" maxlength="200" placeholder="如 /api/admin/sheet-music" />
        </el-form-item>
        <el-form-item label="请求方式">
          <el-select v-model="permForm.method" placeholder="请选择" clearable style="width: 100%">
            <el-option label="GET" value="GET" />
            <el-option label="POST" value="POST" />
            <el-option label="PUT" value="PUT" />
            <el-option label="DELETE" value="DELETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="所属模块">
          <el-tree-select
            v-model="permForm.module_id"
            :data="moduleTree"
            :props="{ label: 'name', children: 'children' }"
            placeholder="请选择模块"
            node-key="id"
            check-strictly
            clearable
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="permVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSavePermission">保存</el-button>
      </template>
    </el-dialog>

    <!-- 新增/编辑模块 -->
    <el-dialog v-model="moduleVisible" :title="moduleForm.id ? '编辑模块' : '新增模块'" width="480px">
      <el-form :model="moduleForm" label-width="90px">
        <el-form-item label="模块名称" required>
          <el-input v-model="moduleForm.name" maxlength="100" />
        </el-form-item>
        <el-form-item label="上级模块">
          <el-tree-select
            v-model="moduleForm.parent_id"
            :data="moduleTree"
            :props="{ label: 'name', children: 'children' }"
            placeholder="请选择上级模块，不选则为根模块"
            node-key="id"
            check-strictly
            clearable
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="moduleVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveModule">保存</el-button>
      </template>
    </el-dialog>

    <!-- 分配权限 -->
    <el-dialog v-model="assignVisible" :title="`分配权限 — ${assignTarget?.name || ''}`" width="640px">
      <el-transfer
        v-model="selectedPermIds"
        :data="transferData"
        :titles="['可选权限', '已选权限']"
      />
      <template #footer>
        <el-button @click="assignVisible = false">取消</el-button>
        <el-button type="primary" :loading="assigning" @click="handleAssignPermissions">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getPermissionTree,
  createModule,
  updateModule,
  deleteModule,
  getAllPermissions,
  createPermission,
  updatePermission,
  deletePermission,
  getModulePermissions,
  assignPermissionsToModule,
} from '@/api/permission'
import type { ModuleTreeNode, Permission } from '@/types/permission'

const moduleTree = ref<ModuleTreeNode[]>([])
const allPermissions = ref<Permission[]>([])
const loading = ref(false)
const saving = ref(false)
const search = ref('')

const currentModule = ref<ModuleTreeNode | null>(null)
const currentModulePermIds = ref<number[]>([])

const permVisible = ref(false)
const permForm = ref<{
  id: number
  name: string
  code: string
  url: string
  method: string
  module_id: number | null
}>({ id: 0, name: '', code: '', url: '', method: '', module_id: null })

const moduleVisible = ref(false)
const moduleForm = ref<{ id: number; name: string; parent_id: number | null; path: string; icon: string }>({
  id: 0,
  name: '',
  parent_id: null,
  path: '',
  icon: '',
})

const assignVisible = ref(false)
const assigning = ref(false)
const assignTarget = ref<ModuleTreeNode | null>(null)
const selectedPermIds = ref<number[]>([])

const transferData = computed(() =>
  allPermissions.value.map((p) => ({ key: p.id, label: p.name, disabled: false }))
)

const filteredPermissions = computed(() => {
  let list = allPermissions.value
  if (currentModule.value) {
    list = list.filter((p) => currentModulePermIds.value.includes(p.id))
  }
  if (search.value) {
    const kw = search.value.toLowerCase()
    list = list.filter(
      (p) => p.name.toLowerCase().includes(kw) || p.code.toLowerCase().includes(kw)
    )
  }
  return list
})

const fetchTree = async () => {
  try {
    const res = await getPermissionTree()
    moduleTree.value = res.result || []
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const fetchPermissions = async () => {
  loading.value = true
  try {
    const res = await getAllPermissions()
    allPermissions.value = res.result || []
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

// ===== 模块 =====

const handleModuleClick = async (node: ModuleTreeNode) => {
  currentModule.value = node
  try {
    const res = await getModulePermissions(node.id)
    currentModulePermIds.value = res.result || []
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const clearModuleFilter = () => {
  currentModule.value = null
  currentModulePermIds.value = []
}

const openCreateModule = () => {
  moduleForm.value = { id: 0, name: '', parent_id: null, path: '', icon: '' }
  moduleVisible.value = true
}

const openEditModule = (node: ModuleTreeNode) => {
  moduleForm.value = {
    id: node.id,
    name: node.name,
    parent_id: node.parent_id || null,
    path: node.path || '',
    icon: node.icon || '',
  }
  moduleVisible.value = true
}

const handleSaveModule = async () => {
  if (!moduleForm.value.name) {
    ElMessage.warning('请填写模块名称')
    return
  }
  if (moduleForm.value.id && moduleForm.value.parent_id === moduleForm.value.id) {
    ElMessage.warning('上级模块不能是自身')
    return
  }
  try {
    saving.value = true
    const data = {
      name: moduleForm.value.name,
      parent_id: moduleForm.value.parent_id,
      path: moduleForm.value.path,
      icon: moduleForm.value.icon,
    }
    if (moduleForm.value.id) {
      await updateModule(moduleForm.value.id, data)
    } else {
      await createModule(data)
    }
    ElMessage.success('保存成功')
    moduleVisible.value = false
    fetchTree()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    saving.value = false
  }
}

const handleDeleteModule = async (node: ModuleTreeNode) => {
  try {
    await ElMessageBox.confirm('确定要删除该模块吗？子模块会一并删除', '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteModule(node.id)
    ElMessage.success('删除成功')
    if (currentModule.value?.id === node.id) {
      clearModuleFilter()
    }
    fetchTree()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

// ===== 权限 =====

const openCreatePermission = () => {
  permForm.value = {
    id: 0,
    name: '',
    code: '',
    url: '',
    method: '',
    module_id: currentModule.value?.id || null,
  }
  permVisible.value = true
}

const openEditPermission = (row: Permission) => {
  permForm.value = {
    id: row.id,
    name: row.name,
    code: row.code,
    url: row.url || '',
    method: row.method || '',
    module_id: currentModule.value?.id || null,
  }
  permVisible.value = true
}

const handleSavePermission = async () => {
  if (!permForm.value.name || !permForm.value.code) {
    ElMessage.warning('请填写权限名称和权限标识')
    return
  }
  try {
    saving.value = true
    const data = {
      name: permForm.value.name,
      code: permForm.value.code,
      url: permForm.value.url,
      method: permForm.value.method,
      module_id: permForm.value.module_id,
    }
    if (permForm.value.id) {
      await updatePermission(permForm.value.id, data)
    } else {
      await createPermission(data)
    }
    ElMessage.success('保存成功')
    permVisible.value = false
    fetchPermissions()
    if (currentModule.value) {
      handleModuleClick(currentModule.value)
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    saving.value = false
  }
}

const handleDeletePermission = async (row: Permission) => {
  try {
    await ElMessageBox.confirm('确定要删除该权限吗？', '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await deletePermission(row.id)
    ElMessage.success('删除成功')
    fetchPermissions()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

// ===== 分配权限 =====

const openAssignDialog = async (node: ModuleTreeNode) => {
  assignTarget.value = node
  selectedPermIds.value = []
  assignVisible.value = true
  try {
    const res = await getModulePermissions(node.id)
    selectedPermIds.value = res.result || []
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const handleAssignPermissions = async () => {
  if (!assignTarget.value) return
  try {
    assigning.value = true
    await assignPermissionsToModule(assignTarget.value.id, selectedPermIds.value)
    ElMessage.success('保存成功')
    assignVisible.value = false
    if (currentModule.value?.id === assignTarget.value.id) {
      handleModuleClick(currentModule.value)
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    assigning.value = false
  }
}

onMounted(() => {
  fetchTree()
  fetchPermissions()
})
</script>

<style scoped>
.permission-manage {
  padding: 20px;
}

.toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.tree-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
}

.module-actions {
  margin-left: 8px;
  opacity: 0;
  transition: opacity 0.3s;
}

:deep(.el-tree-node__content:hover) .module-actions {
  opacity: 1;
}

:deep(.el-transfer) {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

:deep(.el-transfer-panel) {
  width: 42%;
}
</style>
