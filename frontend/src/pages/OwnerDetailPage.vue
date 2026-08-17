<template>
  <div class="page">
    <div class="page__head">
      <div>
        <div class="page__title">{{ stats?.owner.name || '货主详情' }}</div>
        <div class="page__sub">基本信息 + 库存统计 + 历史单据 + 账务概况</div>
      </div>
      <div class="row-actions">
        <router-link to="/owners" class="btn">返回列表</router-link>
        <button v-if="can('ownerManage') && stats?.owner.status === 'Active'" class="btn" @click="suspend">暂停合作</button>
        <button v-if="can('ownerManage') && stats?.owner.status === 'Suspended'" class="btn" @click="activate">恢复合作</button>
      </div>
    </div>

    <div v-if="loading" class="loading-mask">加载中…</div>
    <template v-else-if="stats">
      <div class="grid grid--4">
        <StatCard icon="📦" label="商品数" :value="stats.product_count" />
        <StatCard icon="📥" label="入库单数" :value="stats.inbound_count" />
        <StatCard icon="📤" label="出库单数" :value="stats.outbound_count" />
        <StatCard icon="💰" label="库存金额" :value="`¥${formatAmount(stats.inventory_value)}`" />
      </div>

      <div class="grid grid--2 detail-grid">
        <div class="card">
          <h3 class="card__title">基本信息</h3>
          <div class="kv">
            <div class="kv__row"><span>联系人</span><b>{{ stats.owner.contact_name }}</b></div>
            <div class="kv__row"><span>电话</span><b>{{ stats.owner.phone }}</b></div>
            <div class="kv__row"><span>邮箱</span><b>{{ stats.owner.email || '-' }}</b></div>
            <div class="kv__row"><span>地址</span><b>{{ stats.owner.address || '-' }}</b></div>
            <div class="kv__row"><span>状态</span><b><StatusBadge :status="stats.owner.status" /></b></div>
          </div>
        </div>
        <div class="card">
          <h3 class="card__title">账务概况</h3>
          <div class="kv">
            <div class="kv__row"><span>结算方式</span><b>{{ stats.owner.settlement_text }}</b></div>
            <div class="kv__row"><span>信用额度</span><b>¥{{ formatAmount(stats.owner.credit_limit) }}</b></div>
            <div class="kv__row"><span>当前欠款</span><b>¥{{ formatAmount(stats.owner.current_debt) }}</b></div>
            <div class="kv__row"><span>库存金额</span><b>¥{{ formatAmount(stats.inventory_value) }}</b></div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import StatCard from '../components/StatCard.vue'
import StatusBadge from '../components/StatusBadge.vue'
import { useOwnerStore } from '../stores/owner'
import { useAuth } from '../hooks/useAuth'
import { formatAmount } from '../utils/format'
import { toast } from '../utils/toast'

const route = useRoute()
const store = useOwnerStore()
const { can } = useAuth()

const id = computed(() => Number(route.params.id))
const stats = computed(() => store.stats)
const loading = computed(() => false)

async function load() {
  await store.fetchStats(id.value)
}

async function suspend() {
  try {
    await store.suspend(id.value)
    toast('已暂停合作')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  }
}

async function activate() {
  try {
    await store.activate(id.value)
    toast('合作已恢复')
    load()
  } catch (e) {
    toast((e as Error).message, 'error')
  }
}

onMounted(load)
</script>

<style scoped>
.row-actions { display: flex; gap: 8px; }
.detail-grid { margin-top: 16px; }
.card__title { font-size: 15px; font-weight: 700; margin-bottom: 14px; }
.kv { display: flex; flex-direction: column; gap: 10px; }
.kv__row { display: flex; justify-content: space-between; font-size: 13px; color: var(--text-muted); }
.kv__row b { color: var(--text-main); font-weight: 600; }
</style>
