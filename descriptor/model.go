package descriptor

import (
	"image"
	"time"

	"github.com/google/uuid"
	img "github.com/inokone/photostorage/image"
)

type Format int64

const (
	GP2 Format = iota
	ARW
)

func (d Format) String() string {
	return [...]string{"GP2", "ARW"}[d]
}

var Formats = map[string]Format{
	"gp2": GP2,
	"arw": ARW,
}

type Descriptor struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FileName  string    `gorm:"type:varchar(255);not null"`
	Uploaded  time.Time
	Format    Format
	Width     int
	Height    int
	Thumbnail image.Image `gorm:"-"`
	Tags      []string    `gorm:"type:text[]"`
}

func (p Descriptor) AsGet() (*Get, error) {
	thumbnail, error := img.ExportJpeg(p.Thumbnail)
	if error != nil {
		return nil, error
	}
	return &Get{
		ID:        p.ID.String(),
		FileName:  p.FileName,
		Uploaded:  p.Uploaded,
		Format:    p.Format.String(),
		Width:     p.Width,
		Height:    p.Height,
		Thumbnail: string(thumbnail),
		Tags:      p.Tags,
	}, nil
}

type Get struct {
	ID        string    `json:"id"`
	FileName  string    `json:"filename"`
	Uploaded  time.Time `json:"uploaded"`
	Format    string    `json:"format"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	Thumbnail string    `json:"thumbnail"`
	Tags      []string  `json:"tags"`
}
