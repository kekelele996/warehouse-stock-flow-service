import { get, post } from '../utils/request'
import type { Paginated, PageQuery } from '../types/enums'

export interface InboundItemView {
  id: number
  product_id: number
  product_name: string
  sku: string
  batch_no: string
  expected_qty: number
  actual_qty: number
  qc_result: string
  qc_result_text: string
  bin_location_id: number | null
  bin_code: string
}

export interface InboundOrderView {
  id: number
  order_no: string
  owner_id: number
  owner_name: string
  supplier_name: string
  expected_arrival_date: string | null
  actual_arrival_date: string | null
  status: string
  status_text: string
  qc_inspector_name: string
  keeper_name: string
  remark: string
  items: InboundItemView[]
  created_at: string
  updated_at: string
}

export function apiListInbound(params: PageQuery) {
  return get<Paginated<InboundOrderView>>('/inbound', params)
}

export function apiGetInbound(id: number) {
  return get<InboundOrderView>(`/inbound/${id}`)
}

export function apiCreateInbound(data: Record<string, unknown>) {
  return post<InboundOrderView>('/inbound', data)
}

export function apiReceiveInbound(id: number) {
  return post<InboundOrderView>(`/inbound/${id}/receive`, {})
}

export function apiQCInbound(id: number, items: { item_id: number; actual_qty: number; qc_result: string }[]) {
  return post<InboundOrderView>(`/inbound/${id}/qc`, { items })
}

export function apiShelveInbound(id: number, items: { item_id: number; bin_location_id: number }[]) {
  return post<InboundOrderView>(`/inbound/${id}/shelve`, { items })
}

export function apiCompleteInbound(id: number) {
  return post<InboundOrderView>(`/inbound/${id}/complete`, {})
}
