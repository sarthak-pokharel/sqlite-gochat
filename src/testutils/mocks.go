package testutils

import (
	"github/sarthak-pokharel/sqlite-d1-gochat/src/models"
)

// MockEmitter is a mock implementation of events.Emitter
type MockEmitter struct {
	EmittedEvents []EmittedEvent
	EmitError     error
}

type EmittedEvent struct {
	EventType string
	Payload   map[string]interface{}
	Metadata  map[string]string
}

func NewMockEmitter() *MockEmitter {
	return &MockEmitter{
		EmittedEvents: make([]EmittedEvent, 0),
	}
}

func (m *MockEmitter) Emit(eventType string, payload map[string]interface{}) error {
	if m.EmitError != nil {
		return m.EmitError
	}
	m.EmittedEvents = append(m.EmittedEvents, EmittedEvent{
		EventType: eventType,
		Payload:   payload,
	})
	return nil
}

func (m *MockEmitter) EmitWithMetadata(eventType string, payload map[string]interface{}, metadata map[string]string) error {
	if m.EmitError != nil {
		return m.EmitError
	}
	m.EmittedEvents = append(m.EmittedEvents, EmittedEvent{
		EventType: eventType,
		Payload:   payload,
		Metadata:  metadata,
	})
	return nil
}

func (m *MockEmitter) Close() error {
	return nil
}

// MockOrganizationRepository is a mock implementation of OrganizationRepository
type MockOrganizationRepository struct {
	Organizations       map[int64]*models.Organization
	OrganizationsBySlug map[string]*models.Organization
	NextID              int64
	CreateError         error
	GetError            error
	ListError           error
	UpdateError         error
	DeleteError         error
}

func NewMockOrganizationRepository() *MockOrganizationRepository {
	return &MockOrganizationRepository{
		Organizations:       make(map[int64]*models.Organization),
		OrganizationsBySlug: make(map[string]*models.Organization),
		NextID:              1,
	}
}

func (m *MockOrganizationRepository) Create(req *models.CreateOrganizationRequest) (*models.Organization, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	org := &models.Organization{
		ID:   m.NextID,
		Name: req.Name,
		Slug: req.Slug,
	}
	m.Organizations[org.ID] = org
	m.OrganizationsBySlug[org.Slug] = org
	m.NextID++
	return org, nil
}

func (m *MockOrganizationRepository) GetByID(id int64) (*models.Organization, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	org, ok := m.Organizations[id]
	if !ok {
		return nil, nil
	}
	return org, nil
}

func (m *MockOrganizationRepository) GetBySlug(slug string) (*models.Organization, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	org, ok := m.OrganizationsBySlug[slug]
	if !ok {
		return nil, nil
	}
	return org, nil
}

func (m *MockOrganizationRepository) List(limit, offset int) ([]*models.Organization, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	result := make([]*models.Organization, 0)
	for _, org := range m.Organizations {
		result = append(result, org)
	}
	return result, nil
}

func (m *MockOrganizationRepository) Update(id int64, req *models.UpdateOrganizationRequest) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	org, ok := m.Organizations[id]
	if !ok {
		return nil
	}
	if req.Name != nil {
		org.Name = *req.Name
	}
	return nil
}

func (m *MockOrganizationRepository) Delete(id int64) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}
	org, ok := m.Organizations[id]
	if ok {
		delete(m.OrganizationsBySlug, org.Slug)
		delete(m.Organizations, id)
	}
	return nil
}

// MockConversationRepository is a mock implementation of ConversationRepository
type MockConversationRepository struct {
	Conversations map[int64]*models.Conversation
	NextID        int64
	CreateError   error
	GetError      error
	ListError     error
	UpdateError   error
}

func NewMockConversationRepository() *MockConversationRepository {
	return &MockConversationRepository{
		Conversations: make(map[int64]*models.Conversation),
		NextID:        1,
	}
}

func (m *MockConversationRepository) Create(req *models.CreateConversationRequest) (*models.Conversation, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	conv := &models.Conversation{
		ID:             m.NextID,
		ChannelID:      req.ChannelID,
		ExternalUserID: req.ExternalUserID,
		Status:         models.ConversationStatusOpen,
		Priority:       req.Priority,
	}
	m.Conversations[conv.ID] = conv
	m.NextID++
	return conv, nil
}

func (m *MockConversationRepository) GetByID(id int64) (*models.Conversation, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	conv, ok := m.Conversations[id]
	if !ok {
		return nil, nil
	}
	return conv, nil
}

func (m *MockConversationRepository) GetOrCreateByUser(channelID, externalUserID int64) (*models.Conversation, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	for _, conv := range m.Conversations {
		if conv.ChannelID == channelID && conv.ExternalUserID == externalUserID {
			return conv, nil
		}
	}
	return m.Create(&models.CreateConversationRequest{
		ChannelID:      channelID,
		ExternalUserID: externalUserID,
		Priority:       models.PriorityNormal,
	})
}

func (m *MockConversationRepository) List(channelID int64, status *models.ConversationStatus, limit, offset int) ([]*models.Conversation, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	result := make([]*models.Conversation, 0)
	for _, conv := range m.Conversations {
		if conv.ChannelID == channelID {
			if status == nil || conv.Status == *status {
				result = append(result, conv)
			}
		}
	}
	return result, nil
}

func (m *MockConversationRepository) Update(id int64, req *models.UpdateConversationRequest) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	conv, ok := m.Conversations[id]
	if !ok {
		return nil
	}
	if req.Status != nil {
		conv.Status = *req.Status
	}
	if req.Priority != nil {
		conv.Priority = *req.Priority
	}
	if req.AssignedToExternalID != nil {
		conv.AssignedToExternalID = req.AssignedToExternalID
	}
	return nil
}

func (m *MockConversationRepository) UpdateLastMessage(id int64) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

// MockMessageRepository is a mock implementation of MessageRepository
type MockMessageRepository struct {
	Messages    map[int64]*models.Message
	NextID      int64
	CreateError error
	GetError    error
	ListError   error
	UpdateError error
}

func NewMockMessageRepository() *MockMessageRepository {
	return &MockMessageRepository{
		Messages: make(map[int64]*models.Message),
		NextID:   1,
	}
}

func (m *MockMessageRepository) Create(msg *models.Message) (*models.Message, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	msg.ID = m.NextID
	m.Messages[msg.ID] = msg
	m.NextID++
	return msg, nil
}

func (m *MockMessageRepository) GetByID(id int64) (*models.Message, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	msg, ok := m.Messages[id]
	if !ok {
		return nil, nil
	}
	return msg, nil
}

func (m *MockMessageRepository) ListByConversation(conversationID int64, limit, offset int, before *int64) ([]*models.Message, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	result := make([]*models.Message, 0)
	for _, msg := range m.Messages {
		if msg.ConversationID == conversationID {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (m *MockMessageRepository) UpdateStatus(id int64, status models.MessageStatus) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	msg, ok := m.Messages[id]
	if ok {
		msg.Status = status
	}
	return nil
}

// MockExternalUserRepository is a mock implementation of ExternalUserRepository
type MockExternalUserRepository struct {
	Users       map[int64]*models.ExternalUser
	NextID      int64
	CreateError error
	GetError    error
	UpdateError error
}

func NewMockExternalUserRepository() *MockExternalUserRepository {
	return &MockExternalUserRepository{
		Users:  make(map[int64]*models.ExternalUser),
		NextID: 1,
	}
}

func (m *MockExternalUserRepository) Create(req *models.CreateExternalUserRequest) (*models.ExternalUser, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	user := &models.ExternalUser{
		ID:             m.NextID,
		ChannelID:      req.ChannelID,
		PlatformUserID: req.PlatformUserID,
	}
	if req.DisplayName != nil {
		user.DisplayName = req.DisplayName
	}
	m.Users[user.ID] = user
	m.NextID++
	return user, nil
}

func (m *MockExternalUserRepository) GetByID(id int64) (*models.ExternalUser, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	user, ok := m.Users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (m *MockExternalUserRepository) GetByPlatformUser(channelID int64, platformUserID string) (*models.ExternalUser, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	for _, user := range m.Users {
		if user.ChannelID == channelID && user.PlatformUserID == platformUserID {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockExternalUserRepository) FindOrCreate(req *models.CreateExternalUserRequest) (*models.ExternalUser, error) {
	user, err := m.GetByPlatformUser(req.ChannelID, req.PlatformUserID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	return m.Create(req)
}

func (m *MockExternalUserRepository) Update(id int64, req *models.UpdateExternalUserRequest) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

func (m *MockExternalUserRepository) UpdateLastSeen(id int64) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

// MockWebhookEventRepository is a mock implementation of WebhookEventRepository
type MockWebhookEventRepository struct {
	Events      map[int64]*models.WebhookEvent
	NextID      int64
	CreateError error
	GetError    error
	ListError   error
	UpdateError error
}

func NewMockWebhookEventRepository() *MockWebhookEventRepository {
	return &MockWebhookEventRepository{
		Events: make(map[int64]*models.WebhookEvent),
		NextID: 1,
	}
}

func (m *MockWebhookEventRepository) Create(event *models.WebhookEvent) (*models.WebhookEvent, error) {
	if m.CreateError != nil {
		return nil, m.CreateError
	}
	event.ID = m.NextID
	m.Events[event.ID] = event
	m.NextID++
	return event, nil
}

func (m *MockWebhookEventRepository) GetByID(id int64) (*models.WebhookEvent, error) {
	if m.GetError != nil {
		return nil, m.GetError
	}
	event, ok := m.Events[id]
	if !ok {
		return nil, nil
	}
	return event, nil
}

func (m *MockWebhookEventRepository) ListUnprocessed(channelID int64, limit int) ([]*models.WebhookEvent, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}
	result := make([]*models.WebhookEvent, 0)
	for _, event := range m.Events {
		if event.ChannelID == channelID && event.ProcessedAt == nil {
			result = append(result, event)
		}
	}
	return result, nil
}

func (m *MockWebhookEventRepository) MarkProcessed(id int64) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	return nil
}

func (m *MockWebhookEventRepository) MarkFailed(id int64, errorMsg string) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}
	event, ok := m.Events[id]
	if ok {
		event.Error = &errorMsg
	}
	return nil
}
