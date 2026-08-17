<template>
  <div class="steps">
    <div
      v-for="(step, idx) in steps"
      :key="step.key"
      class="step"
      :class="{ 'step--done': idx < current, 'step--active': idx === current }"
    >
      <div class="step__dot">{{ idx < current ? '✓' : idx + 1 }}</div>
      <span class="step__label">{{ step.label }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FlowStep } from '../hooks/useInboundFlow'

defineProps<{ steps: FlowStep[]; current: number }>()
</script>

<style scoped>
.steps {
  display: flex;
  align-items: center;
  gap: 4px;
  overflow-x: auto;
}
.step {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 8px;
  color: #7a8794;
  font-size: 13px;
  white-space: nowrap;
}
.step__dot {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: #eef1f4;
  color: #7a8794;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
}
.step--active .step__dot { background: #2563eb; color: #fff; }
.step--active { color: #2563eb; font-weight: 600; }
.step--done .step__dot { background: #16a34a; color: #fff; }
.step--done { color: #16a34a; }
</style>
