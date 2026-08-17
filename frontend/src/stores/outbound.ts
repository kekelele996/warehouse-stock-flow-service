import { defineStore } from 'pinia'
import {
  apiListOutbound, apiGetOutbound, apiCreateOutbound, apiPickingOutbound,
  apiCheckingOutbound, apiPackingOutbound, apiShipOutbound, apiCompleteOutbound,
  type OutboundOrderView,
} from '../api/outbound'
import type { PageQuery, Paginated } from '../types/enums'

export const useOutboundStore = defineStore('outbound', {
  state: () => ({
    orders: [] as OutboundOrderView[],
    total: 0,
    loading: false,
    current: null as OutboundOrderView | null,
  }),
  actions: {
    async list(params: PageQuery) {
      this.loading = true
      try {
        const res: Paginated<OutboundOrderView> = await apiListOutbound(params)
        this.orders = res.list
        this.total = res.total
        return res
      } finally {
        this.loading = false
      }
    },
    async fetch(id: number) {
      this.current = await apiGetOutbound(id)
      return this.current
    },
    async create(data: Record<string, unknown>) {
      return apiCreateOutbound(data)
    },
    async picking(id: number, items: { item_id: number; actual_qty: number }[]) {
      this.current = await apiPickingOutbound(id, items)
      return this.current
    },
    async checking(id: number) {
      this.current = await apiCheckingOutbound(id)
      return this.current
    },
    async packing(id: number) {
      this.current = await apiPackingOutbound(id)
      return this.current
    },
    async ship(id: number, trackingNo: string) {
      this.current = await apiShipOutbound(id, trackingNo)
      return this.current
    },
    async complete(id: number) {
      this.current = await apiCompleteOutbound(id)
      return this.current
    },
  },
})
