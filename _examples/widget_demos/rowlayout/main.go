package main

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

/*
The Row Layout is built to position children in a row or column based on their preferred size.
*/
func main() {

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			//Which direction to layout children
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			//Set how much padding before displaying content
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			//Set how far apart to space the children
			widget.RowLayoutOpts.Spacing(15),
		)),
	)

	innerContainer1 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{255, 0, 0, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				//Specify where within the row or column this element should be positioned.
				Position: widget.RowLayoutPositionStart,
				//Should this widget be stretched across the row or column
				Stretch: false,
				//How wide can this element grow to (override preferred widget size)
				MaxWidth: 100,
				//How tall can this element grow to (override preferred widget size)
				MaxHeight: 100,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	rootContainer.AddChild(innerContainer1)

	innerContainer2 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 255, 0, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				//Specify where within the row or column this element should be positioned.
				Position: widget.RowLayoutPositionCenter,
				//Should this widget be stretched across the row or column
				Stretch: true,
				//How wide can this element grow to (override preferred widget size)
				MaxWidth: 200,
				//How tall can this element grow to (override preferred widget size)
				MaxHeight: 100,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	rootContainer.AddChild(innerContainer2)

	innerContainer3 := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				//Specify where within the row or column this element should be positioned.
				Position: widget.RowLayoutPositionEnd,
				//Should this widget be stretched across the row or column
				Stretch: true,
				//How wide can this element grow to (override preferred widget size)
				MaxWidth: 400,
				//How tall can this element grow to (override preferred widget size)
				MaxHeight: 100,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)
	rootContainer.AddChild(innerContainer3)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Row Layout")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

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
	// draw the UI onto the screen
	g.ui.Draw(screen)
}
