package widget

import (
	"image"
)

// TODO: RootLayout should probably reside in "internal" subpackage

type RootLayout struct {
	widget PreferredSizeLocateableWidget
}

func NewRootLayout(w PreferredSizeLocateableWidget) *RootLayout {
	r := RootLayout{
		widget: w,
	}
	return &r
}

func (r *RootLayout) PreferredSize(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) (int, int) {
	return 0, 0
}

func (r *RootLayout) Layout(widgets []PreferredSizeLocateableWidget, rect image.Rectangle) {
	// unused
}

func (r *RootLayout) LayoutRoot(rect image.Rectangle) {
	r.widget.SetLocation(rect)
}
