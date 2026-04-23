package handlers

import (
	"net/http"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type ProfessorPortalHandler struct {
	portalService services.ProfessorPortalService
}

func NewProfessorPortalHandler(portalService services.ProfessorPortalService) *ProfessorPortalHandler {
	return &ProfessorPortalHandler{portalService: portalService}
}

func (h *ProfessorPortalHandler) GetProfile(c *gin.Context) {
	professorID := c.MustGet("user_id").(uint)
	profile, err := h.portalService.GetProfile(professorID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_profile_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, profile)
}

func (h *ProfessorPortalHandler) GetSections(c *gin.Context) {
	professorID := c.MustGet("user_id").(uint)
	sections, err := h.portalService.GetSections(professorID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_sections_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, sections)
}

func (h *ProfessorPortalHandler) GetSectionRoster(c *gin.Context) {
	professorID := c.MustGet("user_id").(uint)
	sectionID, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	roster, err := h.portalService.GetSectionRoster(professorID, sectionID)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "get_roster_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, roster)
}

func (h *ProfessorPortalHandler) UpdateGrade(c *gin.Context) {
	professorID := c.MustGet("user_id").(uint)
	sectionID, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	var input struct {
		StudentID uint    `json:"student_id" binding:"required"`
		Grade     float64 `json:"grade" binding:"required,min=0,max=100"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	enrollment, err := h.portalService.UpdateGrade(professorID, sectionID, input.StudentID, input.Grade)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "update_grade_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, enrollment)
}

func (h *ProfessorPortalHandler) RecordAttendance(c *gin.Context) {
	professorID := c.MustGet("user_id").(uint)
	sectionID, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	var input struct {
		Date            time.Time    `json:"date" binding:"required"`
		StudentPresence map[uint]bool `json:"student_presence" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	if err := h.portalService.RecordAttendance(professorID, sectionID, input.Date, input.StudentPresence); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "record_attendance_failed", err.Error())
		return
	}
	apiutil.NoContent(c)
}

func (h *ProfessorPortalHandler) CreateSectionAnnouncement(c *gin.Context) {
	professorID := c.MustGet("user_id").(uint)
	sectionID, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", err.Error())
		return
	}

	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	announcement, err := h.portalService.CreateSectionAnnouncement(professorID, sectionID, input.Title, input.Content)
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "create_announcement_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusCreated, announcement)
}

func (h *ProfessorPortalHandler) GetSchedule(c *gin.Context) {
	professorID := c.MustGet("user_id").(uint)
	schedule, err := h.portalService.GetSchedule(professorID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "get_schedule_failed", err.Error())
		return
	}
	apiutil.Success(c, http.StatusOK, schedule)
}
