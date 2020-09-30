package widget

import (
	"image"
)

// TODO: RootLayout should probably reside in "internal" subpackage

type RootLayout struct {
	layout   Layouter
	widgets  []PreferredSizeLocateableWidget
	lastRect image.Rectangle
}

func NewRootLayout(w PreferredSizeLocateableWidget) *RootLayout {
	r := RootLayout{
		layout:  NewFillLayout(),
		widgets: []PreferredSizeLocateableWidget{w},
	}
	r.MarkDirty()
	return &r
}

func (r *RootLayout) MarkDirty() {
	r.layout.(Dirtyable).MarkDirty()
}

func (r *RootLayout) PreferredSize(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) (int, int) {
	return 0, 0
}

func (r *RootLayout) Layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) {
	r.layout.Layout(r.widgets, rect)
}

func (r *RootLayout) LayoutRoot(rect image.Rectangle) {
	if rect != r.lastRect {
		r.MarkDirty()
	}

	r.Layout(r.widgets, rect)
}
