import { get, post, put } from '../utils/request'
import type { Paginated, PageQuery } from '../types/enums'

export interface OwnerView {
  id: number
  name: string
  contact_name: string
  phone: string
  email: string
  address: string
  settlement_method: string
  settlement_text: string
  credit_limit: number
  current_debt: number
  status: string
  status_text: string
  created_at: string
}

export interface OwnerStatsView {
  owner: OwnerView
  product_count: number
  inbound_count: number
  outbound_count: number
  inventory_value: number
}

export function apiListOwners(params: PageQuery) {
  return get<Paginated<OwnerView>>('/owners', params)
}

export function apiGetOwner(id: number) {
  return get<OwnerView>(`/owners/${id}`)
}

export function apiGetOwnerStats(id: number) {
  return get<OwnerStatsView>(`/owners/${id}/stats`)
}

export function apiCreateOwner(data: Record<string, unknown>) {
  return post<OwnerView>('/owners', data)
}

export function apiUpdateOwner(id: number, data: Record<string, unknown>) {
  return put<OwnerView>(`/owners/${id}`, data)
}

export function apiAdjustCredit(id: number, creditLimit: number) {
  return put<OwnerView>(`/owners/${id}/credit`, { credit_limit: creditLimit })
}

export function apiSuspendOwner(id: number) {
  return put<OwnerView>(`/owners/${id}/suspend`, {})
}

export function apiActivateOwner(id: number) {
  return put<OwnerView>(`/owners/${id}/activate`, {})
}
