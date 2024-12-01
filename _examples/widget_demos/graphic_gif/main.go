package main

import (
	"bytes"
	_ "embed"
	"image/gif"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed ebitenui.gif
var data []byte

// Like button example, but use image instead of text for the label.
// Game object used by Ebitengine.
type game struct {
	ui *ebitenui.UI
}

func main() {
	g, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	// Construct a new container that serves as the root of the UI hierarchy.
	root := widget.NewContainer(
		// The container will use an anchor layout to layout its single child widget.
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// Since our button is a multi-widget object, add its wrapping container.
	renderer := widget.NewGraphic(
		widget.GraphicOpts.GIF(g),
	)
	root.AddChild(renderer)

	game := game{
		ui: &ebitenui.UI{
			Container: root,
		},
	}

	// Ebitengine setup.
	ebiten.SetWindowSize(600, 400)
	ebiten.SetWindowTitle("Ebiten UI - Graphic GIF")
	if err = ebiten.RunGame(&game); err != nil {
		panic(err)
	}
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *game) Update() error {
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}
