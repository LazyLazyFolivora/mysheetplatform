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

<style scoped lang="scss">
@use "@/styles/variables" as *;

.nav-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 $space-xl;
  height: 60px;
  max-width: $content-max-width;
  margin: 0 auto;
}

.logo a {
  font-family: $font-family-serif;
  font-size: 22px;
  font-weight: 600;
  color: $text-primary;
  letter-spacing: 0.02em;
  text-decoration: none;

  &:hover {
    color: $primary-dark;
  }
}

.nav-links {
  display: flex;
  gap: $space-xl;
}

.nav-link {
  position: relative;
  color: $text-secondary;
  text-decoration: none;
  font-size: 15px;
  padding: 4px 0;
  transition: color $transition-base;

  &::after {
    content: "";
    position: absolute;
    left: 50%;
    bottom: -4px;
    width: 0;
    height: 2px;
    border-radius: 1px;
    background-color: $primary;
    transform: translateX(-50%);
    transition: width $transition-base;
  }

  &:hover {
    color: $text-primary;
  }

  &.router-link-active {
    color: $text-primary;

    &::after {
      width: 18px;
    }
  }
}

.user-actions {
  display: flex;
  align-items: center;
  gap: $space-lg;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  color: $text-primary;
  transition: color $transition-base;

  &:hover {
    color: $primary-dark;
  }
}
</style>
