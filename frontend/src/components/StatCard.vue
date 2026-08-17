<template>
  <div class="stat-card">
    <div class="stat-card__icon">{{ icon }}</div>
    <div class="stat-card__body">
      <div class="stat-card__value">{{ value }}</div>
      <div class="stat-card__label">{{ label }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{ icon: string; label: string; value: string | number }>()
</script>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  background: #fff;
  border: 1px solid #e6eaf0;
  border-radius: 12px;
  padding: 18px 20px;
  box-shadow: 0 1px 3px rgba(16, 42, 67, 0.06);
}
.stat-card__icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #eef4ff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22px;
}
.stat-card__value {
  font-size: 24px;
  font-weight: 700;
  color: #102a43;
  line-height: 1.2;
}
.stat-card__label {
  font-size: 13px;
  color: #7a8794;
  margin-top: 2px;
}
</style>
