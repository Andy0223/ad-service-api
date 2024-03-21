package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaginationParams 从请求参数中获取分页参数
func GetPaginationParams(c *gin.Context) (limit int, offset int) {
	limit = 5  // 默认值
	offset = 0 // 默认值

	limitQuery := c.Query("limit")
	if limitQuery != "" {
		if i, err := strconv.Atoi(limitQuery); err == nil {
			limit = max(1, min(i, 100)) // 确保 limit 在 1 到 100 之间
		}
	}

	offsetQuery := c.Query("offset")
	if offsetQuery != "" {
		if i, err := strconv.Atoi(offsetQuery); err == nil {
			offset = max(0, i) // 确保 offset 不小于 0
		}
	}

	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
