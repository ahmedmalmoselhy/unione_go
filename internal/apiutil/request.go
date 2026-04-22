package apiutil

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseUintParam(c *gin.Context, key string) (uint, error) {
	value := c.Param(key)
	if value == "" {
		return 0, fmt.Errorf("missing path parameter %s", key)
	}

	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid path parameter %s", key)
	}

	return uint(parsed), nil
}

func ParseRequiredUintQuery(c *gin.Context, key string) (uint, error) {
	value := c.Query(key)
	if value == "" {
		return 0, fmt.Errorf("missing query parameter %s", key)
	}

	parsed, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid query parameter %s", key)
	}

	return uint(parsed), nil
}
