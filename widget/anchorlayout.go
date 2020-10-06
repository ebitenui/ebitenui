package widget

import "image"

type anchorLayout struct {
	padding Insets
}

type AnchorLayoutOpt func(a *anchorLayout)

const AnchorLayoutOpts = anchorLayoutOpts(true)

type anchorLayoutOpts bool

type AnchorLayoutPosition int

type AnchorLayoutData struct {
	HorizontalPosition AnchorLayoutPosition
	VerticalPosition   AnchorLayoutPosition
	StretchHorizontal  bool
	StretchVertical    bool
}

const (
	AnchorLayoutPositionStart = AnchorLayoutPosition(iota)
	AnchorLayoutPositionCenter
	AnchorLayoutPositionEnd
)

func NewAnchorLayout(opts ...AnchorLayoutOpt) Layouter {
	a := &anchorLayout{}

	for _, o := range opts {
		o(a)
	}

	return a
}

func (o anchorLayoutOpts) Padding(i Insets) AnchorLayoutOpt {
	return func(a *anchorLayout) {
		a.padding = i
	}
}

func (a *anchorLayout) PreferredSize(widgets []PreferredSizeLocateableWidget) (int, int) {
	px, py := a.padding.Dx(), a.padding.Dy()

	if len(widgets) == 0 {
		return px, py
	}

	w, h := widgets[0].PreferredSize()
	return w + px, h + py
}

func (a *anchorLayout) Layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) {
	if len(widgets) == 0 {
		return
	}

	widget := widgets[0]
	ww, wh := widget.PreferredSize()
	rect = a.padding.Apply(rect)
	wx := 0
	wy := 0

	if ald, ok := widget.GetWidget().LayoutData.(*AnchorLayoutData); ok {
		wx, wy, ww, wh = a.applyLayoutData(ald, wx, wy, ww, wh, rect)
	}

	r := image.Rect(0, 0, ww, wh)
	r = r.Add(image.Point{wx, wy})
	r = r.Add(rect.Min)

	widget.SetLocation(r)
}

func (a *anchorLayout) applyLayoutData(ld *AnchorLayoutData, wx int, wy int, ww int, wh int, rect image.Rectangle) (int, int, int, int) {
	if ld.StretchHorizontal {
		ww = rect.Dx()
	}

	if ld.StretchVertical {
		wh = rect.Dy()
	}

	hPos := ld.HorizontalPosition
	vPos := ld.VerticalPosition

	switch hPos {
	case AnchorLayoutPositionCenter:
		wx = (rect.Dx() - ww) / 2
	case AnchorLayoutPositionEnd:
		wx = rect.Dx() - ww
	}

	switch vPos {
	case AnchorLayoutPositionCenter:
		wy = (rect.Dy() - wh) / 2
	case AnchorLayoutPositionEnd:
		wy = rect.Dy() - wh
	}

	return wx, wy, ww, wh
}
