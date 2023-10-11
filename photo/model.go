package photo

import (
	"time"

	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/descriptor"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

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

func (p Photo) AsResp(baseUrl string) Response {
	desc := p.Desc.AsResp(baseUrl)
	return Response{
		ID:   p.ID.String(),
		Desc: desc,
	}
}

type Response struct {
	ID   string              `json:"id"`
	Desc descriptor.Response `json:"descriptor"`
}
