import { defineStore } from 'pinia'
import { apiDashboardSummary, type DashboardSummary } from '../api/dashboard'

export const useDashboardStore = defineStore('dashboard', {
  state: () => ({
    summary: null as DashboardSummary | null,
    loading: false,
  }),
  actions: {
    async fetch() {
      this.loading = true
      try {
        this.summary = await apiDashboardSummary()
        return this.summary
      } finally {
        this.loading = false
      }
    },
  },
})
