package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
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
	db *sql.DB
}

func NewWebhookEventRepository(db *sql.DB) WebhookEventRepository {
	return &webhookEventRepository{db: db}
}

func (r *webhookEventRepository) Create(event *models.WebhookEvent) (*models.WebhookEvent, error) {
	query := `
		INSERT INTO webhook_events (channel_id, event_type, payload, created_at)
		VALUES (?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, event.ChannelID, event.EventType, event.Payload, event.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create webhook event: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get event ID: %w", err)
	}

	return r.GetByID(id)
}

func (r *webhookEventRepository) GetByID(id int64) (*models.WebhookEvent, error) {
	query := `
		SELECT id, channel_id, event_type, payload, processed, created_at, processed_at, error
		FROM webhook_events
		WHERE id = ?
	`
	event := &models.WebhookEvent{}
	err := r.db.QueryRow(query, id).Scan(
		&event.ID, &event.ChannelID, &event.EventType, &event.Payload,
		&event.Processed, &event.CreatedAt, &event.ProcessedAt, &event.Error,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("webhook event not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook event: %w", err)
	}
	return event, nil
}

func (r *webhookEventRepository) ListUnprocessed(channelID int64, limit int) ([]*models.WebhookEvent, error) {
	query := `
		SELECT id, channel_id, event_type, payload, processed, created_at, processed_at, error
		FROM webhook_events
		WHERE channel_id = ? AND processed = 0
		ORDER BY created_at ASC
		LIMIT ?
	`
	rows, err := r.db.Query(query, channelID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list unprocessed events: %w", err)
	}
	defer rows.Close()

	events := []*models.WebhookEvent{}
	for rows.Next() {
		event := &models.WebhookEvent{}
		if err := rows.Scan(
			&event.ID, &event.ChannelID, &event.EventType, &event.Payload,
			&event.Processed, &event.CreatedAt, &event.ProcessedAt, &event.Error,
		); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *webhookEventRepository) MarkProcessed(id int64) error {
	query := `
		UPDATE webhook_events
		SET processed = 1, processed_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark event as processed: %w", err)
	}
	return nil
}

func (r *webhookEventRepository) MarkFailed(id int64, errorMsg string) error {
	query := `
		UPDATE webhook_events
		SET processed = 1, processed_at = ?, error = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, time.Now(), errorMsg, id)
	if err != nil {
		return fmt.Errorf("failed to mark event as failed: %w", err)
	}
	return nil
}
