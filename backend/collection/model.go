package collection

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/photo"
	"gorm.io/gorm"
)

// Type is an enum for various collections.
type Type string

const (
	// Upload is a collection type for upload batches of photos
	Upload Type = "UPLOAD"
	// Album is a collection type for custom collections created by users
	Album Type = "ALBUM"
)

// Scan is a function to return a `CollectionType` for value
func (ct *Type) Scan(value interface{}) error {
	*ct = Type(value.([]byte))
	return nil
}

// Value is a function to return the SQL value for a `CollectionType`
func (ct Type) Value() (driver.Value, error) {
	return string(ct), nil
}

// Collection is a type for a collection of photos.
type Collection struct {
	ID        uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Owner     user.User     `gorm:"foreignKey:OwnerID"`
	OwnerID   string        `gorm:"index"`
	Type      Type          `gorm:"type:collection_type"`
	Name      string        `gorm:"type:varchar(255)"`
	Tags      []string      `gorm:"type:text[]"`
	Photos    []photo.Photo `gorm:"many2many:collection_photos;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// AsResp is a method of `Collection` to convert to JSON representation.
func (c Collection) AsResp() Resp {
	var photos []photo.Response = make([]photo.Response, len(c.Photos))
	for i, photo := range c.Photos {
		photos[i] = photo.AsResp()
	}
	return Resp{
		ID:        c.ID.String(),
		Name:      c.Name,
		Tags:      c.Tags,
		Photos:    photos,
		CreatedAt: c.CreatedAt,
	}
}

// Resp is a JSON type for photo collection responses
type Resp struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Tags      []string         `json:"tags"`
	Photos    []photo.Response `json:"photos"`
	CreatedAt time.Time        `json:"descriptor"`
}

// ListResp is a JSON representation
type ListResp struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"descriptor"`
}

// CreateAlbum is a JSON type for creating a new album type collection
type CreateAlbum struct {
	User     string   `json:"user"`
	Name     string   `json:"name"`
	Tags     []string `json:"tags"`
	PhotoIDs []string `json:"photos"`
}

// AlbumItems is a JSON type for changing photos of a collection
type AlbumItems struct {
	ID       string   `json:"id"`
	PhotoIDs []string `json:"photos"`
}

// UpdateAlbum is a JSON type for changing photos of a collection
type UpdateAlbum struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
