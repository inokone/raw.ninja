package photo

import (
	"github.com/inokone/photostorage/descriptor"
	"github.com/inokone/photostorage/user"

	"github.com/google/uuid"
)

type Photo struct {
	ID   uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Raw  []byte                `gorm:"-"`
	User user.User             `gorm:"foreignKey:ID"`
	Desc descriptor.Descriptor `gorm:"foreignKey:ID"`
}
