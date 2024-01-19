package onetime

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Access is a type for a one time accessible link for an object with TTL
type Access struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	OriginalID uuid.UUID
	OneTime    bool
	TTL        time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

// CreateAccess is a JSON type for creating a new one time access
type CreateAccess struct {
	OriginalID string `json:"original_id"`
	OneTime    bool   `json:"one_time"`
}

// Resp is a JSON type representing a one time access
type Resp struct {
	ID string `json:"id"`
}
