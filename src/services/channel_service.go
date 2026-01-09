package services

import (
	"github/sarthak-pokharel/sqlite-d1-gochat/src/events"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/repositories"
)

type ChannelService interface {
	Create(req *models.CreateChannelRequest) (*models.ChatChannel, error)
	GetByID(id int64) (*models.ChatChannel, error)
	ListByOrganization(orgID int64, limit, offset int) ([]*models.ChatChannel, error)
	Update(id int64, req *models.UpdateChannelRequest) error
	UpdateStatus(id int64, status models.ChannelStatus) error
	Delete(id int64) error
}

type channelService struct {
	repo    repositories.ChannelRepository
	emitter events.Emitter
}

func NewChannelService(repo repositories.ChannelRepository, emitter events.Emitter) ChannelService {
	return &channelService{
		repo:    repo,
		emitter: emitter,
	}
}

func (s *channelService) Create(req *models.CreateChannelRequest) (*models.ChatChannel, error) {
	channel, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	go s.emitter.Emit(events.EventChannelCreated, map[string]interface{}{
		"channel_id":      channel.ID,
		"organization_id": channel.OrganizationID,
		"platform":        channel.Platform,
		"name":            channel.Name,
	})

	return channel, nil
}

func (s *channelService) GetByID(id int64) (*models.ChatChannel, error) {
	return s.repo.GetByID(id)
}

func (s *channelService) ListByOrganization(orgID int64, limit, offset int) ([]*models.ChatChannel, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.ListByOrganization(orgID, limit, offset)
}

func (s *channelService) Update(id int64, req *models.UpdateChannelRequest) error {
	if err := s.repo.Update(id, req); err != nil {
		return err
	}

	go s.emitter.Emit(events.EventChannelUpdated, map[string]interface{}{
		"channel_id": id,
	})

	return nil
}

func (s *channelService) UpdateStatus(id int64, status models.ChannelStatus) error {
	if err := s.repo.UpdateStatus(id, status); err != nil {
		return err
	}

	go s.emitter.Emit(events.EventChannelUpdated, map[string]interface{}{
		"channel_id": id,
		"status":     status,
	})

	return nil
}

func (s *channelService) Delete(id int64) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	go s.emitter.Emit(events.EventChannelDeleted, map[string]interface{}{
		"channel_id": id,
	})

	return nil
}
