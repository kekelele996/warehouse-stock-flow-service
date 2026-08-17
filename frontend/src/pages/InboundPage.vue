<template>
  <div class="page">
    <div class="page__head">
      <div>
        <div class="page__title">入库管理</div>
        <div class="page__sub">创建入库单 → 收货 → 质检 → 上架 → 完成</div>
      </div>
      <button v-if="can('inboundManage')" class="btn btn--primary" @click="openCreate">+ 创建入库单</button>
    </div>

    <div class="tabs">
      <div
        v-for="tab in tabs"
        :key="tab.key"
        class="tab"
        :class="{ 'tab--active': activeTab === tab.key }"
        @click="switchTab(tab.key)"
      >
        {{ tab.label }}
      </div>
    </div>

    <DataTable :columns="columns" :rows="orders">
      <template #order_no="{ row }"><span class="mono">{{ (row as Order).order_no }}</span></template>
      <template #status="{ row }"><StatusBadge :status="(row as Order).status" /></template>
      <template #expected="{ row }">{{ formatDate((row as Order).expected_arrival_date) }}</template>
      <template #actual="{ row }">{{ formatDate((row as Order).actual_arrival_date) }}</template>
      <template #created_at="{ row }">{{ formatDateTime((row as Order).created_at) }}</template>
      <template #actions="{ row }">
        <div class="row-actions">
          <button class="btn btn--ghost" @click="openDetail((row as Order).id)">详情</button>
          <button v-if="(row as Order).status === 'Pending'" class="btn" @click="receive((row as Order))">收货</button>
        </div>
      </template>
    </DataTable>

    <div class="pager">
      <span>共 {{ total }} 条</span>
      <div class="pager__btns">
        <button class="btn" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
        <button class="btn" :disabled="page * pageSize >= total" @click="changePage(page + 1)">下一页</button>
      </div>
    </div>

    <!-- 创建入库单 -->
    <Modal :open="showCreate" title="创建入库单" width="760px" @close="showCreate = false">
      <div class="grid grid--2">
        <div class="form-row">
          <label class="form-row__label">货主</label>
          <select v-model="createForm.owner_id" class="select">
            <option :value="0" disabled>请选择货主</option>
            <option v-for="o in owners" :key="o.id" :value="o.id">{{ o.name }}</option>
          </select>
        </div>
        <div class="form-row">
          <label class="form-row__label">供应商</label>
          <input v-model="createForm.supplier_name" class="input" placeholder="供应商名称" />
        </div>
        <div class="form-row">
          <label class="form-row__label">预计到货日期</label>
          <input v-model="createForm.expected_arrival_date" type="date" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">备注</label>
          <input v-model="createForm.remark" class="input" placeholder="备注" />
        </div>
      </div>

      <div class="sub-head">
        <span class="sub-head__title">入库明细</span>
        <button class="btn" @click="addItemRow">+ 添加明细</button>
      </div>

      <div v-for="(item, idx) in createForm.items" :key="idx" class="item-row">
        <select v-model="item.product_id" class="select item-row__product">
          <option :value="0" disabled>选择商品</option>
          <option v-for="p in productsFor(createForm.owner_id)" :key="p.id" :value="p.id">
            {{ p.name }}（{{ p.sku }}）
          </option>
        </select>
        <input v-model="item.batch_no" class="input" placeholder="批次号" />
        <input v-model.number="item.expected_qty" type="number" min="1" class="input item-row__qty" placeholder="数量" />
        <button class="btn btn--ghost" @click="removeItemRow(idx)">✕</button>
      </div>

      <div class="modal-actions">
        <button class="btn btn--primary" :disabled="creating" @click="submitCreate">
          {{ creating ? '提交中…' : '提交' }}
        </button>
      </div>
    </Modal>

    <!-- 入库单详情与流程操作 -->
    <Modal :open="showDetail" :title="`入库单 ${current?.order_no || ''}`" width="860px" @close="showDetail = false">
      <template v-if="current">
        <StepIndicator :steps="INBOUND_STEPS" :current="flow.current.value" />
        <div class="detail-grid">
          <div class="detail-item"><span>货主</span><b>{{ current.owner_name }}</b></div>
          <div class="detail-item"><span>供应商</span><b>{{ current.supplier_name || '-' }}</b></div>
          <div class="detail-item"><span>预计到货</span><b>{{ formatDate(current.expected_arrival_date) }}</b></div>
          <div class="detail-item"><span>实际到货</span><b>{{ formatDate(current.actual_arrival_date) }}</b></div>
          <div class="detail-item"><span>质检员</span><b>{{ current.qc_inspector_name || '-' }}</b></div>
          <div class="detail-item"><span>库管员</span><b>{{ current.keeper_name || '-' }}</b></div>
        </div>

        <div class="sub-head"><span class="sub-head__title">明细（质检/上架操作）</span></div>
        <DataTable :columns="detailColumns" :rows="current.items">
          <template #product_name="{ row }">
            <div>{{ (row as Item).product_name }}</div>
            <div class="muted mono">{{ (row as Item).sku }}</div>
          </template>
          <template #expected_qty="{ row }">{{ (row as Item).expected_qty }}</template>
          <template #actual_qty="{ row }">
            <input v-if="flow.canQC" v-model.number="qcForm[String((row as Item).id)]" type="number" min="0" class="input cell-input" :placeholder="String((row as Item).expected_qty)" />
            <span v-else>{{ (row as Item).actual_qty || '-' }}</span>
          </template>
          <template #qc_result="{ row }">
            <select v-if="flow.canQC" v-model="qcResultForm[String((row as Item).id)]" class="select cell-input">
              <option value="Pass">合格</option>
              <option value="Fail">不合格</option>
              <option value="Partial">部分合格</option>
            </select>
            <StatusBadge v-else-if="(row as Item).qc_result" :status="(row as Item).qc_result" />
            <span v-else>-</span>
          </template>
          <template #bin_code="{ row }">
            <select v-if="flow.canShelve" v-model="shelveForm[String((row as Item).id)]" class="select cell-input">
              <option :value="0" disabled>选择库位</option>
              <option v-for="b in bins" :key="b.id" :value="b.id">{{ b.code }}（{{ b.status_text }}）</option>
            </select>
            <span v-else>{{ (row as Item).bin_code || '-' }}</span>
          </template>
        </DataTable>

        <div class="modal-actions flow-actions">
          <button v-if="flow.canReceive" class="btn btn--primary" :disabled="acting" @click="receive(current)">确认收货</button>
          <button v-if="flow.canQC" class="btn btn--primary" :disabled="acting" @click="submitQC">提交质检</button>
          <button v-if="flow.canShelve" class="btn btn--primary" :disabled="acting" @click="submitShelve">提交上架</button>
          <button v-if="flow.canComplete" class="btn btn--primary" :disabled="acting" @click="complete(current)">完成入库</button>
        </div>
      </template>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import DataTable from '../components/DataTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import StepIndicator from '../components/StepIndicator.vue'
import Modal from '../components/Modal.vue'
import { useInboundStore } from '../stores/inbound'
import { useOwnerStore } from '../stores/owner'
import { useProductStore } from '../stores/product'
import { useBinLocationStore } from '../stores/binLocation'
import { useAuth } from '../hooks/useAuth'
import { INBOUND_STEPS, useInboundFlow } from '../hooks/useInboundFlow'
import { INBOUND_STATUSES } from '../types/enums'
import type { InboundOrderView, InboundItemView } from '../api/inbound'
import { formatDate, formatDateTime } from '../utils/format'
import { toast } from '../utils/toast'
import type { Column } from '../components/DataTable.vue'

type Order = InboundOrderView
interface Item extends InboundItemView {}

const inboundStore = useInboundStore()
const ownerStore = useOwnerStore()
const productStore = useProductStore()
const binStore = useBinLocationStore()
const { can } = useAuth()

const orders = computed(() => inboundStore.orders as Order[])
const total = computed(() => inboundStore.total)
const owners = computed(() => ownerStore.owners)
const bins = computed(() => binStore.bins)
const current = computed(() => inboundStore.current as Order | null)

const page = ref(1)
const pageSize = 10
const activeTab = ref('')
const showCreate = ref(false)
const showDetail = ref(false)
const creating = ref(false)
const acting = ref(false)
const qcForm = reactive<Record<string, number>>({})
const qcResultForm = reactive<Record<string, string>>({})
const shelveForm = reactive<Record<string, number>>({})

const tabs = [
  { key: '', label: '全部' },
  ...INBOUND_STATUSES.map((s) => ({ key: s, label: textOf(s) })),
]

const columns: Column[] = [
  { key: 'order_no', label: '单号' },
  { key: 'owner_name', label: '货主' },
  { key: 'supplier_name', label: '供应商' },
  { key: 'status', label: '状态' },
  { key: 'expected', label: '预计到货' },
  { key: 'actual', label: '实际到货' },
  { key: 'created_at', label: '创建时间' },
  { key: 'actions', label: '操作' },
]

const detailColumns: Column[] = [
  { key: 'product_name', label: '商品' },
  { key: 'batch_no', label: '批次号' },
  { key: 'expected_qty', label: '预期数量' },
  { key: 'actual_qty', label: '实收数量' },
  { key: 'qc_result', label: '质检结果' },
  { key: 'bin_code', label: '库位' },
]

const createForm = reactive({
  owner_id: 0,
  supplier_name: '',
  expected_arrival_date: '',
  remark: '',
  items: [] as { product_id: number; batch_no: string; expected_qty: number }[],
})

const flow = useInboundFlow(() => inboundStore.current as InboundOrderView | null)

function textOf(s: string) {
  const map: Record<string, string> = {
    Pending: '待收货', Received: '已收货', QCInProgress: '质检中', Shelved: '已上架', Completed: '已完成',
  }
  return map[s] || s
}

function productsFor(ownerId: number) {
  return productStore.products.filter((p) => p.owner_id === ownerId)
}

function addItemRow() {
  createForm.items.push({ product_id: 0, batch_no: '', expected_qty: 1 })
}

function removeItemRow(idx: number) {
  createForm.items.splice(idx, 1)
}

async function load() {
  const params: Record<string, unknown> = { page: page.value, page_size: pageSize }
  if (activeTab.value) params.status = activeTab.value
  await inboundStore.list(params)
}

async function loadRefs() {
  await Promise.all([
    ownerStore.list({ page: 1, page_size: 100 }),
    productStore.list({ page: 1, page_size: 100 }),
    binStore.list({ page: 1, page_size: 100 }),
  ])
}

function switchTab(key: string) {
  activeTab.value = key
  page.value = 1
  load()
}

function changePage(p: number) {
  page.value = p
  load()
}

function openCreate() {
  createForm.owner_id = 0
  createForm.supplier_name = ''
  createForm.expected_arrival_date = ''
  createForm.remark = ''
  createForm.items = [{ product_id: 0, batch_no: '', expected_qty: 1 }]
  showCreate.value = true
}

async function submitCreate() {
  if (!createForm.owner_id) return toast('请选择货主', 'error')
  const items = createForm.items.filter((i) => i.product_id > 0 && i.expected_qty > 0)
  if (!items.length) return toast('请添加至少一条有效的商品明细', 'error')
  creating.value = true
  try {
    await inboundStore.create({
      owner_id: createForm.owner_id,
      supplier_name: createForm.supplier_name,
      expected_arrival_date: createForm.expected_arrival_date || null,
      remark: createForm.remark,
      items,
    })
    toast('入库单创建成功')
    showCreate.value = false
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    creating.value = false
  }
}

async function openDetail(id: number) {
  await inboundStore.fetch(id)
  resetForms()
  showDetail.value = true
}

function resetForms() {
  for (const key of Object.keys(qcForm)) delete qcForm[key]
  for (const key of Object.keys(qcResultForm)) delete qcResultForm[key]
  for (const key of Object.keys(shelveForm)) delete shelveForm[key]
}

async function receive(order: Order) {
  acting.value = true
  try {
    await inboundStore.receive(order.id)
    toast('收货完成')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    acting.value = false
  }
}

async function submitQC() {
  if (!current.value) return
  const items = current.value.items.map((it) => ({
    item_id: it.id,
    actual_qty: qcForm[String(it.id)] ?? it.expected_qty,
    qc_result: qcResultForm[String(it.id)] || 'Pass',
  }))
  acting.value = true
  try {
    await inboundStore.qc(current.value.id, items)
    toast('质检完成')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    acting.value = false
  }
}

async function submitShelve() {
  if (!current.value) return
  const items = current.value.items
    .filter((it) => shelveForm[String(it.id)])
    .map((it) => ({ item_id: it.id, bin_location_id: shelveForm[String(it.id)] }))
  if (!items.length) return toast('请为每个明细指派库位', 'error')
  acting.value = true
  try {
    await inboundStore.shelve(current.value.id, items)
    toast('上架完成')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    acting.value = false
  }
}

async function complete(order: Order) {
  acting.value = true
  try {
    await inboundStore.complete(order.id)
    toast('入库单已完成')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    acting.value = false
  }
}

onMounted(async () => {
  await loadRefs()
  await load()
})
</script>

<style scoped>
.row-actions { display: flex; gap: 6px; }
.sub-head { display: flex; align-items: center; justify-content: space-between; margin: 16px 0 10px; }
.sub-head__title { font-size: 14px; font-weight: 700; color: var(--text-main); }
.item-row { display: flex; gap: 8px; margin-bottom: 8px; align-items: center; }
.item-row__product { flex: 1; }
.item-row__qty { width: 90px; }
.muted { color: var(--text-muted); }
.cell-input { width: 110px; }
.modal-actions { display: flex; justify-content: flex-end; margin-top: 18px; }
.flow-actions { gap: 10px; }
.detail-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin: 16px 0;
  background: #f8fafc;
  border-radius: 10px;
  padding: 14px 16px;
}
.detail-item { font-size: 13px; color: var(--text-muted); }
.detail-item b { display: block; color: var(--text-main); margin-top: 2px; font-weight: 600; }
</style>
