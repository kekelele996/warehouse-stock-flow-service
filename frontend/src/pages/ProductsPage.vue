<template>
  <div class="page">
    <div class="page__head">
      <div>
        <div class="page__title">商品管理</div>
        <div class="page__sub">维护货主商品档案（SKU/品类/存储要求/单价）</div>
      </div>
      <button class="btn btn--primary" @click="openCreate">+ 创建商品</button>
    </div>

    <div class="toolbar">
      <input v-model="keyword" class="input toolbar__input" placeholder="搜索名称/SKU/条码" @keyup.enter="search" />
      <select v-model="storageFilter" class="select toolbar__select">
        <option value="">全部存储要求</option>
        <option v-for="s in STORAGE_REQUIREMENTS" :key="s" :value="s">{{ storageText(s) }}</option>
      </select>
      <button class="btn" @click="search">查询</button>
    </div>

    <DataTable :columns="columns" :rows="products">
      <template #owner_name="{ row }">{{ (row as Product).owner_name }}</template>
      <template #storage="{ row }"><StatusBadge :status="(row as Product).storage_requirement" /></template>
      <template #price="{ row }">¥{{ formatAmount((row as Product).price) }}</template>
    </DataTable>

    <div class="pager">
      <span>共 {{ total }} 条</span>
      <div class="pager__btns">
        <button class="btn" :disabled="page <= 1" @click="changePage(page - 1)">上一页</button>
        <button class="btn" :disabled="page * pageSize >= total" @click="changePage(page + 1)">下一页</button>
      </div>
    </div>

    <!-- 创建商品 -->
    <Modal :open="showCreate" title="创建商品" width="720px" @close="showCreate = false">
      <div class="grid grid--2">
        <div class="form-row">
          <label class="form-row__label">所属货主 *</label>
          <select v-model="createForm.owner_id" class="select">
            <option :value="0" disabled>请选择货主</option>
            <option v-for="o in owners" :key="o.id" :value="o.id">{{ o.name }}</option>
          </select>
        </div>
        <div class="form-row">
          <label class="form-row__label">商品名称 *</label>
          <input v-model="createForm.name" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">SKU *</label>
          <input v-model="createForm.sku" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">条形码</label>
          <input v-model="createForm.barcode" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">品类</label>
          <input v-model="createForm.category" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">规格</label>
          <input v-model="createForm.spec" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">单位 *</label>
          <input v-model="createForm.unit" class="input" placeholder="袋/箱/桶/件" />
        </div>
        <div class="form-row">
          <label class="form-row__label">存储要求 *</label>
          <select v-model="createForm.storage_requirement" class="select">
            <option v-for="s in STORAGE_REQUIREMENTS" :key="s" :value="s">{{ storageText(s) }}</option>
          </select>
        </div>
        <div class="form-row">
          <label class="form-row__label">保质期（天，可选）</label>
          <input v-model.number="createForm.shelf_life_days" type="number" min="0" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">单价（元）</label>
          <input v-model.number="createForm.price" type="number" min="0" step="0.01" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">体积（m³）</label>
          <input v-model.number="createForm.volume" type="number" min="0" step="0.001" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">重量（kg）</label>
          <input v-model.number="createForm.weight" type="number" min="0" step="0.001" class="input" />
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn btn--primary" :disabled="creating" @click="submitCreate">{{ creating ? '提交中…' : '提交' }}</button>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import DataTable from '../components/DataTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import Modal from '../components/Modal.vue'
import { useProductStore } from '../stores/product'
import { useOwnerStore } from '../stores/owner'
import { STORAGE_REQUIREMENTS } from '../types/enums'
import type { ProductView } from '../api/product'
import { formatAmount } from '../utils/format'
import { toast } from '../utils/toast'
import type { Column } from '../components/DataTable.vue'

type Product = ProductView

const productStore = useProductStore()
const ownerStore = useOwnerStore()

const products = computed(() => productStore.products)
const owners = computed(() => ownerStore.owners)
const total = ref(0)
const page = ref(1)
const pageSize = 10
const keyword = ref('')
const storageFilter = ref('')
const showCreate = ref(false)
const creating = ref(false)

const createForm = reactive({
  owner_id: 0,
  name: '',
  sku: '',
  barcode: '',
  category: '',
  spec: '',
  unit: '',
  shelf_life_days: null as number | null,
  storage_requirement: 'Normal',
  volume: 0,
  weight: 0,
  price: 0,
})

const columns: Column[] = [
  { key: 'name', label: '商品名称' },
  { key: 'sku', label: 'SKU' },
  { key: 'owner_name', label: '货主' },
  { key: 'category', label: '品类' },
  { key: 'spec', label: '规格' },
  { key: 'storage', label: '存储要求' },
  { key: 'price', label: '单价' },
]

function storageText(s: string) {
  const map: Record<string, string> = { Normal: '常温', ColdChain: '冷链', Dangerous: '危险品', Fragile: '易碎' }
  return map[s] || s
}

async function load() {
  const params: Record<string, unknown> = { page: page.value, page_size: pageSize }
  if (keyword.value) params.keyword = keyword.value
  if (storageFilter.value) params.storage_requirement = storageFilter.value
  const res = await productStore.list(params)
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

function openCreate() {
  Object.assign(createForm, {
    owner_id: 0, name: '', sku: '', barcode: '', category: '', spec: '', unit: '',
    shelf_life_days: null, storage_requirement: 'Normal', volume: 0, weight: 0, price: 0,
  })
  showCreate.value = true
}

async function submitCreate() {
  if (!createForm.owner_id) return toast('请选择货主', 'error')
  if (!createForm.name || !createForm.sku || !createForm.unit) return toast('请填写商品名称、SKU 和单位', 'error')
  creating.value = true
  try {
    await productStore.create(createForm)
    toast('商品创建成功')
    showCreate.value = false
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    creating.value = false
  }
}

onMounted(async () => {
  await ownerStore.list({ page: 1, page_size: 100 })
  await load()
})
</script>

<style scoped>
.toolbar { display: flex; gap: 10px; margin-bottom: 16px; align-items: center; }
.toolbar__input { max-width: 280px; }
.toolbar__select { max-width: 180px; }
.modal-actions { display: flex; justify-content: flex-end; margin-top: 18px; }
</style>
