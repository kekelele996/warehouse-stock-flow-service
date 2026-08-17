import { get, post, put } from '../utils/request'
import type { Paginated, PageQuery } from '../types/enums'

export interface BinView {
  id: number
  code: string
  area: string
  rack_no: string
  layer_no: number
  column_no: number
  capacity: number
  occupancy_rate: number
  storage_requirement: string
  storage_text: string
  status: string
  status_text: string
  created_at: string
}

export interface BinContent {
  product_id: number
  product_name: string
  sku: string
  batch_no: string
  quantity: number
  owner_name: string
}

export function apiListBins(params: PageQuery) {
  return get<Paginated<BinView>>('/bin-locations', params)
}

export function apiGetBin(id: number) {
  return get<BinView>(`/bin-locations/${id}`)
}

export function apiGetBinContents(id: number) {
  return get<BinContent[]>(`/bin-locations/${id}/contents`)
}

export function apiCreateBin(data: Record<string, unknown>) {
  return post<BinView>('/bin-locations', data)
}

export function apiBatchCreateBins(data: Record<string, unknown>) {
  return post<{ created: number }>('/bin-locations/batch', data)
}

export function apiUpdateBin(id: number, data: Record<string, unknown>) {
  return put<BinView>(`/bin-locations/${id}`, data)
}

export function apiRecommendBin(productId: number, quantity: number) {
  return post<BinView>('/bin-locations/recommend', { product_id: productId, quantity })
}
