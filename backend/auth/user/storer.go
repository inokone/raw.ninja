package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Writer is the interface for changing `User` in persistence
type Writer interface {
	Store(user *User) error
	Update(user *User) error
	Patch(usr Patch) error
	Delete(email string) error
	SetEnabled(id uuid.UUID, enabled bool) error
}

// Loader is the interface from loading `User` from persistence
type Loader interface {
	ByEmail(email string) (*User, error)
	ByID(id uuid.UUID) (*User, error)
	List() ([]User, error)
}

// Storer is the interface for `User` persistence
type Storer interface {
	Writer
	Loader

	Stats() (Stats, error)
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

// Store is a method of the `GORMStorer` struct. Takes a `User` as parameter and persists it.
func (s *GORMStorer) Store(user *User) error {
	result := s.db.Create(user)
	return result.Error
}

// ByEmail is a method of the `GORMStorer` struct. Takes an email as parameter to load a `User` object from persistence.
func (s *GORMStorer) ByEmail(email string) (*User, error) {
	var user User
	result := s.db.Preload("Role").Where(&User{Email: email}).First(&user)
	return &user, result.Error
}

// ByID is a method of the `GORMStorer` struct. Takes an UUID as parameter to load a `User` object from persistence.
func (s *GORMStorer) ByID(id uuid.UUID) (*User, error) {
	var user User
	result := s.db.Preload("Role").Where(&User{ID: id}).First(&user)
	return &user, result.Error
}

// List is a method of the `GORMStorer` struct. Loads all `User` objects from persistence.
func (s *GORMStorer) List() ([]User, error) {
	var users []User
	result := s.db.Preload("Role").Find(&users)
	return users, result.Error
}

// Delete is a method of the `GORMStorer` struct. Takes an email as parameter and deletes the corresponding `User` from persistence.
func (s *GORMStorer) Delete(email string) error {
	var user User
	result := s.db.Where(&User{Email: email}).Delete(&user)
	return result.Error
}

// Patch is a method of the `GORMStorer` struct. Takes a `Patch` and updates settings for it.
func (s *GORMStorer) Patch(usr Patch) error {
	var (
		persisted *User
		id        uuid.UUID
		err       error
	)

	id, err = uuid.Parse(usr.ID)
	if err != nil {
		return err
	}

	persisted, err = s.ByID(id)
	if err != nil {
		return err
	}

	res := s.db.Model(&persisted).Updates(usr)
	return res.Error
}

// SetEnabled is a method of the `GORMStorer` struct. Takes a user and updates if it is enabled.
func (s *GORMStorer) SetEnabled(id uuid.UUID, enabled bool) error {
	var (
		persisted *User
		err       error
	)

	persisted, err = s.ByID(id)
	if err != nil {
		return err
	}

	res := s.db.Model(&persisted).Select("enabled").Updates(map[string]interface{}{"enabled": enabled})
	return res.Error
}

// Update is a method of the `GORMStorer` struct. Takes a `User` and updates it.
func (s *GORMStorer) Update(usr *User) error {
	res := s.db.Updates(usr)
	return res.Error
}

// Stats is a method of `GORMStorer` for collecting aggregated data on the storer.
func (s *GORMStorer) Stats() (Stats, error) {
	var (
		distribution []RoleUser
		count        int
	)

	res := s.db.Raw("SELECT r.display_name as role, count(u.id) as users FROM roles r JOIN users u on u.role_id = r.role_type GROUP BY r.display_name").Scan(&distribution)
	for _, value := range distribution {
		count += value.Users
	}

	return Stats{
		TotalUsers:   count,
		Distribution: distribution,
	}, res.Error
}
