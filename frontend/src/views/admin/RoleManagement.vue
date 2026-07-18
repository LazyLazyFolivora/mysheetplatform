<template>
  <div class="role-manage">
    <el-card>
      <div class="toolbar">
        <el-button type="primary" @click="openCreate">新增角色</el-button>
        <el-input
          v-model="search"
          placeholder="搜索角色名称"
          style="width: 200px; margin-left: 10px"
          clearable
        />
      </div>
      <el-table :data="filteredRoles" v-loading="loading" style="width: 100%; margin-top: 10px">
        <el-table-column prop="name" label="角色名称" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-button size="small" @click="openEdit(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
            <el-button size="small" type="success" @click="openModuleDialog(row)">分配模块</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑角色 -->
    <el-dialog v-model="editVisible" :title="editForm.id ? '编辑角色' : '新增角色'" width="480px">
      <el-form :model="editForm" label-width="90px">
        <el-form-item label="角色名称" required>
          <el-input v-model="editForm.name" maxlength="50" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editForm.description" maxlength="200" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <!-- 分配权限模块 -->
    <el-dialog v-model="moduleVisible" title="分配权限模块" width="500px">
      <el-tree
        ref="treeRef"
        :data="moduleTree"
        show-checkbox
        check-strictly
        node-key="id"
        :default-checked-keys="checkedModuleIds"
        :props="{ label: 'name', children: 'children' }"
        default-expand-all
      />
      <template #footer>
        <el-button @click="moduleVisible = false">取消</el-button>
        <el-button type="primary" :loading="assigning" @click="handleAssignModules">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { ElTree } from 'element-plus'
import {
  getAllRoles,
  createRole,
  updateRole,
  deleteRole,
  getPermissionTree,
  getAssignedModules,
  assignModules,
} from '@/api/permission'
import type { Role, ModuleTreeNode } from '@/types/permission'

const roles = ref<Role[]>([])
const loading = ref(false)
const search = ref('')

const editVisible = ref(false)
const saving = ref(false)
const editForm = ref({ id: 0, name: '', description: '' })

const moduleVisible = ref(false)
const assigning = ref(false)
const moduleTree = ref<ModuleTreeNode[]>([])
const checkedModuleIds = ref<number[]>([])
const currentRole = ref<Role | null>(null)
const treeRef = ref<InstanceType<typeof ElTree>>()

const filteredRoles = computed(() => {
  if (!search.value) return roles.value
  return roles.value.filter((r) => r.name.includes(search.value))
})

const fetchRoles = async () => {
  loading.value = true
  try {
    const res = await getAllRoles()
    roles.value = res.result || []
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    loading.value = false
  }
}

const openCreate = () => {
  editForm.value = { id: 0, name: '', description: '' }
  editVisible.value = true
}

const openEdit = (row: Role) => {
  editForm.value = { id: row.id, name: row.name, description: row.description || '' }
  editVisible.value = true
}

const handleSave = async () => {
  if (!editForm.value.name) {
    ElMessage.warning('请填写角色名称')
    return
  }
  try {
    saving.value = true
    const data = { name: editForm.value.name, description: editForm.value.description }
    if (editForm.value.id) {
      await updateRole(editForm.value.id, data)
    } else {
      await createRole(data)
    }
    ElMessage.success('保存成功')
    editVisible.value = false
    fetchRoles()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    saving.value = false
  }
}

const handleDelete = async (row: Role) => {
  try {
    await ElMessageBox.confirm(`确定要删除角色「${row.name}」吗？`, '提示', { type: 'warning' })
  } catch {
    return
  }
  try {
    await deleteRole(row.id)
    ElMessage.success('删除成功')
    fetchRoles()
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const openModuleDialog = async (row: Role) => {
  currentRole.value = row
  checkedModuleIds.value = []
  moduleVisible.value = true
  try {
    const [treeRes, assignedRes] = await Promise.all([getPermissionTree(), getAssignedModules(row.id)])
    moduleTree.value = treeRes.result || []
    checkedModuleIds.value = assignedRes.result || []
    treeRef.value?.setCheckedKeys(checkedModuleIds.value)
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

const handleAssignModules = async () => {
  if (!currentRole.value || !treeRef.value) return
  const checked = treeRef.value.getCheckedKeys(false) as number[]
  try {
    assigning.value = true
    await assignModules(currentRole.value.id, checked)
    ElMessage.success('分配成功')
    moduleVisible.value = false
  } catch {
    // 错误提示已由 request 拦截器统一处理
  } finally {
    assigning.value = false
  }
}

onMounted(() => {
  fetchRoles()
})
</script>

<style scoped>
.role-manage {
  padding: 20px;
}

.toolbar {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}
</style>
