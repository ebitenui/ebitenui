package widget

import (
	"image"
	"testing"

	"github.com/matryer/is"
)

func TestRowLayout_Layout(t *testing.T) {
	is := is.New(t)

	l := newRowLayout(t,
		RowLayoutOpts.WithPadding(Insets{
			Top:    10,
			Left:   20,
			Right:  30,
			Bottom: 40,
		}),
		RowLayoutOpts.WithSpacing(7))

	widgets := []HasWidget{
		newSimpleWidget(10, 10, &RowLayoutData{}),
		newSimpleWidget(20, 20, &RowLayoutData{
			Position: RowLayoutPositionCenter,
		}),
		newSimpleWidget(30, 30, &RowLayoutData{
			Position: RowLayoutPositionEnd,
		}),
		newSimpleWidget(40, 40, &RowLayoutData{
			Stretch: true,
		}),
		// TODO: MaxWidth, MaxHeight
	}

	l.Layout(widgets, image.Rect(25, 25, 200, 200))

	is.Equal(widgets[0].GetWidget().Rect, image.Rect(45, 35, 55, 45))
	is.Equal(widgets[1].GetWidget().Rect, image.Rect(62, 87, 82, 107))
	is.Equal(widgets[2].GetWidget().Rect, image.Rect(89, 130, 119, 160))
	is.Equal(widgets[3].GetWidget().Rect, image.Rect(126, 35, 166, 160))
}

func newRowLayout(t *testing.T, opts ...RowLayoutOpt) Layouter {
	t.Helper()
	l := NewRowLayout(opts...)
	l.(Dirtyable).MarkDirty()
	return l
}
