package constants

// InboundStatus 入库单状态。
const (
	InboundStatusPending      = "Pending"      // 待收货
	InboundStatusReceived     = "Received"     // 已收货
	InboundStatusQCInProgress = "QCInProgress" // 质检中
	InboundStatusShelved      = "Shelved"      // 已上架
	InboundStatusCompleted    = "Completed"    // 已完成
)

// InboundStatusList 入库单状态集合。
var InboundStatusList = []string{
	InboundStatusPending,
	InboundStatusReceived,
	InboundStatusQCInProgress,
	InboundStatusShelved,
	InboundStatusCompleted,
}

// OutboundStatus 出库单状态。
const (
	OutboundStatusPending   = "Pending"   // 待拣货
	OutboundStatusPicking   = "Picking"   // 拣货中
	OutboundStatusChecking  = "Checking"  // 复核中
	OutboundStatusPacking   = "Packing"   // 打包中
	OutboundStatusShipped   = "Shipped"   // 已发货
	OutboundStatusCompleted = "Completed" // 已完成
)

// OutboundStatusList 出库单状态集合。
var OutboundStatusList = []string{
	OutboundStatusPending,
	OutboundStatusPicking,
	OutboundStatusChecking,
	OutboundStatusPacking,
	OutboundStatusShipped,
	OutboundStatusCompleted,
}

// StorageRequirement 存储要求。
const (
	StorageNormal     = "Normal"     // 常温
	StorageColdChain  = "ColdChain"  // 冷链
	StorageDangerous  = "Dangerous"  // 危险品
	StorageFragile    = "Fragile"    // 易碎
)

// StorageRequirementList 存储要求集合。
var StorageRequirementList = []string{StorageNormal, StorageColdChain, StorageDangerous, StorageFragile}

// QCResult 质检结果。
const (
	QCResultPass    = "Pass"    // 合格
	QCResultFail    = "Fail"    // 不合格
	QCResultPartial = "Partial" // 部分合格
)

// QCResultList 质检结果集合。
var QCResultList = []string{QCResultPass, QCResultFail, QCResultPartial}

// SettlementMethod 结算方式。
const (
	SettlementMonthly   = "Monthly"   // 月结
	SettlementPerOrder  = "PerOrder"  // 单结
	SettlementPrepaid   = "Prepaid"   // 预付
)

// SettlementMethodList 结算方式集合。
var SettlementMethodList = []string{SettlementMonthly, SettlementPerOrder, SettlementPrepaid}

// OwnerStatus 货主状态。
const (
	OwnerStatusActive    = "Active"    // 合作中
	OwnerStatusSuspended = "Suspended" // 已暂停
)

// OwnerStatusList 货主状态集合。
var OwnerStatusList = []string{OwnerStatusActive, OwnerStatusSuspended}

// BinLocationStatus 库位状态。
const (
	BinStatusAvailable   = "Available"   // 可用
	BinStatusOccupied    = "Occupied"    // 占用
	BinStatusReserved    = "Reserved"    // 预留
	BinStatusMaintenance = "Maintenance" // 维护
)

// BinStatusList 库位状态集合。
var BinStatusList = []string{BinStatusAvailable, BinStatusOccupied, BinStatusReserved, BinStatusMaintenance}

// BinAreaList 仓库区域集合。
var BinAreaList = []string{"A", "B", "C", "D"}

// InboundTransitions 入库单状态机：key 为当前状态，value 为允许流转的下一个状态。
var InboundTransitions = map[string][]string{
	InboundStatusPending:      {InboundStatusReceived},
	InboundStatusReceived:     {InboundStatusQCInProgress},
	InboundStatusQCInProgress: {InboundStatusShelved},
	InboundStatusShelved:      {InboundStatusCompleted},
	InboundStatusCompleted:    {},
}

// OutboundTransitions 出库单状态机：key 为当前状态，value 为允许流转的下一个状态。
var OutboundTransitions = map[string][]string{
	OutboundStatusPending:   {OutboundStatusPicking, OutboundStatusChecking},
	OutboundStatusPicking:   {OutboundStatusChecking},
	OutboundStatusChecking:  {OutboundStatusPacking},
	OutboundStatusPacking:   {OutboundStatusShipped},
	OutboundStatusShipped:   {OutboundStatusCompleted},
	OutboundStatusCompleted: {},
}

// CanTransitInbound 判断入库单状态是否可流转。
func CanTransitInbound(from, to string) bool {
	for _, next := range InboundTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

// CanTransitOutbound 判断出库单状态是否可流转。
func CanTransitOutbound(from, to string) bool {
	for _, next := range OutboundTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

// Contains 判断字符串切片是否包含目标值。
func Contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
