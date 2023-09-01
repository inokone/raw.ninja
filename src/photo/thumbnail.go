package photo

import (
	"image"
	"math"

	"golang.org/x/image/draw"
)

const thumbWidth float64 = 200
const thumbHeight float64 = 200

func thumbnail(original image.Image) (image.Image, error) {
	ratio := math.Max(float64(original.Bounds().Size().X)/thumbWidth, float64(original.Bounds().Size().Y)/thumbHeight)
	width := int(ratio * float64(original.Bounds().Size().X))
	height := int(ratio * float64(original.Bounds().Size().Y))
	// Create a new RGBA image with a white background
	result := image.NewRGBA(image.Rect(0, 0, width, height))
	// Resize
	draw.NearestNeighbor.Scale(result, result.Rect, original, original.Bounds(), draw.Over, nil)
	return result, nil
}
