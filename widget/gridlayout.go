package widget

import (
	"image"
	"math"
)

type gridLayout struct {
	columns       int
	padding       Insets
	columnSpacing int
	rowSpacing    int
	columnStretch []bool
	rowStretch    []bool
}

type GridLayoutOpt func(g *gridLayout)

type GridLayoutData struct {
	MaxWidth           int
	MaxHeight          int
	HorizontalPosition GridLayoutPosition
	VerticalPosition   GridLayoutPosition
}

type GridLayoutPosition int

const (
	GridLayoutPositionStart = GridLayoutPosition(iota)
	GridLayoutPositionCenter
	GridLayoutPositionEnd
)

const GridLayoutOpts = gridLayoutOpts(true)

type gridLayoutOpts bool

func NewGridLayout(opts ...GridLayoutOpt) Layouter {
	g := &gridLayout{}

	for _, o := range opts {
		o(g)
	}

	return g
}

func (o gridLayoutOpts) Columns(c int) GridLayoutOpt {
	return func(g *gridLayout) {
		g.columns = c
	}
}

func (o gridLayoutOpts) Padding(p Insets) GridLayoutOpt {
	return func(g *gridLayout) {
		g.padding = p
	}
}

func (o gridLayoutOpts) Spacing(columnSpacing int, rowSpacing int) GridLayoutOpt {
	return func(g *gridLayout) {
		g.columnSpacing = columnSpacing
		g.rowSpacing = rowSpacing
	}
}

func (o gridLayoutOpts) Stretch(columns []bool, rows []bool) GridLayoutOpt {
	return func(g *gridLayout) {
		g.columnStretch = columns
		g.rowStretch = rows
	}
}

func (g *gridLayout) PreferredSize(widgets []PreferredSizeLocateableWidget) (int, int) {
	colWidths, rowHeights := g.preferredColumnWidthsAndRowHeights(widgets)
	return g.padding.Dx() + g.columnSpacing*(len(colWidths)-1) + sumInts(colWidths),
		g.padding.Dy() + g.rowSpacing*(len(rowHeights)-1) + sumInts(rowHeights)
}

func (g *gridLayout) Layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) {
	rect = g.padding.Apply(rect)

	colWidths, rowHeights := g.preferredColumnWidthsAndRowHeights(widgets)

	stretchedColWidth, stretchedRowHeight, firstStretchedColWidth, firstStretchedRowHeight := g.stretchedCellSizes(colWidths, rowHeights, rect)

	// part 2: layout

	c, r := 0, 0
	x, y := 0, 0
	firstStretchedCol, firstStretchedRow := true, true
	for _, w := range widgets {
		var cw int
		var ch int

		if g.columnStretched(c) {
			if firstStretchedCol {
				cw = firstStretchedColWidth
				firstStretchedCol = false
			} else {
				cw = stretchedColWidth
			}
		} else {
			cw = colWidths[c]
		}

		if g.rowStretched(r) {
			if firstStretchedRow {
				ch = firstStretchedRowHeight
				firstStretchedRow = false
			} else {
				ch = stretchedRowHeight
			}
		} else {
			ch = rowHeights[r]
		}

		ww, wh := cw, ch
		wx, wy := x, y

		ld := w.GetWidget().LayoutData
		if gld, ok := ld.(*GridLayoutData); ok {
			wx, wy, ww, wh = g.applyLayoutData(gld, wx, wy, ww, wh, x, y, cw, ch)
		}

		w.SetLocation(image.Rect(rect.Min.X+wx, rect.Min.Y+wy, rect.Min.X+wx+ww, rect.Min.Y+wy+wh))

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

func (g *gridLayout) stretchedCellSizes(colWidths []int, rowHeights []int, rect image.Rectangle) (int, int, int, int) {
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

func (g *gridLayout) columnStretched(c int) bool {
	return g.columnStretch != nil && g.columnStretch[c]
}

func (g *gridLayout) rowStretched(r int) bool {
	return g.rowStretch != nil && g.rowStretch[r]
}

func (g *gridLayout) preferredColumnWidthsAndRowHeights(widgets []PreferredSizeLocateableWidget) ([]int, []int) {
	colWidths := make([]int, g.columns)
	rowHeights := make([]int, int(math.Ceil(float64(len(widgets))/float64(g.columns))))

	c := 0
	r := 0
	for _, w := range widgets {
		ww, wh := w.PreferredSize()

		ld := w.GetWidget().LayoutData
		if gld, ok := ld.(*GridLayoutData); ok {
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

func (g *gridLayout) applyLayoutData(ld *GridLayoutData, wx int, wy int, ww int, wh int, x int, y int, cw int, ch int) (int, int, int, int) {
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
		wy = x + (ch-wh)/2
	case GridLayoutPositionEnd:
		wy = y + ch - wh
	}

	return wx, wy, ww, wh
}

func (g *gridLayout) applyMaxSize(ld *GridLayoutData, ww int, wh int) (int, int) {
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
