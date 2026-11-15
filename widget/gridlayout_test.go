package widget

import (
	"image"
	"testing"

	"github.com/matryer/is"
)

func TestGridLayout_PositionsInCell(t *testing.T) {
	is := is.New(t)

	l := NewGridLayout(GridLayoutOpts.Columns(2), GridLayoutOpts.Spacing(10, 10), GridLayoutOpts.DefaultStretch(true, true))
	at := func(h, v GridLayoutPosition) *simpleWidget {
		return newSimpleWidget(10, 10, GridLayoutData{MaxWidth: 10, MaxHeight: 10, HorizontalPosition: h, VerticalPosition: v})
	}
	// Cells are 50x30; the second column starts at 60, the second row at 40.
	cells := []PreferredSizeLocateableWidget{
		newSimpleWidget(50, 30, nil),
		at(GridLayoutPositionCenter, GridLayoutPositionCenter),
		at(GridLayoutPositionEnd, GridLayoutPositionEnd),
		at(GridLayoutPositionStart, GridLayoutPositionStart),
	}

	l.Layout(cells, image.Rect(0, 0, 110, 70))

	is.Equal(cells[0].GetWidget().Rect, image.Rect(0, 0, 50, 30))
	is.Equal(cells[1].GetWidget().Rect, image.Rect(80, 10, 90, 20)) // centered
	is.Equal(cells[2].GetWidget().Rect, image.Rect(40, 60, 50, 70)) // at the end
	is.Equal(cells[3].GetWidget().Rect, image.Rect(60, 40, 70, 50)) // at the start
}

// Stretched columns and rows fill the layout; cells in a row share its
// height, and cells in a column its width.
func TestGridLayout_StretchedCellsFill(t *testing.T) {
	is := is.New(t)

	l := NewGridLayout(GridLayoutOpts.Columns(3), GridLayoutOpts.DefaultStretch(true, true))
	cells := make([]PreferredSizeLocateableWidget, 9)
	for i := range cells {
		cells[i] = newSimpleWidget(10, 10, nil)
	}

	l.Layout(cells, image.Rect(0, 0, 100, 100))

	// 100 does not divide by 3: the first column and row take the extra pixel.
	edges := []int{0, 34, 67, 100}
	for i, w := range cells {
		c, r := i%3, i/3
		is.Equal(w.GetWidget().Rect, image.Rect(edges[c], edges[r], edges[c+1], edges[r+1]))
	}
}

// A row ends where its cells do, though its last widget is hidden.
func TestGridLayout_RowAfterHiddenWidget(t *testing.T) {
	is := is.New(t)

	l := NewGridLayout(GridLayoutOpts.Columns(2), GridLayoutOpts.DefaultStretch(false, true))
	cells := make([]PreferredSizeLocateableWidget, 4)
	for i := range cells {
		cells[i] = newSimpleWidget(10, 10, nil)
	}
	cells[1].GetWidget().SetVisibility(Visibility_Hide)

	l.Layout(cells, image.Rect(0, 0, 100, 100))

	is.Equal(cells[0].GetWidget().Rect, image.Rect(0, 0, 10, 50))
	is.Equal(cells[2].GetWidget().Rect, image.Rect(0, 50, 10, 100))
}

// Stretched columns share what the others leave, the first taking the
// extra pixel.
func TestGridLayout_StretchedColumnsShareTheRest(t *testing.T) {
	is := is.New(t)

	l := NewGridLayout(GridLayoutOpts.Columns(3), GridLayoutOpts.Spacing(5, 0), GridLayoutOpts.Stretch([]bool{false, true, true}, nil))
	cells := []PreferredSizeLocateableWidget{
		newSimpleWidget(20, 10, nil),
		newSimpleWidget(10, 10, nil),
		newSimpleWidget(10, 10, nil),
	}

	l.Layout(cells, image.Rect(0, 0, 101, 50))

	is.Equal(cells[0].GetWidget().Rect, image.Rect(0, 0, 20, 10))
	is.Equal(cells[1].GetWidget().Rect, image.Rect(25, 0, 61, 10))
	is.Equal(cells[2].GetWidget().Rect, image.Rect(66, 0, 101, 10))
}

// A grid layout measures and lays out without allocating once it has
// measured, and a layout shared by nested containers measures each right.
func TestGridLayout_MeasuresWithoutAllocating(t *testing.T) {
	is := is.New(t)

	l := NewGridLayout(GridLayoutOpts.Columns(2))
	cells := []PreferredSizeLocateableWidget{
		newSimpleWidget(20, 10, nil),
		newSimpleWidget(10, 10, nil),
		newSimpleWidget(10, 30, nil),
	}
	w, h := l.PreferredSize(cells)
	is.Equal(w, 30)
	is.Equal(h, 40)
	is.Equal(testing.AllocsPerRun(10, func() { l.PreferredSize(cells) }), 0.0)
	is.Equal(testing.AllocsPerRun(10, func() { l.Layout(cells, image.Rect(0, 0, 30, 40)) }), 0.0)

	inner := NewContainer(ContainerOpts.Layout(l))
	inner.AddChild(newSimpleWidget(10, 10, nil), newSimpleWidget(10, 10, nil))
	outer := NewContainer(ContainerOpts.Layout(l))
	outer.AddChild(newSimpleWidget(40, 10, nil), inner)
	w, h = outer.PreferredSize()
	is.Equal(w, 60) // 40, and the inner container's 20 measured with the same layout
	is.Equal(h, 10)
}
