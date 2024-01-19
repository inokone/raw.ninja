package onetime

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Storer is the interface for `OneTimeAccess` persistence
type Storer interface {
	Loader
	Writer
}

// Loader is an interface for loading `OneTimeAccess` entities from persistence.
type Loader interface {
	ByID(id uuid.UUID) (*Access, error)
}

// Writer is an interface for persistence of `OneTimeAccess` entities.
type Writer interface {
	Store(ota *Access) error
}

// ExpiredAccess is an error for requesting an access with expired TTL
type ExpiredAccess struct {
	ID uuid.UUID
}

// Error is the string representation of an `InvalidPhotoID` error
func (e ExpiredAccess) Error() string { return fmt.Sprintf("Expired acces [%v]", e.ID) }

// GORMStorer is the `Storer` implementation based on GORM library.
type GORMStorer struct {
	db *gorm.DB
}

// NewGORMStorer creates a new `GORMStorer` instance based on the GORM library.
func NewGORMStorer(db *gorm.DB) *GORMStorer {
	return &GORMStorer{
		db: db,
	}
}

// ByID is a method of `GORMStorer` for loading a single `OneTimeAccess` entity by ID provided as parameter.
func (s *GORMStorer) ByID(id uuid.UUID) (*Access, error) {
	var o Access
	result := s.db.First(&o, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	if o.TTL.Before(time.Now()) {
		return nil, ExpiredAccess{id} // expired TTL, return error
	}
	if o.OneTime {
		result = s.db.Delete(&o, "id = ?", id) // one time usage, delete after retrieving
	}
	return &o, result.Error
}

// Store is a method of the `GORMStorer` struct. Takes a `OneTimeAccess` as parameter and persists it.
func (s *GORMStorer) Store(o *Access) error {
	result := s.db.Create(o)
	return result.Error
}
