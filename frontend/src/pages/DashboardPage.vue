<template>
  <div class="page">
    <div class="page__head">
      <div>
        <div class="page__title">仓库总览</div>
        <div class="page__sub">今日收货/发货、库位占用率、货主库存金额与待处理任务</div>
      </div>
      <button class="btn" :disabled="loading" @click="load">刷新</button>
    </div>

    <div v-if="loading" class="loading-mask">加载中…</div>
    <template v-else-if="summary">
      <div class="grid grid--4">
        <StatCard icon="📥" label="今日收货（单）" :value="summary.today_inbound_count" />
        <StatCard icon="📤" label="今日发货（单）" :value="summary.today_outbound_count" />
        <StatCard icon="⏳" label="待处理入库" :value="summary.pending_inbound" />
        <StatCard icon="📦" label="待处理出库" :value="summary.pending_outbound" />
      </div>

      <div class="grid grid--2 dashboard-grid">
        <div class="card">
          <h3 class="card__title">库位占用率</h3>
          <div class="card__ring">
            <OccupancyRing :value="summary.bin_occupancy_rate" />
            <div class="ring-meta">
              <div>总库位：{{ summary.total_bin_count }} 个</div>
              <div>占用库位：{{ summary.occupied_bin_count }} 个</div>
              <div>平均占用率：{{ percent(summary.bin_occupancy_rate) }}</div>
            </div>
          </div>
        </div>

        <div class="card">
          <h3 class="card__title">各货主库存金额 TOP10</h3>
          <div v-if="summary.owner_top10.length" class="top10">
            <div v-for="(item, i) in summary.owner_top10" :key="item.owner_id" class="top10__row">
              <span class="top10__rank">{{ i + 1 }}</span>
              <span class="top10__name">{{ item.owner_name }}</span>
              <div class="top10__bar-wrap">
                <div class="top10__bar" :style="{ width: barWidth(item.total_value) }"></div>
              </div>
              <span class="top10__value">¥{{ formatAmount(item.total_value) }}</span>
            </div>
          </div>
          <EmptyState v-else title="暂无库存数据" />
        </div>
      </div>

      <div class="card">
        <h3 class="card__title">待处理任务</h3>
        <DataTable :columns="taskColumns" :rows="summary.pending_tasks">
          <template #type="{ row }">
            <StatusBadge :status="(row as PendingTask).type === 'inbound' ? 'Received' : 'Picking'" :text="(row as PendingTask).type === 'inbound' ? '入库' : '出库'" />
          </template>
          <template #order_no="{ row }">
            <span class="mono">{{ (row as PendingTask).order_no }}</span>
          </template>
          <template #status="{ row }">
            <StatusBadge :status="(row as PendingTask).status" />
          </template>
        </DataTable>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import StatCard from '../components/StatCard.vue'
import OccupancyRing from '../components/OccupancyRing.vue'
import StatusBadge from '../components/StatusBadge.vue'
import DataTable from '../components/DataTable.vue'
import EmptyState from '../components/EmptyState.vue'
import { useDashboardStore } from '../stores/dashboard'
import { formatAmount, percent } from '../utils/format'
import type { PendingTaskItem } from '../api/dashboard'
import type { Column } from '../components/DataTable.vue'

const store = useDashboardStore()
const summary = computed(() => store.summary)
const loading = computed(() => store.loading)

interface PendingTask extends PendingTaskItem {}

const taskColumns: Column[] = [
  { key: 'type', label: '类型' },
  { key: 'title', label: '任务' },
  { key: 'order_no', label: '单号' },
  { key: 'status', label: '状态' },
]

const maxTopValue = computed(() => {
  const items = summary.value?.owner_top10 || []
  return items.length ? Math.max(...items.map((i) => i.total_value), 1) : 1
})

function barWidth(v: number) {
  return `${Math.max(2, Math.round((v / maxTopValue.value) * 100))}%`
}

function load() {
  store.fetch().catch((e) => console.error(e))
}

onMounted(load)
</script>

<style scoped>
.dashboard-grid { margin-top: 16px; }
.card__title { font-size: 15px; font-weight: 700; margin-bottom: 16px; color: var(--text-main); }
.card__ring { display: flex; align-items: center; gap: 32px; }
.ring-meta { font-size: 13px; color: var(--text-sub); line-height: 2; }
.top10 { display: flex; flex-direction: column; gap: 10px; }
.top10__row { display: flex; align-items: center; gap: 10px; font-size: 13px; }
.top10__rank {
  width: 22px; height: 22px; border-radius: 6px; background: #eef1f4; color: var(--text-muted);
  display: flex; align-items: center; justify-content: center; font-size: 12px; font-weight: 700;
}
.top10__name { width: 120px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: var(--text-sub); }
.top10__bar-wrap { flex: 1; height: 10px; background: #eef1f4; border-radius: 999px; overflow: hidden; }
.top10__bar { height: 100%; background: linear-gradient(90deg, #2563eb, #60a5fa); border-radius: 999px; }
.top10__value { width: 110px; text-align: right; font-weight: 600; color: var(--text-main); }
</style>
