package image

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	raw "github.com/inokone/golibraw"
	"github.com/rs/zerolog/log"
)

// Importer is an interface for importing RAW camera images.
type Importer interface {
	Image(raw []byte) (*image.Image, error)

	Describe(raw []byte) (*Metadata, error)

	Thumbnail(raw []byte) ([]byte, error)
}

// LibrawImporter is an implementation of `Importer` using LibRAW library.
type LibrawImporter struct {
	tempDir string
}

// NewLibrawImporter creates a new `LibrawImporter` instance.
func NewLibrawImporter() Importer {
	return LibrawImporter{
		tempDir: "photostore",
	}
}

// Image is a method of `LibrawImporter` for importing a RAW image byte array into an `image.Image`
func (p LibrawImporter) Image(rawBytes []byte) (*image.Image, error) {
	path, err := p.tempFile("image", rawBytes)
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
func (p LibrawImporter) Describe(rawBytes []byte) (*Metadata, error) {
	path, err := p.tempFile("desc", rawBytes)
	defer removeTempFile(path)

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

// Thumbnail is a methof of `LibrawImporter` for extracting existing thumbnail image from the RAW image byte array.
// If the RAW image does not contain a thumbnail, this function generates one from the RAW image.
func (p LibrawImporter) Thumbnail(rawBytes []byte) ([]byte, error) {
	start := time.Now()
	path, err := p.tempFile("raw", rawBytes)
	defer removeTempFile(path)
	if err != nil {
		return nil, fmt.Errorf("thumbnail extract error [%v]", err)
	}
	log.Debug().Dur("Elapsed time", time.Since(start)).Msg("Temp file created for raw.")

	exportPath := p.tempPath("thumb")
	defer removeTempFile(exportPath)
	log.Debug().Dur("Elapsed time", time.Since(start)).Msg("Temp path created for thumb.")

	err = raw.ExtractThumbnail(path, exportPath)
	if err == nil {
		return os.ReadFile(exportPath)
	}
	log.Debug().AnErr("Thumbnail extraction", err).Msg("Failed to extract thumbnail")
	log.Debug().Dur("Elapsed time", time.Since(start)).Msg("Thumbnail extraction finished.")
	// most likely we have no thumbnail embedded in the RAW image, let's create one
	img, err := p.Image(rawBytes)
	if err != nil {
		return nil, err
	}
	log.Debug().Dur("Elapsed time", time.Since(start)).Msg("Image bytes loaded.")
	*img, err = Thumbnail(*img)
	if err != nil {
		return nil, err
	}
	log.Debug().Dur("Elapsed time", time.Since(start)).Msg("Thumbnail generated.")
	res, err := ExportJpeg(*img)
	log.Debug().Dur("Elapsed time", time.Since(start)).Msg("JPEG thumbnail exported.")
	return res, err
}

func (p LibrawImporter) tempFile(target string, content []byte) (string, error) {
	f, err := os.CreateTemp("", target+"_*")
	if err != nil {
		return "", err
	}
	defer closeTempFile(f)

	_, err = f.Write(content)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func (p LibrawImporter) tempPath(target string) string {
	tempDir := os.TempDir()
	return filepath.Join(tempDir, target+"_"+uuid.New().String())
}

func removeTempFile(path string) {
	if err := os.Remove(path); err != nil {
		log.Warn().AnErr("Cause", err).Str("Path", path).Msg("Could not remove temp file.")
	}
}

func closeTempFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Warn().AnErr("Cause", err).Str("Name", f.Name()).Msg("Could not close temp file.")
	}
}
