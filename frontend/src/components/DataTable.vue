<template>
  <div class="table-wrap">
    <table class="data-table">
      <thead>
        <tr>
          <th v-for="col in columns" :key="String(col.key)" class="data-table__th">{{ col.label }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(row, i) in rows" :key="i">
          <td v-for="col in columns" :key="String(col.key)" class="data-table__td">
            <slot :name="String(col.key)" :row="row" :value="(row as Record<string, unknown>)[String(col.key)]">
              {{ (row as Record<string, unknown>)[String(col.key)] }}
            </slot>
          </td>
        </tr>
      </tbody>
    </table>
    <div v-if="!rows.length" class="table-wrap__empty">
      <EmptyState />
    </div>
  </div>
</template>

<script setup lang="ts">
import EmptyState from './EmptyState.vue'

export interface Column {
  key: string | number
  label: string
}

defineProps<{ columns: Column[]; rows: unknown[] }>()
</script>

<style scoped>
.table-wrap {
  overflow-x: auto;
  background: #fff;
  border: 1px solid #e6eaf0;
  border-radius: 12px;
}
.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}
.data-table__th {
  text-align: left;
  padding: 12px 14px;
  color: #5f6b76;
  font-weight: 600;
  background: #f8fafc;
  border-bottom: 1px solid #e6eaf0;
  white-space: nowrap;
}
.data-table__td {
  padding: 12px 14px;
  border-bottom: 1px solid #f0f3f7;
  color: #334e68;
  vertical-align: middle;
}
.table-wrap__empty { border-top: 1px solid #f0f3f7; }
</style>
