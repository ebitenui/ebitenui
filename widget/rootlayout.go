package widget

import (
	"image"
)

type RootLayout struct {
	layout   Layouter
	widgets  []HasWidget
	lastRect image.Rectangle
}

func NewRootLayout(w HasWidget) *RootLayout {
	r := RootLayout{
		layout:  NewFillLayout(),
		widgets: []HasWidget{w},
	}
	r.MarkDirty()
	return &r
}

func (r *RootLayout) MarkDirty() {
	r.layout.(Dirtyable).MarkDirty()
}

func (r *RootLayout) PreferredSize(widgets []HasWidget, rect image.Rectangle) (width int, height int) {
	return 0, 0
}

func (r *RootLayout) Layout(widgets []HasWidget, rect image.Rectangle) {
	r.layout.Layout(r.widgets, rect)
}

func (r *RootLayout) LayoutRoot(rect image.Rectangle) {
	if rect != r.lastRect {
		r.MarkDirty()
	}

	r.Layout(r.widgets, rect)
}
