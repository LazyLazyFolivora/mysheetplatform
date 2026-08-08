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

<style scoped lang="scss">
@use "@/styles/variables" as *;

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
  color: $text-secondary;
  font-size: 13px;

  a {
    color: $text-secondary;
    text-decoration: none;
    margin: 0 5px;
    transition: color $transition-base;

    &:hover {
      color: $primary-dark;
      text-decoration: underline;
    }
  }
}
</style>
