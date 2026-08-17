<template>
  <div class="layout">
    <aside class="sidebar">
      <div class="sidebar__brand">
        <div class="sidebar__logo">📦</div>
        <div>
          <div class="sidebar__title">WMSFlow</div>
          <div class="sidebar__sub">仓储出入库协同</div>
        </div>
      </div>
      <nav class="sidebar__nav">
        <router-link
          v-for="item in menus"
          :key="item.path"
          :to="item.path"
          class="sidebar__item"
          active-class="sidebar__item--active"
        >
          <span class="sidebar__icon">{{ item.icon }}</span>
          <span>{{ item.title }}</span>
        </router-link>
      </nav>
    </aside>
    <div class="main">
      <header class="topbar">
        <div class="topbar__title">{{ route.meta.title || 'WMSFlow' }}</div>
        <div class="topbar__right">
          <span class="topbar__role">{{ roleText }}</span>
          <span class="topbar__user">{{ auth.user?.name || auth.user?.username }}</span>
          <button class="btn btn--ghost" @click="logout">退出</button>
        </div>
      </header>
      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuth } from '../hooks/useAuth'
import { ROLE_TEXT } from '../constants/enums'

const route = useRoute()
const router = useRouter()
const { store: auth, menus } = useAuth()

const roleText = computed(() => ROLE_TEXT[auth.role] || auth.role)

function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout {
  display: flex;
  height: 100%;
}
.sidebar {
  width: 220px;
  background: var(--sidebar-bg);
  color: #cbd5e1;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}
.sidebar__brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 18px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}
.sidebar__logo {
  width: 38px;
  height: 38px;
  border-radius: 10px;
  background: var(--sidebar-active);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
}
.sidebar__title { font-size: 16px; font-weight: 700; color: #fff; }
.sidebar__sub { font-size: 11px; color: #8aa0b8; }
.sidebar__nav { padding: 12px 10px; display: flex; flex-direction: column; gap: 4px; flex: 1; overflow-y: auto; }
.sidebar__item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  color: #cbd5e1;
  font-size: 13px;
  font-weight: 500;
  transition: background 0.15s ease;
}
.sidebar__item:hover { background: rgba(255, 255, 255, 0.08); color: #fff; }
.sidebar__item--active { background: var(--sidebar-active); color: #fff; }
.sidebar__icon { font-size: 16px; }
.main { flex: 1; display: flex; flex-direction: column; min-width: 0; }
.topbar {
  height: 56px;
  background: #fff;
  border-bottom: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}
.topbar__title { font-size: 16px; font-weight: 700; color: var(--text-main); }
.topbar__right { display: flex; align-items: center; gap: 12px; font-size: 13px; }
.topbar__role {
  background: #eef4ff;
  color: var(--primary);
  padding: 4px 10px;
  border-radius: 999px;
  font-weight: 600;
}
.topbar__user { color: var(--text-sub); font-weight: 600; }
.content { flex: 1; overflow-y: auto; }
</style>
