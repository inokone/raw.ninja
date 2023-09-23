package image

import (
	"fmt"
	"image"
	"os"

	"github.com/google/uuid"
	raw "github.com/inokone/golibraw"
)

type Importer interface {
	Image(raw []byte) (*image.Image, error)

	Describe(raw []byte) (*Metadata, error)

	Thumbnail(raw []byte) ([]byte, error)
}

type LibrawImporter struct {
	uuid uuid.UUID
}

func (p LibrawImporter) Image(rawBytes []byte) (*image.Image, error) {
	path, err := p.tempFile(rawBytes)
	defer os.Remove(path)

	if err != nil {
		return nil, fmt.Errorf("RAW import error [%v]", err)
	}
	result, err := raw.ImportRaw(path)
	if err != nil {
		return nil, fmt.Errorf("RAW import error [%v]", err)
	}
	return &result, nil
}

func (p LibrawImporter) Describe(rawBytes []byte) (*Metadata, error) {
	path, err := p.tempFile(rawBytes)
	defer os.Remove(path)

	if err != nil {
		return nil, fmt.Errorf("metadata extract error [%v]", err)
	}
	metadata, err := raw.ExtractMetadata(path)
	if err != nil {
		return nil, fmt.Errorf("metadata extract error [%v]", err)
	}
	return &Metadata{
		Height:    metadata.Height,
		Width:     metadata.Width,
		Timestamp: metadata.Timestamp,
		DataSize:  metadata.DataSize,
		Camera: Camera{
			Make:     metadata.Camera.Make,
			Model:    metadata.Camera.Model,
			Software: metadata.Camera.Software,
		},
		Lens: Lens{
			Make:  metadata.Lens.Make,
			Model: metadata.Lens.Model,
		},
		ISO:      metadata.ISO,
		Aperture: metadata.Aperture,
		Shutter:  metadata.Shutter,
	}, nil
}

func (p LibrawImporter) Thumbnail(rawBytes []byte) ([]byte, error) {
	path, err := p.tempFile(rawBytes)
	defer os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("thumbnail extract error [%v]", err)
	}

	exportPath, err := p.tempFile(make([]byte, 0))
	defer os.Remove(path)
	if err != nil {
		return nil, fmt.Errorf("thumbnail extract error [%v]", err)
	}

	err = raw.ExtractThumbnail(path, exportPath)
	if err == nil {
		return os.ReadFile(path)
	}

	// most likely we have no thumbnail embedded in the RAW image, let's create one
	img, err := p.Image(rawBytes)
	if err != nil {
		return nil, err
	}
	*img, err = Thumbnail(*img)
	if err != nil {
		return nil, err
	}
	return ExportJpeg(*img)
}

func (l LibrawImporter) tempFile(content []byte) (string, error) {
	f, err := os.CreateTemp("", fmt.Sprintf("%v_*", l.uuid.String()))
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func NewImporter() Importer {
	return LibrawImporter{
		uuid: uuid.New(),
	}
}
