package widget

import "image"

type rowLayout struct {
	direction Direction
	padding   Insets
	spacing   int
	dirty     bool
}

type RowLayoutOpt func(f *rowLayout)

type RowLayoutData struct {
	Position  RowLayoutPosition
	Stretch   bool
	MaxWidth  int
	MaxHeight int
}

type RowLayoutPosition int

const (
	RowLayoutPositionStart = RowLayoutPosition(iota)
	RowLayoutPositionCenter
	RowLayoutPositionEnd
)

const RowLayoutOpts = rowLayoutOpts(true)

type rowLayoutOpts bool

func NewRowLayout(opts ...RowLayoutOpt) Layouter {
	r := &rowLayout{}

	for _, o := range opts {
		o(r)
	}

	return r
}

func (o rowLayoutOpts) Direction(d Direction) RowLayoutOpt {
	return func(r *rowLayout) {
		r.direction = d
	}
}

func (o rowLayoutOpts) Padding(i Insets) RowLayoutOpt {
	return func(r *rowLayout) {
		r.padding = i
	}
}

func (o rowLayoutOpts) Spacing(s int) RowLayoutOpt {
	return func(f *rowLayout) {
		f.spacing = s
	}
}

func (f *rowLayout) PreferredSize(widgets []HasWidget) (int, int) {
	r := image.Rectangle{}
	f.layout(widgets, image.Rectangle{}, false, func(w HasWidget, wr image.Rectangle) {
		r = r.Union(wr)
	})
	return r.Dx() + f.padding.Dx(), r.Dy() + f.padding.Dy()
}

func (r *rowLayout) Layout(widgets []HasWidget, rect image.Rectangle) {
	if !r.dirty {
		return
	}

	defer func() {
		r.dirty = false
	}()

	r.layout(widgets, rect, true, func(w HasWidget, wr image.Rectangle) {
		w.(Locateable).SetLocation(wr)
	})
}

func (r *rowLayout) layout(widgets []HasWidget, rect image.Rectangle, usePosition bool, locationFunc func(w HasWidget, wr image.Rectangle)) {
	if len(widgets) == 0 {
		return
	}

	rect = r.padding.Apply(rect)
	x, y := 0, 0

	for _, widget := range widgets {
		wx, wy := x, y

		var ww int
		var wh int
		if p, ok := widget.(PreferredSizer); ok {
			ww, wh = p.PreferredSize()
		} else {
			ww, wh = 50, 50
		}

		ld := widget.GetWidget().LayoutData
		if rld, ok := ld.(*RowLayoutData); ok {
			if usePosition && rld.Stretch {
				if r.direction == DirectionHorizontal {
					wh = rect.Dy()
				} else {
					ww = rect.Dx()
				}
			}

			if rld.MaxWidth > 0 && ww > rld.MaxWidth {
				ww = rld.MaxWidth
			}

			if rld.MaxHeight > 0 && wh > rld.MaxHeight {
				wh = rld.MaxHeight
			}

			if usePosition {
				switch rld.Position {
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
			}
		}

		if _, ok := widget.(Locateable); ok {
			locationFunc(widget, image.Rect(rect.Min.X+wx, rect.Min.Y+wy, rect.Min.X+wx+ww, rect.Min.Y+wy+wh))
		}

		if r.direction == DirectionHorizontal {
			x += ww + r.spacing
		} else {
			y += wh + r.spacing
		}
	}
}

func (r *rowLayout) MarkDirty() {
	r.dirty = true
}
