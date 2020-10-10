package image

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/matryer/is"
)

func TestBufferedImage_Image(t *testing.T) {
	is := is.New(t)

	b := &BufferedImage{}
	b.Width, b.Height = 100, 100
	i := b.Image()
	w, h := i.Size()
	is.Equal(w, 100)
	is.Equal(h, 100)

	b.Width, b.Height = 150, 70
	i = b.Image()
	w, h = i.Size()
	is.Equal(w, 150)
	is.Equal(h, 70)
}

func TestMaskedRenderBuffer_Draw(t *testing.T) {
	is := is.New(t)

	b := NewMaskedRenderBuffer()
	screen := newImageEmptySize(100, 100, t)

	draw := false
	drawMask := false

	b.Draw(screen, func(buf *ebiten.Image) {
		draw = true
	}, func(buf *ebiten.Image) {
		drawMask = true
	})

	is.True(draw)
	is.True(drawMask)
}
