import { get, post } from '../utils/request'
import type { Paginated, PageQuery } from '../types/enums'

export interface OutboundItemView {
  id: number
  product_id: number
  product_name: string
  sku: string
  bin_location_id: number | null
  bin_code: string
  expected_qty: number
  actual_qty: number
}

export interface OutboundOrderView {
  id: number
  order_no: string
  owner_id: number
  owner_name: string
  receiver_name: string
  receiver_address: string
  required_ship_date: string | null
  actual_ship_date: string | null
  status: string
  status_text: string
  picker_name: string
  checker_name: string
  tracking_no: string
  items: OutboundItemView[]
  created_at: string
  updated_at: string
}

export function apiListOutbound(params: PageQuery) {
  return get<Paginated<OutboundOrderView>>('/outbound', params)
}

export function apiGetOutbound(id: number) {
  return get<OutboundOrderView>(`/outbound/${id}`)
}

export function apiCreateOutbound(data: Record<string, unknown>) {
  return post<OutboundOrderView>('/outbound', data)
}

export function apiPickingOutbound(id: number, items: { item_id: number; actual_qty: number }[]) {
  return post<OutboundOrderView>(`/outbound/${id}/picking`, { items })
}

export function apiCheckingOutbound(id: number) {
  return post<OutboundOrderView>(`/outbound/${id}/checking`, {})
}

export function apiPackingOutbound(id: number) {
  return post<OutboundOrderView>(`/outbound/${id}/packing`, {})
}

export function apiShipOutbound(id: number, trackingNo: string) {
  return post<OutboundOrderView>(`/outbound/${id}/ship`, { tracking_no: trackingNo })
}

export function apiCompleteOutbound(id: number) {
  return post<OutboundOrderView>(`/outbound/${id}/complete`, {})
}
