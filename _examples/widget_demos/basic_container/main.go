package main

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game object used by ebitengine
type game struct {
	//Our main UI object
	ui *ebitenui.UI
}

func main() {
	// construct a new container that serves as the root of the UI hierarchy
	// The root container will fill the entire window set up by ebitengine
	rootContainer := widget.NewContainer(
		// The container will use a plain color as its background. This is not required if you wish
		// the container to be transparent.
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
		// Containers have the concept of a Layout. This is how children of this container should be
		// displayed within the bounds of this container.
		// The container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	//From here you can add other containers and widgets to the container.
	//Please see the button widget example to continue.
	//rootContainer.AddChild(button)

	// construct the UI object
	// This ui object has an Update method and Draw method that must be
	// called in the ebitengine Update and Draw methods as seen below.
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Basic Container")

	game := game{
		ui: &ui,
	}

	// run Ebiten main loop
	err := ebiten.RunGame(&game)
	if err != nil {
		log.Println(err)
	}
}

// Layout implements Game.
func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Update implements Game.
func (g *game) Update() error {
	// update the UI
	g.ui.Update()
	return nil
}

// Draw implements Ebiten's Draw method.
func (g *game) Draw(screen *ebiten.Image) {
	// draw the UI onto the screen. Typically you will want this to be the
	// last line in the Draw call to ensure the UI is drawn on top of everything else.
	g.ui.Draw(screen)
}
