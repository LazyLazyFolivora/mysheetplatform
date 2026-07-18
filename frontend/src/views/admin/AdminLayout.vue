<template>
  <div class="admin-layout">
    <el-container>
      <el-aside width="200px" class="admin-aside">
        <div class="logo">
          <h2>管理后台</h2>
        </div>
        <el-menu
          :default-active="$route.path"
          class="admin-menu"
          @select="handleSelect"
        >
          <el-menu-item index="/admin">
            <el-icon><DataLine /></el-icon>
            <span>仪表盘</span>
          </el-menu-item>
          <el-menu-item index="/admin/users">
            <el-icon><User /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/sheet-music">
            <el-icon><Document /></el-icon>
            <span>乐谱管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/permissions">
            <el-icon><Lock /></el-icon>
            <span>权限管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/roles">
            <el-icon><UserFilled /></el-icon>
            <span>角色管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/contact-messages">
            <el-icon><Message /></el-icon>
            <span>联系消息</span>
          </el-menu-item>
          <el-menu-item index="/admin/friend-links">
            <el-icon><Link /></el-icon>
            <span>友情链接</span>
          </el-menu-item>
          <el-menu-item index="/admin/banned-ip">
            <el-icon><Lock /></el-icon>
            <span>IP封禁管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/system-config">
            <el-icon><Setting /></el-icon>
            <span>系统配置</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-container class="admin-main">
        <el-header height="60px" class="admin-header">
          <div class="header-right">
            <el-dropdown>
              <span class="user-info">
                <el-avatar :size="32" :src="userStore.userInfo?.avatar || ''" />
                <span class="username">{{ userStore.userInfo?.nickname || userStore.userInfo?.username || '管理员' }}</span>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        <el-main>
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import {
  DataLine,
  User,
  Document,
  Lock,
  UserFilled,
  Message,
  Link,
  Setting
} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const handleLogout = () => {
  userStore.logout()
  router.push('/admin/login')
}

const handleSelect = (index: string) => {
  router.push(index)
}
</script>

<style scoped>
.admin-layout {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.el-container {
  height: 100%;
}

.admin-aside {
  background-color: #001529;
  color: #fff;
  overflow-x: hidden;
  overflow-y: auto;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.admin-menu {
  border-right: none;
  background-color: transparent;
}

.admin-menu .el-menu-item {
  color: rgba(255, 255, 255, 0.65);
}

.admin-menu .el-menu-item:hover {
  color: #fff;
  background-color: rgba(255, 255, 255, 0.1);
}

.admin-menu .el-menu-item.is-active {
  color: #fff;
  background-color: #1890ff;
}

.admin-main {
  background-color: #f0f2f5;
}

.admin-header {
  background-color: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 20px;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin-left: 8px;
  color: #333;
}

.el-main {
  padding: 20px;
  overflow-y: auto;
}
</style>
