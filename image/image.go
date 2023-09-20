package image

import (
	"image"
	"image/color"
)

type LibrawImage struct {
	img **uint16
}

func (i LibrawImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (i LibrawImage) Bounds() image.Rectangle {
	return
}

func (i LibrawImage) At(x, y int) color.Color {

}
