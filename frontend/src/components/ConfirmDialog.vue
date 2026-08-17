<template>
  <div v-if="open" class="overlay" @click.self="$emit('close')">
    <div class="dialog">
      <div class="dialog__title">{{ title }}</div>
      <div class="dialog__body">{{ message }}</div>
      <div class="dialog__actions">
        <button class="btn btn--ghost" @click="$emit('close')">取消</button>
        <button class="btn btn--danger" :disabled="loading" @click="$emit('confirm')">
          {{ loading ? '处理中…' : confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{ open: boolean; title?: string; message?: string; confirmText?: string; loading?: boolean }>(),
  { title: '确认操作', message: '确定要执行该操作吗？', confirmText: '确认', loading: false }
)
defineEmits<{ (e: 'confirm'): void; (e: 'close'): void }>()
</script>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(16, 42, 67, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.dialog {
  width: 400px;
  max-width: calc(100vw - 32px);
  background: #fff;
  border-radius: 14px;
  padding: 24px;
  box-shadow: 0 20px 60px rgba(16, 42, 67, 0.2);
}
.dialog__title { font-size: 17px; font-weight: 700; color: #102a43; }
.dialog__body { margin-top: 10px; font-size: 14px; color: #486581; line-height: 1.6; }
.dialog__actions { margin-top: 22px; display: flex; justify-content: flex-end; gap: 10px; }
</style>
