package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Writer is the interface for changing `AuthenticationState` in persistence
type Writer interface {
	Store(state AuthenticationState) error
	Update(state AuthenticationState) error
}

// Loader is the interface from loading `AuthenticationState` from persistence
type Loader interface {
	ByUser(userID uuid.UUID) (AuthenticationState, error)
}

// Storer is the interface for `AuthenticationState` persistence
type Storer interface {
	Writer
	Loader
}

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

// Store is a method of the `GORMStorer` struct. Takes a `AuthenticationState` as parameter and persists it.
func (s *GORMStorer) Store(state AuthenticationState) error {
	result := s.db.Create(&state)
	return result.Error
}

// Update is a method of the `GORMStorer` struct. Takes a `AuthenticationState` as parameter and updates it.
func (s *GORMStorer) Update(state AuthenticationState) error {
	result := s.db.Updates(&state)
	return result.Error
}

// ByUser is a method of the `GORMStorer` struct. Takes a userID as parameter to load a `AuthenticationState` object from persistence.
func (s *GORMStorer) ByUser(userID uuid.UUID) (AuthenticationState, error) {
	var state AuthenticationState
	result := s.db.Where(&AuthenticationState{UserID: userID}).First(&state)
	return state, result.Error
}
