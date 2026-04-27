package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ahmedmalmoselhy/unione_go/internal/models"
	"gorm.io/gorm"
)

type WebhookService struct {
	db *gorm.DB
}

func NewWebhookService(db *gorm.DB) *WebhookService {
	return &WebhookService{db: db}
}

func (s *WebhookService) CreateWebhook(url, event string) (*models.Webhook, error) {
	// Generate a random secret
	secret := generateRandomString(32)

	webhook := &models.Webhook{
		URL:    url,
		Event:  event,
		Secret: secret,
	}

	err := s.db.Create(webhook).Error
	return webhook, err
}

func (s *WebhookService) ListWebhooks() ([]models.Webhook, error) {
	var webhooks []models.Webhook
	err := s.db.Find(&webhooks).Error
	return webhooks, err
}

func (s *WebhookService) DeleteWebhook(id uint) error {
	return s.db.Delete(&models.Webhook{}, id).Error
}

func (s *WebhookService) Trigger(event string, payload interface{}) {
	go func() {
		var webhooks []models.Webhook
		// Find active webhooks for this event or all events (*)
		err := s.db.Where("is_active = ? AND (event = ? OR event = ?)", true, event, "*").Find(&webhooks).Error
		if err != nil {
			log.Printf("Webhook trigger error (finding webhooks): %v", err)
			return
		}

		for _, wh := range webhooks {
			s.deliver(wh, event, payload)
		}
	}()
}

func (s *WebhookService) deliver(wh models.Webhook, event string, payload interface{}) {
	payloadBytes, err := json.Marshal(map[string]interface{}{
		"event":     event,
		"payload":   payload,
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		log.Printf("Webhook deliver error (marshaling payload): %v", err)
		return
	}

	start := time.Now()
	
	req, err := http.NewRequest("POST", wh.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("Webhook deliver error (creating request): %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "UniOne-Webhook/1.0")
	req.Header.Set("X-UniOne-Event", event)
	
	// Sign the payload
	signature := signPayload(payloadBytes, wh.Secret)
	req.Header.Set("X-UniOne-Signature", signature)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	duration := time.Since(start).Milliseconds()

	delivery := models.WebhookDelivery{
		WebhookID: wh.ID,
		Event:     event,
		Payload:   string(payloadBytes),
		Duration:  duration,
	}

	if err != nil {
		delivery.ErrorMessage = err.Error()
		delivery.StatusCode = 0
		s.handleFailure(wh)
	} else {
		defer resp.Body.Close()
		delivery.StatusCode = resp.StatusCode
		bodyBytes, _ := io.ReadAll(resp.Body)
		delivery.ResponseBody = string(bodyBytes)

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			s.handleSuccess(wh)
		} else {
			s.handleFailure(wh)
		}
	}

	if err := s.db.Create(&delivery).Error; err != nil {
		log.Printf("Webhook deliver error (saving delivery): %v", err)
	}
}

func (s *WebhookService) handleSuccess(wh models.Webhook) {
	if wh.FailCount > 0 {
		s.db.Model(&wh).Update("fail_count", 0)
	}
}

func (s *WebhookService) handleFailure(wh models.Webhook) {
	newFailCount := wh.FailCount + 1
	updates := map[string]interface{}{
		"fail_count": newFailCount,
	}
	
	// Auto-disable after 10 failures
	if newFailCount >= 10 {
		updates["is_active"] = false
		log.Printf("Webhook %d auto-disabled after %d failures", wh.ID, newFailCount)
	}

	s.db.Model(&wh).Updates(updates)
}

func signPayload(payload []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func generateRandomString(n int) string {
	// Simple random string for secret
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

func (s *WebhookService) GetDeliveries(webhookID uint, limit, offset int) ([]models.WebhookDelivery, int64, error) {
	var deliveries []models.WebhookDelivery
	var total int64

	query := s.db.Model(&models.WebhookDelivery{})
	if webhookID > 0 {
		query = query.Where("webhook_id = ?", webhookID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at desc").Limit(limit).Offset(offset).Find(&deliveries).Error
	return deliveries, total, err
}
