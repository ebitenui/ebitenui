package widget

import (
	"image"
	"math"
)

// GridLayout layouts widgets in a grid fashion, with columns or rows optionally being stretched.
//
// Widget.LayoutData of widgets being layouted by GridLayout need to be of type GridLayoutData.
type GridLayout struct {
	columns       int
	padding       Insets
	columnSpacing int
	rowSpacing    int
	columnStretch []bool
	rowStretch    []bool
}

// GridLayoutOpt is a function that configures g.
type GridLayoutOpt func(g *GridLayout)

type GridLayoutOptions struct {
}

// GridLayoutData specifies layout settings for a widget.
type GridLayoutData struct {
	// MaxWidth specifies the maximum width.
	MaxWidth int

	// MaxHeight specifies the maximum height..
	MaxHeight int

	// HorizontalPosition specifies the horizontal anchoring position inside the grid cell.
	HorizontalPosition GridLayoutPosition

	// VerticalPosition specifies the vertical anchoring position inside the grid cell.
	VerticalPosition GridLayoutPosition
}

// GridLayoutPosition is the type used to specify an anchoring position.
type GridLayoutPosition int

const (
	// GridLayoutPositionStart is the anchoring position for "left" (in the horizontal direction) or "top" (in the vertical direction.)
	GridLayoutPositionStart = GridLayoutPosition(iota)

	// GridLayoutPositionStart is the center anchoring position.
	GridLayoutPositionCenter

	// GridLayoutPositionEnd is the anchoring position for "right" (in the horizontal direction) or "bottom" (in the vertical direction.)
	GridLayoutPositionEnd
)

// GridLayoutOpts contains functions that configure a GridLayout.
var GridLayoutOpts GridLayoutOptions

// NewGridLayout constructs a new GridLayout, configured by opts.
func NewGridLayout(opts ...GridLayoutOpt) *GridLayout {
	g := &GridLayout{}

	for _, o := range opts {
		o(g)
	}

	return g
}

// Columns configures a grid layout to use c columns.
func (o GridLayoutOptions) Columns(c int) GridLayoutOpt {
	return func(g *GridLayout) {
		g.columns = c
	}
}

// Padding configures a grid layout to use padding i.
func (o GridLayoutOptions) Padding(i Insets) GridLayoutOpt {
	return func(g *GridLayout) {
		g.padding = i
	}
}

// Spacing configures a grid layout to separate columns by spacing c and rows by spacing r.
func (o GridLayoutOptions) Spacing(c int, r int) GridLayoutOpt {
	return func(g *GridLayout) {
		g.columnSpacing = c
		g.rowSpacing = r
	}
}

// Stretch configures a grid layout to stretch columns according to c and rows according to r.
// The number of elements of c and r must correspond with the number of columns and rows in the
// layout.
func (o GridLayoutOptions) Stretch(c []bool, r []bool) GridLayoutOpt {
	return func(g *GridLayout) {
		g.columnStretch = c
		g.rowStretch = r
	}
}

// PreferredSize implements Layouter.
func (g *GridLayout) PreferredSize(widgets []PreferredSizeLocateableWidget) (int, int) {
	colWidths, rowHeights := g.preferredColumnWidthsAndRowHeights(widgets)
	return g.padding.Dx() + g.columnSpacing*(len(colWidths)-1) + sumInts(colWidths),
		g.padding.Dy() + g.rowSpacing*(len(rowHeights)-1) + sumInts(rowHeights)
}

// Layout implements Layouter.
func (g *GridLayout) CalcLayout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) []image.Rectangle {
	res := make([]image.Rectangle, 0)
	g.doLayout(widgets, rect, func(pslw PreferredSizeLocateableWidget, r image.Rectangle) {
		res = append(res, r)
	})
	return res
}
func (g *GridLayout) Layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) {
	g.doLayout(widgets, rect, func(pslw PreferredSizeLocateableWidget, r image.Rectangle) {
		pslw.SetLocation(r)
	})
}
func (g *GridLayout) doLayout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle, locater LocationFunction) {
	rect = g.padding.Apply(rect)

	colWidths, rowHeights := g.preferredColumnWidthsAndRowHeights(widgets)
	stretchedColWidth, stretchedRowHeight, firstStretchedColWidth, firstStretchedRowHeight := g.stretchedCellSizes(colWidths, rowHeights, rect)

	c, r := 0, 0
	x, y := 0, 0
	firstStretchedCol, firstStretchedRow := true, true
	for _, w := range widgets {
		if w.GetWidget().Visibility == Visibility_Hide {
			continue
		}

		cw := colWidths[c]
		if g.columnStretched(c) {
			cw = stretchedColWidth
			if firstStretchedCol {
				cw = firstStretchedColWidth
				firstStretchedCol = false
			}
		}

		ch := rowHeights[r]
		if g.rowStretched(r) {
			ch = stretchedRowHeight
			if firstStretchedRow {
				ch = firstStretchedRowHeight
				firstStretchedRow = false
			}
		}

		ww, wh := cw, ch
		wx, wy := x, y

		ld := w.GetWidget().LayoutData
		if gld, ok := ld.(GridLayoutData); ok {
			wx, wy, ww, wh = g.applyLayoutData(gld, wx, wy, ww, wh, x, y, cw, ch)
		}

		locater(w, image.Rect(rect.Min.X+wx, rect.Min.Y+wy, rect.Min.X+wx+ww, rect.Min.Y+wy+wh))

		c++
		x += cw + g.columnSpacing

		if c >= g.columns {
			c = 0
			r++
			x = 0
			y += ch + g.rowSpacing
			firstStretchedCol = true
		}
	}
}

func (g *GridLayout) stretchedCellSizes(colWidths []int, rowHeights []int, rect image.Rectangle) (int, int, int, int) {
	stretchedColWidth, stretchedRowHeight := 0, 0

	remainingWidth := rect.Dx() - g.columnSpacing*(len(colWidths)-1)
	remainingHeight := rect.Dy() - g.rowSpacing*(len(rowHeights)-1)
	stretchedCols, stretchedRows := 0, 0

	for c, cw := range colWidths {
		if g.columnStretched(c) {
			stretchedCols++
		} else {
			remainingWidth -= cw
		}
	}

	for r, rh := range rowHeights {
		if g.rowStretched(r) {
			stretchedRows++
		} else {
			remainingHeight -= rh
		}
	}

	if stretchedCols > 0 {
		stretchedColWidth = int(math.Floor(float64(remainingWidth) / float64(stretchedCols)))
	}
	if stretchedRows > 0 {
		stretchedRowHeight = int(math.Floor(float64(remainingHeight) / float64(stretchedRows)))
	}

	firstStretchedColWidth := stretchedColWidth + (remainingWidth - stretchedColWidth*stretchedCols)
	firstStretchedRowHeight := stretchedRowHeight + (remainingHeight - stretchedRowHeight*stretchedRows)

	return stretchedColWidth, stretchedRowHeight, firstStretchedColWidth, firstStretchedRowHeight
}

func (g *GridLayout) columnStretched(c int) bool {
	return g.columnStretch != nil && c < len(g.columnStretch) && g.columnStretch[c]
}

func (g *GridLayout) rowStretched(r int) bool {
	return g.rowStretch != nil && r < len(g.rowStretch) && g.rowStretch[r]
}

func (g *GridLayout) preferredColumnWidthsAndRowHeights(widgets []PreferredSizeLocateableWidget) ([]int, []int) {
	colWidths := make([]int, g.columns)
	rowHeights := make([]int, int(math.Ceil(float64(len(widgets))/float64(g.columns))))

	c := 0
	r := 0
	for _, w := range widgets {
		ww, wh := w.PreferredSize()

		ld := w.GetWidget().LayoutData
		if gld, ok := ld.(GridLayoutData); ok {
			ww, wh = g.applyMaxSize(gld, ww, wh)
		}

		if ww > colWidths[c] {
			colWidths[c] = ww
		}

		if wh > rowHeights[r] {
			rowHeights[r] = wh
		}

		c++
		if c >= g.columns {
			c = 0
			r++
		}
	}

	return colWidths, rowHeights
}

func (g *GridLayout) applyLayoutData(ld GridLayoutData, wx int, wy int, ww int, wh int, x int, y int, cw int, ch int) (int, int, int, int) {
	if ld.MaxWidth > 0 && ww > ld.MaxWidth {
		ww = ld.MaxWidth
	}

	if ld.MaxHeight > 0 && wh > ld.MaxHeight {
		wh = ld.MaxHeight
	}

	switch ld.HorizontalPosition {
	case GridLayoutPositionCenter:
		wx = x + (cw-ww)/2
	case GridLayoutPositionEnd:
		wx = x + cw - ww
	}

	switch ld.VerticalPosition {
	case GridLayoutPositionCenter:
		wy = y + (ch-wh)/2
	case GridLayoutPositionEnd:
		wy = y + ch - wh
	}

	return wx, wy, ww, wh
}

func (g *GridLayout) applyMaxSize(ld GridLayoutData, ww int, wh int) (int, int) {
	if ld.MaxWidth > 0 && ww > ld.MaxWidth {
		ww = ld.MaxWidth
	}

	if ld.MaxHeight > 0 && wh > ld.MaxHeight {
		wh = ld.MaxHeight
	}

	return ww, wh
}

func sumInts(ints []int) int {
	s := 0
	for _, i := range ints {
		s += i
	}
	return s
}
