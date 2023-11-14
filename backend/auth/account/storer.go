package account

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Writer is the interface for changing `Account` in persistence
type Writer interface {
	Store(state *Account) error
	Update(state *Account) error
}

// Loader is the interface from loading `Account` from persistence
type Loader interface {
	ByUser(userID uuid.UUID) (Account, error)
	ByConfirmToken(token string) (Account, error)
	ByRecoveryToken(token string) (Account, error)
}

// Storer is the interface for `Account` persistence
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

// Store is a method of the `GORMStorer` struct. Takes a `Account` as parameter and persists it.
func (s *GORMStorer) Store(state *Account) error {
	result := s.db.Create(state)
	return result.Error
}

// Update is a method of the `GORMStorer` struct. Takes a `Account` as parameter and updates it.
func (s *GORMStorer) Update(state *Account) error {
	result := s.db.Updates(state)
	return result.Error
}

// ByUser is a method of the `GORMStorer` struct. Takes a userID as parameter to load a `Account` object from persistence.
func (s *GORMStorer) ByUser(userID uuid.UUID) (Account, error) {
	var state Account
	result := s.db.Where(&Account{UserID: userID}).First(&state)
	return state, result.Error
}

// ByConfirmToken is a method of the `GORMStorer` struct. Takes a confirmation token as parameter to load a `Account` object from persistence.
func (s *GORMStorer) ByConfirmToken(token string) (Account, error) {
	var state Account
	result := s.db.Where(&Account{ConfirmationToken: token}).First(&state)
	return state, result.Error
}

// ByRecoveryToken is a method of the `GORMStorer` struct. Takes a recovery token as parameter to load a `Account` object from persistence.
func (s *GORMStorer) ByRecoveryToken(token string) (Account, error) {
	var state Account
	result := s.db.Where(&Account{RecoveryToken: token}).First(&state)
	return state, result.Error
}
