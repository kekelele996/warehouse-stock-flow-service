import { defineStore } from 'pinia'
import { apiListProducts, apiCreateProduct, apiUpdateProduct, type ProductView } from '../api/product'
import type { PageQuery, Paginated } from '../types/enums'

export const useProductStore = defineStore('product', {
  state: () => ({
    products: [] as ProductView[],
    total: 0,
    loading: false,
  }),
  actions: {
    async list(params: PageQuery) {
      this.loading = true
      try {
        const res: Paginated<ProductView> = await apiListProducts(params)
        this.products = res.list
        this.total = res.total
        return res
      } finally {
        this.loading = false
      }
    },
    async create(data: Record<string, unknown>) {
      const product = await apiCreateProduct(data)
      return product
    },
    async update(id: number, data: Record<string, unknown>) {
      return apiUpdateProduct(id, data)
    },
  },
})
