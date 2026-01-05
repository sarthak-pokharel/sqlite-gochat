package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
)

type ChannelRepository interface {
	Create(req *models.CreateChannelRequest) (*models.ChatChannel, error)
	GetByID(id int64) (*models.ChatChannel, error)
	ListByOrganization(orgID int64, limit, offset int) ([]*models.ChatChannel, error)
	Update(id int64, req *models.UpdateChannelRequest) error
	UpdateStatus(id int64, status models.ChannelStatus) error
	Delete(id int64) error
}

type channelRepository struct {
	db *sql.DB
}

func NewChannelRepository(db *sql.DB) ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) Create(req *models.CreateChannelRequest) (*models.ChatChannel, error) {
	query := `
		INSERT INTO chat_channels (
			organization_id, platform, name, account_identifier,
			status, webhook_secret, access_token, config, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.Exec(query,
		req.OrganizationID, req.Platform, req.Name, req.AccountIdentifier,
		models.ChannelStatusPending, req.WebhookSecret, req.AccessToken,
		req.Config, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get channel ID: %w", err)
	}

	return r.GetByID(id)
}

func (r *channelRepository) GetByID(id int64) (*models.ChatChannel, error) {
	query := `
		SELECT id, organization_id, platform, name, account_identifier, status,
		       webhook_secret, access_token, config, created_at, updated_at,
		       last_message_at, is_active
		FROM chat_channels
		WHERE id = ?
	`
	ch := &models.ChatChannel{}
	err := r.db.QueryRow(query, id).Scan(
		&ch.ID, &ch.OrganizationID, &ch.Platform, &ch.Name,
		&ch.AccountIdentifier, &ch.Status, &ch.WebhookSecret,
		&ch.AccessToken, &ch.Config, &ch.CreatedAt, &ch.UpdatedAt,
		&ch.LastMessageAt, &ch.IsActive,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("channel not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}
	return ch, nil
}

func (r *channelRepository) ListByOrganization(orgID int64, limit, offset int) ([]*models.ChatChannel, error) {
	query := `
		SELECT id, organization_id, platform, name, account_identifier, status,
		       webhook_secret, config, created_at, updated_at,
		       last_message_at, is_active
		FROM chat_channels
		WHERE organization_id = ? AND is_active = 1
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}
	defer rows.Close()

	channels := []*models.ChatChannel{}
	for rows.Next() {
		ch := &models.ChatChannel{}
		err := rows.Scan(
			&ch.ID, &ch.OrganizationID, &ch.Platform, &ch.Name,
			&ch.AccountIdentifier, &ch.Status, &ch.WebhookSecret,
			&ch.Config, &ch.CreatedAt, &ch.UpdatedAt,
			&ch.LastMessageAt, &ch.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}

		ch.AccessToken = nil
		channels = append(channels, ch)
	}
	return channels, nil
}

func (r *channelRepository) Update(id int64, req *models.UpdateChannelRequest) error {
	query := "UPDATE chat_channels SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.Name != nil {
		query += ", name = ?"
		args = append(args, *req.Name)
	}
	if req.Status != nil {
		query += ", status = ?"
		args = append(args, *req.Status)
	}
	if req.WebhookSecret != nil {
		query += ", webhook_secret = ?"
		args = append(args, *req.WebhookSecret)
	}
	if req.AccessToken != nil {
		query += ", access_token = ?"
		args = append(args, *req.AccessToken)
	}
	if req.Config != nil {
		query += ", config = ?"
		args = append(args, *req.Config)
	}
	if req.IsActive != nil {
		query += ", is_active = ?"
		args = append(args, *req.IsActive)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update channel: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("channel not found")
	}

	return nil
}

func (r *channelRepository) UpdateStatus(id int64, status models.ChannelStatus) error {
	query := "UPDATE chat_channels SET status = ?, updated_at = ? WHERE id = ?"
	result, err := r.db.Exec(query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update channel status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("channel not found")
	}

	return nil
}

func (r *channelRepository) Delete(id int64) error {
	query := "UPDATE chat_channels SET is_active = 0 WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete channel: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("channel not found")
	}

	return nil
}
