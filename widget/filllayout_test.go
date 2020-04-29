package widget

import (
	"image"
	"testing"

	"github.com/matryer/is"
)

func TestFillLayout_Layout_Padding(t *testing.T) {
	is := is.New(t)

	l := newFillLayout(t,
		FillLayoutOpts.WithPadding(Insets{
			Top:    10,
			Left:   20,
			Right:  30,
			Bottom: 40,
		}))

	b := newButton(t)
	l.Layout([]HasWidget{b}, image.Rect(25, 25, 100, 100))

	is.Equal(b.GetWidget().Rect, image.Rect(45, 35, 70, 60))
}

func newFillLayout(t *testing.T, opts ...FillLayoutOpt) Layouter {
	t.Helper()
	l := NewFillLayout(opts...)
	l.(Dirtyable).MarkDirty()
	return l
}
