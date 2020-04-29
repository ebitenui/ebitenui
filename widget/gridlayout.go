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
	dirty         bool
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

func (o gridLayoutOpts) WithColumns(c int) GridLayoutOpt {
	return func(g *gridLayout) {
		g.columns = c
	}
}

func (o gridLayoutOpts) WithPadding(p Insets) GridLayoutOpt {
	return func(g *gridLayout) {
		g.padding = p
	}
}

func (o gridLayoutOpts) WithSpacing(columnSpacing int, rowSpacing int) GridLayoutOpt {
	return func(g *gridLayout) {
		g.columnSpacing = columnSpacing
		g.rowSpacing = rowSpacing
	}
}

func (o gridLayoutOpts) WithStretch(columns []bool, rows []bool) GridLayoutOpt {
	return func(g *gridLayout) {
		g.columnStretch = columns
		g.rowStretch = rows
	}
}

func (g *gridLayout) PreferredSize(widgets []HasWidget) (int, int) {
	colWidths, rowHeights := g.preferredColumnWidthsAndRowHeights(widgets)
	return g.padding.Dx() + g.columnSpacing*(len(colWidths)-1) + sumInts(colWidths),
		g.padding.Dy() + g.rowSpacing*(len(rowHeights)-1) + sumInts(rowHeights)
}

func (g *gridLayout) Layout(widgets []HasWidget, rect image.Rectangle) {
	if !g.dirty {
		return
	}

	defer func() {
		g.dirty = false
	}()

	rect = g.padding.Apply(rect)

	colWidths, rowHeights := g.preferredColumnWidthsAndRowHeights(widgets)
	remainingWidth := rect.Dx() - g.columnSpacing*(len(colWidths)-1)
	remainingHeight := rect.Dy() - g.rowSpacing*(len(rowHeights)-1)

	// part 1: calculate stretched column widths/row heights

	stretchedCols, stretchedRows := 0, 0

	if g.columnStretch != nil {
		for c, cw := range colWidths {
			if g.columnStretch[c] {
				stretchedCols++
			} else {
				remainingWidth -= cw
			}
		}
	}

	if g.rowStretch != nil {
		for r, rh := range rowHeights {
			if g.rowStretch[r] {
				stretchedRows++
			} else {
				remainingHeight -= rh
			}
		}
	}

	stretchedColWidth, stretchedRowHeight := 0, 0
	if stretchedCols > 0 {
		stretchedColWidth = int(math.Floor(float64(remainingWidth) / float64(stretchedCols)))
	}
	if stretchedRows > 0 {
		stretchedRowHeight = int(math.Floor(float64(remainingHeight) / float64(stretchedRows)))
	}

	firstStretchedColWidth := stretchedColWidth + (remainingWidth - stretchedColWidth*stretchedCols)
	firstStretchedRowHeight := stretchedRowHeight + (remainingHeight - stretchedRowHeight*stretchedRows)

	// part 2: layout

	c, r := 0, 0
	x, y := 0, 0
	firstStretchedCol, firstStretchedRow := true, true
	for _, w := range widgets {
		var cw int
		var ch int

		if g.columnStretch != nil && g.columnStretch[c] {
			if firstStretchedCol {
				cw = firstStretchedColWidth
				firstStretchedCol = false
			} else {
				cw = stretchedColWidth
			}
		} else {
			cw = colWidths[c]
		}

		if g.rowStretch != nil && g.rowStretch[r] {
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
			if gld.MaxWidth > 0 && ww > gld.MaxWidth {
				ww = gld.MaxWidth
			}

			if gld.MaxHeight > 0 && wh > gld.MaxHeight {
				wh = gld.MaxHeight
			}

			switch gld.HorizontalPosition {
			case GridLayoutPositionCenter:
				wx = x + (cw-ww)/2
			case GridLayoutPositionEnd:
				wx = x + cw - ww
			}

			switch gld.VerticalPosition {
			case GridLayoutPositionCenter:
				wy = x + (ch-wh)/2
			case GridLayoutPositionEnd:
				wy = y + ch - wh
			}
		}

		if l, ok := w.(Locateable); ok {
			l.SetLocation(image.Rect(rect.Min.X+wx, rect.Min.Y+wy, rect.Min.X+wx+ww, rect.Min.Y+wy+wh))
		}

		c++
		x += cw + g.columnSpacing

		if c >= g.columns {
			c = 0
			r++
			x = 0
			y += ch + g.rowSpacing
		}
	}
}

func (g *gridLayout) preferredColumnWidthsAndRowHeights(widgets []HasWidget) ([]int, []int) {
	colWidths := make([]int, g.columns)
	rowHeights := make([]int, int(math.Ceil(float64(len(widgets))/float64(g.columns))))

	c := 0
	r := 0
	for _, w := range widgets {
		var ww int
		var wh int
		if p, ok := w.(PreferredSizer); ok {
			ww, wh = p.PreferredSize()
		} else {
			ww, wh = 50, 50
		}

		ld := w.GetWidget().LayoutData
		if gld, ok := ld.(*GridLayoutData); ok {
			if gld.MaxWidth > 0 && ww > gld.MaxWidth {
				ww = gld.MaxWidth
			}

			if gld.MaxHeight > 0 && wh > gld.MaxHeight {
				wh = gld.MaxHeight
			}
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

func sumInts(ints []int) int {
	s := 0
	for _, i := range ints {
		s += i
	}
	return s
}

func (g *gridLayout) MarkDirty() {
	g.dirty = true
}
