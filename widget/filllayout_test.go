package widget

import (
	"image"
	"testing"

	"github.com/matryer/is"
)

func TestFillLayout_PreferredSize(t *testing.T) {
	is := is.New(t)

	padding := Insets{
		Top:    10,
		Left:   20,
		Right:  30,
		Bottom: 40,
	}

	l := newFillLayout(t,
		FillLayoutOpts.WithPadding(padding))

	wi := newSimpleWidget(35, 45, nil)

	w, h := l.PreferredSize([]HasWidget{wi})

	is.Equal(w, wi.preferredWidth+padding.Dx())
	is.Equal(h, wi.preferredHeight+padding.Dy())
}

func TestFillLayout_Layout(t *testing.T) {
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
