<template>
  <div class="about">
    <h1>关于作者</h1>
    <el-card class="about-card">
      <div v-if="content" class="about-content" v-html="content"></div>
      <div v-else v-loading="loading" class="about-content" :style="{ minHeight: '200px' }"></div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import DOMPurify from 'dompurify'
import { getPublicSystemConfig } from '@/api/system-config'

const content = ref('')
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await getPublicSystemConfig('about_us')
    content.value = DOMPurify.sanitize(res.result?.config_val || '')
  } catch {
    content.value = ''
  } finally {
    loading.value = false
  }
})
</script>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.about {
  padding: 20px;
  max-width: $content-max-width;
  margin: 0 auto;
}
.about-card {
  max-width: 800px;
  margin: 20px auto 0;
  border-color: $border-light;
  border-radius: $radius-xl;
}
.about-content {
  padding: 20px;
  line-height: 1.8;
}
.about-content :deep(h2) {
  margin: 0 0 12px 0;
  color: $text-primary;
}
.about-content :deep(p) {
  margin: 8px 0;
  color: $text-secondary;
}
.about-content :deep(ul) {
  margin: 8px 0;
  padding-left: 20px;
  color: $text-secondary;
}
.about-content :deep(a) {
  color: $primary;
  text-decoration: none;
  transition: color $transition-base;
}
.about-content :deep(a:hover) {
  color: $primary-hover;
  text-decoration: underline;
}
</style>
