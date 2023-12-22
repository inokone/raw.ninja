package collection

import (
	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"gorm.io/gorm"
)

const (
	listQuery = `SELECT c.id as ID, c.name as Name, c.tags as Tags, c.created_at as Created, count(p.photo_id) as Photos, c.thumbnail_id as Thumbnail
	FROM collections c 
	LEFT JOIN collection_photos p ON c.id = p.collection_id 
	WHERE user_id = ? and type = ?
	GROUP BY c.id
	ORDER by c.created_at DESC`
)

// Writer is the interface for changing `Collection` in persistence
type Writer interface {
	Store(c *Collection) error
	Update(c *Collection) error
	Delete(id uuid.UUID) error
}

// Loader is the interface from loading `Collection` from persistence
type Loader interface {
	ByUserAndType(usr *user.User, ct Type) ([]ListItem, error)
	ByID(id uuid.UUID) (*Collection, error)
}

// Storer is the interface for `Collection` persistence
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

// Store is a method of the `GORMStorer` struct. Takes a `Collection` as parameter and persists it.
func (s *GORMStorer) Store(collection *Collection) error {
	result := s.db.Create(collection)
	return result.Error
}

// Update is a method of the `GORMStorer` struct. Takes a `Collection` and updates it.
func (s *GORMStorer) Update(collection *Collection) error {
	res := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(collection)
	return res.Error
}

// ByID is a method of the `GORMStorer` struct. Takes an UUID as parameter to load a `Collection` object from persistence.
func (s *GORMStorer) ByID(id uuid.UUID) (*Collection, error) {
	var collection Collection
	result := s.db.Preload("Photos.Desc.Metadata").First(&collection, "id = ?", id.String())
	return &collection, result.Error
}

// ByUserAndType is a method of the `GORMStorer` struct. Takes a user and a type as parameters to lists `Collection` objects as `ListResp` from persistence.
// The returned
func (s *GORMStorer) ByUserAndType(usr *user.User, ct Type) ([]ListItem, error) {
	var collection []ListItem
	res := s.db.Raw(listQuery, usr.ID, ct).Scan(&collection)
	return collection, res.Error
}

// Delete is a method of the `GORMStorer` struct. Takes an id as parameter and deletes the corresponding `Collection` from persistence.
func (s *GORMStorer) Delete(id uuid.UUID) error {
	var collection Collection
	result := s.db.Where(&Collection{ID: id}).Delete(&collection)
	return result.Error
}
