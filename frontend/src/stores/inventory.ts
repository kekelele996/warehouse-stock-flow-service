import { defineStore } from 'pinia'
import { apiListInventory, apiInventorySummary, type InventoryView, type InventorySummaryItem } from '../api/inventory'
import type { PageQuery, Paginated } from '../types/enums'

export const useInventoryStore = defineStore('inventory', {
  state: () => ({
    items: [] as InventoryView[],
    total: 0,
    loading: false,
    summary: [] as InventorySummaryItem[],
  }),
  actions: {
    async list(params: PageQuery) {
      this.loading = true
      try {
        const res: Paginated<InventoryView> = await apiListInventory(params)
        this.items = res.list
        this.total = res.total
        return res
      } finally {
        this.loading = false
      }
    },
    async fetchSummary() {
      this.summary = await apiInventorySummary()
      return this.summary
    },
  },
})
