package descriptor

import (
	"strings"
	"time"

	"github.com/google/uuid"
	img "github.com/inokone/photostorage/image"
	"gorm.io/gorm"
)

type Format string

func ParseFormat(s string) Format {
	return Format(strings.ToLower(strings.TrimSpace(s)))
}

type Descriptor struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FileName   string    `gorm:"type:varchar(255);not null"`
	Uploaded   time.Time
	Format     Format
	Thumbnail  []byte       `gorm:"-"`
	Tags       []string     `gorm:"type:text[]"`
	Metadata   img.Metadata `gorm:"foreignKey:MetadataID"`
	MetadataID string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func (p Descriptor) AsResp() Response {
	return Response{
		ID:        p.ID.String(),
		FileName:  p.FileName,
		Uploaded:  p.Uploaded,
		Format:    string(p.Format),
		Thumbnail: string(p.Thumbnail),
		Metadata:  p.Metadata.AsResp(),
		Tags:      p.Tags,
	}
}

type Response struct {
	ID        string       `json:"id"`
	FileName  string       `json:"filename"`
	Uploaded  time.Time    `json:"uploaded"`
	Format    string       `json:"format"`
	Thumbnail string       `json:"thumbnail"`
	Metadata  img.Response `json:"metadata"`
	Tags      []string     `json:"tags"`
}
