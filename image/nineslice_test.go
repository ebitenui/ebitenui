package image

import (
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
