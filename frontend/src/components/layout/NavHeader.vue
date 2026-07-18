<template>
  <div class="nav-header">
    <div class="logo">
      <router-link to="/">Folivora懒懒</router-link>
    </div>
    <div class="nav-links">
      <router-link to="/sheet-music" class="nav-link">乐谱展示</router-link>
      <router-link to="/about" class="nav-link">关于作者</router-link>
      <router-link to="/contact" class="nav-link">联系作者</router-link>
      <router-link to="/friend-links" class="nav-link">友情链接</router-link>
      <router-link v-if="userStore.isAdmin" to="/admin" class="nav-link">管理后台</router-link>
    </div>
    <div class="user-actions">
      <template v-if="userStore.isLoggedIn">
        <el-dropdown>
          <span class="user-info">
            {{ userStore.userInfo?.nickname || userStore.userInfo?.username }}
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>
                <router-link to="/profile">个人中心</router-link>
              </el-dropdown-item>
              <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>
      <template v-else>
        <router-link to="/login" class="nav-link">登录</router-link>
        <router-link to="/register" class="nav-link">注册</router-link>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ArrowDown } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const handleLogout = () => {
  userStore.logout()
  router.push('/')
}
</script>

<style scoped>
.nav-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  height: 60px;
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.logo a {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
  text-decoration: none;
}

.nav-links {
  display: flex;
  gap: 20px;
}

.nav-link {
  color: #606266;
  text-decoration: none;
  font-size: 16px;
}

.nav-link:hover {
  color: #409EFF;
}

.user-actions {
  display: flex;
  align-items: center;
  gap: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: #606266;
}

.user-info:hover {
  color: #409EFF;
}
</style>
