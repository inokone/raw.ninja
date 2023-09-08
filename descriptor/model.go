package descriptor

import (
	"image"
	"time"

	"github.com/google/uuid"
)

type Format int64

const (
	GP2 Format = iota
	ARW
)

type Descriptor struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FileName  string
	Uploaded  time.Time
	Format    Format
	Width     int
	Height    int
	Thumbnail image.Image `gorm:"-"`
	Tags      []string
}
