package widget

import (
	"testing"

	"github.com/ebitenui/ebitenui/event"
	"github.com/matryer/is"
)

func TestGraphic_PreferredSize(t *testing.T) {
	is := is.New(t)

	i := newImageEmptySize(47, 11, t)
	g := newGraphic(t, GraphicOpts.Image(i))
	w, h := g.PreferredSize()
	is.Equal(w, i.Bounds().Dx())
	is.Equal(h, i.Bounds().Dy())
}

func newGraphic(t *testing.T, opts ...GraphicOpt) *Graphic {
	t.Helper()

	g := NewGraphic(opts...)
	event.ExecuteDeferred()
	render(g, t)
	return g
}
