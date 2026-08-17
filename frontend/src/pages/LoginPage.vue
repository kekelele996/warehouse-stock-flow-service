<template>
  <div class="login">
    <div class="login__panel">
      <div class="login__brand">
        <div class="login__logo">📦</div>
        <div>
          <div class="login__name">WMSFlow</div>
          <div class="login__slogan">供应链仓储出入库协同管理系统</div>
        </div>
      </div>

      <template v-if="mode === 'login'">
        <h2 class="login__title">欢迎回来</h2>
        <div class="form-row">
          <label class="form-row__label">用户名</label>
          <input v-model="form.username" class="input" placeholder="admin" autocomplete="username" />
        </div>
        <div class="form-row">
          <label class="form-row__label">密码</label>
          <input v-model="form.password" type="password" class="input" placeholder="wmsflow@123" autocomplete="current-password" @keyup.enter="submit" />
        </div>
        <button class="btn btn--primary login__btn" :disabled="loading" @click="submit">
          {{ loading ? '登录中…' : '登 录' }}
        </button>
        <div class="login__switch">
          还没有账号？
          <a @click="mode = 'register'">注册货主账号</a>
        </div>
        <div class="login__hint">
          演示账号：admin / wmsflow@123；manager / wmsflow@123；owner01 / wmsflow@123
        </div>
      </template>

      <template v-else>
        <h2 class="login__title">注册货主账号</h2>
        <div class="form-row">
          <label class="form-row__label">用户名</label>
          <input v-model="reg.username" class="input" placeholder="3-64位" />
        </div>
        <div class="form-row">
          <label class="form-row__label">密码</label>
          <input v-model="reg.password" type="password" class="input" placeholder="至少6位" />
        </div>
        <div class="form-row">
          <label class="form-row__label">姓名</label>
          <input v-model="reg.name" class="input" placeholder="联系人姓名" />
        </div>
        <div class="form-row">
          <label class="form-row__label">公司名称</label>
          <input v-model="reg.owner_name" class="input" placeholder="货主公司名称" />
        </div>
        <div class="form-row">
          <label class="form-row__label">联系人</label>
          <input v-model="reg.contact_name" class="input" placeholder="业务联系人" />
        </div>
        <div class="form-row">
          <label class="form-row__label">联系电话</label>
          <input v-model="reg.phone" class="input" placeholder="手机号" />
        </div>
        <button class="btn btn--primary login__btn" :disabled="loading" @click="register">
          {{ loading ? '注册中…' : '注 册' }}
        </button>
        <div class="login__switch">
          已有账号？
          <a @click="mode = 'login'">返回登录</a>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { toast } from '../utils/toast'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const mode = ref<'login' | 'register'>('login')
const loading = ref(false)
const form = reactive({ username: 'admin', password: 'wmsflow@123' })
const reg = reactive({ username: '', password: '', name: '', owner_name: '', contact_name: '', phone: '' })

async function submit() {
  if (!form.username || !form.password) {
    toast('请输入用户名和密码', 'error')
    return
  }
  loading.value = true
  try {
    await auth.login(form.username, form.password)
    toast('登录成功')
    router.push((route.query.redirect as string) || '/dashboard')
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    loading.value = false
  }
}

async function register() {
  if (!reg.username || !reg.password || !reg.name || !reg.owner_name || !reg.contact_name || !reg.phone) {
    toast('请完整填写注册信息', 'error')
    return
  }
  loading.value = true
  try {
    await auth.register(reg)
    toast('注册成功，请登录')
    mode.value = 'login'
    form.username = reg.username
    form.password = reg.password
  } catch (e) {
    toast((e as Error).message, 'error')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0f2440 0%, #1d4ed8 60%, #2563eb 100%);
}
.login__panel {
  width: 420px;
  max-width: calc(100vw - 32px);
  background: #fff;
  border-radius: 18px;
  padding: 36px 36px 28px;
  box-shadow: 0 24px 64px rgba(6, 18, 36, 0.35);
}
.login__brand { display: flex; align-items: center; gap: 12px; margin-bottom: 28px; }
.login__logo {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #eef4ff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}
.login__name { font-size: 20px; font-weight: 800; color: var(--text-main); }
.login__slogan { font-size: 12px; color: var(--text-muted); }
.login__title { font-size: 20px; font-weight: 700; margin-bottom: 20px; }
.login__btn { width: 100%; margin-top: 6px; }
.login__switch { margin-top: 16px; font-size: 13px; color: var(--text-muted); text-align: center; }
.login__switch a { cursor: pointer; color: var(--primary); font-weight: 600; }
.login__hint { margin-top: 18px; background: #f6f8fb; border-radius: 8px; padding: 10px 12px; font-size: 12px; color: var(--text-muted); line-height: 1.6; }
</style>
