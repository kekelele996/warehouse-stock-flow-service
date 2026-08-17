import { get } from '../utils/request'

export interface OwnerTopItem {
  owner_id: number
  owner_name: string
  total_value: number
}

export interface PendingTaskItem {
  type: string
  title: string
  order_no: string
  status: string
}

export interface DashboardSummary {
  today_inbound_count: number
  today_outbound_count: number
  pending_inbound: number
  pending_outbound: number
  bin_occupancy_rate: number
  total_bin_count: number
  occupied_bin_count: number
  owner_top10: OwnerTopItem[]
  pending_tasks: PendingTaskItem[]
}

export function apiDashboardSummary() {
  return get<DashboardSummary>('/dashboard/summary')
}
