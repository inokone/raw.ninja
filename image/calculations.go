package image

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"os"

	"log"

	"golang.org/x/image/draw"

	raw "github.com/MRHT-SRProject/LibRawGo/librawgo"
)

type RawProcessor interface {
	Process(raw []byte) (image.Image, error)
}

func importRaw(path string) (image.Image, error) {
	lr := raw.Libraw_init(0)
	if int(raw.LIBRAW_SUCCESS) != raw.Libraw_open_file(lr, path) {
		return nil, fmt.Errorf("Libraw import error at path [%v]", path)
	}
	return LibrawImage{
		img: lr.GetImage(),
	}, nil
}

// Type for converting Canon raw images to Go image.Image
type Cr2Processor struct{}

func (p *Cr2Processor) Process(raw []byte) (image.Image, error) {
	file, err := os.CreateTemp("", "raw_import_")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())
	_, err = file.Write(raw)
	if err != nil {
		return nil, err
	}
	result, err := importRaw(file.Name())
	if err != nil {
		log.Printf("Image processing failed with cause: %v", err)
	}
	return result, err
}

func Factory(extension string) (RawProcessor, error) {
	if extension == "cr2" {
		return new(Cr2Processor), nil
	}
	return nil, image.ErrFormat
}

const (
	thumbWidth  float64 = 200
	thumbHeight float64 = 200
)

func Thumbnail(original image.Image) (image.Image, error) {
	ratio := math.Max(float64(original.Bounds().Size().X)/thumbWidth, float64(original.Bounds().Size().Y)/thumbHeight)
	width := int(ratio * float64(original.Bounds().Size().X))
	height := int(ratio * float64(original.Bounds().Size().Y))
	// Create a new RGBA image with a white background
	result := image.NewRGBA(image.Rect(0, 0, width, height))
	// Resize
	draw.NearestNeighbor.Scale(result, result.Rect, original, original.Bounds(), draw.Over, nil)
	return result, nil
}

func ExportJpeg(image image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)
	return buf.Bytes(), err
}

func ImportJpeg(b []byte) (image.Image, error) {
	return jpeg.Decode(bytes.NewReader(b))
}
