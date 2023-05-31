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

func main() {
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// construct a progressbar
	progressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			//Set the required anchor layout data to determine where in
			//the container to place the progressbar
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			//Set the minimum size for the progressbar.
			//This is necessary if you wish to have the progressbar be larger than
			//the provided track image. In this exampe since we are using NineSliceColor
			//which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(200, 20),
		),
		widget.ProgressBarOpts.Images(
			//Set the track images (Idle, Hover, Disabled)
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			//Set the progress images (Idle, Hover, Disabled)
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		//Set the min, max, and current values
		widget.ProgressBarOpts.Values(0, 10, 7),
		//Set how much of the track is displayed when the bar is overlayed.
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    2,
			Bottom: 2,
			Left:   0,
			Right:  0,
		}),
	)
	/*
		To update the progressbar programmatically you can use
		progressbar.SetCurrent(value)
		progressbar.GetCurrent()
		progressbar.Min = 5
		progressbar.Max = 10
	*/
	// add the progressbar as a child of the container
	rootContainer.AddChild(progressbar)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Progressbar")

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
