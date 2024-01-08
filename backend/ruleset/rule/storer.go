package rule

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Writer is an interface for persistence of `Rule` entities.
type Writer interface {
	Store(rule *Rule) error

	Update(rule *Rule) error

	Delete(id uuid.UUID) error
}

// Loader is an interface for loading `Rule` entities from persistence.
type Loader interface {
	ByID(id uuid.UUID) (*Rule, error)

	ByUser(userID uuid.UUID) ([]Rule, error)
}

// Storer is the interface for `Rule` persistence
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

// ByID is a method of `GORMStorer` for loading a single `Rule` entity by ID provided as parameter.
func (s *GORMStorer) ByID(id uuid.UUID) (*Rule, error) {
	var r Rule
	result := s.db.Preload("Rules").First(&r, "id = ?", id)
	return &r, result.Error
}

// ByUser is a method of `GORMStorer` for loading all `Rule`s of a user specified by the ID as a parameter.
func (s *GORMStorer) ByUser(userID uuid.UUID) ([]Rule, error) {
	var r []Rule
	result := s.db.Where("user_id = ?", userID).Order("created_at ASC").Find(&r)
	return r, result.Error
}

// Store is a method of the `GORMStorer` struct. Takes a `Rule` as parameter and persists it.
func (s *GORMStorer) Store(r *Rule) error {
	result := s.db.Create(r)
	return result.Error
}

// Update is a method of the `GORMStorer` struct. Takes a `Rule` and updates it.
func (s *GORMStorer) Update(r *Rule) error {
	var (
		persisted *Rule
		err       error
	)

	persisted, err = s.ByID(r.ID)
	if err != nil {
		return err
	}

	res := s.db.Model(&persisted).Updates(r)
	return res.Error
}

// Delete is a method of `GORMStorer` for deleting a single `Rule` entity by ID provided as parameter.
func (s *GORMStorer) Delete(id uuid.UUID) error {
	var r Rule
	result := s.db.Delete(&r, "id = ?", id)
	return result.Error
}
