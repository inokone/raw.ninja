package role

import (
	"gorm.io/gorm"
)

// Writer is the interface for changing `Role` in persistence
type Writer interface {
	Update(role ProfileRole) error
}

// Loader is the interface from loading `Role` from persistence
type Loader interface {
	List() ([]Role, error)
}

// Storer is the interface for `Role` persistence
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

// Update is a method of the `GORMStorer` struct. Takes a `Role` and updates settings (quota and display name) for it.
func (s *GORMStorer) Update(role ProfileRole) error {
	var persisted Role
	result := s.db.Where(&Role{RoleType: role.RoleType}).First(&persisted)
	if result.Error != nil {
		return result.Error
	}

	persisted.Quota = role.Quota
	persisted.DisplayName = role.DisplayName

	result = s.db.Updates(persisted)

	return result.Error
}

// List is a method of the `GORMStorer` struct. Loads all `Role` objects from persistence.
func (s *GORMStorer) List() ([]Role, error) {
	var roles []Role
	result := s.db.Find(&roles)
	return roles, result.Error
}
