<template>
  <div class="notfound">
    <h1>404</h1>
    <p>页面不存在</p>
    <router-link to="/dashboard" class="btn btn--primary">返回首页</router-link>
  </div>
</template>

<script setup lang="ts"></script>

<style scoped>
.notfound { text-align: center; padding: 120px 20px; }
.notfound h1 { font-size: 64px; color: var(--primary); }
.notfound p { color: var(--text-muted); margin: 12px 0 24px; }
</style>
