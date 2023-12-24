package image

import (
	"bytes"
	"image"
	"image/jpeg"
	"math"

	"golang.org/x/image/draw"
)

const (
	thumbWidth  float64 = 1000
	thumbHeight float64 = 1000
)

// Thumbnail is a function to generate a thumbnail image of max size [`thumbWidth`, `thumbHeight`] for the image provided as a parameter.
func Thumbnail(original image.Image) (image.Image, error) {
	result := canvas(original.Bounds().Size().X, original.Bounds().Size().Y)
	draw.NearestNeighbor.Scale(result, result.Rect, original, original.Bounds(), draw.Over, nil)
	return result, nil
}

func canvas(width int, height int) *image.RGBA {
	ratio := math.Min(thumbWidth/float64(width), thumbHeight/float64(height))
	newWidth := int(ratio * float64(width))
	newHeight := int(ratio * float64(height))
	result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	return result
}

// ExportJpeg is a function to export the image provided as parameter as a byte array in JPEG format.
func ExportJpeg(image image.Image) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, image, nil)
	return buf.Bytes(), err
}

// ImportJpeg is a function to import a byte array in JPEG format into an Image object.
func ImportJpeg(b []byte) (image.Image, error) {
	return jpeg.Decode(bytes.NewReader(b))
}
