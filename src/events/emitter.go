package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Event types that will be emitted to NestJS
const (
	EventNewChatRequest      = "chat.request.new"
	EventChatAccepted        = "chat.request.accepted"
	EventChatRejected        = "chat.request.rejected"
	EventNewMessage          = "chat.message.new"
	EventMessageDelivered    = "chat.message.delivered"
	EventMessageRead         = "chat.message.read"
	EventConversationCreated = "chat.conversation.created"
	EventConversationUpdated = "chat.conversation.updated"
	EventUserOnline          = "chat.user.online"
	EventUserOffline         = "chat.user.offline"
)

// Event represents a generic event structure
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Source    string                 `json:"source"`
	Payload   map[string]interface{} `json:"payload"`
	Metadata  map[string]string      `json:"metadata,omitempty"`
}

// Emitter handles event emission to Redis
type Emitter interface {
	Emit(eventType string, payload map[string]interface{}) error
	EmitWithMetadata(eventType string, payload map[string]interface{}, metadata map[string]string) error
	Close() error
}

type redisEmitter struct {
	client  *redis.Client
	enabled bool
	source  string
}

// NewRedisEmitter creates a new Redis event emitter
func NewRedisEmitter(redisURL string, enabled bool) (Emitter, error) {
	if !enabled {
		return &noopEmitter{}, nil
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &redisEmitter{
		client:  client,
		enabled: true,
		source:  "go-chat-service",
	}, nil
}

func (e *redisEmitter) Emit(eventType string, payload map[string]interface{}) error {
	return e.EmitWithMetadata(eventType, payload, nil)
}

func (e *redisEmitter) EmitWithMetadata(eventType string, payload map[string]interface{}, metadata map[string]string) error {
	if !e.enabled {
		return nil
	}

	event := Event{
		ID:        fmt.Sprintf("%s-%d", eventType, time.Now().UnixNano()),
		Type:      eventType,
		Timestamp: time.Now(),
		Source:    e.source,
		Payload:   payload,
		Metadata:  metadata,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Publish to event type channel
	if err := e.client.Publish(ctx, eventType, data).Err(); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

func (e *redisEmitter) Close() error {
	if e.client != nil {
		return e.client.Close()
	}
	return nil
}

// noopEmitter is used when Redis is disabled
type noopEmitter struct{}

func (e *noopEmitter) Emit(eventType string, payload map[string]interface{}) error {
	return nil
}

func (e *noopEmitter) EmitWithMetadata(eventType string, payload map[string]interface{}, metadata map[string]string) error {
	return nil
}

func (e *noopEmitter) Close() error {
	return nil
}
