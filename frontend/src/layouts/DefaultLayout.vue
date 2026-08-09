<template>
  <el-container class="app-container">
    <!-- 顶部导航栏 -->
    <el-header height="60px">
      <nav-header />
    </el-header>

    <!-- 主要内容区 -->
    <el-main>
      <router-view />
    </el-main>

    <!-- 页脚 -->
    <el-footer height="60px">
      <nav-footer />
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import NavHeader from '@/components/layout/NavHeader.vue'
import NavFooter from '@/components/layout/NavFooter.vue'
</script>

<style scoped lang="scss">
@use "@/styles/variables" as *;

.app-container {
  min-height: 100vh;
  background-color: $bg-page;
}

.el-header {
  position: fixed;
  width: 100%;
  z-index: 1000;
  padding: 0;
  background-color: rgba(255, 255, 255, 0.82);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid $border-light;
}

.el-main {
  margin-top: 60px;
  padding: 20px;
  background-color: $bg-page;
}

.el-footer {
  background-color: transparent;
  border-top: 1px solid $border-light;
}
</style>
