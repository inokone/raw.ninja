package ruleset

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Writer is an interface for persistence of `RuleSet` entities.
type Writer interface {
	Store(set *RuleSet) error

	Update(set *RuleSet) error

	Delete(id uuid.UUID) error
}

// Loader is an interface for loading `RuleSet` entities from persistence.
type Loader interface {
	ByID(id uuid.UUID) (*RuleSet, error)

	ByUser(userID uuid.UUID) ([]RuleSet, error)
}

// Storer is the interface for `RuleSet` persistence
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

// ByID is a method of `GORMStorer` for loading a single `RuleSet` entity by ID provided as parameter.
func (s *GORMStorer) ByID(id uuid.UUID) (*RuleSet, error) {
	var r RuleSet
	result := s.db.Preload("Rules", func(db *gorm.DB) *gorm.DB {
		return db.Order("rules.timing ASC")
	}).First(&r, "id = ?", id)
	return &r, result.Error
}

// ByUser is a method of `GORMStorer` for loading all `RuleSet`s of a user specified by the ID as a parameter.
func (s *GORMStorer) ByUser(userID uuid.UUID) ([]RuleSet, error) {
	var r []RuleSet
	result := s.db.Preload("Rules", func(db *gorm.DB) *gorm.DB {
		return db.Order("rules.timing ASC")
	}).Where("user_id = ?", userID).Order("created_at ASC").Find(&r)
	return r, result.Error
}

// Store is a method of the `GORMStorer` struct. Takes a `RuleSet` as parameter and persists it.
func (s *GORMStorer) Store(rs *RuleSet) error {
	result := s.db.Create(rs)
	return result.Error
}

// Update is a method of the `GORMStorer` struct. Takes a `RuleSet` and updates it.
func (s *GORMStorer) Update(rs *RuleSet) error {
	res := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&rs).Where("ID = ?", rs.ID).Updates(rs)
	return res.Error
}

// Delete is a method of the `GORMStorer` struct. Takes an id as parameter and deletes the corresponding `RuleSet` from persistence.
func (s *GORMStorer) Delete(id uuid.UUID) error {
	var collection RuleSet
	result := s.db.Where(&RuleSet{ID: id}).Delete(&collection)
	return result.Error
}
