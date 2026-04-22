package middlewares

import (
	"fmt"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}

func GetRole(c *gin.Context) (models.Role, bool) {
	value, exists := c.Get("role")
	if !exists {
		return "", false
	}

	role, ok := value.(models.Role)
	return role, ok
}

func GetFacultyID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("faculty_id")
	if !exists {
		return 0, false
	}

	facultyID, ok := value.(uint)
	return facultyID, ok
}

func GetDepartmentID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("department_id")
	if !exists {
		return 0, false
	}

	departmentID, ok := value.(uint)
	return departmentID, ok
}

func claimUint(claims map[string]any, key string) (uint, error) {
	raw, exists := claims[key]
	if !exists {
		return 0, fmt.Errorf("claim %s missing", key)
	}

	switch v := raw.(type) {
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case uint:
		return v, nil
	default:
		return 0, fmt.Errorf("claim %s has invalid type %T", key, raw)
	}
}
