package services

import (
	"github/sarthak-pokharel/sqlite-d1-gochat/src/events"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/repositories"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/utils"
)

type ConversationService interface {
	GetByID(id int64) (*models.Conversation, error)
	ListByChannel(channelID int64, status *models.ConversationStatus, limit, offset int) ([]*models.Conversation, error)
	Assign(conversationID int64, assigneeID string) error
	UpdateStatus(conversationID int64, status models.ConversationStatus) error
	UpdatePriority(conversationID int64, priority models.ConversationPriority) error
}

type conversationService struct {
	repo    repositories.ConversationRepository
	emitter events.Emitter
}

func NewConversationService(repo repositories.ConversationRepository, emitter events.Emitter) ConversationService {
	return &conversationService{
		repo:    repo,
		emitter: emitter,
	}
}

func (s *conversationService) GetByID(id int64) (*models.Conversation, error) {
	return s.repo.GetByID(id)
}

func (s *conversationService) ListByChannel(channelID int64, status *models.ConversationStatus, limit, offset int) ([]*models.Conversation, error) {
	return s.repo.List(channelID, status, utils.NormalizeLimit(limit), utils.NormalizeOffset(offset))
}

func (s *conversationService) Assign(conversationID int64, assigneeID string) error {
	if err := s.repo.Update(conversationID, &models.UpdateConversationRequest{
		AssignedToExternalID: &assigneeID,
	}); err != nil {
		return err
	}

	go s.emitter.Emit("conversation.assigned", map[string]interface{}{
		"conversation_id": conversationID,
		"assignee_id":     assigneeID,
	})

	return nil
}

func (s *conversationService) UpdateStatus(conversationID int64, status models.ConversationStatus) error {
	if err := s.repo.Update(conversationID, &models.UpdateConversationRequest{
		Status: &status,
	}); err != nil {
		return err
	}

	go s.emitter.Emit(events.EventConversationUpdated, map[string]interface{}{
		"conversation_id": conversationID,
		"status":          status,
	})

	return nil
}

func (s *conversationService) UpdatePriority(conversationID int64, priority models.ConversationPriority) error {
	if err := s.repo.Update(conversationID, &models.UpdateConversationRequest{
		Priority: &priority,
	}); err != nil {
		return err
	}

	go s.emitter.Emit(events.EventConversationUpdated, map[string]interface{}{
		"conversation_id": conversationID,
		"priority":        priority,
	})

	return nil
}
