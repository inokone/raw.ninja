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
	WHERE user_id = ? and type = ? and c.deleted_at IS NULL
	GROUP BY c.id
	ORDER by c.created_at DESC`

	storerStatsQuery = `SELECT count(*) as collection_count FROM collections WHERE type = ? and deleted_at IS NULL`

	searchQuery = `SELECT c.id as ID, c.name as Name, c.tags as Tags, c.created_at as Created, count(p.photo_id) as Photos, c.thumbnail_id as Thumbnail
	FROM collections c 
	LEFT JOIN collection_photos p ON c.id = p.collection_id 
	WHERE user_id = ? and type = ? and c.deleted_at IS NULL and c.name LIKE ? or ? LIKE ANY(c.tags)
	GROUP BY c.id
	ORDER by c.created_at DESC`

	statQuery = `SELECT c.created_at as Created, c.type as Type, count(p.photo_id) as Photos
	FROM collections c 
	LEFT JOIN collection_photos p ON c.id = p.collection_id 
	WHERE user_id = ?
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
	Stats(usr *user.User) ([]Stat, error)
	StorerStats() (int, int, error)
	Search(usrID uuid.UUID, ct Type, query string) ([]ListItem, error)
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
	var resID string
	res := s.db.Raw("DELETE from collection_photos WHERE collection_id = ? RETURNING collection_id", collection.ID).Scan(&resID)
	if res.Error != nil {
		return res.Error
	}
	res = s.db.Updates(collection)
	return res.Error
}

// ByID is a method of the `GORMStorer` struct. Takes an UUID as parameter to load a `Collection` object from persistence.
func (s *GORMStorer) ByID(id uuid.UUID) (*Collection, error) {
	var collection Collection
	result := s.db.Preload("Photos.Desc.Metadata").Preload("RuleSet.Rules").First(&collection, "id = ?", id.String())
	return &collection, result.Error
}

// ByUserAndType is a method of the `GORMStorer` struct. Takes a user and a type as parameters to lists `Collection` objects as `ListResp` from persistence.
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

// Stats is a method of the `GORMStorer` struct. Takes a user to list `Collection` objects with types and photo counts and creation date.
func (s *GORMStorer) Stats(usr *user.User) ([]Stat, error) {
	var collection []Stat
	res := s.db.Raw(statQuery, usr.ID).Scan(&collection)
	return collection, res.Error
}

// StorerStats is a method of the `GORMStorer` struct. List number of albums and uploads.
func (s *GORMStorer) StorerStats() (int, int, error) {
	var (
		albums  int
		uploads int
	)
	if res := s.db.Raw(storerStatsQuery, Album).Scan(&albums); res.Error != nil {
		return 0, 0, res.Error
	}
	if res := s.db.Raw(storerStatsQuery, Upload).Scan(&uploads); res.Error != nil {
		return 0, 0, res.Error
	}
	return albums, uploads, nil
}

// Search is a method of the `GORMStorer` struct. Takes a user and query string to search `Collection` objects matching these criteria.
func (s *GORMStorer) Search(usrID uuid.UUID, ct Type, query string) ([]ListItem, error) {
	var (
		q          = "%" + query + "%"
		collection []ListItem
	)
	res := s.db.Raw(searchQuery, usrID, ct, q, q).Scan(&collection)
	return collection, res.Error
}
