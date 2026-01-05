package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
	"github/sarthak-pokharel/sqlite-d1-gochat/src/repositories"
)


type WebhookService interface {
	ProcessWebhook(channelID int64, eventType string, payload interface{}) error
}

type webhookService struct {
	eventRepo  repositories.WebhookEventRepository
	msgService MessageService
}

func NewWebhookService(
	eventRepo repositories.WebhookEventRepository,
	msgService MessageService,
) WebhookService {
	return &webhookService{
		eventRepo:  eventRepo,
		msgService: msgService,
	}
}

func (s *webhookService) ProcessWebhook(channelID int64, eventType string, payload interface{}) error {
	
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	webhookEvent := &models.WebhookEvent{
		ChannelID: channelID,
		EventType: eventType,
		Payload:   string(payloadJSON),
		CreatedAt: time.Now(),
	}

	savedEvent, err := s.eventRepo.Create(webhookEvent)
	if err != nil {
		return fmt.Errorf("failed to store webhook event: %w", err)
	}

	
	var processErr error
	switch eventType {
	case "message":
		processErr = s.processMessageEvent(channelID, payload)
	case "status_update":
		processErr = s.processStatusUpdate(channelID, payload)
	default:
		processErr = fmt.Errorf("unknown event type: %s", eventType)
	}

	
	if processErr != nil {
		_ = s.eventRepo.MarkFailed(savedEvent.ID, processErr.Error())
		return processErr
	}

	_ = s.eventRepo.MarkProcessed(savedEvent.ID)
	return nil
}

func (s *webhookService) processMessageEvent(channelID int64, payload interface{}) error {
	
	data, ok := payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload format")
	}

	
	platformMsgID, _ := data["message_id"].(string)
	platformUserID, _ := data["user_id"].(string)
	userDisplayName, _ := data["user_name"].(string)
	content, _ := data["content"].(string)
	msgTypeStr, _ := data["message_type"].(string)

	if platformUserID == "" || content == "" {
		return fmt.Errorf("missing required fields")
	}

	msgType := models.MessageType(msgTypeStr)
	if msgType == "" {
		msgType = models.MessageTypeText
	}

	
	_, err := s.msgService.ProcessIncomingMessage(&ProcessIncomingMessageRequest{
		ChannelID:         channelID,
		PlatformMessageID: platformMsgID,
		PlatformUserID:    platformUserID,
		UserDisplayName:   userDisplayName,
		Content:           content,
		MessageType:       msgType,
	})

	return err
}

func (s *webhookService) processStatusUpdate(channelID int64, payload interface{}) error {
	
	data, ok := payload.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payload format")
	}

	messageID, ok := data["message_id"].(float64)
	if !ok {
		return fmt.Errorf("missing message_id")
	}

	status, _ := data["status"].(string)

	switch status {
	case "delivered":
		return s.msgService.MarkDelivered(int64(messageID))
	case "read":
		return s.msgService.MarkRead(int64(messageID))
	}

	return nil
}
