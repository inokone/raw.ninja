package photo

import (
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/image"
)

type Repository struct {
	DB *gorm.DB
	Ir image.Repository
}

func (s *Repository) Create(photo Photo) (uuid.UUID, error) {
	s.DB.Save(&photo)
	err := s.Ir.Create(photo.ID.String(), photo.Raw, photo.Desc.Thumbnail)
	return photo.ID, err
}

func (s *Repository) Delete(id string) error {
	var photo Photo
	s.DB.Delete(&photo, "id = ?", id)
	return nil
}

func (s *Repository) Update(photo Photo) error {
	s.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&photo)
	return nil
}

func (s *Repository) Get(id string) (Photo, error) {
	var photo Photo
	result := s.DB.Preload("Desc.Metadata").First(&photo, "id = ?", id)
	return photo, result.Error
}

func (s *Repository) All(userID string) ([]Photo, error) {
	var photos []Photo
	result := s.DB.Preload("Desc.Metadata").Where("user_id = ?", userID).Order("created_at ASC").Find(&photos)
	return photos, result.Error
}

func (s *Repository) Favorites(userID string) ([]Photo, error) {
	var photos []Photo
	result := s.DB.Preload("Desc.Metadata").Joins("JOIN descriptors ON descriptors.id = photos.desc_id").Where("photos.user_id = ?", userID).Where("descriptors.favorite = true").Order("photos.created_at DESC").Find(&photos)
	return photos, result.Error
}

func (s *Repository) Search(userID string, searchText string) ([]Photo, error) {
	var photos []Photo
	result := s.DB.Preload("Desc.Metadata").Joins("JOIN descriptors ON descriptors.id = photos.desc_id").Where("photos.user_id = ?", userID).Where("descriptors.file_name LIKE ?", "%"+searchText+"%").Order("photos.created_at ASC").Find(&photos)
	return photos, result.Error
}

func (s *Repository) Raw(id string) ([]byte, error) {
	return s.Ir.Image(id)
}

func (s *Repository) Thumbnail(id string) ([]byte, error) {
	return s.Ir.Thumbnail(id)
}

func (s *Repository) Statistics(userID string) (UserStats, error) {
	var photos, favorites int
	s.DB.Raw("SELECT count(id) FROM photos WHERE user_id = ?", userID).Scan(&photos)
	s.DB.Raw("SELECT count(p.id) FROM photos p JOIN descriptors d ON d.id = p.desc_id WHERE p.user_id = ? and d.favorite = true", userID).Scan(&favorites)

	photoList, err := s.All(userID)
	if err != nil {
		return UserStats{}, err
	}
	photoIDs := make([]string, len(photoList))

	for idx, photo := range photoList {
		photoIDs[idx] = photo.ID.String()
	}

	usedSpace, err := s.Ir.UsedSpace(photoIDs)
	if err != nil {
		return UserStats{}, err
	}

	return UserStats{
		Photos:    photos,
		Favorites: favorites,
		UsedSpace: usedSpace,
	}, nil
}
