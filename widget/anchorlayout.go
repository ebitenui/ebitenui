package widget

import "image"

// AnchorLayout layouts widgets anchored to either a corner or edge of a rectangle,
// optionally stretching it in one or both directions.
//
// AnchorLayout will layout all widgets  in the container to the specified locations regardless of overlap.
// The widgets in the container will be drawn in the order they were added to the container.
//
// Widget.LayoutData of widgets being layouted by AnchorLayout need to be of type AnchorLayoutData.
type AnchorLayout struct {
	padding *Insets
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

	// Sets the padding for the child.
	Padding *Insets
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

	if a.padding == nil {
		a.padding = &Insets{}
	}

	return a
}

// Padding configures an anchor layout to use padding i. This affects all children.
func (o AnchorLayoutOptions) Padding(i *Insets) AnchorLayoutOpt {
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
	for idx := range widgets {
		widget := widgets[idx]
		if widget.GetWidget().GetVisibility() == Visibility_Hide {
			continue
		}

		ww, wh := widget.PreferredSize()
		wrect := a.padding.Apply(rect)
		wx := 0
		wy := 0

		if ald, ok := widget.GetWidget().LayoutData.(AnchorLayoutData); ok {
			if ald.Padding != nil {
				wrect = ald.Padding.Apply(wrect)
			}
			wx, wy, ww, wh = a.applyLayoutData(ald, wx, wy, ww, wh, wrect)
		}

		r := image.Rect(0, 0, ww, wh)
		r = r.Add(image.Point{wx, wy})
		r = r.Add(wrect.Min)

		widget.SetLocation(r)
	}
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
	case AnchorLayoutPositionStart:
		// Do nothing
	}

	switch vPos {
	case AnchorLayoutPositionCenter:
		wy = (rect.Dy() - wh) / 2
	case AnchorLayoutPositionEnd:
		wy = rect.Dy() - wh
	case AnchorLayoutPositionStart:
		// Do nothing
	}

	return wx, wy, ww, wh
}
