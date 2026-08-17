<template>
  <div class="page">
    <div class="page__head">
      <div>
        <div class="page__title">库位管理</div>
        <div class="page__sub">按区域查看库位网格，颜色表示占用率；支持批量创建</div>
      </div>
      <button v-if="can('binManage')" class="btn btn--primary" @click="showBatch = true">+ 批量创建库位</button>
    </div>

    <div class="tabs">
      <div class="tab" :class="{ 'tab--active': area === '' }" @click="switchArea('')">全部</div>
      <div v-for="a in AREAS" :key="a" class="tab" :class="{ 'tab--active': area === a }" @click="switchArea(a)">{{ a }} 区</div>
    </div>

    <div v-if="loading" class="loading-mask">加载中…</div>
    <div v-else-if="bins.length" class="bin-grid">
      <div
        v-for="bin in bins"
        :key="bin.id"
        class="bin-card"
        :class="`bin-card--${bin.status.toLowerCase()}`"
        @click="openDetail(bin)"
      >
        <div class="bin-card__head">
          <span class="bin-card__code mono">{{ bin.code }}</span>
          <StatusBadge :status="bin.status" />
        </div>
        <div class="bin-card__body">
          <OccupancyRing :value="bin.occupancy_rate" :size="64" :color="ringColor(bin.occupancy_rate)" />
          <div class="bin-card__meta">
            <div>容量 {{ bin.capacity }} m³</div>
            <div>存储：{{ bin.storage_text }}</div>
            <div>占用 {{ bin.occupancy_rate.toFixed(1) }}%</div>
          </div>
        </div>
      </div>
    </div>
    <EmptyState v-else title="当前区域暂无库位" />

    <div class="pager">
      <span>共 {{ total }} 个</span>
      <div class="pager__btns">
        <button class="btn" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
        <button class="btn" :disabled="page * pageSize >= total" @click="changePage(page + 1)">下一页</button>
      </div>
    </div>

    <!-- 库位详情 -->
    <Modal :open="showDetail" :title="`库位 ${current?.code || ''}`" width="560px" @close="showDetail = false">
      <template v-if="current">
        <div class="detail-grid">
          <div class="detail-item"><span>区域</span><b>{{ current.area }} 区</b></div>
          <div class="detail-item"><span>货架/层/列</span><b>{{ current.rack_no }} / {{ current.layer_no }} / {{ current.column_no }}</b></div>
          <div class="detail-item"><span>容量</span><b>{{ current.capacity }} m³</b></div>
          <div class="detail-item"><span>占用率</span><b>{{ current.occupancy_rate.toFixed(1) }}%</b></div>
          <div class="detail-item"><span>存储要求</span><b>{{ current.storage_text }}</b></div>
          <div class="detail-item"><span>状态</span><b><StatusBadge :status="current.status" /></b></div>
        </div>
        <div class="sub-head"><span class="sub-head__title">存放商品</span></div>
        <DataTable :columns="contentColumns" :rows="contents">
          <template #product_name="{ row }">{{ (row as BinContent).product_name }}</template>
          <template #owner_name="{ row }">{{ (row as BinContent).owner_name }}</template>
        </DataTable>
      </template>
    </Modal>

    <!-- 批量创建 -->
    <Modal :open="showBatch" title="批量创建库位" width="600px" @close="showBatch = false">
      <div class="grid grid--2">
        <div class="form-row">
          <label class="form-row__label">区域</label>
          <select v-model="batch.area" class="select">
            <option v-for="a in AREAS" :key="a" :value="a">{{ a }} 区</option>
          </select>
        </div>
        <div class="form-row">
          <label class="form-row__label">货架号</label>
          <input v-model="batch.rack_no" class="input" placeholder="如 R01" />
        </div>
        <div class="form-row">
          <label class="form-row__label">层范围</label>
          <div class="range">
            <input v-model.number="batch.layer_start" type="number" class="input" />
            <span>~</span>
            <input v-model.number="batch.layer_end" type="number" class="input" />
          </div>
        </div>
        <div class="form-row">
          <label class="form-row__label">列范围</label>
          <div class="range">
            <input v-model.number="batch.column_start" type="number" class="input" />
            <span>~</span>
            <input v-model.number="batch.column_end" type="number" class="input" />
          </div>
        </div>
        <div class="form-row">
          <label class="form-row__label">容量（m³）</label>
          <input v-model.number="batch.capacity" type="number" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">存储要求</label>
          <select v-model="batch.storage_requirement" class="select">
            <option v-for="s in STORAGE_REQUIREMENTS" :key="s" :value="s">{{ storageText(s) }}</option>
          </select>
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn btn--primary" :disabled="creating" @click="submitBatch">
          {{ creating ? '创建中…' : '批量创建' }}
        </button>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import StatusBadge from '../components/StatusBadge.vue'
import OccupancyRing from '../components/OccupancyRing.vue'
import DataTable from '../components/DataTable.vue'
import EmptyState from '../components/EmptyState.vue'
import Modal from '../components/Modal.vue'
import { useBinLocationStore } from '../stores/binLocation'
import { useAuth } from '../hooks/useAuth'
import { BIN_AREAS, STORAGE_REQUIREMENTS } from '../types/enums'
import type { BinView, BinContent } from '../api/binLocation'
import { toast } from '../utils/toast'
import type { Column } from '../components/DataTable.vue'

const AREAS = BIN_AREAS

const store = useBinLocationStore()
const { can } = useAuth()

const bins = computed(() => store.bins)
const total = computed(() => store.total)
const loading = computed(() => store.loading)
const current = computed(() => store.current as BinView | null)
const contents = computed(() => store.contents as BinContent[])

const page = ref(1)
const pageSize = 30
const area = ref('')
const showDetail = ref(false)
const showBatch = ref(false)
const creating = ref(false)

const batch = reactive({
  area: 'A',
  rack_no: 'R02',
  layer_start: 1,
  layer_end: 2,
  column_start: 1,
  column_end: 3,
  capacity: 8,
  storage_requirement: 'Normal',
})

const contentColumns: Column[] = [
  { key: 'product_name', label: '商品' },
  { key: 'sku', label: 'SKU' },
  { key: 'batch_no', label: '批次号' },
  { key: 'quantity', label: '数量' },
  { key: 'owner_name', label: '货主' },
]

function storageText(s: string) {
  const map: Record<string, string> = { Normal: '常温', ColdChain: '冷链', Dangerous: '危险品', Fragile: '易碎' }
  return map[s] || s
}

function ringColor(rate: number) {
  if (rate >= 80) return '#dc2626'
  if (rate >= 50) return '#f59e0b'
  return '#16a34a'
}

async function load() {
  const params: Record<string, unknown> = { page: page.value, page_size: pageSize }
  if (area.value) params.area = area.value
  await store.list(params)
}

function switchArea(a: string) {
  area.value = a
  page.value = 1
  load()
}

function changePage(p: number) {
  page.value = p
  load()
}

async function openDetail(bin: BinView) {
  await store.fetch(bin.id)
  await store.fetchContents(bin.id)
  showDetail.value = true
}

async function submitBatch() {
  creating.value = true
  try {
    const res = await store.batchCreate(batch)
    toast(`已创建 ${res.created} 个库位`)
    showBatch.value = false
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    creating.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.bin-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(230px, 1fr));
  gap: 14px;
}
.bin-card {
  background: #fff;
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 14px;
  cursor: pointer;
  transition: box-shadow 0.15s ease, transform 0.15s ease;
}
.bin-card:hover { box-shadow: 0 6px 18px rgba(16, 42, 67, 0.1); transform: translateY(-2px); }
.bin-card--maintenance { opacity: 0.65; }
.bin-card__head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 10px; }
.bin-card__code { font-size: 13px; font-weight: 700; color: var(--text-main); }
.bin-card__body { display: flex; align-items: center; gap: 12px; }
.bin-card__meta { font-size: 12px; color: var(--text-muted); line-height: 1.8; }
.detail-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 14px;
  background: #f8fafc;
  border-radius: 10px;
  padding: 14px 16px;
}
.detail-item { font-size: 13px; color: var(--text-muted); }
.detail-item b { display: block; color: var(--text-main); margin-top: 2px; font-weight: 600; }
.sub-head { margin: 14px 0 10px; }
.sub-head__title { font-size: 14px; font-weight: 700; color: var(--text-main); }
.range { display: flex; align-items: center; gap: 8px; }
.modal-actions { display: flex; justify-content: flex-end; margin-top: 18px; }
</style>
