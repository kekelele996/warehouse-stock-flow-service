import { defineStore } from 'pinia'
import { apiListAudit, type AuditView } from '../api/audit'
import type { PageQuery, Paginated } from '../types/enums'

export const useAuditStore = defineStore('audit', {
  state: () => ({
    logs: [] as AuditView[],
    total: 0,
    loading: false,
  }),
  actions: {
    async list(params: PageQuery) {
      this.loading = true
      try {
        const res: Paginated<AuditView> = await apiListAudit(params)
        this.logs = res.list
        this.total = res.total
        return res
      } finally {
        this.loading = false
      }
    },
  },
})
