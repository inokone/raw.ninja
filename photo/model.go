package photo

import (
	"time"

	"github.com/inokone/photostorage/auth"
	"github.com/inokone/photostorage/descriptor"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type Photo struct {
	ID        uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Raw       []byte                `gorm:"-"`
	User      auth.User             `gorm:"foreignKey:ID"`
	Desc      descriptor.Descriptor `gorm:"foreignKey:ID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (p Photo) AsResp() (*Response, error) {
	desc, error := p.Desc.AsResp()
	if error != nil {
		return nil, error
	}
	return &Response{
		ID:   p.ID.String(),
		Desc: *desc,
	}, nil
}

type Response struct {
	ID   string              `json:"id"`
	Desc descriptor.Response `json:"descriptor"`
}
