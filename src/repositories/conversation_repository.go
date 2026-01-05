package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
)

// ConversationRepository interface
type ConversationRepository interface {
	Create(req *models.CreateConversationRequest) (*models.Conversation, error)
	GetByID(id int64) (*models.Conversation, error)
	GetOrCreateByUser(channelID, externalUserID int64) (*models.Conversation, error)
	List(channelID int64, status *models.ConversationStatus, limit, offset int) ([]*models.Conversation, error)
	Update(id int64, req *models.UpdateConversationRequest) error
	UpdateLastMessage(id int64) error
}

type conversationRepository struct {
	db *sql.DB
}

func NewConversationRepository(db *sql.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) Create(req *models.CreateConversationRequest) (*models.Conversation, error) {
	priority := req.Priority
	if priority == "" {
		priority = models.PriorityNormal
	}

	query := `
		INSERT INTO conversations (
			channel_id, external_user_id, status, priority, subject, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.Exec(query,
		req.ChannelID, req.ExternalUserID, models.ConversationStatusOpen,
		priority, req.Subject, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation ID: %w", err)
	}

	return r.GetByID(id)
}

func (r *conversationRepository) GetByID(id int64) (*models.Conversation, error) {
	query := `
		SELECT id, channel_id, external_user_id, assigned_to_external_id,
		       status, priority, subject, first_message_at, last_message_at,
		       resolved_at, created_at, updated_at, metadata
		FROM conversations
		WHERE id = ?
	`
	conv := &models.Conversation{}
	err := r.db.QueryRow(query, id).Scan(
		&conv.ID, &conv.ChannelID, &conv.ExternalUserID,
		&conv.AssignedToExternalID, &conv.Status, &conv.Priority,
		&conv.Subject, &conv.FirstMessageAt, &conv.LastMessageAt,
		&conv.ResolvedAt, &conv.CreatedAt, &conv.UpdatedAt, &conv.Metadata,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("conversation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	return conv, nil
}

func (r *conversationRepository) GetOrCreateByUser(channelID, externalUserID int64) (*models.Conversation, error) {
	// Check for existing open conversation
	query := `
		SELECT id, channel_id, external_user_id, assigned_to_external_id,
		       status, priority, subject, first_message_at, last_message_at,
		       resolved_at, created_at, updated_at, metadata
		FROM conversations
		WHERE channel_id = ? AND external_user_id = ? AND status IN ('open', 'pending')
		ORDER BY created_at DESC
		LIMIT 1
	`
	conv := &models.Conversation{}
	err := r.db.QueryRow(query, channelID, externalUserID).Scan(
		&conv.ID, &conv.ChannelID, &conv.ExternalUserID,
		&conv.AssignedToExternalID, &conv.Status, &conv.Priority,
		&conv.Subject, &conv.FirstMessageAt, &conv.LastMessageAt,
		&conv.ResolvedAt, &conv.CreatedAt, &conv.UpdatedAt, &conv.Metadata,
	)
	if err == nil {
		return conv, nil
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// No open conversation, create new
	return r.Create(&models.CreateConversationRequest{
		ChannelID:      channelID,
		ExternalUserID: externalUserID,
	})
}

func (r *conversationRepository) List(channelID int64, status *models.ConversationStatus, limit, offset int) ([]*models.Conversation, error) {
	query := `
		SELECT id, channel_id, external_user_id, assigned_to_external_id,
		       status, priority, subject, first_message_at, last_message_at,
		       resolved_at, created_at, updated_at, metadata
		FROM conversations
		WHERE channel_id = ?
	`
	args := []interface{}{channelID}

	if status != nil {
		query += " AND status = ?"
		args = append(args, *status)
	}

	query += " ORDER BY updated_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}
	defer rows.Close()

	conversations := []*models.Conversation{}
	for rows.Next() {
		conv := &models.Conversation{}
		err := rows.Scan(
			&conv.ID, &conv.ChannelID, &conv.ExternalUserID,
			&conv.AssignedToExternalID, &conv.Status, &conv.Priority,
			&conv.Subject, &conv.FirstMessageAt, &conv.LastMessageAt,
			&conv.ResolvedAt, &conv.CreatedAt, &conv.UpdatedAt, &conv.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan conversation: %w", err)
		}
		conversations = append(conversations, conv)
	}
	return conversations, nil
}

func (r *conversationRepository) Update(id int64, req *models.UpdateConversationRequest) error {
	query := "UPDATE conversations SET updated_at = ?"
	args := []interface{}{time.Now()}

	if req.AssignedToExternalID != nil {
		query += ", assigned_to_external_id = ?"
		args = append(args, *req.AssignedToExternalID)
	}
	if req.Status != nil {
		query += ", status = ?"
		args = append(args, *req.Status)
		if *req.Status == models.ConversationStatusResolved || *req.Status == models.ConversationStatusClosed {
			query += ", resolved_at = ?"
			args = append(args, time.Now())
		}
	}
	if req.Priority != nil {
		query += ", priority = ?"
		args = append(args, *req.Priority)
	}
	if req.Subject != nil {
		query += ", subject = ?"
		args = append(args, *req.Subject)
	}
	if req.Metadata != nil {
		query += ", metadata = ?"
		args = append(args, *req.Metadata)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update conversation: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

func (r *conversationRepository) UpdateLastMessage(id int64) error {
	now := time.Now()
	query := `
		UPDATE conversations
		SET last_message_at = ?, updated_at = ?,
		    first_message_at = COALESCE(first_message_at, ?)
		WHERE id = ?
	`
	_, err := r.db.Exec(query, now, now, now, id)
	return err
}
