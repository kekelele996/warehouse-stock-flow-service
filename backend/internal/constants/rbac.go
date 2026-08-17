package constants

// UserRole 系统角色。
const (
	RoleAdmin            = "Admin"            // 系统管理员
	RoleWarehouseManager = "WarehouseManager" // 仓库经理
	RoleQCInspector      = "QCInspector"      // 质检员
	RolePicker           = "Picker"           // 拣货员
	RoleChecker          = "Checker"          // 复核员
	RoleOwner            = "Owner"            // 货主
)

// AllRoles 全部角色集合。
var AllRoles = []string{RoleAdmin, RoleWarehouseManager, RoleQCInspector, RolePicker, RoleChecker, RoleOwner}

// UserStatus 用户状态。
const (
	UserStatusActive   = "Active"
	UserStatusDisabled = "Disabled"
)

// Permission 权限点。
const (
	PermOwnerManage     = "owner:manage"     // 货主管理
	PermProductManage   = "product:manage"   // 商品管理
	PermInboundManage   = "inbound:manage"   // 入库管理
	PermOutboundManage  = "outbound:manage"  // 出库管理
	PermBinManage       = "bin:manage"       // 库位管理
	PermAuditView       = "audit:view"       // 操作日志查看
	PermDashboardView   = "dashboard:view"   // 总览查看
)

// RolePermissions 角色 → 权限点映射。
var RolePermissions = map[string][]string{
	RoleAdmin:            {PermOwnerManage, PermProductManage, PermInboundManage, PermOutboundManage, PermBinManage, PermAuditView, PermDashboardView},
	RoleWarehouseManager: {PermOwnerManage, PermProductManage, PermInboundManage, PermOutboundManage, PermBinManage, PermAuditView, PermDashboardView},
	RoleQCInspector:      {PermInboundManage, PermOutboundManage, PermBinManage, PermDashboardView},
	RolePicker:           {PermOutboundManage, PermBinManage, PermDashboardView},
	RoleChecker:          {PermOutboundManage, PermDashboardView},
	RoleOwner:            {PermProductManage, PermInboundManage, PermOutboundManage, PermDashboardView},
}

// HasPermission 判断角色是否拥有权限点。
func HasPermission(role, perm string) bool {
	for _, p := range RolePermissions[role] {
		if p == perm {
			return true
		}
	}
	return false
}
