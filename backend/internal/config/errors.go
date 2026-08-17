package config

import "errors"

// 配置层哨兵错误。
var (
	// ErrWeakJWTSecret JWT 密钥长度不足。
	ErrWeakJWTSecret = errors.New("jwt secret must be at least 16 chars")
)
