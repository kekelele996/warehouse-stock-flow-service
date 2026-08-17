// 角色与菜单可见性配置（RBAC 前端部分）。
import type { UserRole } from '../types/enums'

export interface MenuItem {
  path: string
  title: string
  icon: string
  roles: UserRole[]
}

export const MENU_ITEMS: MenuItem[] = [
  { path: '/dashboard', title: '仓库总览', icon: '📊', roles: ['Admin', 'WarehouseManager', 'QCInspector', 'Picker', 'Checker', 'Owner'] },
  { path: '/inbound', title: '入库管理', icon: '📥', roles: ['Admin', 'WarehouseManager', 'QCInspector', 'Owner'] },
  { path: '/outbound', title: '出库管理', icon: '📤', roles: ['Admin', 'WarehouseManager', 'Picker', 'Checker', 'Owner'] },
  { path: '/bin-locations', title: '库位管理', icon: '🗄️', roles: ['Admin', 'WarehouseManager', 'QCInspector', 'Picker'] },
  { path: '/products', title: '商品管理', icon: '📦', roles: ['Admin', 'WarehouseManager', 'Owner'] },
  { path: '/owners', title: '货主管理', icon: '🏢', roles: ['Admin', 'WarehouseManager', 'Owner'] },
  { path: '/audit', title: '操作日志', icon: '📜', roles: ['Admin', 'WarehouseManager'] },
]

export function canAccess(role: UserRole | string | undefined, item: MenuItem): boolean {
  if (!role) return false
  return item.roles.includes(role as UserRole)
}

export const ACTION_ROLES: Record<string, UserRole[]> = {
  ownerManage: ['Admin', 'WarehouseManager'],
  inboundManage: ['Admin', 'WarehouseManager', 'QCInspector', 'Owner'],
  outboundManage: ['Admin', 'WarehouseManager', 'Picker', 'Checker', 'Owner'],
  binManage: ['Admin', 'WarehouseManager'],
  inboundReceive: ['Admin', 'WarehouseManager'],
  inboundQC: ['Admin', 'WarehouseManager', 'QCInspector'],
  inboundShelve: ['Admin', 'WarehouseManager'],
  outboundPicking: ['Admin', 'WarehouseManager', 'Picker'],
  outboundChecking: ['Admin', 'WarehouseManager', 'Checker'],
  outboundShip: ['Admin', 'WarehouseManager'],
}
