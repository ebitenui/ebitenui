package main

import (
	"image"
	"sync"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"
	"github.com/blizzy78/ebitenui/widget"

	"github.com/hajimehoshi/ebiten"
)

// UI encapsulates a complete user interface that can be rendered onto the screen.
// There should only be exactly one UI per application.
type UI struct {
	Container *widget.Container

	init     sync.Once
	layout   *widget.RootLayout
	lastRect image.Rectangle
}

// Update updates u. This function should be called in the Ebiten Update function.
func (u *UI) Update() {
	u.init.Do(u.initUI)

	event.FireDeferredEvents()
	input.Update()
}

// Draw renders u onto screen, with rect as the area reserved for rendering.
// This function should be called in the Ebiten Draw function.
func (u *UI) Draw(screen *ebiten.Image, rect image.Rectangle) {
	u.init.Do(u.initUI)

	input.Draw()

	defer func() {
		u.lastRect = rect
	}()

	if rect != u.lastRect {
		u.Container.RequestRelayout()
	}

	input.SetupInputLayersWithDeferred(u.Container)
	u.layout.LayoutRoot(rect)
	widget.RenderWithDeferred(u.Container, screen)
}

func (u *UI) initUI() {
	u.layout = widget.NewRootLayout(u.Container)
}
