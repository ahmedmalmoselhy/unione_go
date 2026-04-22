package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Authorization header missing")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Invalid authorization format")
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Failed to parse token claims")
			return
		}

		userID, err := claimUint(claims, "sub")
		if err != nil {
			apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Invalid token subject")
			return
		}

		roleValue, ok := claims["role"]
		if !ok {
			apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Missing token role")
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Invalid token role")
			return
		}

		c.Set("user_id", userID)
		c.Set("role", models.Role(role))

		if _, exists := claims["faculty_id"]; exists {
			facultyID, err := claimUint(claims, "faculty_id")
			if err != nil {
				apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Invalid faculty scope")
				return
			}
			c.Set("faculty_id", facultyID)
		}

		if _, exists := claims["department_id"]; exists {
			departmentID, err := claimUint(claims, "department_id")
			if err != nil {
				apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Invalid department scope")
				return
			}
			c.Set("department_id", departmentID)
		}

		c.Next()
	}
}

func RequireFacultyScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := GetRole(c)
		if exists && role == models.RoleAdmin {
			c.Next()
			return
		}

		userFacultyID, exists := GetFacultyID(c)
		if !exists {
			apiutil.AbortError(c, http.StatusForbidden, "forbidden", "Missing faculty scope")
			return
		}

		reqFacultyID := c.Param("faculty_id")
		if fmt.Sprintf("%d", userFacultyID) != reqFacultyID {
			apiutil.AbortError(c, http.StatusForbidden, "forbidden", "Unauthorized faculty scope")
			return
		}

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := GetRole(c)
		if !exists {
			apiutil.AbortError(c, http.StatusUnauthorized, "unauthorized", "Authentication required")
			return
		}

		for _, role := range roles {
			if userRole == models.Role(role) {
				c.Next()
				return
			}
		}

		apiutil.AbortError(c, http.StatusForbidden, "forbidden", "Insufficient permissions")
	}
}
