package importer

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"

	img "github.com/inokone/photostorage/image"
	"github.com/rs/zerolog/log"
	"github.com/rwcarlsen/goexif/exif"
)

// DefaultImporter is an implementation of `Importer` using Libvips anf Goexif libraries.
type DefaultImporter struct{}

// NewDefaultImporter creates a new `DefaultImporter` instance, setting up format regsitrations.
func NewDefaultImporter() DefaultImporter {
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
	return DefaultImporter{}
}

// Image is a method of `DefaultImporter` for importing an image byte array into an `image.Image`
func (i DefaultImporter) Image(raw []byte) (*image.Image, error) {
	im, _, err := image.Decode(bytes.NewReader(raw))
	return &im, err
}

// Describe is a method of `DefaultImporter` for importing EXIF metadata from the image
func (i DefaultImporter) Describe(raw []byte) (*img.Metadata, error) {
	var (
		err error
		m   *exif.Exif
		b   []byte
		js  string
	)

	m, err = exif.Decode(bytes.NewReader(raw))
	if err != nil {
		return i.noExif(raw)
	}
	b, err = m.MarshalJSON()
	if err != nil {
		return nil, err
	}

	js = string(b)
	log.Debug().Str("data", js).Msg("EXIF")

	return &img.Metadata{
		Width:  asInt(m, exif.PixelXDimension),
		Height: asInt(m, exif.PixelYDimension),
		Camera: img.Camera{
			Make:     asString(m, exif.Make),
			Model:    asString(m, exif.Model),
			Software: asString(m, exif.Software),
		},
		Lens: img.Lens{
			Make:  asString(m, exif.LensMake),
			Model: asString(m, exif.LensModel),
		},
		Aperture:  asFloat(m, exif.FNumber),
		Shutter:   asApex(m, exif.ShutterSpeedValue),
		ISO:       asInt(m, exif.ISOSpeedRatings),
		DataSize:  int64(len(raw)),
		Timestamp: asTime(m),
	}, nil
}

func (i DefaultImporter) noExif(raw []byte) (*img.Metadata, error) {
	im, _, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	return &img.Metadata{
		Width:    im.Width,
		Height:   im.Height,
		DataSize: int64(len(raw)),
	}, nil
}

func asApex(m *exif.Exif, f exif.FieldName) float64 {
	focal, err := m.Get(f)
	if err != nil {
		return 0
	}
	numer, denom, err := focal.Rat2(0)
	if err != nil {
		return 0
	}
	res := math.Pow(2, -float64(numer)/float64(denom))
	if math.IsNaN(res) {
		return 0
	}
	return res
}

func asFloat(m *exif.Exif, f exif.FieldName) float64 {
	focal, err := m.Get(f)
	if err != nil {
		return 0
	}
	numer, denom, err := focal.Rat2(0)
	if err != nil {
		return 0
	}
	res := float64(numer) / float64(denom)
	if math.IsNaN(res) {
		return 0
	}
	return res
}

func asInt(m *exif.Exif, f exif.FieldName) int {
	t, err := m.Get(f)
	if err != nil {
		return 0
	}
	i, err := t.Int(0)
	if err != nil {
		return 0
	}
	return i
}

func asString(m *exif.Exif, f exif.FieldName) string {
	t, err := m.Get(f)
	if err != nil {
		return ""
	}
	i := t.String()
	return i
}

func asTime(m *exif.Exif) int64 {
	time, err := m.DateTime()
	if err != nil {
		return 0
	}
	return time.Unix()
}

// Thumbnail is a method of `DefaultImporter` for generating a thumbnail image byte array for an image bye array.
func (i DefaultImporter) Thumbnail(raw []byte) (*img.ThumbnailImg, error) {
	var (
		im  *image.Image
		err error
		f   string
		tn  image.Image
		res []byte
	)
	im, err = i.Image(raw)
	if err != nil {
		f, _ = tempFile("forensics", raw)
		log.Warn().Str("path", f).Msg("Image import failed, writing forensics file.")
		return nil, err
	}
	tn, err = img.Thumbnail(*im)
	if err != nil {
		return nil, err
	}
	res, err = img.ExportJpeg(tn)
	if err != nil {
		return nil, err
	}
	return &img.ThumbnailImg{
		Image:  res,
		Width:  tn.Bounds().Dx(),
		Height: tn.Bounds().Dy(),
	}, nil
}
