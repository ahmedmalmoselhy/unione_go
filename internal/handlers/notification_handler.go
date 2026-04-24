package handlers

import (
	"net/http"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/middlewares"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifSvc services.NotificationService
}

func NewNotificationHandler(notifSvc services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifSvc: notifSvc}
}

func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		apiutil.Error(c, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	unreadOnly := c.Query("unread") == "1"

	notifications, unreadCount, err := h.notifSvc.GetNotifications(userID, unreadOnly)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "fetch_notifications_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{
		"data": notifications,
		"meta": gin.H{
			"unread_count": unreadCount,
		},
	})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, err := apiutil.ParseID(c.Param("id"))
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", "Invalid notification ID")
		return
	}

	userID, exists := middlewares.GetUserID(c)
	if !exists {
		apiutil.Error(c, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	if err := h.notifSvc.MarkAsRead(id, userID); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "mark_read_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{"message": "Marked as read"})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, exists := middlewares.GetUserID(c)
	if !exists {
		apiutil.Error(c, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	if err := h.notifSvc.MarkAllAsRead(userID); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "mark_all_read_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{"message": "All marked as read"})
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id, err := apiutil.ParseID(c.Param("id"))
	if err != nil {
		apiutil.Error(c, http.StatusBadRequest, "invalid_id", "Invalid notification ID")
		return
	}

	userID, exists := middlewares.GetUserID(c)
	if !exists {
		apiutil.Error(c, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}

	if err := h.notifSvc.DeleteNotification(id, userID); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "delete_failed", err.Error())
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{"message": "Notification deleted"})
}
