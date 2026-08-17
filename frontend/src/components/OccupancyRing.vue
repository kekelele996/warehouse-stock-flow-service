<template>
  <div class="ring" :style="{ width: size + 'px', height: size + 'px' }">
    <svg :width="size" :height="size" viewBox="0 0 36 36">
      <circle cx="18" cy="18" r="15.9" fill="none" class="ring__bg" stroke-width="3.4" />
      <circle
        cx="18" cy="18" r="15.9" fill="none"
        class="ring__fg"
        stroke-width="3.4"
        :stroke="color"
        :stroke-dasharray="`${clamped} 100`"
        stroke-linecap="round"
        transform="rotate(-90 18 18)"
      />
    </svg>
    <div class="ring__label">
      <div class="ring__value">{{ value }}</div>
      <div class="ring__caption">{{ caption }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{ value: number; caption?: string; size?: number; color?: string }>(),
  { caption: '占用率', size: 120, color: '#2563eb' }
)

const clamped = computed(() => Math.max(0, Math.min(100, props.value)))
const value = computed(() => `${clamped.value.toFixed(1)}%`)
</script>

<style scoped>
.ring {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}
.ring__bg { stroke: #e8edf3; }
.ring__fg { transition: stroke-dasharray 0.4s ease; }
.ring__label {
  position: absolute;
  text-align: center;
}
.ring__value {
  font-size: 18px;
  font-weight: 700;
  color: #102a43;
}
.ring__caption {
  font-size: 11px;
  color: #7a8794;
}
</style>
