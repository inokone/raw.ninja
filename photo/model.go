package photo

import (
	"time"

	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/descriptor"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

// Photo is a struct representing a photo object including image, thumbnail and metadata.
type Photo struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Raw       []byte    `gorm:"-"`
	UserID    string
	User      auth.User `gorm:"foreignKey:UserID"`
	DescID    string
	Desc      descriptor.Descriptor `gorm:"foreignKey:DescID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// AsResp is a method of the `Photo` struct. It converts a `Photo` object into a `Response` object.
func (p Photo) AsResp(baseURL string) Response {
	desc := p.Desc.AsResp(baseURL)
	return Response{
		ID:   p.ID.String(),
		Desc: desc,
	}
}

// Response is the JSON representation of `Photo` when retrieving from the application
type Response struct {
	ID   string              `json:"id"`
	Desc descriptor.Response `json:"descriptor"`
}

// UserStats is aggregated data on the photos of a user.
type UserStats struct {
	ID        uuid.UUID
	Photos    int
	Favorites int
	UsedSpace int64
}
