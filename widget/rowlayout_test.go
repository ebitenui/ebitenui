package widget

import (
	"image"
	"testing"

	"github.com/matryer/is"
)

func TestRowLayout_PreferredSize(t *testing.T) {
	is := is.New(t)

	spacing := 7

	padding := Insets{
		Top:    10,
		Left:   20,
		Right:  30,
		Bottom: 40,
	}

	l := newRowLayout(t,
		RowLayoutOpts.Padding(padding),
		RowLayoutOpts.Spacing(spacing))

	widgets := []PreferredSizeLocateableWidget{
		newSimpleWidget(10, 10, nil),
		newSimpleWidget(20, 20, RowLayoutData{
			Position: RowLayoutPositionCenter,
		}),
		newSimpleWidget(30, 30, RowLayoutData{
			Position: RowLayoutPositionEnd,
		}),
		newSimpleWidget(40, 40, RowLayoutData{
			Stretch: true,
		}),
		newSimpleWidget(200, 200, RowLayoutData{
			MaxHeight: 45,
		}),
		newSimpleWidget(200, 200, RowLayoutData{
			Position:  RowLayoutPositionCenter,
			MaxHeight: 45,
		}),
		newSimpleWidget(200, 200, RowLayoutData{
			Position:  RowLayoutPositionEnd,
			MaxHeight: 45,
		}),
	}

	expectedWidth, expectedHeight := 0, 0
	for i, wi := range widgets {
		s := wi.(*simpleWidget)
		w, h := s.preferredWidth, s.preferredHeight

		expectedWidth += w

		if i > 0 {
			expectedWidth += spacing
		}

		if rld, ok := s.GetWidget().LayoutData.(RowLayoutData); ok {
			if rld.MaxHeight > 0 && h > rld.MaxHeight {
				h = rld.MaxHeight
			}
		}

		if h > expectedHeight {
			expectedHeight = h
		}
	}
	expectedWidth += padding.Dx()
	expectedHeight += padding.Dy()

	w, h := l.PreferredSize(widgets)

	is.Equal(w, expectedWidth)
	is.Equal(h, expectedHeight)
}

func TestRowLayout_Layout(t *testing.T) {
	is := is.New(t)

	l := newRowLayout(t,
		RowLayoutOpts.Padding(Insets{
			Top:    10,
			Left:   20,
			Right:  30,
			Bottom: 40,
		}),
		RowLayoutOpts.Spacing(7))

	widgets := []PreferredSizeLocateableWidget{
		newSimpleWidget(10, 10, nil),
		newSimpleWidget(20, 20, RowLayoutData{
			Position: RowLayoutPositionCenter,
		}),
		newSimpleWidget(30, 30, RowLayoutData{
			Position: RowLayoutPositionEnd,
		}),
		newSimpleWidget(40, 40, RowLayoutData{
			Stretch: true,
		}),
		newSimpleWidget(200, 200, RowLayoutData{
			MaxHeight: 45,
		}),
		newSimpleWidget(200, 200, RowLayoutData{
			Position:  RowLayoutPositionCenter,
			MaxHeight: 45,
		}),
		newSimpleWidget(200, 200, RowLayoutData{
			Position:  RowLayoutPositionEnd,
			MaxHeight: 45,
		}),
	}

	l.Layout(widgets, image.Rect(25, 25, 200, 200))

	expected := []image.Rectangle{
		image.Rect(45, 35, 55, 45),
		image.Rect(62, 87, 82, 107),
		image.Rect(89, 130, 119, 160),
		image.Rect(126, 35, 166, 160),
		image.Rect(173, 35, 373, 80),
		image.Rect(380, 75, 580, 120),
		image.Rect(587, 115, 787, 160),
	}

	for i, r := range expected {
		is.Equal(widgets[i].GetWidget().Rect, r)
	}
}

func newRowLayout(t *testing.T, opts ...RowLayoutOpt) Layouter {
	t.Helper()
	l := NewRowLayout(opts...)
	return l
}
