package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Writer is the interface for changing `User` in persistence
type Writer interface {
	Store(user User) error
	Delete(email string) error
}

// Loader is the interface from loading `User` from persistence
type Loader interface {
	ByEmail(email string) (User, error)
	ByID(id uuid.UUID) (User, error)
}

// Storer is the interface for `User` persistence
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

// Store is is a method of the `GORMStorer` struct. Takes a `User` as parameter and persists it.
func (s *GORMStorer) Store(user User) error {
	result := s.db.Create(&user)
	return result.Error
}

// ByEmail is is a method of the `GORMStorer` struct. Takes an email as parameter to load a `User` object from persistence.
func (s *GORMStorer) ByEmail(email string) (User, error) {
	var user User
	result := s.db.Where(&User{Email: email}).First(&user)
	return user, result.Error
}

// ByID is is a method of the `GORMStorer` struct. Takes an UUID as parameter to load a `User` object from persistence.
func (s *GORMStorer) ByID(id uuid.UUID) (User, error) {
	var user User
	result := s.db.Where(&User{ID: id}).First(&user)
	return user, result.Error
}

// Delete is is a method of the `GORMStorer` struct. Takes an email as parameter and deletes the corresponding `User` from persistence.
func (s *GORMStorer) Delete(email string) error {
	var user User
	result := s.db.Where(&User{Email: email}).Delete(&user)
	return result.Error
}
