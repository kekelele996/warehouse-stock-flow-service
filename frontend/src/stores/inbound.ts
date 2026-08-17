import { defineStore } from 'pinia'
import {
  apiListInbound, apiGetInbound, apiCreateInbound, apiReceiveInbound,
  apiQCInbound, apiShelveInbound, apiCompleteInbound, type InboundOrderView,
} from '../api/inbound'
import type { PageQuery, Paginated } from '../types/enums'

export const useInboundStore = defineStore('inbound', {
  state: () => ({
    orders: [] as InboundOrderView[],
    total: 0,
    loading: false,
    current: null as InboundOrderView | null,
  }),
  actions: {
    async list(params: PageQuery) {
      this.loading = true
      try {
        const res: Paginated<InboundOrderView> = await apiListInbound(params)
        this.orders = res.list
        this.total = res.total
        return res
      } finally {
        this.loading = false
      }
    },
    async fetch(id: number) {
      this.current = await apiGetInbound(id)
      return this.current
    },
    async create(data: Record<string, unknown>) {
      return apiCreateInbound(data)
    },
    async receive(id: number) {
      this.current = await apiReceiveInbound(id)
      return this.current
    },
    async qc(id: number, items: { item_id: number; actual_qty: number; qc_result: string }[]) {
      this.current = await apiQCInbound(id, items)
      return this.current
    },
    async shelve(id: number, items: { item_id: number; bin_location_id: number }[]) {
      this.current = await apiShelveInbound(id, items)
      return this.current
    },
    async complete(id: number) {
      this.current = await apiCompleteInbound(id)
      return this.current
    },
  },
})
