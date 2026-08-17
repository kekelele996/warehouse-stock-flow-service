<template>
  <div class="empty">
    <div class="empty__icon">🗂️</div>
    <div class="empty__title">{{ title }}</div>
    <div v-if="description" class="empty__desc">{{ description }}</div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{ title?: string; description?: string }>(),
  { title: '暂无数据', description: '' }
)
</script>

<style scoped>
.empty {
  text-align: center;
  padding: 48px 16px;
  color: #7a8794;
}
.empty__icon { font-size: 40px; margin-bottom: 10px; }
.empty__title { font-size: 15px; font-weight: 600; color: #486581; }
.empty__desc { font-size: 13px; margin-top: 4px; }
</style>
