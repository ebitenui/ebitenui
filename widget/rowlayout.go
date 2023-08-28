package widget

import "image"

// RowLayout layouts widgets in either a single row or a single column,
// optionally stretching them in the other direction.
//
// Widget.LayoutData of widgets being layouted by RowLayout need to be of type RowLayoutData.
type RowLayout struct {
	direction Direction
	padding   Insets
	spacing   int
}

type RowLayoutOptions struct {
}

// RowLayoutOpt is a function that configures r.
type RowLayoutOpt func(r *RowLayout)

// RowLayoutData specifies layout settings for a widget.
type RowLayoutData struct {
	// Position specifies the anchoring position for the direction that is not the primary direction of the layout.
	Position RowLayoutPosition

	// Stretch specifies whether to stretch in the direction that is not the primary direction of the layout.
	Stretch bool

	// MaxWidth specifies the maximum width.
	MaxWidth int

	// MaxHeight specifies the maximum height.
	MaxHeight int
}

// RowLayoutPosition is the type used to specify an anchoring position.
type RowLayoutPosition int

const (
	// RowLayoutPositionStart is the anchoring position for "left" (in the horizontal direction) or "top" (in the vertical direction.)
	RowLayoutPositionStart = RowLayoutPosition(iota)

	// RowLayoutPositionCenter is the center anchoring position.
	RowLayoutPositionCenter

	// RowLayoutPositionEnd is the anchoring position for "right" (in the horizontal direction) or "bottom" (in the vertical direction.)
	RowLayoutPositionEnd
)

// RowLayoutOpts contains functions that configure a RowLayout.
var RowLayoutOpts RowLayoutOptions

// NewRowLayout constructs a new RowLayout, configured by opts.
func NewRowLayout(opts ...RowLayoutOpt) *RowLayout {
	r := &RowLayout{}

	for _, o := range opts {
		o(r)
	}

	return r
}

// Direction configures a row layout to layout widgets in the primary direction d. This will also switch the meaning
// of any widget's RowLayoutData.Position and RowLayoutData.Stretch to the other direction.
func (o RowLayoutOptions) Direction(d Direction) RowLayoutOpt {
	return func(r *RowLayout) {
		r.direction = d
	}
}

// Padding configures a row layout to use padding i.
func (o RowLayoutOptions) Padding(i Insets) RowLayoutOpt {
	return func(r *RowLayout) {
		r.padding = i
	}
}

// Spacing configures a row layout to separate widgets by spacing s.
func (o RowLayoutOptions) Spacing(s int) RowLayoutOpt {
	return func(f *RowLayout) {
		f.spacing = s
	}
}

// PreferredSize implements Layouter.
func (r *RowLayout) PreferredSize(widgets []PreferredSizeLocateableWidget) (int, int) {
	rect := image.Rectangle{}
	r.layout(widgets, image.Rectangle{}, false, func(_ PreferredSizeLocateableWidget, wr image.Rectangle) {
		rect = rect.Union(wr)
	})
	return rect.Dx() + r.padding.Dx(), rect.Dy() + r.padding.Dy()
}

// Layout implements Layouter.
func (r *RowLayout) Layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) {
	r.layout(widgets, rect, true, func(w PreferredSizeLocateableWidget, wr image.Rectangle) {
		w.SetLocation(wr)
	})
}

func (r *RowLayout) layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle, usePosition bool, locationFunc func(w PreferredSizeLocateableWidget, wr image.Rectangle)) {
	if len(widgets) == 0 {
		return
	}

	rect = r.padding.Apply(rect)
	x, y := 0, 0

	for _, widget := range widgets {
		if widget.GetWidget().Visibility == Visibility_Hide {
			continue
		}

		wx, wy := x, y
		ww, wh := widget.PreferredSize()

		ld := widget.GetWidget().LayoutData
		if rld, ok := ld.(RowLayoutData); ok {
			wx, wy, ww, wh = r.applyLayoutData(rld, wx, wy, ww, wh, usePosition, rect, x, y)
		}

		wr := image.Rect(0, 0, ww, wh)
		wr = wr.Add(rect.Min)
		wr = wr.Add(image.Point{wx, wy})
		locationFunc(widget, wr)

		if r.direction == DirectionHorizontal {
			x += ww + r.spacing
		} else {
			y += wh + r.spacing
		}
	}
}

func (r *RowLayout) applyLayoutData(ld RowLayoutData, wx int, wy int, ww int, wh int, usePosition bool, rect image.Rectangle, x int, y int) (int, int, int, int) {
	if usePosition {
		ww, wh = r.applyStretch(ld, ww, wh, rect)
	}

	ww, wh = r.applyMaxSize(ld, ww, wh)

	if usePosition {
		wx, wy = r.applyPosition(ld, wx, wy, ww, wh, rect, x, y)
	}

	return wx, wy, ww, wh
}

func (r *RowLayout) applyStretch(ld RowLayoutData, ww int, wh int, rect image.Rectangle) (int, int) {
	if !ld.Stretch {
		return ww, wh
	}

	if r.direction == DirectionHorizontal {
		wh = rect.Dy()
	} else {
		ww = rect.Dx()
	}

	return ww, wh
}

func (r *RowLayout) applyMaxSize(ld RowLayoutData, ww int, wh int) (int, int) {
	if ld.MaxWidth > 0 && ww > ld.MaxWidth {
		ww = ld.MaxWidth
	}

	if ld.MaxHeight > 0 && wh > ld.MaxHeight {
		wh = ld.MaxHeight
	}

	return ww, wh
}

func (r *RowLayout) applyPosition(ld RowLayoutData, wx int, wy int, ww int, wh int, rect image.Rectangle, x int, y int) (int, int) {
	switch ld.Position {
	case RowLayoutPositionCenter:
		if r.direction == DirectionHorizontal {
			wy = y + (rect.Dy()-wh)/2
		} else {
			wx = x + (rect.Dx()-ww)/2
		}

	case RowLayoutPositionEnd:
		if r.direction == DirectionHorizontal {
			wy = y + rect.Dy() - wh
		} else {
			wx = x + rect.Dx() - ww
		}
	}

	return wx, wy
}
