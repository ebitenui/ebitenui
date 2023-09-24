package widget

import (
	"image"
)

// StackedLayout lays out multiple widgets stacked on top of each other in the order they are added as children
// Each child will have the dimensions equal to the max size of all the StackedLayout's children
//
// Note: Events will propogate to each layer. e.g. if you have overlapping buttons, a click event over both will trigger both events.
//
// Widget.LayoutData of widgets being layed out by StackedLayout should be left empty.
type StackedLayout struct {
	padding Insets
}

// StackedLayoutOpt is a function that configures a.
type StackedLayoutOpt func(a *StackedLayout)

type StackedLayoutOptions struct {
}

// StackedLayoutData specifies layout settings for a widget.
type StackedLayoutData struct {
}

// StackedLayoutOpts contains functions that configure an StackedLayout.
var StackedLayoutOpts StackedLayoutOptions

// NewStackedLayout constructs a new StackedLayout, configured by opts.
func NewStackedLayout(opts ...StackedLayoutOpt) *StackedLayout {
	a := &StackedLayout{}

	for _, o := range opts {
		o(a)
	}

	return a
}

// Padding configures an Stacked layout to use padding i.
func (o StackedLayoutOptions) Padding(i Insets) StackedLayoutOpt {
	return func(a *StackedLayout) {
		a.padding = i
	}
}

// PreferredSize implements Layouter.
func (a *StackedLayout) PreferredSize(widgets []PreferredSizeLocateableWidget) (int, int) {
	px, py := a.padding.Dx(), a.padding.Dy()

	if len(widgets) == 0 {
		return px, py
	}
	var w, h int
	for idx, widget := range widgets {
		if widget.GetWidget().Visibility == Visibility_Hide {
			continue
		}

		w1, h1 := widgets[idx].PreferredSize()
		if w1 > w {
			w = w1
		}
		if h1 > h {
			h = h1
		}
	}
	return w + px, h + py
}

// Layout implements Layouter.
func (a *StackedLayout) Layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) {
	if len(widgets) == 0 {
		return
	}
	rect = a.padding.Apply(rect)
	for idx := range widgets {
		widgets[idx].SetLocation(rect)
	}
}
