package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaginationParams 从请求参数中获取分页参数
func GetPaginationParams(c *gin.Context) (limit int, offset int, err error) {
	limit = 5  // Default value
	offset = 0 // Default value

	limitQuery := c.Query("limit")
	if limitQuery != "" {
		var i int
		i, err = strconv.Atoi(limitQuery)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid limit value: %v", err)
		}
		limit = max(1, min(i, 100)) // Ensure limit is between 1 and 100
	}

	offsetQuery := c.Query("offset")
	if offsetQuery != "" {
		var i int
		i, err = strconv.Atoi(offsetQuery)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid offset value: %v", err)
		}
		offset = max(0, i) // Ensure offset is not less than 0
	}

	return limit, offset, nil
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
