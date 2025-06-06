package descriptor

import (
	"strings"
	"time"

	"github.com/google/uuid"
	img "github.com/inokone/photostorage/image"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Format is a representation of imge format
type Format string

// ParseFormat parses a string - usually file extension - and returns the image format
func ParseFormat(s string) Format {
	return Format(strings.ToLower(strings.TrimSpace(s)))
}

// Descriptor is a collection of metadata for a photo
type Descriptor struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	FileName    string    `gorm:"type:varchar(255);index;not null"`
	Uploaded    time.Time `gorm:"index"`
	Format      Format
	Tags        pq.StringArray `gorm:"type:text[]"`
	Favorite    bool           `gorm:"index"`
	Rating      int8           `gorm:"index"`
	Metadata    img.Metadata   `gorm:"foreignKey:MetadataID"`
	MetadataID  uuid.UUID
	ThumbWidth  int
	ThumbHeight int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

// AsResp converts `Descriptor` entity to a `Response“ entity
func (p Descriptor) AsResp() Response {
	return Response{
		ID:          p.ID.String(),
		FileName:    p.FileName,
		Uploaded:    p.Uploaded,
		Format:      string(p.Format),
		Metadata:    p.Metadata.AsResp(),
		Tags:        p.Tags,
		Favorite:    p.Favorite,
		ThumbWidth:  p.ThumbWidth,
		ThumbHeight: p.ThumbHeight,
		Rating:      p.Rating,
	}
}

// Response entity is a REST response representation of a `Descriptor“.
type Response struct {
	ID          string       `json:"id"`
	FileName    string       `json:"filename"`
	Uploaded    time.Time    `json:"uploaded"`
	Format      string       `json:"format"`
	Thumbnail   string       `json:"thumbnail"`
	ThumbWidth  int          `json:"thumbnail_width"`
	ThumbHeight int          `json:"thumbnail_height"`
	Metadata    img.Response `json:"metadata"`
	Tags        []string     `json:"tags"`
	Favorite    bool         `json:"favorite"`
	Rating      int8         `json:"rating"`
}
