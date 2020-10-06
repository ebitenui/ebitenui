package widget

import (
	"image"
	"testing"

	"github.com/matryer/is"
)

func TestAnchorLayout_PreferredSize_Fill(t *testing.T) {
	is := is.New(t)

	padding := Insets{
		Top:    10,
		Left:   20,
		Right:  30,
		Bottom: 40,
	}

	l := newAnchorLayout(t, AnchorLayoutOpts.Padding(padding))

	wi := newSimpleWidget(35, 45, &AnchorLayoutData{
		StretchHorizontal: true,
		StretchVertical:   true,
	})

	w, h := l.PreferredSize([]PreferredSizeLocateableWidget{wi})

	is.Equal(w, wi.preferredWidth+padding.Dx())
	is.Equal(h, wi.preferredHeight+padding.Dy())
}

func TestAnchorLayout_Layout_Fill(t *testing.T) {
	is := is.New(t)

	l := newAnchorLayout(t, AnchorLayoutOpts.Padding(Insets{
		Top:    10,
		Left:   20,
		Right:  30,
		Bottom: 40,
	}))

	b := newButton(t, ButtonOpts.WidgetOpts(WidgetOpts.LayoutData(&AnchorLayoutData{
		StretchHorizontal: true,
		StretchVertical:   true,
	})))
	l.Layout([]PreferredSizeLocateableWidget{b}, image.Rect(25, 25, 100, 100))

	is.Equal(b.GetWidget().Rect, image.Rect(45, 35, 70, 60))
}

func newAnchorLayout(t *testing.T, opts ...AnchorLayoutOpt) Layouter {
	t.Helper()
	l := NewAnchorLayout(opts...)
	return l
}
