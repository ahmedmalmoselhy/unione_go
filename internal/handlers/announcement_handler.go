package handlers

import (
	"net/http"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/middlewares"
	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	notifSvc services.NotificationService
}

func NewAnnouncementHandler(notifSvc services.NotificationService) *AnnouncementHandler {
	return &AnnouncementHandler{notifSvc: notifSvc}
}

type CreateAnnouncementInput struct {
	Title     string                  `json:"title" binding:"required"`
	Content   string                  `json:"content" binding:"required"`
	Type      models.AnnouncementType `json:"type" binding:"required"`
	FacultyID *uint                   `json:"faculty_id"`
	SectionID *uint                   `json:"section_id"`
}

func (h *AnnouncementHandler) CreateAnnouncement(c *gin.Context) {
	var input CreateAnnouncementInput
	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "validation_error", err.Error())
		return
	}

	userID, exists := middlewares.GetUserID(c)
	if !exists {
		apiutil.Error(c, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	announcement := &models.Announcement{
		Title:     input.Title,
		Content:   input.Content,
		Type:      input.Type,
		FacultyID: input.FacultyID,
		SectionID: input.SectionID,
		AuthorID:  userID,
	}

	if err := h.notifSvc.CreateAnnouncement(announcement); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "create_announcement_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusCreated, announcement)
}

func (h *AnnouncementHandler) ListAnnouncements(c *gin.Context) {
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		apiutil.Error(c, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	announcements, err := h.notifSvc.ListAnnouncements(userID)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "fetch_announcements_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, announcements)
}

func (h *AnnouncementHandler) MarkRead(c *gin.Context) {
	announcementID, err := apiutil.ParseUintParam(c, "id")
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", "Invalid announcement ID")
		return
	}

	userID, exists := middlewares.GetUserID(c)
	if !exists {
		apiutil.Error(c, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	if err := h.notifSvc.MarkAnnouncementRead(announcementID, userID); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "mark_read_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{"message": "Marked as read"})
}
