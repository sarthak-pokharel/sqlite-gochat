package repositories

import (
	"fmt"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"

	"gorm.io/gorm"
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
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) ChannelRepository {
	return &channelRepository{db: db}
}

func (r *channelRepository) Create(req *models.CreateChannelRequest) (*models.ChatChannel, error) {
	channel := &models.ChatChannel{
		OrganizationID:    req.OrganizationID,
		Platform:          req.Platform,
		Name:              req.Name,
		AccountIdentifier: req.AccountIdentifier,
		Status:            models.ChannelStatusPending,
		WebhookSecret:     req.WebhookSecret,
		AccessToken:       req.AccessToken,
		Config:            req.Config,
		IsActive:          true,
	}

	if err := r.db.Create(channel).Error; err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	return channel, nil
}

func (r *channelRepository) GetByID(id int64) (*models.ChatChannel, error) {
	var channel models.ChatChannel
	if err := r.db.First(&channel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("channel not found")
		}
		return nil, fmt.Errorf("failed to get channel: %w", err)
	}
	return &channel, nil
}

func (r *channelRepository) ListByOrganization(orgID int64, limit, offset int) ([]*models.ChatChannel, error) {
	var channels []*models.ChatChannel
	err := r.db.Where("organization_id = ? AND is_active = ?", orgID, 1).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&channels).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list channels: %w", err)
	}

	// Remove access tokens from response
	for _, ch := range channels {
		ch.AccessToken = nil
	}

	return channels, nil
}

func (r *channelRepository) Update(id int64, req *models.UpdateChannelRequest) error {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.WebhookSecret != nil {
		updates["webhook_secret"] = *req.WebhookSecret
	}
	if req.AccessToken != nil {
		updates["access_token"] = *req.AccessToken
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	result := r.db.Model(&models.ChatChannel{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update channel: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("channel not found")
	}

	return nil
}

func (r *channelRepository) UpdateStatus(id int64, status models.ChannelStatus) error {
	result := r.db.Model(&models.ChatChannel{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update channel status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("channel not found")
	}
	return nil
}

func (r *channelRepository) Delete(id int64) error {
	result := r.db.Model(&models.ChatChannel{}).Where("id = ?", id).Update("is_active", 0)
	if result.Error != nil {
		return fmt.Errorf("failed to delete channel: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("channel not found")
	}
	return nil
}
