package main

import (
	"fmt"
	"image/gif"
	"log"
	"os"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// Like button example, but use image instead of text for the label.

// Game object used by Ebitengine.
type game struct {
	ui *ebitenui.UI
}

func main() {
	gif := loadGIF()

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	gifW := widget.NewGraphic(
		widget.GraphicOpts.GIF(gif),
	)

	// since our button is a multi-widget object, add its wrapping container
	rootContainer.AddChild(gifW)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(600, 400)
	ebiten.SetWindowTitle("Ebiten UI - Graphic GIF")

	game := game{ui: &ui}

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
	// draw the UI onto the screen
	g.ui.Draw(screen)
}

func loadGIF() *gif.GIF {
	f, err := os.Open("./ebiten-ui.gif")
	if err != nil {
		panic(fmt.Errorf("Failing to load the gif is most likely due to the execution path, it's intended to be from the root of the project like 'go run ./_examples/widget_demos/graphic_gif/': %w", err))
	}
	defer f.Close()

	g, err := gif.DecodeAll(f)
	if err != nil {
		panic(err)
	}
	return g
}
