package repository

import "sync/atomic"

var seq int64

// dbTxSeq 返回递增序号，用于生成唯一测试数据。
func dbTxSeq() int64 { return atomic.AddInt64(&seq, 1) }
