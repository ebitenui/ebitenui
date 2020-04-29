package main

import (
	"image"
	"sync"

	"github.com/blizzy78/ebitenui/event"
	"github.com/blizzy78/ebitenui/input"
	"github.com/blizzy78/ebitenui/widget"

	"github.com/hajimehoshi/ebiten"
)

type UI struct {
	Container *widget.Container

	init     sync.Once
	layout   *widget.RootLayout
	lastRect image.Rectangle
}

func (u *UI) Update() {
	u.init.Do(u.initUI)

	event.FireDeferredEvents()
	input.Update()
}

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
