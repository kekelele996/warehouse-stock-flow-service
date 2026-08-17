package util

import (
	"fmt"
	"strings"
	"time"

	"github.com/wmsflow/wmsflow/internal/constants"
)

// FormatDate 格式化日期为 YYYY-MM-DD。
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatDateTime 格式化为 YYYY-MM-DD HH:MM:SS。
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatTimePtr 空指针返回空串。
func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatDateTime(*t)
}

// InboundStatusText 入库单状态中文文本。
func InboundStatusText(status string) string {
	switch status {
	case constants.InboundStatusPending:
		return "待收货"
	case constants.InboundStatusReceived:
		return "已收货"
	case constants.InboundStatusQCInProgress:
		return "质检中"
	case constants.InboundStatusShelved:
		return "已上架"
	case constants.InboundStatusCompleted:
		return "已完成"
	default:
		return status
	}
}

// OutboundStatusText 出库单状态中文文本。
func OutboundStatusText(status string) string {
	switch status {
	case constants.OutboundStatusPending:
		return "待拣货"
	case constants.OutboundStatusPicking:
		return "复核中"
	case constants.OutboundStatusChecking:
		return "复核中"
	case constants.OutboundStatusPacking:
		return "打包中"
	case constants.OutboundStatusShipped:
		return "已发货"
	case constants.OutboundStatusCompleted:
		return "已完成"
	default:
		return status
	}
}

// StorageRequirementText 存储要求中文文本。
func StorageRequirementText(req string) string {
	switch req {
	case constants.StorageNormal:
		return "常温"
	case constants.StorageColdChain:
		return "冷链"
	case constants.StorageDangerous:
		return "危险品"
	case constants.StorageFragile:
		return "易碎"
	default:
		return req
	}
}

// QCResultText 质检结果中文文本。
func QCResultText(result string) string {
	switch result {
	case constants.QCResultPass:
		return "合格"
	case constants.QCResultFail:
		return "不合格"
	case constants.QCResultPartial:
		return "部分合格"
	default:
		return result
	}
}

// SettlementMethodText 结算方式中文文本。
func SettlementMethodText(method string) string {
	switch method {
	case constants.SettlementMonthly:
		return "月结"
	case constants.SettlementPerOrder:
		return "单结"
	case constants.SettlementPrepaid:
		return "预付"
	default:
		return method
	}
}

// OwnerStatusText 货主状态中文文本。
func OwnerStatusText(status string) string {
	if status == constants.OwnerStatusSuspended {
		return "已暂停"
	}
	return "合作中"
}

// BinStatusText 库位状态中文文本。
func BinStatusText(status string) string {
	switch status {
	case constants.BinStatusAvailable:
		return "可用"
	case constants.BinStatusOccupied:
		return "占用"
	case constants.BinStatusReserved:
		return "预留"
	case constants.BinStatusMaintenance:
		return "维护"
	default:
		return status
	}
}

// RoleText 角色中文文本。
func RoleText(role string) string {
	switch role {
	case constants.RoleAdmin:
		return "系统管理员"
	case constants.RoleWarehouseManager:
		return "仓库经理"
	case constants.RoleQCInspector:
		return "质检员"
	case constants.RolePicker:
		return "拣货员"
	case constants.RoleChecker:
		return "复核员"
	case constants.RoleOwner:
		return "货主"
	default:
		return role
	}
}

// JoinNames 拼接名称列表。
func JoinNames(names []string) string {
	return strings.Join(names, "、")
}

// FormatAmount 格式化金额，保留两位小数。
func FormatAmount(amount float64) string {
	return fmt.Sprintf("%.2f", amount)
}
