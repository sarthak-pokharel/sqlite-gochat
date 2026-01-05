package repositories

import (
	"fmt"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"

	"gorm.io/gorm"
)

// WebhookEventRepository interface
type WebhookEventRepository interface {
	Create(event *models.WebhookEvent) (*models.WebhookEvent, error)
	GetByID(id int64) (*models.WebhookEvent, error)
	ListUnprocessed(channelID int64, limit int) ([]*models.WebhookEvent, error)
	MarkProcessed(id int64) error
	MarkFailed(id int64, errorMsg string) error
}

type webhookEventRepository struct {
	db *gorm.DB
}

func NewWebhookEventRepository(db *gorm.DB) WebhookEventRepository {
	return &webhookEventRepository{db: db}
}

func (r *webhookEventRepository) Create(event *models.WebhookEvent) (*models.WebhookEvent, error) {
	if err := r.db.Create(event).Error; err != nil {
		return nil, fmt.Errorf("failed to create webhook event: %w", err)
	}
	return event, nil
}

func (r *webhookEventRepository) GetByID(id int64) (*models.WebhookEvent, error) {
	var event models.WebhookEvent
	if err := r.db.First(&event, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("webhook event not found")
		}
		return nil, fmt.Errorf("failed to get webhook event: %w", err)
	}
	return &event, nil
}

func (r *webhookEventRepository) ListUnprocessed(channelID int64, limit int) ([]*models.WebhookEvent, error) {
	var events []*models.WebhookEvent
	err := r.db.Where("channel_id = ? AND processed = ?", channelID, 0).
		Order("created_at ASC").
		Limit(limit).
		Find(&events).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list unprocessed events: %w", err)
	}
	return events, nil
}

func (r *webhookEventRepository) MarkProcessed(id int64) error {
	return r.db.Model(&models.WebhookEvent{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed":    1,
			"processed_at": gorm.Expr("CURRENT_TIMESTAMP"),
		}).Error
}

func (r *webhookEventRepository) MarkFailed(id int64, errorMsg string) error {
	return r.db.Model(&models.WebhookEvent{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"processed":    1,
			"processed_at": gorm.Expr("CURRENT_TIMESTAMP"),
			"error":        errorMsg,
		}).Error
}
