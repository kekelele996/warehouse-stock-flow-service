<template>
  <div v-if="open" class="overlay" @click.self="$emit('close')">
    <div class="modal" :style="{ width: width }">
      <div class="modal__head">
        <div class="modal__title">{{ title }}</div>
        <button class="modal__close" @click="$emit('close')">✕</button>
      </div>
      <div class="modal__body">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ open: boolean; title: string; width?: string }>(), { width: '640px' })
defineEmits<{ (e: 'close'): void }>()
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
.modal {
  max-width: calc(100vw - 32px);
  max-height: calc(100vh - 64px);
  background: #fff;
  border-radius: 14px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(16, 42, 67, 0.2);
}
.modal__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #eef1f4;
}
.modal__title { font-size: 16px; font-weight: 700; color: #102a43; }
.modal__close {
  border: none;
  background: none;
  font-size: 14px;
  color: #7a8794;
  cursor: pointer;
}
.modal__body { padding: 20px; overflow-y: auto; }
</style>
