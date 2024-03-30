package common

import "github.com/google/uuid"

// EventType is a type for events
type EventType string

const (
	// snsMessagingType represents AWS SNS messaging type in configuration
	snsMessagingType = "sns"
	// audit is a type for events that are used for persisting audit trail on the app
	audit EventType = "raw-ninja.audit"
	// lifecycle is a type for events that are used for enforcing lifecycle rules
	lifecycle EventType = "raw-ninja.lifecycle"
)

// Event is the interface for message bus events.
type Event interface{}

// auditEvent is an event that is used for persisting audit trail on the app
type auditEvent struct {
	CorrID     uuid.UUID         `json:"correlation_id"`
	EType      EventType         `json:"type"`
	UserID     string            `json:"user_id"`
	Action     string            `json:"action"`
	TargetIDs  []string          `json:"target_ids"`
	TargetType string            `json:"target_type"`
	Metdata    map[string]string `json:"metadata"`
	Outcome    string            `json:"outcome"`
}

// NewAuditEvent creates an event for audit purposes
func NewAuditEvent(userID string, action string, targetIDs []string, targetType string, metadata map[string]string, outcome string) Event {
	return &auditEvent{
		CorrID:     uuid.New(),
		EType:      audit,
		UserID:     userID,
		Action:     action,
		TargetIDs:  targetIDs,
		TargetType: targetType,
		Metdata:    metadata,
		Outcome:    outcome,
	}
}

// lifecycleEvent is an event that is used for enforcing lifecycle rules
type lifecycleEvent struct {
	CorrID  uuid.UUID `json:"correlatrion_id"`
	EType   EventType `json:"type"`
	UserID  string    `json:"user_id"`
	Action  string    `json:"action"`
	AlbumID string    `json:"album_id"`
	Photos  []string  `json:"photos"`
}

// NewLifecycleEvent creates an event for lyfecycle enforcing purposes
func NewLifecycleEvent(userID string, action string, albumID string, photos []string) Event {
	return &lifecycleEvent{
		CorrID:  uuid.New(),
		EType:   lifecycle,
		UserID:  userID,
		Action:  action,
		AlbumID: albumID,
		Photos:  photos,
	}
}
