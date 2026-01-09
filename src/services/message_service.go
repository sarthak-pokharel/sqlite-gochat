package services

import (
	"fmt"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/events"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/repositories"
)

type MessageService interface {
	ProcessIncomingMessage(req *ProcessIncomingMessageRequest) (*models.Message, error)
	SendOutgoingMessage(req *SendOutgoingMessageRequest) (*models.Message, error)
	GetMessageHistory(conversationID int64, limit, offset int, before *int64) ([]*models.Message, error)
	MarkDelivered(messageID int64) error
	MarkRead(messageID int64) error
}

type ProcessIncomingMessageRequest struct {
	ChannelID         int64
	PlatformMessageID string
	PlatformUserID    string
	UserDisplayName   string
	UserPhone         *string
	UserEmail         *string
	Content           string
	MessageType       models.MessageType
	MediaURL          *string
	Metadata          *string
}

type SendOutgoingMessageRequest struct {
	ConversationID int64
	Content        string
	MessageType    models.MessageType
	MediaURL       *string
	SenderID       *int64
	Metadata       *string
}

type messageService struct {
	messageRepo      repositories.MessageRepository
	conversationRepo repositories.ConversationRepository
	externalUserRepo repositories.ExternalUserRepository
	emitter          events.Emitter
}

func NewMessageService(
	messageRepo repositories.MessageRepository,
	conversationRepo repositories.ConversationRepository,
	externalUserRepo repositories.ExternalUserRepository,
	emitter events.Emitter,
) MessageService {
	return &messageService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		externalUserRepo: externalUserRepo,
		emitter:          emitter,
	}
}

func (s *messageService) ProcessIncomingMessage(req *ProcessIncomingMessageRequest) (*models.Message, error) {

	user, err := s.externalUserRepo.FindOrCreate(&models.CreateExternalUserRequest{
		ChannelID:      req.ChannelID,
		PlatformUserID: req.PlatformUserID,
		DisplayName:    &req.UserDisplayName,
		PhoneNumber:    req.UserPhone,
		Email:          req.UserEmail,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find/create user: %w", err)
	}

	conversation, err := s.conversationRepo.GetOrCreateByUser(req.ChannelID, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create conversation: %w", err)
	}

	message := &models.Message{
		ConversationID:    conversation.ID,
		PlatformMessageID: &req.PlatformMessageID,
		SenderType:        models.SenderExternal,
		SenderID:          &user.ID,
		Content:           req.Content,
		MessageType:       req.MessageType,
		MediaURL:          req.MediaURL,
		Direction:         models.DirectionInbound,
		Status:            models.MessageStatusReceived,
		CreatedAt:         time.Now(),
		Metadata:          req.Metadata,
	}

	savedMessage, err := s.messageRepo.Create(message)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	if err := s.conversationRepo.UpdateLastMessage(conversation.ID); err != nil {

		fmt.Printf("Warning: failed to update conversation last message: %v\n", err)
	}

	if err := s.externalUserRepo.UpdateLastSeen(user.ID); err != nil {

		fmt.Printf("Warning: failed to update user last seen: %v\n", err)
	}

	go s.emitter.Emit(events.EventNewMessage, map[string]interface{}{
		"message_id":       savedMessage.ID,
		"conversation_id":  conversation.ID,
		"channel_id":       req.ChannelID,
		"external_user_id": user.ID,
		"content":          req.Content,
		"message_type":     req.MessageType,
		"direction":        models.DirectionInbound,
		"timestamp":        savedMessage.CreatedAt,
	})

	return savedMessage, nil
}

func (s *messageService) SendOutgoingMessage(req *SendOutgoingMessageRequest) (*models.Message, error) {

	conversation, err := s.conversationRepo.GetByID(req.ConversationID)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}
	if conversation == nil {
		return nil, fmt.Errorf("conversation not found")
	}

	message := &models.Message{
		ConversationID: req.ConversationID,
		SenderType:     models.SenderInternal,
		SenderID:       req.SenderID,
		Content:        req.Content,
		MessageType:    req.MessageType,
		MediaURL:       req.MediaURL,
		Direction:      models.DirectionOutbound,
		Status:         models.MessageStatusSent,
		CreatedAt:      time.Now(),
		Metadata:       req.Metadata,
	}

	savedMessage, err := s.messageRepo.Create(message)
	if err != nil {
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	if err := s.conversationRepo.UpdateLastMessage(conversation.ID); err != nil {
		fmt.Printf("Warning: failed to update conversation: %v\n", err)
	}

	go s.emitter.Emit(events.EventNewMessage, map[string]interface{}{
		"message_id":      savedMessage.ID,
		"conversation_id": conversation.ID,
		"channel_id":      conversation.ChannelID,
		"content":         req.Content,
		"message_type":    req.MessageType,
		"direction":       models.DirectionOutbound,
		"timestamp":       savedMessage.CreatedAt,
	})

	return savedMessage, nil
}

func (s *messageService) GetMessageHistory(conversationID int64, limit, offset int, before *int64) ([]*models.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	return s.messageRepo.ListByConversation(conversationID, limit, offset, before)
}

func (s *messageService) MarkDelivered(messageID int64) error {
	if err := s.messageRepo.UpdateStatus(messageID, models.MessageStatusDelivered); err != nil {
		return err
	}

	go s.emitter.Emit(events.EventMessageDelivered, map[string]interface{}{
		"message_id": messageID,
		"status":     models.MessageStatusDelivered,
	})

	return nil
}

func (s *messageService) MarkRead(messageID int64) error {
	if err := s.messageRepo.UpdateStatus(messageID, models.MessageStatusRead); err != nil {
		return err
	}

	go s.emitter.Emit(events.EventMessageRead, map[string]interface{}{
		"message_id": messageID,
		"status":     models.MessageStatusRead,
	})

	return nil
}
