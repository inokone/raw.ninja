package photo

import (
	"time"

	"github.com/inokone/photostorage/auth/user"
	"github.com/inokone/photostorage/image"
	"github.com/inokone/photostorage/photo/descriptor"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

// Photo is a struct representing a photo object including image, thumbnail and metadata.
type Photo struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Raw       []byte    `gorm:"-"`
	Thumbnail []byte    `gorm:"-"`
	UserID    string    `gorm:"index"`
	User      user.User `gorm:"foreignKey:UserID"`
	DescID    string
	Desc      descriptor.Descriptor `gorm:"foreignKey:DescID"`
	UsedSpace int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// AsResp is a method of the `Photo` struct. It converts a `Photo` object into a `Response` object.
func (p Photo) AsResp() Response {
	desc := p.Desc.AsResp()
	return Response{
		ID:   p.ID.String(),
		Desc: desc,
	}
}

// Response is the JSON representation of `Photo` when retrieving from the application
type Response struct {
	ID        string                  `json:"id"`
	Desc      descriptor.Response     `json:"descriptor"`
	Raw       *image.PresignedRequest `json:"raw"`
	Thumbnail *image.PresignedRequest `json:"thumbnail"`
}

// UserStats is aggregated data on the photos of a user.
type UserStats struct {
	ID        uuid.UUID
	Photos    int
	Favorites int
	UsedSpace int64
}

// Stats is aggregated data on the storer.
type Stats struct {
	Photos    int
	Favorites int
	UsedSpace int64
}

// UploadSuccess is a JSON response type for upload results.
type UploadSuccess struct {
	PhotoIDs []string `json:"photo_ids"`
	UserID   string   `json:"user_id"`
}
