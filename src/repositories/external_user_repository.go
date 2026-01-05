package repositories

import (
	"fmt"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"

	"gorm.io/gorm"
)

// ExternalUserRepository interface
type ExternalUserRepository interface {
	Create(req *models.CreateExternalUserRequest) (*models.ExternalUser, error)
	GetByID(id int64) (*models.ExternalUser, error)
	GetByPlatformUser(channelID int64, platformUserID string) (*models.ExternalUser, error)
	FindOrCreate(req *models.CreateExternalUserRequest) (*models.ExternalUser, error)
	Update(id int64, req *models.UpdateExternalUserRequest) error
	UpdateLastSeen(id int64) error
}

type externalUserRepository struct {
	db *gorm.DB
}

func NewExternalUserRepository(db *gorm.DB) ExternalUserRepository {
	return &externalUserRepository{db: db}
}

func (r *externalUserRepository) Create(req *models.CreateExternalUserRequest) (*models.ExternalUser, error) {
	user := &models.ExternalUser{
		ChannelID:        req.ChannelID,
		PlatformUserID:   req.PlatformUserID,
		PlatformUsername: req.PlatformUsername,
		DisplayName:      req.DisplayName,
		PhoneNumber:      req.PhoneNumber,
		Email:            req.Email,
		AvatarURL:        req.AvatarURL,
		Metadata:         req.Metadata,
		IsBlocked:        false,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create external user: %w", err)
	}

	return user, nil
}

func (r *externalUserRepository) GetByID(id int64) (*models.ExternalUser, error) {
	var user models.ExternalUser
	if err := r.db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *externalUserRepository) GetByPlatformUser(channelID int64, platformUserID string) (*models.ExternalUser, error) {
	var user models.ExternalUser
	err := r.db.Where("channel_id = ? AND platform_user_id = ?", channelID, platformUserID).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

func (r *externalUserRepository) FindOrCreate(req *models.CreateExternalUserRequest) (*models.ExternalUser, error) {
	// Try to find existing user
	user, err := r.GetByPlatformUser(req.ChannelID, req.PlatformUserID)
	if err == nil {
		// User exists, update last seen
		_ = r.UpdateLastSeen(user.ID)
		return user, nil
	}

	// User doesn't exist, create new
	return r.Create(req)
}

func (r *externalUserRepository) Update(id int64, req *models.UpdateExternalUserRequest) error {
	updates := make(map[string]interface{})

	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}
	if req.Metadata != nil {
		updates["metadata"] = *req.Metadata
	}
	if req.IsBlocked != nil {
		updates["is_blocked"] = *req.IsBlocked
	}

	if len(updates) == 0 {
		return nil // Nothing to update
	}

	result := r.db.Model(&models.ExternalUser{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *externalUserRepository) UpdateLastSeen(id int64) error {
	return r.db.Model(&models.ExternalUser{}).Where("id = ?", id).
		Update("last_seen_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}
