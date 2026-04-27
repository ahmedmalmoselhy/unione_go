package handlers

import (
	"net/http"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := "UP"
	dbStatus := "UP"

	sqlDB, err := h.db.DB()
	if err != nil {
		status = "DEGRADED"
		dbStatus = "DOWN (Connection Error)"
	} else {
		err = sqlDB.Ping()
		if err != nil {
			status = "DEGRADED"
			dbStatus = "DOWN (Ping Error)"
		}
	}

	apiutil.Success(c, http.StatusOK, gin.H{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"checks": gin.H{
			"database": dbStatus,
		},
	})
}
