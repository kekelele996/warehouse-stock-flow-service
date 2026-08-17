<template>
  <div class="page">
    <div class="page__head">
      <div>
        <div class="page__title">货主管理</div>
        <div class="page__sub">货主档案、信用额度与账务概况</div>
      </div>
      <button v-if="can('ownerManage')" class="btn btn--primary" @click="openCreate">+ 添加货主</button>
    </div>

    <div class="toolbar">
      <input v-model="keyword" class="input toolbar__input" placeholder="搜索货主名称/联系人/电话" @keyup.enter="search" />
      <select v-model="statusFilter" class="select toolbar__select">
        <option value="">全部状态</option>
        <option value="Active">合作中</option>
        <option value="Suspended">已暂停</option>
      </select>
      <button class="btn" @click="search">查询</button>
    </div>

    <DataTable :columns="columns" :rows="owners">
      <template #name="{ row }">
        <router-link :to="`/owners/${(row as Owner).id}`" class="link">{{ (row as Owner).name }}</router-link>
      </template>
      <template #settlement="{ row }">{{ (row as Owner).settlement_text }}</template>
      <template #credit_limit="{ row }">¥{{ formatAmount((row as Owner).credit_limit) }}</template>
      <template #current_debt="{ row }">¥{{ formatAmount((row as Owner).current_debt) }}</template>
      <template #status="{ row }"><StatusBadge :status="(row as Owner).status" /></template>
      <template #actions="{ row }">
        <div class="row-actions">
          <router-link :to="`/owners/${(row as Owner).id}`" class="btn btn--ghost">详情</router-link>
          <button v-if="can('ownerManage')" class="btn" @click="adjustCredit((row as Owner))">调额度</button>
          <button v-if="can('ownerManage') && (row as Owner).status === 'Active'" class="btn" @click="suspend((row as Owner))">暂停</button>
          <button v-if="can('ownerManage') && (row as Owner).status === 'Suspended'" class="btn" @click="activate((row as Owner))">恢复</button>
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

    <!-- 创建货主 -->
    <Modal :open="showCreate" title="添加货主" width="640px" @close="showCreate = false">
      <div class="grid grid--2">
        <div class="form-row">
          <label class="form-row__label">货主名称 *</label>
          <input v-model="createForm.name" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">联系人 *</label>
          <input v-model="createForm.contact_name" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">联系电话 *</label>
          <input v-model="createForm.phone" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">邮箱</label>
          <input v-model="createForm.email" class="input" />
        </div>
        <div class="form-row">
          <label class="form-row__label">结算方式 *</label>
          <select v-model="createForm.settlement_method" class="select">
            <option value="Monthly">月结</option>
            <option value="PerOrder">单结</option>
            <option value="Prepaid">预付</option>
          </select>
        </div>
        <div class="form-row">
          <label class="form-row__label">信用额度 *</label>
          <input v-model.number="createForm.credit_limit" type="number" min="0" class="input" />
        </div>
        <div class="form-row form-row--full">
          <label class="form-row__label">地址</label>
          <input v-model="createForm.address" class="input" />
        </div>
      </div>
      <div class="modal-actions">
        <button class="btn btn--primary" :disabled="creating" @click="submitCreate">{{ creating ? '提交中…' : '提交' }}</button>
      </div>
    </Modal>

    <!-- 调整额度 -->
    <Modal :open="showCredit" :title="`调整额度 - ${creditOwner?.name || ''}`" width="420px" @close="showCredit = false">
      <div class="form-row">
        <label class="form-row__label">信用额度</label>
        <input v-model.number="creditValue" type="number" min="0" class="input" />
      </div>
      <div class="modal-actions">
        <button class="btn btn--primary" :disabled="creating" @click="submitCredit">保存</button>
      </div>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import DataTable from '../components/DataTable.vue'
import StatusBadge from '../components/StatusBadge.vue'
import Modal from '../components/Modal.vue'
import { useOwnerStore } from '../stores/owner'
import { useAuth } from '../hooks/useAuth'
import type { OwnerView } from '../api/owner'
import { formatAmount } from '../utils/format'
import { toast } from '../utils/toast'
import type { Column } from '../components/DataTable.vue'

type Owner = OwnerView

const store = useOwnerStore()
const { can } = useAuth()

const owners = computed(() => store.owners)
const total = ref(0)
const page = ref(1)
const pageSize = 10
const keyword = ref('')
const statusFilter = ref('')
const showCreate = ref(false)
const showCredit = ref(false)
const creating = ref(false)
const creditOwner = ref<Owner | null>(null)
const creditValue = ref(0)

const createForm = reactive({
  name: '',
  contact_name: '',
  phone: '',
  email: '',
  address: '',
  settlement_method: 'Monthly',
  credit_limit: 100000,
})

const columns: Column[] = [
  { key: 'name', label: '货主名称' },
  { key: 'contact_name', label: '联系人' },
  { key: 'phone', label: '电话' },
  { key: 'settlement', label: '结算方式' },
  { key: 'credit_limit', label: '信用额度' },
  { key: 'current_debt', label: '当前欠款' },
  { key: 'status', label: '状态' },
  { key: 'actions', label: '操作' },
]

async function load() {
  const params: Record<string, unknown> = { page: page.value, page_size: pageSize }
  if (keyword.value) params.keyword = keyword.value
  if (statusFilter.value) params.status = statusFilter.value
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

function openCreate() {
  Object.assign(createForm, { name: '', contact_name: '', phone: '', email: '', address: '', settlement_method: 'Monthly', credit_limit: 100000 })
  showCreate.value = true
}

async function submitCreate() {
  if (!createForm.name || !createForm.contact_name || !createForm.phone) {
    return toast('请填写货主名称、联系人和电话', 'error')
  }
  creating.value = true
  try {
    await store.create(createForm)
    toast('货主创建成功')
    showCreate.value = false
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    creating.value = false
  }
}

function adjustCredit(owner: Owner) {
  creditOwner.value = owner
  creditValue.value = owner.credit_limit
  showCredit.value = true
}

async function submitCredit() {
  if (!creditOwner.value) return
  creating.value = true
  try {
    await store.adjustCredit(creditOwner.value.id, creditValue.value)
    toast('信用额度已更新')
    showCredit.value = false
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    creating.value = false
  }
}

async function suspend(owner: Owner) {
  try {
    await store.suspend(owner.id)
    toast('已暂停合作')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  }
}

async function activate(owner: Owner) {
  try {
    await store.activate(owner.id)
    toast('合作已恢复')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  }
}

onMounted(load)
</script>

<style scoped>
.toolbar { display: flex; gap: 10px; margin-bottom: 16px; align-items: center; }
.toolbar__input { max-width: 280px; }
.toolbar__select { max-width: 160px; }
.row-actions { display: flex; gap: 6px; align-items: center; }
.link { font-weight: 600; }
.form-row--full { grid-column: 1 / -1; }
.modal-actions { display: flex; justify-content: flex-end; margin-top: 18px; }
</style>
