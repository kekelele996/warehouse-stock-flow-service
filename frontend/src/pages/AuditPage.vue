<template>
  <div class="page">
    <div class="page__head">
      <div>
        <div class="page__title">操作日志</div>
        <div class="page__sub">收货、质检、上架、拣货、复核、发货等全链路操作记录</div>
      </div>
      <button class="btn" :disabled="loading" @click="load">刷新</button>
    </div>

    <div class="toolbar">
      <select v-model="moduleFilter" class="select toolbar__select">
        <option value="">全部模块</option>
        <option v-for="m in modules" :key="m" :value="m">{{ moduleText(m) }}</option>
      </select>
      <button class="btn" @click="search">查询</button>
    </div>

    <DataTable :columns="columns" :rows="logs">
      <template #created_at="{ row }">{{ formatDateTime((row as Audit).created_at) }}</template>
      <template #role="{ row }"><StatusBadge :status="(row as Audit).role" :text="(row as Audit).role_text" /></template>
      <template #detail="{ row }">
        <span class="detail-text" :title="(row as Audit).detail">{{ (row as Audit).detail }}</span>
      </template>
    </DataTable>

    <div class="pager">
      <span>共 {{ total }} 条</span>
      <div class="pager__btns">
        <button class="btn" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
        <button class="btn" :disabled="page * pageSize >= total" @click="changePage(page + 1)">下一页</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import DataTable from '../components/DataTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { useAuditStore } from '../stores/audit'
import type { AuditView } from '../api/audit'
import { formatDateTime } from '../utils/format'
import type { Column } from '../components/DataTable.vue'

type Audit = AuditView

const store = useAuditStore()

const logs = computed(() => store.logs)
const loading = computed(() => store.loading)
const total = ref(0)
const page = ref(1)
const pageSize = 15
const moduleFilter = ref('')

const modules = ['inbound', 'outbound', 'owner', 'product', 'bin', 'auth']

const columns: Column[] = [
  { key: 'created_at', label: '时间' },
  { key: 'username', label: '操作人' },
  { key: 'role', label: '角色' },
  { key: 'module', label: '模块' },
  { key: 'action', label: '动作' },
  { key: 'entity_id', label: '单据/ID' },
  { key: 'detail', label: '详情' },
]

function moduleText(m: string) {
  const map: Record<string, string> = { inbound: '入库', outbound: '出库', owner: '货主', product: '商品', bin: '库位', auth: '认证' }
  return map[m] || m
}

async function load() {
  const params: Record<string, unknown> = { page: page.value, page_size: pageSize }
  if (moduleFilter.value) params.module = moduleFilter.value
  const res = await store.list(params)
  total.value = res.total
}

function search() {
  page.value = 1
  load()
}

function changePage(p: number) {
  page.value = p
  load()
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 10px; margin-bottom: 16px; }
.toolbar__select { max-width: 180px; }
.detail-text { max-width: 360px; display: inline-block; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; vertical-align: bottom; }
</style>
