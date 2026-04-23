package handlers

import (
	"net/http"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type StudentPortalHandler struct {
	portalService services.StudentPortalService
}

func NewStudentPortalHandler(portalService services.StudentPortalService) *StudentPortalHandler {
	return &StudentPortalHandler{portalService: portalService}
}

func (h *StudentPortalHandler) GetProfile(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	profile, err := h.portalService.GetProfile(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_profile_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, profile)
}

func (h *StudentPortalHandler) GetEnrollments(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	enrollments, err := h.portalService.GetEnrollments(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_enrollments_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, enrollments)
}

func (h *StudentPortalHandler) GetTranscript(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	transcript, err := h.portalService.GetTranscript(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_transcript_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, transcript)
}

func (h *StudentPortalHandler) GetAcademicHistory(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	history, err := h.portalService.GetAcademicHistory(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_academic_history_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, history)
}

func (h *StudentPortalHandler) GetSchedule(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	schedule, err := h.portalService.GetSchedule(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_schedule_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, schedule)
}

func (h *StudentPortalHandler) ExportScheduleICS(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	ics, err := h.portalService.ExportScheduleICS(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "export_ics_failed", err.Error())
		return
	}
	c.Header("Content-Type", "text/calendar")
	c.Header("Content-Disposition", "attachment; filename=schedule.ics")
	c.Data(http.StatusOK, "text/calendar", ics)
}

func (h *StudentPortalHandler) GetAttendanceSummary(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	summary, err := h.portalService.GetAttendanceSummary(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_attendance_summary_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, summary)
}

func (h *StudentPortalHandler) JoinWaitlist(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	var input struct {
		SectionID uint `json:"section_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	waitlist, err := h.portalService.JoinWaitlist(studentID, input.SectionID)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "join_waitlist_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusCreated, waitlist)
}

func (h *StudentPortalHandler) GetWaitlist(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	waitlist, err := h.portalService.GetWaitlist(studentID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_waitlist_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, waitlist)
}

func (h *StudentPortalHandler) LeaveWaitlist(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	id, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	if err := h.portalService.LeaveWaitlist(studentID, id); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "leave_waitlist_failed", err.Error())
		return
	}
	apiutil.NoContent(c)
}

func (h *StudentPortalHandler) RateCourse(c *gin.Context) {
	studentID := c.MustGet("user_id").(uint)
	courseID, err := apiutil.ParseUintParam(c, "course_id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_course_id", err.Error())
		return
	}

	var input struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	rating, err := h.portalService.RateCourse(studentID, courseID, input.Rating, input.Comment)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "rate_course_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusCreated, rating)
}
