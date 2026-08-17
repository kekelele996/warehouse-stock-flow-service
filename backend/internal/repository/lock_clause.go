package repository

import "gorm.io/gorm/clause"

// LockClause 返回 FOR UPDATE 行锁子句。
func LockClause() clause.Locking {
	return clause.Locking{Strength: "UPDATE"}
}
