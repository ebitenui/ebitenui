package main

import (
	_ "image/png"
	"log"

	"github.com/blizzy78/ebitenui/_demo/app"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Ebiten UI Demo")
	ebiten.SetWindowResizable(true)
	ebiten.SetScreenClearedEveryFrame(false)

	ui, closeUI, err := app.CreateUI()
	if err != nil {
		log.Fatal(err)
	}

	defer closeUI()
	instance := app.NewApp(ui)
	err = ebiten.RunGame(instance)
	if err != nil {
		log.Print(err)
	}
}
