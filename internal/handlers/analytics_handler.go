package handlers

import (
	"net/http"
	"strconv"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetDashboardStats(c *gin.Context) {
	// Scoping logic (faculty_id, department_id from query if admin/employee)
	var facultyID, departmentID *uint

	fIDStr := c.Query("faculty_id")
	if fIDStr != "" {
		id, err := strconv.ParseUint(fIDStr, 10, 32)
		if err == nil {
			val := uint(id)
			facultyID = &val
		}
	}

	dIDStr := c.Query("department_id")
	if dIDStr != "" {
		id, err := strconv.ParseUint(dIDStr, 10, 32)
		if err == nil {
			val := uint(id)
			departmentID = &val
		}
	}

	stats, err := h.analyticsService.GetDashboardStats(facultyID, departmentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "analytics_error", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetEnrollmentTrends(c *gin.Context) {
	months := 12
	monthsStr := c.DefaultQuery("months", "12")
	if m, err := strconv.Atoi(monthsStr); err == nil {
		months = m
	}

	trends, err := h.analyticsService.GetEnrollmentTrends(months)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "analytics_error", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, trends)
}

func (h *AnalyticsHandler) PredictStudentPerformance(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	studentID, err := strconv.ParseUint(studentIDStr, 10, 32)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", "Invalid student ID")
		return
	}

	prediction, err := h.analyticsService.PredictStudentPerformance(uint(studentID))
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "analytics_error", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, prediction)
}

func (h *AnalyticsHandler) GetCourseDemand(c *gin.Context) {
	demand, err := h.analyticsService.GetCourseDemand()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "analytics_error", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, demand)
}

func (h *AnalyticsHandler) GetProfessorWorkload(c *gin.Context) {
	workload, err := h.analyticsService.GetProfessorWorkload()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "analytics_error", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, workload)
}

func (h *AnalyticsHandler) GetAttendanceAnalytics(c *gin.Context) {
	analytics, err := h.analyticsService.GetAttendanceAnalytics()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "analytics_error", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, analytics)
}
