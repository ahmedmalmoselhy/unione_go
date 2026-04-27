package handlers

import (
	"net/http"
	"strconv"

	"github.com/ahmedmalmoselhy/unione_go/internal/apiutil"
	"github.com/ahmedmalmoselhy/unione_go/internal/services"
	"github.com/gin-gonic/gin"
)

type GovernanceHandler struct {
	auditService   *services.AuditService
	webhookService *services.WebhookService
}

func NewGovernanceHandler(auditService *services.AuditService, webhookService *services.WebhookService) *GovernanceHandler {
	return &GovernanceHandler{
		auditService:   auditService,
		webhookService: webhookService,
	}
}

// Audit Logs

func (h *GovernanceHandler) GetAuditLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, total, err := h.auditService.GetLogs(limit, offset)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch audit logs")
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{
		"items":  logs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// Webhooks

func (h *GovernanceHandler) CreateWebhook(c *gin.Context) {
	var input struct {
		URL   string `json:"url" binding:"required,url"`
		Event string `json:"event" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		apiutil.Error(c, http.StatusBadRequest, "INVALID_INPUT", err.Error())
		return
	}

	webhook, err := h.webhookService.CreateWebhook(input.URL, input.Event)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to create webhook")
		return
	}

	apiutil.Success(c, http.StatusCreated, gin.H{
		"message": "Webhook created successfully",
		"webhook": webhook,
		"secret":  webhook.Secret, // Only shown once
	})
}

func (h *GovernanceHandler) ListWebhooks(c *gin.Context) {
	webhooks, err := h.webhookService.ListWebhooks()
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch webhooks")
		return
	}

	apiutil.Success(c, http.StatusOK, webhooks)
}

func (h *GovernanceHandler) DeleteWebhook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.webhookService.DeleteWebhook(uint(id)); err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to delete webhook")
		return
	}

	apiutil.Message(c, http.StatusOK, "Webhook deleted successfully")
}

func (h *GovernanceHandler) GetWebhookDeliveries(c *gin.Context) {
	webhookID, _ := strconv.Atoi(c.Param("id"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	deliveries, total, err := h.webhookService.GetDeliveries(uint(webhookID), limit, offset)
	if err != nil {
		apiutil.Error(c, http.StatusInternalServerError, "DB_ERROR", "Failed to fetch deliveries")
		return
	}

	apiutil.Success(c, http.StatusOK, gin.H{
		"items":  deliveries,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}
