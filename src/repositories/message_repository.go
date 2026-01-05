package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
)

// MessageRepository interface
type MessageRepository interface {
	Create(msg *models.Message) (*models.Message, error)
	GetByID(id int64) (*models.Message, error)
	ListByConversation(conversationID int64, limit, offset int, before *int64) ([]*models.Message, error)
	UpdateStatus(id int64, status models.MessageStatus) error
}

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(msg *models.Message) (*models.Message, error) {
	query := `
		INSERT INTO messages (
			conversation_id, platform_message_id, sender_type, sender_id,
			content, message_type, media_url, direction, status,
			created_at, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		msg.ConversationID, msg.PlatformMessageID, msg.SenderType,
		msg.SenderID, msg.Content, msg.MessageType, msg.MediaURL,
		msg.Direction, msg.Status, msg.CreatedAt, msg.Metadata,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get message ID: %w", err)
	}

	return r.GetByID(id)
}

func (r *messageRepository) GetByID(id int64) (*models.Message, error) {
	query := `
		SELECT id, conversation_id, platform_message_id, sender_type, sender_id,
		       content, message_type, media_url, direction, status, created_at,
		       delivered_at, read_at, metadata
		FROM messages
		WHERE id = ?
	`
	msg := &models.Message{}
	err := r.db.QueryRow(query, id).Scan(
		&msg.ID, &msg.ConversationID, &msg.PlatformMessageID,
		&msg.SenderType, &msg.SenderID, &msg.Content, &msg.MessageType,
		&msg.MediaURL, &msg.Direction, &msg.Status, &msg.CreatedAt,
		&msg.DeliveredAt, &msg.ReadAt, &msg.Metadata,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("message not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return msg, nil
}

func (r *messageRepository) ListByConversation(conversationID int64, limit, offset int, before *int64) ([]*models.Message, error) {
	query := `
		SELECT id, conversation_id, platform_message_id, sender_type, sender_id,
		       content, message_type, media_url, direction, status, created_at,
		       delivered_at, read_at, metadata
		FROM messages
		WHERE conversation_id = ?
	`
	args := []interface{}{conversationID}

	if before != nil {
		query += " AND id < ?"
		args = append(args, *before)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer rows.Close()

	messages := []*models.Message{}
	for rows.Next() {
		msg := &models.Message{}
		err := rows.Scan(
			&msg.ID, &msg.ConversationID, &msg.PlatformMessageID,
			&msg.SenderType, &msg.SenderID, &msg.Content, &msg.MessageType,
			&msg.MediaURL, &msg.Direction, &msg.Status, &msg.CreatedAt,
			&msg.DeliveredAt, &msg.ReadAt, &msg.Metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *messageRepository) UpdateStatus(id int64, status models.MessageStatus) error {
	query := "UPDATE messages SET status = ?"
	args := []interface{}{status}

	now := time.Now()
	if status == models.MessageStatusDelivered {
		query += ", delivered_at = ?"
		args = append(args, now)
	} else if status == models.MessageStatusRead {
		query += ", read_at = ?"
		args = append(args, now)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}
