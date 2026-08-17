import { get, post, put } from '../utils/request'
import type { Paginated, PageQuery } from '../types/enums'

export interface ProductView {
  id: number
  owner_id: number
  owner_name: string
  name: string
  sku: string
  barcode: string
  category: string
  spec: string
  unit: string
  shelf_life_days: number | null
  storage_requirement: string
  storage_text: string
  volume: number
  weight: number
  price: number
  created_at: string
}

export function apiListProducts(params: PageQuery) {
  return get<Paginated<ProductView>>('/products', params)
}

export function apiListProductsByOwner(ownerId: number, params: PageQuery) {
  return get<Paginated<ProductView>>(`/products/by-owner/${ownerId}`, params)
}

export function apiGetProduct(id: number) {
  return get<ProductView>(`/products/${id}`)
}

export function apiCreateProduct(data: Record<string, unknown>) {
  return post<ProductView>('/products', data)
}

export function apiUpdateProduct(id: number, data: Record<string, unknown>) {
  return put<ProductView>(`/products/${id}`, data)
}
