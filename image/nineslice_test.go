package image

import (
	"image"
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

func TestNineSlice_MinSize(t *testing.T) {
	is := is.New(t)

	n := NewNineSlice(newImageEmptySize(20, 20, t), [3]int{3, 10, 7}, [3]int{2, 16, 2})
	w, h := n.MinSize()
	is.Equal(w, 10)
	is.Equal(h, 4)

	n = NewNineSliceColor(color.White)
	w, h = n.MinSize()
	is.Equal(w, 0)
	is.Equal(h, 0)

	n = NewNineSliceColor(color.Transparent)
	w, h = n.MinSize()
	is.Equal(w, 0)
	is.Equal(h, 0)
}

func newImageEmptySize(width int, height int, t *testing.T) *ebiten.Image {
	t.Helper()
	return ebiten.NewImage(width, height)
}

func Test_NewNineSliceColor(t *testing.T) {
	is := is.New(t)

	n := NewNineSliceColor(nil)
	is.Equal(n.transparent, true)

	n = NewNineSliceColor(color.RGBA{0, 0, 0, 0})
	is.Equal(n.transparent, true)
}

func Test_NewNineSliceBorder(t *testing.T) {
    is := is.New(t)

    img := newImageEmptySize(20, 20, t)
    n := NewNineSliceBorder(img, 3)

    is.Equal(n.widths[0], 3)
    is.Equal(n.widths[1], 20 - 3 * 2)
    is.Equal(n.widths[2], 3)
    is.Equal(n.heights[0], 3)
    is.Equal(n.heights[1], 20 - 3 * 2)
    is.Equal(n.heights[2], 3)
}

// A nine-slice is drawn when any of it lands on the image, and skipped when
// none does; its options are set once to place it, then once for each tile.
func TestNineSlice_SkipsDrawingOffTheImage(t *testing.T) {
	is := is.New(t)

	n := NewNineSliceSimple(ebiten.NewImage(12, 12), 4, 4)
	screen := ebiten.NewImageWithOptions(image.Rect(100, 100, 200, 200), nil)
	calls := func(opts func(*ebiten.DrawImageOptions)) int {
		count := 0
		n.Draw(screen, 20, 20, func(o *ebiten.DrawImageOptions) {
			count++
			opts(o)
		})
		return count
	}
	at := func(x, y float64) func(*ebiten.DrawImageOptions) {
		return func(o *ebiten.DrawImageOptions) { o.GeoM.Translate(x, y) }
	}

	is.Equal(calls(at(120, 120)), 10) // inside
	is.Equal(calls(at(81, 81)), 10)   // a pixel in, at the top left
	is.Equal(calls(at(199, 199)), 10) // a pixel in, at the bottom right
	is.Equal(calls(at(80, 120)), 1)   // touching the left edge
	is.Equal(calls(at(120, 200)), 1)  // touching the bottom edge
	is.Equal(calls(at(0, 0)), 1)      // outside
	is.Equal(calls(func(o *ebiten.DrawImageOptions) {
		o.GeoM.Scale(-1, 1) // mirrored onto the left of x
		o.GeoM.Translate(110, 120)
	}), 10)
}
