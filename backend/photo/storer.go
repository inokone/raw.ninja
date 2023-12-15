package photo

import (
	"github.com/google/uuid"
	_ "github.com/lib/pq" // Postgres driver package for GORM, no need to have a name
	"gorm.io/gorm"
)

// Writer is an interface for persistence of `Photo` entities.
type Writer interface {
	Store(photo *Photo) (uuid.UUID, error)
	Update(photo *Photo) error
	Delete(id string) error
}

// Loader is an interface for loading `Photo` entities from persistence.
type Loader interface {
	Load(id string) (*Photo, error)
	All(userID string) ([]Photo, error)
}

// Searcher is an interface for searching `Photo` entities by various filters in persistence.
type Searcher interface {
	Search(userID string, searchText string) ([]Photo, error)
	Favorites(userID string) ([]Photo, error)
}

// Storer is an interface for types that can store `Photo`s.
type Storer interface {
	Writer
	Loader
	Searcher

	UserStats(userID string) (UserStats, error)
	Stats() (Stats, error)
}

// GORMStorer is an implementation of `Storer` interface based on GORM library.
type GORMStorer struct {
	db *gorm.DB
}

// NewGORMStorer creates a new `GORMStorer` instance based on the GORM library and image persistence provided in parameters.
func NewGORMStorer(db *gorm.DB) *GORMStorer {
	return &GORMStorer{db: db}
}

// Store is a method of `GORMStorer` for persisting a `Photo` entity.
func (s *GORMStorer) Store(photo *Photo) (uuid.UUID, error) {
	result := s.db.Save(&photo)
	return photo.ID, result.Error
}

// Delete is a method of `GORMStorer` for deleting a `Photo` entity from persistence.
func (s *GORMStorer) Delete(id string) error {
	var (
		photo  Photo
		result *gorm.DB
	)
	result = s.db.Delete(&photo, "id = ?", id)
	return result.Error
}

// Update is a method of `GORMStorer` for updating metadata of a `Photo` entity in persistence.
func (s *GORMStorer) Update(photo *Photo) error {
	result := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&photo)
	return result.Error
}

// Load is a method of `GORMStorer` for loading a single `Photo` entity by ID provided as parameter.
func (s *GORMStorer) Load(id string) (*Photo, error) {
	var photo Photo
	result := s.db.Preload("Desc.Metadata").First(&photo, "id = ?", id)
	return &photo, result.Error
}

// All is a method of `GORMStorer` for loading a all `Photo`s of a user specified by the ID as a parameter.
func (s *GORMStorer) All(userID string) ([]Photo, error) {
	var photos []Photo
	result := s.db.Preload("Desc.Metadata").Where("user_id = ?", userID).Order("created_at ASC").Find(&photos)
	return photos, result.Error
}

// Favorites is a method of `GORMStorer` for loading a favorite `Photo`s of a user specified by the ID as a parameter.
func (s *GORMStorer) Favorites(userID string) ([]Photo, error) {
	var photos []Photo
	result := s.db.Preload(
		"Desc.Metadata").Joins(
		"JOIN descriptors ON descriptors.id = photos.desc_id").Where(
		"photos.user_id = ?", userID).Where(
		"descriptors.favorite = true").Order(
		"photos.created_at DESC").Find(&photos)
	return photos, result.Error
}

// Search is a method of `GORMStorer` for loading `Photo`s of a user specified by the ID as a parameter,
// that has filename matching the search text parameter.
func (s *GORMStorer) Search(userID string, searchText string) ([]Photo, error) {
	var photos []Photo
	result := s.db.Preload(
		"Desc.Metadata").Joins(
		"JOIN descriptors ON descriptors.id = photos.desc_id").Where(
		"photos.user_id = ?", userID).Where(
		"descriptors.file_name LIKE ?", "%"+searchText+"%").Order(
		"photos.created_at ASC").Find(&photos)
	return photos, result.Error
}

// UserStats is a method of `GORMStorer` for collecting aggregated data on the photos of the user specified by the ID in the parameter.
func (s *GORMStorer) UserStats(userID string) (UserStats, error) {
	var (
		photos, favorites int
		usedSpace         int64
	)
	res := s.db.Raw("SELECT count(id) FROM photos WHERE user_id = ?", userID).Scan(&photos)
	if res.Error != nil {
		return UserStats{}, res.Error
	}

	res = s.db.Raw("SELECT coalesce(sum(coalesce(used_space, 0)),0) FROM photos WHERE user_id = ?", userID).Scan(&usedSpace)
	if res.Error != nil {
		return UserStats{}, res.Error
	}

	res = s.db.Raw("SELECT count(p.id) FROM photos p JOIN descriptors d ON d.id = p.desc_id WHERE p.user_id = ? and d.favorite = true", userID).Scan(&favorites)
	if res.Error != nil {
		return UserStats{}, res.Error
	}

	photoList, err := s.All(userID)
	if err != nil {
		return UserStats{}, err
	}

	photoIDs := make([]string, len(photoList))
	for idx, photo := range photoList {
		photoIDs[idx] = photo.ID.String()
	}

	return UserStats{
		Photos:    photos,
		Favorites: favorites,
		UsedSpace: usedSpace,
	}, nil
}

// Stats is a method of `GORMStorer` for collecting aggregated data on the storer.
func (s *GORMStorer) Stats() (Stats, error) {
	var (
		photos, favorites int
		usedSpace         int64
	)

	res := s.db.Raw("SELECT count(id) FROM photos").Scan(&photos)
	if res.Error != nil {
		return Stats{}, res.Error
	}

	res = s.db.Raw("SELECT coalesce(sum(coalesce(used_space, 0)), 0) FROM photos").Scan(&usedSpace)
	if res.Error != nil {
		return Stats{}, res.Error
	}

	res = s.db.Raw("SELECT count(id) FROM descriptors WHERE favorite = true").Scan(&favorites)
	if res.Error != nil {
		return Stats{}, res.Error
	}

	return Stats{
		Photos:    photos,
		Favorites: favorites,
		UsedSpace: usedSpace,
	}, nil
}
