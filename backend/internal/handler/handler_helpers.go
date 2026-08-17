package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/wmsflow/wmsflow/pkg/pagination"
)

// parsePage 解析分页参数。
func parsePage(c *gin.Context) pagination.Page {
	return pagination.Parse(c)
}
