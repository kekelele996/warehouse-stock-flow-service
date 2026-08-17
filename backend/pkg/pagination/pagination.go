package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/internal/constants"
)

// Page 分页参数。
type Page struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// Parse 从查询参数解析 page 与 page_size，带默认值与上限。
func Parse(c *gin.Context) Page {
	page := constants.DefaultPage
	pageSize := constants.DefaultPageSize
	if v, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.Atoi(c.DefaultQuery("page_size", "10")); err == nil && v > 0 {
		pageSize = v
	}
	if pageSize > constants.MaxPageSize {
		pageSize = constants.MaxPageSize
	}
	return Page{Page: page, PageSize: pageSize}
}

// Offset 计算 SQL offset。
func (p Page) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Result 分页返回结构。
type Result struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}
