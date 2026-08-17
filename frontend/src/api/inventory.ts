import { get } from '../utils/request'
import type { Paginated, PageQuery } from '../types/enums'

export interface InventoryView {
  id: number
  product_id: number
  product_name: string
  sku: string
  owner_id: number
  owner_name: string
  bin_location_id: number
  bin_code: string
  batch_no: string
  quantity: number
  unit_price: number
  total_value: number
  updated_at: string
}

export interface InventorySummaryItem {
  owner_id: number
  owner_name: string
  total_qty: number
  total_value: number
}

export function apiListInventory(params: PageQuery) {
  return get<Paginated<InventoryView>>('/inventory', params)
}

export function apiInventorySummary() {
  return get<InventorySummaryItem[]>('/inventory/summary')
}
