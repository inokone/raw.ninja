package image

import (
	"testing"
)

type CanvasTest struct {
	inX  int
	inY  int
	outX int
	outY int
}

var tests = []CanvasTest{
	{200, 200, 200, 200},
	{300, 400, 150, 200},
	{100, 50, 200, 100},
	{2344, 1540, 200, 131},
}

func TestCanvas(t *testing.T) {
	for _, test := range tests {
		actual := canvas(test.inX, test.inY)
		if actual.Bounds().Size().X != test.outX || actual.Bounds().Size().Y != test.outY {
			t.Errorf("Canvas(%v, %v) = (%v, %v); want (%v, %v)", test.inX, test.inY, actual.Bounds().Size().X, actual.Bounds().Size().Y, test.outX, test.outY)
		}
	}
}
