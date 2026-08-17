// 与后端 backend/internal/constants/enums.go 对应的业务枚举。

export type InboundStatus = 'Pending' | 'Received' | 'QCInProgress' | 'Shelved' | 'Completed'
export type OutboundStatus = 'Pending' | 'Picking' | 'Checking' | 'Packing' | 'Shipped' | 'Completed'
export type StorageRequirement = 'Normal' | 'ColdChain' | 'Dangerous' | 'Fragile'
export type QCResult = 'Pass' | 'Fail' | 'Partial'
export type SettlementMethod = 'Monthly' | 'PerOrder' | 'Prepaid'
export type OwnerStatus = 'Active' | 'Suspended'
export type BinLocationStatus = 'Available' | 'Occupied' | 'Reserved' | 'Maintenance'
export type UserRole = 'Admin' | 'WarehouseManager' | 'QCInspector' | 'Picker' | 'Checker' | 'Owner'

export const INBOUND_STATUSES: InboundStatus[] = ['Pending', 'Received', 'QCInProgress', 'Shelved', 'Completed']
export const OUTBOUND_STATUSES: OutboundStatus[] = ['Pending', 'Picking', 'Checking', 'Packing', 'Shipped', 'Completed']
export const STORAGE_REQUIREMENTS: StorageRequirement[] = ['Normal', 'ColdChain', 'Dangerous', 'Fragile']
export const QC_RESULTS: QCResult[] = ['Pass', 'Fail', 'Partial']
export const SETTLEMENT_METHODS: SettlementMethod[] = ['Monthly', 'PerOrder', 'Prepaid']
export const BIN_STATUSES: BinLocationStatus[] = ['Available', 'Occupied', 'Reserved', 'Maintenance']
export const BIN_AREAS: string[] = ['A', 'B', 'C', 'D']

export interface Paginated<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

export interface PageQuery {
  page?: number
  page_size?: number
  [key: string]: unknown
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}
