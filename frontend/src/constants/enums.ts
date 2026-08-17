// 前端业务枚举常量（与 src/types/enums.ts 及后端 constants/enums.go 对应，多处引用）。

export const ENUM_TEXT: Record<string, string> = {
  // InboundStatus
  Pending: '待收货',
  Received: '已收货',
  QCInProgress: '质检中',
  Shelved: '已上架',
  Completed: '已完成',
  // OutboundStatus
  Picking: '拣货中',
  Checking: '复核中',
  Packing: '打包中',
  Shipped: '已发货',
  // StorageRequirement
  Normal: '常温',
  ColdChain: '冷链',
  Dangerous: '危险品',
  Fragile: '易碎',
  // QCResult
  Pass: '合格',
  Fail: '不合格',
  Partial: '部分合格',
  // SettlementMethod
  Monthly: '月结',
  PerOrder: '单结',
  Prepaid: '预付',
  // OwnerStatus
  Active: '合作中',
  Suspended: '已暂停',
  // BinLocationStatus
  Available: '可用',
  Occupied: '占用',
  Reserved: '预留',
  Maintenance: '维护',
}

export const STATUS_COLOR: Record<string, string> = {
  Pending: 'gray',
  Received: 'blue',
  QCInProgress: 'orange',
  Shelved: 'cyan',
  Completed: 'green',
  Picking: 'blue',
  Checking: 'orange',
  Packing: 'purple',
  Shipped: 'cyan',
  Pass: 'green',
  Fail: 'red',
  Partial: 'orange',
  Active: 'green',
  Suspended: 'red',
  Available: 'green',
  Occupied: 'orange',
  Reserved: 'blue',
  Maintenance: 'gray',
}

export const ROLE_TEXT: Record<string, string> = {
  Admin: '系统管理员',
  WarehouseManager: '仓库经理',
  QCInspector: '质检员',
  Picker: '拣货员',
  Checker: '复核员',
  Owner: '货主',
}
