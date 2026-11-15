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

// A graphic draws Image while it is set, otherwise its images for its
// state, and picking one allocates nothing.
func TestGraphic_DrawsItsImage(t *testing.T) {
	is := is.New(t)

	idle, disabled, single := newImageEmptySize(2, 2, t), newImageEmptySize(2, 2, t), newImageEmptySize(3, 3, t)
	g := newGraphic(t, GraphicOpts.Images(&GraphicImage{Idle: idle, Disabled: disabled}))
	is.Equal(g.image(), idle)
	g.GetWidget().Disabled = true
	is.Equal(g.image(), disabled)

	g.Image = single
	is.Equal(g.image(), single)
	is.Equal(testing.AllocsPerRun(10, func() { g.image() }), 0.0)
	render(g, t)
	g.Image = nil
	is.Equal(g.image(), disabled) // its images again
}

func newGraphic(t *testing.T, opts ...GraphicOpt) *Graphic {
	t.Helper()

	g := NewGraphic(opts...)
	event.ExecuteDeferred()
	render(g, t)
	return g
}
