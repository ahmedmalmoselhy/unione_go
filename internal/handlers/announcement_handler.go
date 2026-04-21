package handlers

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("user_id")
	
	announcement := &models.Announcement{
		Title:     input.Title,
		Content:   input.Content,
		Type:      input.Type,
		FacultyID: input.FacultyID,
		SectionID: input.SectionID,
		AuthorID:  userID.(uint),
	}

	if err := h.notifSvc.CreateAnnouncement(announcement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, announcement)
}
