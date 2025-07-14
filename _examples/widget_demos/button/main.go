package main

import (
	"image/color"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Engine struct {
	ui      *ebitenui.UI
	enabled bool
	update  func(bool)
}

func (engine *Engine) Init() {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	c1 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	c2 := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(10),
		)),
	)

	colors := []color.RGBA{
		{R: 255, G: 0, B: 0, A: 255},   // Red
		{R: 0, G: 255, B: 0, A: 255},   // Green
		{R: 0, G: 0, B: 255, A: 255},   // Blue
		{R: 255, G: 255, B: 0, A: 255}, // Yellow
		{R: 255, G: 165, B: 0, A: 255}, // Orange
		{R: 128, G: 0, B: 128, A: 255}, // Purple
	}

	setupContainer := func(set bool) {
		if set {
			for _, col := range colors {
				img := ebiten.NewImage(50, 50)
				img.Fill(col)
				c1.AddChild(widget.NewGraphic(
					widget.GraphicOpts.Image(img),
				))
			}
		} else {
			c1.RemoveChildren()
		}
	}

	setupContainer(true)

	for range len(colors) {
		img := ebiten.NewImage(50, 50)
		img.Fill(color.White)
		c2.AddChild(widget.NewGraphic(
			widget.GraphicOpts.Image(img),
		))
	}

	rootContainer.AddChild(c1, c2)

	ui := &ebitenui.UI{
		Container: rootContainer,
	}

	engine.ui = ui
	engine.enabled = true
	engine.update = setupContainer
}

func (engine *Engine) Update() error {
	keys := inpututil.AppendJustPressedKeys(nil)

	for _, key := range keys {
		switch key {
		case ebiten.KeyEscape, ebiten.KeyCapsLock:
			return ebiten.Termination
		case ebiten.KeyT:
			engine.enabled = !engine.enabled
			engine.update(engine.enabled)
		}
	}

	engine.ui.Update()

	return nil
}

func (engine *Engine) Draw(screen *ebiten.Image) {
	engine.ui.Draw(screen)
}

func (engine *Engine) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, Ebiten!")

	var engine Engine
	engine.Init()

	if err := ebiten.RunGame(&engine); err != nil {
		panic(err)
	}
}
