package util

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateOrderNo 生成业务单号：<prefix>-YYYYMMDD-XXXX。
func GenerateOrderNo(prefix string) string {
	now := time.Now()
	n := rand.Intn(10000)
	return fmt.Sprintf("%s-%s-%04d", prefix, now.Format("20060102"), n)
}
