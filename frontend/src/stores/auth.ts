import { defineStore } from 'pinia'
import { apiLogin, apiMe, apiRegister, type UserView } from '../api/auth'
import { TOKEN_KEY } from '../utils/request'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem(TOKEN_KEY) || '',
    user: JSON.parse(localStorage.getItem('wmsflow_user') || 'null') as UserView | null,
  }),
  getters: {
    isLoggedIn: (s) => !!s.token,
    role: (s) => s.user?.role || '',
  },
  actions: {
    async login(username: string, password: string) {
      const res = await apiLogin(username, password)
      this.token = res.token
      this.user = res.user
      localStorage.setItem(TOKEN_KEY, res.token)
      localStorage.setItem('wmsflow_user', JSON.stringify(res.user))
      return res
    },
    async register(data: { username: string; password: string; name: string; owner_name: string; contact_name: string; phone: string }) {
      return apiRegister(data)
    },
    async fetchMe() {
      this.user = await apiMe()
      localStorage.setItem('wmsflow_user', JSON.stringify(this.user))
    },
    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem(TOKEN_KEY)
      localStorage.removeItem('wmsflow_user')
    },
  },
})
