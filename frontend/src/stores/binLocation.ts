import { defineStore } from 'pinia'
import {
  apiListBins, apiGetBin, apiGetBinContents, apiCreateBin, apiBatchCreateBins,
  apiUpdateBin, apiRecommendBin, type BinView, type BinContent,
} from '../api/binLocation'
import type { PageQuery, Paginated } from '../types/enums'

export const useBinLocationStore = defineStore('binLocation', {
  state: () => ({
    bins: [] as BinView[],
    total: 0,
    loading: false,
    current: null as BinView | null,
    contents: [] as BinContent[],
  }),
  actions: {
    async list(params: PageQuery) {
      this.loading = true
      try {
        const res: Paginated<BinView> = await apiListBins(params)
        this.bins = res.list
        this.total = res.total
        return res
      } finally {
        this.loading = false
      }
    },
    async fetch(id: number) {
      this.current = await apiGetBin(id)
      return this.current
    },
    async fetchContents(id: number) {
      this.contents = await apiGetBinContents(id)
      return this.contents
    },
    async create(data: Record<string, unknown>) {
      return apiCreateBin(data)
    },
    async batchCreate(data: Record<string, unknown>) {
      return apiBatchCreateBins(data)
    },
    async update(id: number, data: Record<string, unknown>) {
      return apiUpdateBin(id, data)
    },
    async recommend(productId: number, quantity: number) {
      return apiRecommendBin(productId, quantity)
    },
  },
})
