package widget

import "image"

// AnchorLayout layouts a single widget anchored to either a corner or edge of a rectangle,
// optionally stretching it in one or both directions.
//
// AnchorLayout will only layout the first widget in a container and ignore all other widgets.
//
// Widget.LayoutData of widgets being layouted by AnchorLayout need to be of type AnchorLayoutData.
type AnchorLayout struct {
	padding Insets
}

// AnchorLayoutOpt is a function that configures a.
type AnchorLayoutOpt func(a *AnchorLayout)

type AnchorLayoutOptions struct {
}

// AnchorLayoutPosition is the type used to specify an anchoring position.
type AnchorLayoutPosition int

// AnchorLayoutData specifies layout settings for a widget.
type AnchorLayoutData struct {
	// HorizontalPosition specifies the horizontal anchoring position.
	HorizontalPosition AnchorLayoutPosition

	// VerticalPosition specifies the vertical anchoring position.
	VerticalPosition AnchorLayoutPosition

	// StretchHorizontal specifies whether to stretch in the horizontal direction.
	StretchHorizontal bool

	// StretchVertical specifies whether to stretch in the vertical direction.
	StretchVertical bool
}

const (
	// AnchorLayoutPositionStart is the anchoring position for "left" (in the horizontal direction) or "top" (in the vertical direction.)
	AnchorLayoutPositionStart = AnchorLayoutPosition(iota)

	// AnchorLayoutPositionCenter is the center anchoring position.
	AnchorLayoutPositionCenter

	// AnchorLayoutPositionEnd is the anchoring position for "right" (in the horizontal direction) or "bottom" (in the vertical direction.)
	AnchorLayoutPositionEnd
)

// AnchorLayoutOpts contains functions that configure an AnchorLayout.
var AnchorLayoutOpts AnchorLayoutOptions

// NewAnchorLayout constructs a new AnchorLayout, configured by opts.
func NewAnchorLayout(opts ...AnchorLayoutOpt) *AnchorLayout {
	a := &AnchorLayout{}

	for _, o := range opts {
		o(a)
	}

	return a
}

// Padding configures an anchor layout to use padding i.
func (o AnchorLayoutOptions) Padding(i Insets) AnchorLayoutOpt {
	return func(a *AnchorLayout) {
		a.padding = i
	}
}

// PreferredSize implements Layouter.
func (a *AnchorLayout) PreferredSize(widgets []PreferredSizeLocateableWidget) (int, int) {
	px, py := a.padding.Dx(), a.padding.Dy()

	if len(widgets) == 0 {
		return px, py
	}

	w, h := widgets[0].PreferredSize()
	return w + px, h + py
}

// Layout implements Layouter.
func (a *AnchorLayout) Layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) {
	if len(widgets) == 0 {
		return
	}

	widget := widgets[0]
	if widget.GetWidget().Visibility == Visibility_Hide {
		return
	}

	ww, wh := widget.PreferredSize()
	rect = a.padding.Apply(rect)
	wx := 0
	wy := 0

	if ald, ok := widget.GetWidget().LayoutData.(AnchorLayoutData); ok {
		wx, wy, ww, wh = a.applyLayoutData(ald, wx, wy, ww, wh, rect)
	}

	r := image.Rect(0, 0, ww, wh)
	r = r.Add(image.Point{wx, wy})
	r = r.Add(rect.Min)

	widget.SetLocation(r)
}

func (a *AnchorLayout) applyLayoutData(ld AnchorLayoutData, wx int, wy int, ww int, wh int, rect image.Rectangle) (int, int, int, int) {
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
