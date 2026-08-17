<template>
  <span class="badge" :class="`badge--${color}`">{{ text }}</span>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ENUM_TEXT, STATUS_COLOR } from '../constants/enums'

const props = defineProps<{ status: string; text?: string }>()

const color = computed(() => STATUS_COLOR[props.status] || 'gray')
const text = computed(() => props.text || ENUM_TEXT[props.status] || props.status)
</script>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  line-height: 20px;
  white-space: nowrap;
}
.badge--green { background: #e6f7ed; color: #1a7f4b; }
.badge--red { background: #fdeceb; color: #c0392b; }
.badge--orange { background: #fff4e5; color: #c77700; }
.badge--blue { background: #e8f1fd; color: #1f6fc9; }
.badge--cyan { background: #e4f8f8; color: #0e7c7c; }
.badge--purple { background: #f1ecfb; color: #6c4bb8; }
.badge--gray { background: #eef1f4; color: #5f6b76; }
</style>
