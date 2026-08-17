import { get } from '../utils/request'
import type { Paginated, PageQuery } from '../types/enums'

export interface AuditView {
  id: number
  user_id: number
  username: string
  role: string
  role_text: string
  module: string
  action: string
  entity_type: string
  entity_id: string
  detail: string
  ip: string
  created_at: string
}

export function apiListAudit(params: PageQuery) {
  return get<Paginated<AuditView>>('/audit', params)
}
