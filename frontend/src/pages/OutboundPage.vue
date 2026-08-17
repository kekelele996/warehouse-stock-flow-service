<template>
  <div class="page">
    <div class="page__head">
      <div>
        <div class="page__title">出库管理</div>
        <div class="page__sub">创建出库单 → 拣货 → 复核 → 打包 → 发货 → 完成</div>
      </div>
      <button v-if="can('outboundManage')" class="btn btn--primary" @click="openCreate">+ 创建出库单</button>
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
      <template #status="{ row }"><StatusBadge :status="(row as Order).status" :text="outboundText((row as Order).status)" /></template>
      <template #receiver="{ row }">{{ (row as Order).receiver_name }}</template>
      <template #required="{ row }">{{ formatDate((row as Order).required_ship_date) }}</template>
      <template #tracking="{ row }"><span class="mono">{{ (row as Order).tracking_no || '-' }}</span></template>
      <template #actions="{ row }">
        <div class="row-actions">
          <button class="btn btn--ghost" @click="openDetail((row as Order).id)">详情</button>
          <button v-if="(row as Order).status === 'Pending'" class="btn" @click="pick((row as Order))">拣货</button>
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

    <!-- 创建出库单 -->
    <Modal :open="showCreate" title="创建出库单" width="800px" @close="showCreate = false">
      <div class="grid grid--2">
        <div class="form-row">
          <label class="form-row__label">货主</label>
          <select v-model="createForm.owner_id" class="select" @change="onOwnerChange">
            <option :value="0" disabled>请选择货主</option>
            <option v-for="o in owners" :key="o.id" :value="o.id">{{ o.name }}</option>
          </select>
        </div>
        <div class="form-row">
          <label class="form-row__label">收货方</label>
          <input v-model="createForm.receiver_name" class="input" placeholder="收货方名称" />
        </div>
        <div class="form-row">
          <label class="form-row__label">收货地址</label>
          <input v-model="createForm.receiver_address" class="input" placeholder="收货地址" />
        </div>
        <div class="form-row">
          <label class="form-row__label">要求发货日期</label>
          <input v-model="createForm.required_ship_date" type="date" class="input" />
        </div>
      </div>

      <div class="sub-head">
        <span class="sub-head__title">出库明细（系统推荐库位）</span>
        <button class="btn" @click="addItemRow">+ 添加明细</button>
      </div>

      <div v-for="(item, idx) in createForm.items" :key="idx" class="item-row">
        <select v-model="item.product_id" class="select item-row__product" @change="onProductChange(item)">
          <option :value="0" disabled>选择商品</option>
          <option v-for="p in productsFor(createForm.owner_id)" :key="p.id" :value="p.id">
            {{ p.name }}（{{ p.sku }}）
          </option>
        </select>
        <select v-model="item.bin_location_id" class="select item-row__bin">
          <option :value="0" disabled>{{ item.bin_location_id ? '库位' : '选择/推荐库位' }}</option>
          <option v-for="b in binsFor(item.product_id)" :key="b.id" :value="b.id">{{ b.code }}</option>
        </select>
        <input v-model.number="item.expected_qty" type="number" min="1" class="input item-row__qty" placeholder="数量" />
        <button class="btn btn--ghost" @click="removeItemRow(idx)">✕</button>
      </div>

      <div class="modal-actions">
        <button class="btn btn--primary" :disabled="creating" @click="submitCreate">
          {{ creating ? '提交中…' : '提交' }}
        </button>
      </div>
    </Modal>

    <!-- 出库单详情与流程操作 -->
    <Modal :open="showDetail" :title="`出库单 ${current?.order_no || ''}`" width="880px" @close="showDetail = false">
      <template v-if="current">
        <StepIndicator :steps="OUTBOUND_STEPS" :current="flow.currentStep.value" />
        <div class="detail-grid">
          <div class="detail-item"><span>货主</span><b>{{ current.owner_name }}</b></div>
          <div class="detail-item"><span>收货方</span><b>{{ current.receiver_name }}</b></div>
          <div class="detail-item"><span>地址</span><b>{{ current.receiver_address || '-' }}</b></div>
          <div class="detail-item"><span>要求发货</span><b>{{ formatDate(current.required_ship_date) }}</b></div>
          <div class="detail-item"><span>拣货员</span><b>{{ current.picker_name || '-' }}</b></div>
          <div class="detail-item"><span>复核员</span><b>{{ current.checker_name || '-' }}</b></div>
        </div>

        <div class="sub-head"><span class="sub-head__title">明细（拣货数量）</span></div>
        <DataTable :columns="detailColumns" :rows="current.items">
          <template #product_name="{ row }">
            <div>{{ (row as Item).product_name }}</div>
            <div class="muted mono">{{ (row as Item).sku }}</div>
          </template>
          <template #bin_code="{ row }"><span class="mono">{{ (row as Item).bin_code }}</span></template>
          <template #expected_qty="{ row }">{{ (row as Item).expected_qty }}</template>
          <template #actual_qty="{ row }">
            <input v-if="flow.canPicking" v-model.number="pickForm[String((row as Item).id)]" type="number" min="0" class="input cell-input" :placeholder="String((row as Item).expected_qty)" />
            <span v-else>{{ (row as Item).actual_qty || '-' }}</span>
          </template>
        </DataTable>

        <div v-if="flow.canShip" class="ship-row">
          <div class="form-row ship-row__input">
            <label class="form-row__label">快递单号</label>
            <input v-model="trackingNo" class="input" placeholder="例如 SF1234567890" />
          </div>
        </div>

        <div class="modal-actions flow-actions">
          <button v-if="flow.canPicking" class="btn btn--primary" :disabled="acting" @click="submitPicking">提交拣货</button>
          <button v-if="flow.canChecking" class="btn btn--primary" :disabled="acting" @click="checking(current)">确认复核</button>
          <button v-if="flow.canPacking" class="btn btn--primary" :disabled="acting" @click="packing(current)">确认打包</button>
          <button v-if="flow.canShip" class="btn btn--primary" :disabled="acting" @click="ship(current)">确认发货</button>
          <button v-if="flow.canComplete" class="btn btn--primary" :disabled="acting" @click="complete(current)">完成出库</button>
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
import { useOutboundStore } from '../stores/outbound'
import { useOwnerStore } from '../stores/owner'
import { useProductStore } from '../stores/product'
import { useBinLocationStore } from '../stores/binLocation'
import { useAuth } from '../hooks/useAuth'
import { OUTBOUND_STEPS, useOutboundFlow } from '../hooks/useInboundFlow'
import { OUTBOUND_STATUSES } from '../types/enums'
import type { OutboundOrderView, OutboundItemView } from '../api/outbound'
import { formatDate } from '../utils/format'
import { toast } from '../utils/toast'
import type { Column } from '../components/DataTable.vue'

type Order = OutboundOrderView
interface Item extends OutboundItemView {}

const outboundStore = useOutboundStore()
const ownerStore = useOwnerStore()
const productStore = useProductStore()
const binStore = useBinLocationStore()
const { can } = useAuth()

const orders = computed(() => outboundStore.orders as Order[])
const total = computed(() => outboundStore.total)
const owners = computed(() => ownerStore.owners)
const bins = computed(() => binStore.bins)
const current = computed(() => outboundStore.current as Order | null)

const page = ref(1)
const pageSize = 10
const activeTab = ref('')
const showCreate = ref(false)
const showDetail = ref(false)
const creating = ref(false)
const acting = ref(false)
const trackingNo = ref('')
const pickForm = reactive<Record<string, number>>({})

const tabs = [
  { key: '', label: '全部' },
  ...OUTBOUND_STATUSES.map((s) => ({ key: s, label: textOf(s) })),
]

const columns: Column[] = [
  { key: 'order_no', label: '单号' },
  { key: 'owner_name', label: '货主' },
  { key: 'receiver', label: '收货方' },
  { key: 'status', label: '状态' },
  { key: 'required', label: '要求发货' },
  { key: 'tracking', label: '快递单号' },
  { key: 'actions', label: '操作' },
]

const detailColumns: Column[] = [
  { key: 'product_name', label: '商品' },
  { key: 'bin_code', label: '库位' },
  { key: 'expected_qty', label: '应拣数量' },
  { key: 'actual_qty', label: '实拣数量' },
]

const createForm = reactive({
  owner_id: 0,
  receiver_name: '',
  receiver_address: '',
  required_ship_date: '',
  items: [] as { product_id: number; bin_location_id: number; expected_qty: number }[],
})

const flow = useOutboundFlow(() => outboundStore.current as OutboundOrderView | null)

function outboundText(s: string) {
  return textOf(s)
}

function textOf(s: string) {
  const map: Record<string, string> = {
    Pending: '待拣货', Picking: '拣货中', Checking: '复核中', Packing: '打包中', Shipped: '已发货', Completed: '已完成',
  }
  return map[s] || s
}

function productsFor(ownerId: number) {
  return productStore.products.filter((p) => p.owner_id === ownerId)
}

function binsFor(productId: number) {
  if (!productId) return bins.value
  const product = productStore.products.find((p) => p.id === productId)
  if (!product) return bins.value
  return bins.value.filter((b) => b.storage_requirement === product.storage_requirement && b.status !== 'Maintenance')
}

function addItemRow() {
  createForm.items.push({ product_id: 0, bin_location_id: 0, expected_qty: 1 })
}

function removeItemRow(idx: number) {
  createForm.items.splice(idx, 1)
}

function onOwnerChange() {
  createForm.items = []
}

async function onProductChange(item: { product_id: number; bin_location_id: number; expected_qty: number }) {
  item.bin_location_id = 0
  if (!item.product_id) return
  try {
    const bin = await binStore.recommend(item.product_id, item.expected_qty || 1)
    item.bin_location_id = bin.id
    toast(`已推荐库位 ${bin.code}`)
  } catch {
    item.bin_location_id = 0
  }
}

async function load() {
  const params: Record<string, unknown> = { page: page.value, page_size: pageSize }
  if (activeTab.value) params.status = activeTab.value
  await outboundStore.list(params)
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
  createForm.receiver_name = ''
  createForm.receiver_address = ''
  createForm.required_ship_date = ''
  createForm.items = [{ product_id: 0, bin_location_id: 0, expected_qty: 1 }]
  showCreate.value = true
}

async function submitCreate() {
  if (!createForm.owner_id) return toast('请选择货主', 'error')
  if (!createForm.receiver_name) return toast('请填写收货方', 'error')
  const items = createForm.items.filter((i) => i.product_id > 0 && i.bin_location_id > 0 && i.expected_qty > 0)
  if (!items.length) return toast('请添加至少一条有效明细并确认库位', 'error')
  creating.value = true
  try {
    await outboundStore.create({
      owner_id: createForm.owner_id,
      receiver_name: createForm.receiver_name,
      receiver_address: createForm.receiver_address,
      required_ship_date: createForm.required_ship_date || null,
      items,
    })
    toast('出库单创建成功')
    showCreate.value = false
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    creating.value = false
  }
}

async function openDetail(id: number) {
  await outboundStore.fetch(id)
  for (const key of Object.keys(pickForm)) delete pickForm[key]
  trackingNo.value = ''
  showDetail.value = true
}

async function pick(order: Order) {
  const detail = await outboundStore.fetch(order.id)
  if (!detail) return
  const items = detail.items.map((it) => ({ item_id: it.id, actual_qty: it.expected_qty }))
  acting.value = true
  try {
    await outboundStore.picking(order.id, items)
    toast('拣货完成')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    acting.value = false
  }
}

async function submitPicking() {
  if (!current.value) return
  const items = current.value.items.map((it) => ({
    item_id: it.id,
    actual_qty: pickForm[String(it.id)] ?? it.expected_qty,
  }))
  acting.value = true
  try {
    await outboundStore.picking(current.value.id, items)
    toast('拣货完成')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    acting.value = false
  }
}

async function checking(order: Order) {
  acting.value = true
  try {
    await outboundStore.checking(order.id)
    toast('复核完成')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    acting.value = false
  }
}

async function packing(order: Order) {
  acting.value = true
  try {
    await outboundStore.packing(order.id)
    toast('打包完成')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    acting.value = false
  }
}

async function ship(order: Order) {
  if (!trackingNo.value) return toast('请填写快递单号', 'error')
  acting.value = true
  try {
    await outboundStore.ship(order.id, trackingNo.value)
    toast('发货完成')
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
    await outboundStore.complete(order.id)
    toast('出库单已完成')
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
.item-row__product { flex: 1.4; }
.item-row__bin { flex: 1; }
.item-row__qty { width: 90px; }
.muted { color: var(--text-muted); }
.cell-input { width: 110px; }
.modal-actions { display: flex; justify-content: flex-end; margin-top: 18px; }
.flow-actions { gap: 10px; }
.ship-row { display: flex; justify-content: flex-end; margin-top: 14px; }
.ship-row__input { width: 320px; }
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
