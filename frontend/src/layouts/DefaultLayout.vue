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

<style scoped>
.app-container {
  min-height: 100vh;
}

.el-header {
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: fixed;
  width: 100%;
  z-index: 1000;
}

.el-main {
  margin-top: 60px;
  padding: 20px;
  background-color: #f5f7fa;
}

.el-footer {
  background-color: #fff;
  border-top: 1px solid #e6e6e6;
}
</style>
