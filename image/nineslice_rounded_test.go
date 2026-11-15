package image

import (
	"image/color"
	"testing"

	"github.com/matryer/is"
)

// pixelAt returns the premultiplied pixel at x, y of square pixels of a
// side.
func pixelAt(pix []byte, side, x, y int) [4]byte {
	i := 4 * (y*side + x)
	return [4]byte{pix[i], pix[i+1], pix[i+2], pix[i+3]}
}

// A rounded rectangle is filled in its middle, clear at the tip of a
// corner, partly covered along the curve, and edged by its border along
// its sides; square corners are filled to the tip.
func TestRoundedPixels(t *testing.T) {
	is := is.New(t)

	fill, edge := color.NRGBA{R: 40, G: 52, B: 80, A: 255}, color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	const radius, width = 6, 2
	side := 2*radius + 1
	pix := roundedPixels(fill, radius, edge, width, side)

	is.Equal(pixelAt(pix, side, radius, radius), [4]byte{40, 52, 80, 255})    // the middle: fill
	is.Equal(pixelAt(pix, side, 0, 0), [4]byte{})                             // the corner's tip: clear
	is.Equal(pixelAt(pix, side, 0, radius), [4]byte{255, 255, 255, 255})      // the left side: border
	is.Equal(pixelAt(pix, side, radius, side-1), [4]byte{255, 255, 255, 255}) // the bottom side: border
	a := pixelAt(pix, side, 1, 1)[3]
	is.True(a > 0 && a < 255) // along the curve: partly covered

	square := roundedPixels(fill, 0, nil, 0, 1)
	is.Equal(pixelAt(square, 1, 0, 0), [4]byte{40, 52, 80, 255}) // no radius, no border: fill

	edged := roundedPixels(fill, 0, edge, width, 2*width+1)
	is.Equal(pixelAt(edged, 2*width+1, 0, 0), [4]byte{255, 255, 255, 255}) // a square corner: border to the tip
	is.Equal(pixelAt(edged, 2*width+1, width, width), [4]byte{40, 52, 80, 255})
}

// The nine-slice keeps its corners and stretches its middle.
func TestNewNineSliceRounded(t *testing.T) {
	is := is.New(t)

	n := NewNineSliceRounded(color.White, 6, color.Black, 2)
	is.Equal(n.widths, [3]int{6, 1, 6})
	is.Equal(n.heights, [3]int{6, 1, 6})
	w, h := n.MinSize()
	is.Equal(w, 12)
	is.Equal(h, 12)
}
