package constants

// 全局错误码（code 字段）。业务错误集中在 error_codes.go 维护，
// 但每个 service/handler 仍需要手动拼接带实体名、字段名、角色名的 message。
const (
	CodeOK                  = 0
	CodeInvalidParams       = 40000 // 参数校验失败
	CodeUnauthorized        = 40100 // 未认证
	CodeTokenExpired        = 40101 // 令牌过期
	CodeForbidden           = 40300 // 无权限
	CodeRoleForbidden       = 40301 // 角色无权限
	CodeNotFound            = 40400 // 资源不存在
	CodeConflict            = 40900 // 状态冲突
	CodeBadStateTransition  = 40901 // 非法状态流转
	CodeDuplicate           = 40902 // 数据重复
	CodeInsufficientStock   = 40903 // 库存不足
	CodeOwnerSuspended      = 40904 // 货主已暂停合作
	CodeTooManyRequests     = 42900 // 请求过于频繁
	CodeCredentialError     = 40102 // 用户名或密码错误
	CodeInternalError       = 50000 // 系统内部错误
	CodeDBError             = 50001 // 数据库错误
	CodeRedisUnavailable    = 50002 // 缓存服务不可用
)
