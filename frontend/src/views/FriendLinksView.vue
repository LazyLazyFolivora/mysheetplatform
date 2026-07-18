<template>
  <div class="friend-links">
    <h1>友情链接</h1>
    <el-card class="links-card">
      <el-empty v-if="!links.length" description="暂无友情链接" />
      <el-row v-else :gutter="20">
        <el-col :span="8" v-for="link in links" :key="link.id">
          <el-card class="link-card">
            <div class="link-content">
              <img v-if="link.logo" :src="link.logo" class="link-logo" />
              <div class="link-info">
                <h3>{{ link.name }}</h3>
                <a :href="link.url" target="_blank" rel="noopener">访问网站</a>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/utils/request'
import type { ApiResponse } from '@/types/api'

interface FriendLink {
  id: number
  name: string
  url: string
  logo?: string
  sort_order: number
}

const links = ref<FriendLink[]>([])

const fetchLinks = async () => {
  try {
    const res = await request.get<ApiResponse<FriendLink[]>>('/friend-links')
    if (Array.isArray(res.result)) {
      links.value = res.result
    }
  } catch {
    // 错误提示已由 request 拦截器统一处理
  }
}

onMounted(() => {
  fetchLinks()
})
</script>

<style scoped>
.friend-links {
  padding: 20px;
}
.links-card {
  margin-top: 20px;
}
.link-card {
  margin-bottom: 20px;
}
.link-content {
  display: flex;
  align-items: center;
  gap: 10px;
}
.link-logo {
  width: 50px;
  height: 50px;
  object-fit: cover;
}
.link-info {
  flex: 1;
}
.link-info h3 {
  margin: 0;
}
.link-info a {
  color: #409eff;
  text-decoration: none;
}
</style>
