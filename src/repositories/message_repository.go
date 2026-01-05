package repositories

import (
	"fmt"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"

	"gorm.io/gorm"
)

// MessageRepository interface
type MessageRepository interface {
	Create(msg *models.Message) (*models.Message, error)
	GetByID(id int64) (*models.Message, error)
	ListByConversation(conversationID int64, limit, offset int, before *int64) ([]*models.Message, error)
	UpdateStatus(id int64, status models.MessageStatus) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(msg *models.Message) (*models.Message, error) {
	if err := r.db.Create(msg).Error; err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}
	return msg, nil
}

func (r *messageRepository) GetByID(id int64) (*models.Message, error) {
	var msg models.Message
	if err := r.db.First(&msg, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return &msg, nil
}

func (r *messageRepository) ListByConversation(conversationID int64, limit, offset int, before *int64) ([]*models.Message, error) {
	var messages []*models.Message

	query := r.db.Where("conversation_id = ?", conversationID)

	if before != nil {
		query = query.Where("id < ?", *before)
	}

	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}

	return messages, nil
}

func (r *messageRepository) UpdateStatus(id int64, status models.MessageStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == models.MessageStatusDelivered {
		updates["delivered_at"] = gorm.Expr("CURRENT_TIMESTAMP")
	} else if status == models.MessageStatusRead {
		updates["read_at"] = gorm.Expr("CURRENT_TIMESTAMP")
	}

	result := r.db.Model(&models.Message{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update message status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}
