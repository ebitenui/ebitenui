package widget

import (
	img "image"
	"testing"

	"github.com/matryer/is"
)

// ScrollVisible brings a row to the top, the middle or the bottom of the
// view, and stops at the content's end.
func TestScrollContainer_ScrollVisible(t *testing.T) {
	is := is.New(t)

	content := NewContainer(ContainerOpts.Layout(NewRowLayout(RowLayoutOpts.Direction(DirectionVertical))))
	var rows []*simpleWidget
	for range 10 {
		r := newSimpleWidget(50, 20, nil)
		rows = append(rows, r)
		content.AddChild(r)
	}
	ns := newNineSliceEmpty(t)
	s := NewScrollContainer(ScrollContainerOpts.Content(content), ScrollContainerOpts.Image(&ScrollContainerImage{Idle: ns, Mask: ns}))
	s.SetLocation(img.Rect(0, 0, 50, 50)) // a view of 2.5 rows over 10
	render(s, t)

	for _, c := range []struct {
		row    int
		p      float64
		wantAt int // the row's top, from the view's
	}{
		{5, 0, 0},
		{5, 0.5, 15},
		{5, 1, 30},
		{9, 0, 30}, // the content's end stops it
	} {
		s.ScrollVisible(rows[c.row], c.p)
		is.True(s.ScrollTop >= 0 && s.ScrollTop <= 1) // a fraction of the scroll
		render(s, t)
		is.Equal(rows[c.row].GetWidget().Rect.Min.Y, c.wantAt) // where the row shows
	}

	s.ScrollVisible(rows[0], 0)
	render(s, t)
	outside := newSimpleWidget(50, 20, nil)
	outside.SetLocation(img.Rect(0, 120, 50, 140))
	s.ScrollVisible(outside, 0) // not in the content
	is.Equal(s.ScrollTop, 0.0)
}
