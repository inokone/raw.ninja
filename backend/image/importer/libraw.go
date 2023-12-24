package importer

import (
	"fmt"
	"image"
	"math"
	"os"

	raw "github.com/inokone/golibraw"
	pi "github.com/inokone/photostorage/image"
	"github.com/rs/zerolog/log"
)

// LibrawImporter is an implementation of `Importer` using LibRAW library.
type LibrawImporter struct {
	tempDir string
	def     DefaultImporter
}

// NewLibrawImporter creates a new `LibrawImporter` instance.
func NewLibrawImporter() Importer {
	return LibrawImporter{
		tempDir: "photostore",
		def:     NewDefaultImporter(),
	}
}

// Image is a method of `LibrawImporter` for importing a RAW image byte array into an `image.Image`
func (p LibrawImporter) Image(rawBytes []byte) (*image.Image, error) {
	path, err := tempFile("image", rawBytes)
	defer removeTempFile(path)

	if err != nil {
		return nil, fmt.Errorf("RAW import error [%v]", err)
	}
	result, err := raw.ImportRaw(path)
	if err != nil {
		return nil, fmt.Errorf("RAW import error [%v]", err)
	}
	return &result, nil
}

// Describe is a method of `LibrawImporter` for importing the description from the RAW image byte array.
func (p LibrawImporter) Describe(rawBytes []byte) (*pi.Metadata, error) {
	path, err := tempFile("desc", rawBytes)
	defer removeTempFile(path)

	if err != nil {
		return nil, fmt.Errorf("metadata extract error [%v]", err)
	}
	metadata, err := raw.ExtractMetadata(path)
	if err != nil {
		return nil, fmt.Errorf("metadata extract error [%v]", err)
	}
	if math.IsNaN(metadata.Aperture) {
		metadata.Aperture = 0
	}
	if math.IsNaN(metadata.Shutter) {
		metadata.Shutter = 0
	}
	return &pi.Metadata{
		Height:    metadata.Height,
		Width:     metadata.Width,
		Timestamp: metadata.Timestamp,
		DataSize:  metadata.DataSize,
		Camera: pi.Camera{
			Make:     metadata.Camera.Make,
			Model:    metadata.Camera.Model,
			Software: metadata.Camera.Software,
		},
		Lens: pi.Lens{
			Make:  metadata.Lens.Make,
			Model: metadata.Lens.Model,
		},
		ISO:      metadata.ISO,
		Aperture: metadata.Aperture,
		Shutter:  metadata.Shutter,
	}, nil
}

// Thumbnail is a methof of `LibrawImporter` for extracting existing thumbnail image from the RAW image byte array.
// If the RAW image does not contain a thumbnail, this function generates one from the RAW image.
func (p LibrawImporter) Thumbnail(rawBytes []byte) (*pi.ThumbnailImg, error) {
	path, err := tempFile("raw", rawBytes)
	defer removeTempFile(path)
	if err != nil {
		return nil, fmt.Errorf("thumbnail extract error [%v]", err)
	}
	exportPath := tempPath("thumb")
	defer removeTempFile(exportPath)
	err = raw.ExtractThumbnail(path, exportPath)
	if err == nil {
		rs, err := os.ReadFile(exportPath)
		if err != nil {
			return nil, fmt.Errorf("thumbnail extract error [%v]", err)
		}
		im, err := p.def.Image(rs)
		if err != nil {
			return nil, fmt.Errorf("thumbnail extract error [%v]", err)
		}
		return &pi.ThumbnailImg{
			Image:  rs,
			Width:  (*im).Bounds().Dx(),
			Height: (*im).Bounds().Dy(),
		}, nil
	}
	log.Debug().AnErr("Thumbnail extraction", err).Msg("Failed to extract thumbnail")
	// most likely we have no thumbnail embedded in the RAW image, let's create one
	img, err := p.Image(rawBytes)
	if err != nil {
		return nil, err
	}
	*img, err = pi.Thumbnail(*img)
	if err != nil {
		return nil, err
	}
	res, err := pi.ExportJpeg(*img)
	if err != nil {
		return nil, err
	}
	return &pi.ThumbnailImg{
		Image:  res,
		Width:  (*img).Bounds().Dx(),
		Height: (*img).Bounds().Dy(),
	}, nil
}
