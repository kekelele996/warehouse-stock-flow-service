import { defineStore } from 'pinia'
import {
  apiListOwners, apiGetOwner, apiGetOwnerStats, apiCreateOwner,
  apiUpdateOwner, apiAdjustCredit, apiSuspendOwner, apiActivateOwner,
  type OwnerView, type OwnerStatsView,
} from '../api/owner'
import type { PageQuery, Paginated } from '../types/enums'

export const useOwnerStore = defineStore('owner', {
  state: () => ({
    owners: [] as OwnerView[],
    total: 0,
    loading: false,
    current: null as OwnerView | null,
    stats: null as OwnerStatsView | null,
  }),
  actions: {
    async list(params: PageQuery) {
      this.loading = true
      try {
        const res = await apiListOwners(params)
        this.owners = res.list
        this.total = res.total
        return res
      } finally {
        this.loading = false
      }
    },
    async fetch(id: number) {
      this.current = await apiGetOwner(id)
      return this.current
    },
    async fetchStats(id: number) {
      this.stats = await apiGetOwnerStats(id)
      return this.stats
    },
    async create(data: Record<string, unknown>) {
      const owner = await apiCreateOwner(data)
      return owner
    },
    async update(id: number, data: Record<string, unknown>) {
      return apiUpdateOwner(id, data)
    },
    async adjustCredit(id: number, creditLimit: number) {
      return apiAdjustCredit(id, creditLimit)
    },
    async suspend(id: number) {
      return apiSuspendOwner(id)
    },
    async activate(id: number) {
      return apiActivateOwner(id)
    },
  },
})
