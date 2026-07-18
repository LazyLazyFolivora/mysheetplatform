<template>
  <div class="about">
    <h1>关于作者</h1>
    <el-card class="about-card">
      <div class="about-content">
        <h2>Folivora懒懒</h2>
        <p>热爱音乐与编曲，专注于乐谱扒谱与改编。</p>
        <ul>
          <li>QQ：1015553652</li>
          <li>
            bilibili：
            <a href="https://space.bilibili.com/436626418" target="_blank" rel="noopener">Folivora懒懒</a>
          </li>
        </ul>
        <p>本站提供原创/改编乐谱的展示与下载，欢迎交流。</p>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.about {
  padding: 20px;
}
.about-card {
  max-width: 800px;
  margin: 0 auto;
}
.about-content {
  padding: 20px;
  line-height: 1.8;
}
.about-content h2 {
  margin: 0 0 12px 0;
  color: #333;
}
.about-content p {
  margin: 8px 0;
  color: #666;
}
.about-content ul {
  margin: 8px 0;
  padding-left: 20px;
  color: #666;
}
.about-content a {
  color: #409eff;
  text-decoration: none;
}
.about-content a:hover {
  text-decoration: underline;
}
</style>
