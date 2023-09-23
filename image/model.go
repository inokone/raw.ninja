package image

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Camera struct {
	Make     string
	Model    string
	Software string
	Colors   uint
}

type Lens struct {
	Make   string
	Model  string
	Serial string
}

type Metadata struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Timestamp int64
	Width     int
	Height    int
	DataSize  int64
	Camera    Camera `gorm:"embedded;embeddedPrefix:camera_"`
	Lens      Lens   `gorm:"embedded;embeddedPrefix:lens_"`
	ISO       int
	Aperture  float64
	Shutter   float64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type Response struct {
	Width       int `json:"width"`
	Height      int `json:"height"`
	DataSize    int64
	ISO         int
	Aperture    float64
	Shutter     float64
	Timestamp   time.Time
	CameraMake  string
	CameraModel string
	CameraSW    string
	Colors      uint
	LensMake    string
	LensModel   string
}

func (m Metadata) AsResp() Response {
	return Response{
		Width:       m.Width,
		Height:      m.Height,
		DataSize:    m.DataSize,
		ISO:         m.ISO,
		Aperture:    m.Aperture,
		Shutter:     m.Shutter,
		Timestamp:   time.Unix(m.Timestamp, 0),
		CameraMake:  m.Camera.Make,
		CameraModel: m.Camera.Model,
		CameraSW:    m.Camera.Software,
		Colors:      m.Camera.Colors,
		LensMake:    m.Lens.Make,
		LensModel:   m.Lens.Model,
	}
}
