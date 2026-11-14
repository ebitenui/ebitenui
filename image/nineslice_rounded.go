package image

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// NewNineSliceRounded returns a nine-slice of a rectangle filled with fill,
// its corners rounded by radius pixels, and edged along its outline by a
// border width pixels wide in edge; a width of 0 draws no border, and a
// radius of 0 keeps the corners square. The outline is smoothed over a
// pixel. It stretches to any size at least as large as its two corners.
func NewNineSliceRounded(fill color.Color, radius int, edge color.Color, width int) *NineSlice {
	corner := max(radius, width)
	size := 2*corner + 1
	img := ebiten.NewImage(size, size)
	img.WritePixels(roundedPixels(fill, radius, edge, width, size))
	return NewNineSlice(img, [3]int{corner, 1, corner}, [3]int{corner, 1, corner})
}

// roundedPixels returns the premultiplied RGBA pixels of a size by size
// square drawn as NewNineSliceRounded describes.
func roundedPixels(fill color.Color, radius int, edge color.Color, width int, size int) []byte {
	pix := make([]byte, 4*size*size)
	fr, fg, fb, fa := fill.RGBA()
	er, eg, eb, ea := fr, fg, fb, fa
	if edge != nil {
		er, eg, eb, ea = edge.RGBA()
	}
	half, r := float64(size)/2, float64(radius)
	for y := range size {
		for x := range size {
			// The signed distance to the outline, negative inside.
			qx := math.Abs(float64(x)+0.5-half) - (half - r)
			qy := math.Abs(float64(y)+0.5-half) - (half - r)
			d := math.Hypot(math.Max(qx, 0), math.Max(qy, 0)) + math.Min(math.Max(qx, qy), 0) - r
			cover := clamp01(0.5 - d)
			// How far past the border the pixel is: 0 on it, 1 inside it.
			inside := 1.0
			if width > 0 {
				inside = clamp01(-d - float64(width) + 0.5)
			}
			mix := func(e, f uint32) byte {
				return byte(math.Round(cover * (float64(e)*(1-inside) + float64(f)*inside) / 0x101))
			}
			i := 4 * (y*size + x)
			pix[i], pix[i+1], pix[i+2], pix[i+3] = mix(er, fr), mix(eg, fg), mix(eb, fb), mix(ea, fa)
		}
	}
	return pix
}

func clamp01(v float64) float64 { return math.Min(math.Max(v, 0), 1) }
