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
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	DataSize    int64     `json:"data_size"`
	ISO         int       `json:"ISO"`
	Aperture    float64   `json:"aperture"`
	Shutter     float64   `json:"shutter"`
	Timestamp   time.Time `json:"timestamp"`
	CameraMake  string    `json:"camera_make"`
	CameraModel string    `json:"camera_model"`
	CameraSW    string    `json:"camera_sw"`
	Colors      uint      `json:"colors"`
	LensMake    string    `json:"lens_make"`
	LensModel   string    `json:"lens_model"`
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
