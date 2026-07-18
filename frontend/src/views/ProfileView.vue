<template>
  <div class="profile">
    <h1>个人资料</h1>
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <span>基本信息</span>
        </div>
      </template>
      <el-form :model="profile" label-width="100px">
        <el-form-item label="头像">
          <img v-if="profile.avatar" :src="profile.avatar" class="avatar" />
          <el-avatar v-else :size="80">{{ profile.nickname?.charAt(0) || profile.username?.charAt(0) }}</el-avatar>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="profile.username" disabled />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="profile.nickname" disabled />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="profile.email" disabled />
        </el-form-item>
        <el-form-item label="积分">
          <el-input :model-value="String(profile.points ?? 0)" disabled />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import type { CurrentUser } from '@/types/user'

const userStore = useUserStore()

const profile = ref<Partial<CurrentUser>>({
  username: '',
  nickname: '',
  email: '',
  avatar: '',
  points: 0
})

onMounted(() => {
  if (userStore.userInfo) {
    profile.value = { ...userStore.userInfo }
  }
})
</script>

<style scoped>
.profile {
  padding: 20px;
}

.profile-card {
  max-width: 600px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
}
</style>
