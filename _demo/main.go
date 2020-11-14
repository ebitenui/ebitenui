package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	_ "image/png"

	"github.com/blizzy78/ebitenui"
	"github.com/blizzy78/ebitenui/_demo/demo"
)

type game struct {
	ui *ebitenui.UI
}

func main() {
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Ebiten UI Demo")
	ebiten.SetWindowResizable(true)
	ebiten.SetScreenClearedEveryFrame(false)

	ui, closeUI, err := demo.CreateUI()
	if err != nil {
		log.Fatal(err)
	}

	defer closeUI()

	game := game{
		ui: ui,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Print(err)
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
