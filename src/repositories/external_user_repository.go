package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
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
	db *sql.DB
}

func NewExternalUserRepository(db *sql.DB) ExternalUserRepository {
	return &externalUserRepository{db: db}
}

func (r *externalUserRepository) Create(req *models.CreateExternalUserRequest) (*models.ExternalUser, error) {
	query := `
		INSERT INTO external_users (
			channel_id, platform_user_id, platform_username, display_name,
			phone_number, email, avatar_url, metadata, first_seen_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.Exec(query,
		req.ChannelID, req.PlatformUserID, req.PlatformUsername,
		req.DisplayName, req.PhoneNumber, req.Email,
		req.AvatarURL, req.Metadata, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create external user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	return r.GetByID(id)
}

func (r *externalUserRepository) GetByID(id int64) (*models.ExternalUser, error) {
	query := `
		SELECT id, channel_id, platform_user_id, platform_username, display_name,
		       phone_number, email, avatar_url, metadata, first_seen_at,
		       last_seen_at, is_blocked
		FROM external_users
		WHERE id = ?
	`
	user := &models.ExternalUser{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.ChannelID, &user.PlatformUserID,
		&user.PlatformUsername, &user.DisplayName, &user.PhoneNumber,
		&user.Email, &user.AvatarURL, &user.Metadata, &user.FirstSeenAt,
		&user.LastSeenAt, &user.IsBlocked,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (r *externalUserRepository) GetByPlatformUser(channelID int64, platformUserID string) (*models.ExternalUser, error) {
	query := `
		SELECT id, channel_id, platform_user_id, platform_username, display_name,
		       phone_number, email, avatar_url, metadata, first_seen_at,
		       last_seen_at, is_blocked
		FROM external_users
		WHERE channel_id = ? AND platform_user_id = ?
	`
	user := &models.ExternalUser{}
	err := r.db.QueryRow(query, channelID, platformUserID).Scan(
		&user.ID, &user.ChannelID, &user.PlatformUserID,
		&user.PlatformUsername, &user.DisplayName, &user.PhoneNumber,
		&user.Email, &user.AvatarURL, &user.Metadata, &user.FirstSeenAt,
		&user.LastSeenAt, &user.IsBlocked,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
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
	query := "UPDATE external_users SET "
	args := []interface{}{}
	updates := []string{}

	if req.DisplayName != nil {
		updates = append(updates, "display_name = ?")
		args = append(args, *req.DisplayName)
	}
	if req.PhoneNumber != nil {
		updates = append(updates, "phone_number = ?")
		args = append(args, *req.PhoneNumber)
	}
	if req.Email != nil {
		updates = append(updates, "email = ?")
		args = append(args, *req.Email)
	}
	if req.AvatarURL != nil {
		updates = append(updates, "avatar_url = ?")
		args = append(args, *req.AvatarURL)
	}
	if req.Metadata != nil {
		updates = append(updates, "metadata = ?")
		args = append(args, *req.Metadata)
	}
	if req.IsBlocked != nil {
		updates = append(updates, "is_blocked = ?")
		args = append(args, *req.IsBlocked)
	}

	if len(updates) == 0 {
		return nil // Nothing to update
	}

	query += updates[0]
	for i := 1; i < len(updates); i++ {
		query += ", " + updates[i]
	}
	query += " WHERE id = ?"
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *externalUserRepository) UpdateLastSeen(id int64) error {
	query := "UPDATE external_users SET last_seen_at = ? WHERE id = ?"
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}
