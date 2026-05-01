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

// GetDashboardStats godoc
// @Summary Get dashboard statistics
// @Description Returns overview stats, enrollment status, grade distribution, GPA distribution, and section fill rates. Scoped by faculty or department if provided.
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param faculty_id query int false "Filter by faculty ID"
// @Param department_id query int false "Filter by department ID"
// @Success 200 {object} models.DashboardStats
// @Failure 401 {object} apiutil.errorResponse
// @Failure 500 {object} apiutil.errorResponse
// @Router /admin/dashboard/stats [get]
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

// GetEnrollmentTrends godoc
// @Summary Get enrollment trends
// @Description Returns enrollment data over time (monthly) for the specified period.
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param months query int false "Period in months (default 12)"
// @Success 200 {object} models.EnrollmentTrendsResponse
// @Failure 401 {object} apiutil.errorResponse
// @Failure 500 {object} apiutil.errorResponse
// @Router /admin/analytics/enrollment-trends [get]
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

// PredictStudentPerformance godoc
// @Summary Predict student performance
// @Description Returns a performance prediction for a specific student based on their GPA and historical data.
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param student_id path int true "Student ID"
// @Success 200 {object} models.StudentPerformancePrediction
// @Failure 400 {object} apiutil.errorResponse
// @Failure 401 {object} apiutil.errorResponse
// @Failure 500 {object} apiutil.errorResponse
// @Router /admin/analytics/student-performance/{student_id} [get]
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

// GetCourseDemand godoc
// @Summary Get course demand analysis
// @Description Returns analysis of course demand based on current enrollments and capacities.
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.CourseDemandResponse
// @Failure 401 {object} apiutil.errorResponse
// @Failure 500 {object} apiutil.errorResponse
// @Router /admin/analytics/course-demand [get]
func (h *AnalyticsHandler) GetCourseDemand(c *gin.Context) {
	demand, err := h.analyticsService.GetCourseDemand()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "analytics_error", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, demand)
}

// GetProfessorWorkload godoc
// @Summary Get professor workload analysis
// @Description Returns workload stats for professors, including section counts and hours.
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.ProfessorWorkloadResponse
// @Failure 401 {object} apiutil.errorResponse
// @Failure 500 {object} apiutil.errorResponse
// @Router /admin/analytics/professor-workload [get]
func (h *AnalyticsHandler) GetProfessorWorkload(c *gin.Context) {
	workload, err := h.analyticsService.GetProfessorWorkload()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "analytics_error", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, workload)
}

// GetAttendanceAnalytics godoc
// @Summary Get attendance analytics
// @Description Returns overall attendance rate and status breakdown.
// @Tags analytics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.AttendanceAnalyticsResponse
// @Failure 401 {object} apiutil.errorResponse
// @Failure 500 {object} apiutil.errorResponse
// @Router /admin/analytics/attendance [get]
func (h *AnalyticsHandler) GetAttendanceAnalytics(c *gin.Context) {
	analytics, err := h.analyticsService.GetAttendanceAnalytics()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "analytics_error", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, analytics)
}
