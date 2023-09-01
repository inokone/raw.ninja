package photo

import (
	"io"
	"log"

	"image"

	"github.com/nf/cr2"
)

type RawProcessor interface {
	Process(raw io.Reader) (image.Image, error)
}

// Type for converting Canon raw images to Go image.Image
type Cr2Processor struct{}

func (p Cr2Processor) Process(raw io.Reader) (image.Image, error) {
	result, err := cr2.Decode(raw)
	if err != nil {
		log.Printf("Image processing failed with cause: %v", err)
	}
	return result, err
}
