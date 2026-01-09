package testutils

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
