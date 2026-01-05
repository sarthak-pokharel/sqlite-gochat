package repositories

import (
	"fmt"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"

	"gorm.io/gorm"
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
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db: db}
}

func (r *conversationRepository) Create(req *models.CreateConversationRequest) (*models.Conversation, error) {
	priority := req.Priority
	if priority == "" {
		priority = models.PriorityNormal
	}

	conv := &models.Conversation{
		ChannelID:      req.ChannelID,
		ExternalUserID: req.ExternalUserID,
		Status:         models.ConversationStatusOpen,
		Priority:       priority,
		Subject:        req.Subject,
	}

	if err := r.db.Create(conv).Error; err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	return conv, nil
}

func (r *conversationRepository) GetByID(id int64) (*models.Conversation, error) {
	var conv models.Conversation
	if err := r.db.First(&conv, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("conversation not found")
		}
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}
	return &conv, nil
}

func (r *conversationRepository) GetOrCreateByUser(channelID, externalUserID int64) (*models.Conversation, error) {
	// Check for existing open conversation
	var conv models.Conversation
	err := r.db.Where("channel_id = ? AND external_user_id = ? AND status IN ?",
		channelID, externalUserID, []models.ConversationStatus{models.ConversationStatusOpen, models.ConversationStatusPending}).
		Order("created_at DESC").
		First(&conv).Error

	if err == nil {
		return &conv, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// No open conversation, create new
	return r.Create(&models.CreateConversationRequest{
		ChannelID:      channelID,
		ExternalUserID: externalUserID,
	})
}

func (r *conversationRepository) List(channelID int64, status *models.ConversationStatus, limit, offset int) ([]*models.Conversation, error) {
	var conversations []*models.Conversation

	query := r.db.Where("channel_id = ?", channelID)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("updated_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&conversations).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list conversations: %w", err)
	}

	return conversations, nil
}

func (r *conversationRepository) Update(id int64, req *models.UpdateConversationRequest) error {
	updates := make(map[string]interface{})

	if req.AssignedToExternalID != nil {
		updates["assigned_to_external_id"] = *req.AssignedToExternalID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		if *req.Status == models.ConversationStatusResolved || *req.Status == models.ConversationStatusClosed {
			updates["resolved_at"] = gorm.Expr("CURRENT_TIMESTAMP")
		}
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.Metadata != nil {
		updates["metadata"] = *req.Metadata
	}

	result := r.db.Model(&models.Conversation{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update conversation: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("conversation not found")
	}

	return nil
}

func (r *conversationRepository) UpdateLastMessage(id int64) error {
	return r.db.Model(&models.Conversation{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_message_at":  gorm.Expr("CURRENT_TIMESTAMP"),
			"first_message_at": gorm.Expr("COALESCE(first_message_at, CURRENT_TIMESTAMP)"),
		}).Error
}
