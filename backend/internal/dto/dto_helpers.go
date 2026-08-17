package dto

import "github.com/wmsflow/wmsflow/internal/util"

// settlementText 复用 util formatter 的结算方式文本。
func settlementText(method string) string { return util.SettlementMethodText(method) }

// ownerStatusText 复用 util formatter 的货主状态文本。
func ownerStatusText(status string) string { return util.OwnerStatusText(status) }

// storageText 复用 util formatter 的存储要求文本。
func storageText(req string) string { return util.StorageRequirementText(req) }

// binStatusText 复用 util formatter 的库位状态文本。
func binStatusText(status string) string {
	return util.BinStatusText(status)
}

// inboundStatusText 复用 util formatter 的入库单状态文本。
func inboundStatusText(status string) string { return util.InboundStatusText(status) }

// outboundStatusText 复用 util formatter 的出库单状态文本。
func outboundStatusText(status string) string { return util.OutboundStatusText(status) }

// qcResultText 复用 util formatter 的质检结果文本。
func qcResultText(result string) string { return util.QCResultText(result) }

// roleText 复用 util formatter 的角色文本。
func roleText(role string) string { return util.RoleText(role) }
