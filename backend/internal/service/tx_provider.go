package service

import (
	"database/sql"

	"gorm.io/gorm"
)

// txProvider 事务执行接口：便于单元测试注入 fake 事务。
type txProvider interface {
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error
}
