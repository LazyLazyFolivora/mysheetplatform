<template>
  <div class="nav-footer">
    <div class="footer-content">
      <div class="copyright">
        QQ：1015553652&nbsp;&nbsp;&nbsp;&nbsp;
        <a href="https://space.bilibili.com/436626418" target="_blank" rel="noopener">bilibili：Folivora懒懒</a>
      </div>
      <div class="copyright">
        © 2025 Folivora懒懒. All rights reserved.
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.nav-footer {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.footer-content {
  text-align: center;
}

.copyright {
  color: #999;
  font-size: 14px;
}
.copyright a {
  color: #999;
  text-decoration: none;
  margin: 0 5px;
  transition: color 0.2s;
}
.copyright a:hover {
  color: #409EFF;
  text-decoration: underline;
}
</style>
