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
	// Construct a new container that serves as the root of the UI hierarchy.
	rootContainer := widget.NewContainer(
		// The container will use a plain color as its background.
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// The container will use an anchor layout to lay out its single child.
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// Construct a container to hold the progress bars.
	progressBarsContainer := widget.NewContainer(
		// The container will use a vertical row layout to lay out the progress
		// bars in a vertical row.
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
		)),
		// Set the required anchor layout data to determine where in the root
		// container to place the progress bars.
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
	)

	// Construct a horizontal progress bar.
	hProgressbar := widget.NewProgressBar(
		widget.ProgressBarOpts.WidgetOpts(
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this exampe since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(200, 20),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the progress images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 0, 255, 255}),
			},
		),
		// Set the min, max, and current values.
		widget.ProgressBarOpts.Values(0, 10, 7),
		// Set how much of the track is displayed when the bar is overlayed.
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Top:    2,
			Bottom: 2,
		}),
	)
	// Construct a vertical inverted progress bar.
	vProgressbar := widget.NewProgressBar(
		// Set the direction of the progress bar to vertical.
		widget.ProgressBarOpts.Direction(widget.DirectionVertical),
		// Invert the progress bar, meaning here it will fill from the bottom to the top
		// since it’s vertical.
		widget.ProgressBarOpts.Inverted(true),
		widget.ProgressBarOpts.WidgetOpts(
			// Set the minimum size for the progress bar.
			// This is necessary if you wish to have the progress bar be larger than
			// the provided track image. In this example since we are using NineSliceColor
			// which is 1px x 1px we must set a minimum size.
			widget.WidgetOpts.MinSize(20, 200),

			// Set the progress bar in the middle of the row container cell
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ProgressBarOpts.Images(
			// Set the track images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the progress images (Idle, Hover, Disabled).
			&widget.ProgressBarImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{0, 255, 0, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{0, 255, 0, 255}),
			},
		),
		// Set the min, max, and current values.
		widget.ProgressBarOpts.Values(0, 10, 4),
		// Set how much of the track is displayed when the bar is overlayed.
		widget.ProgressBarOpts.TrackPadding(widget.Insets{
			Left:  2,
			Right: 2,
		}),
	)
	/*
		To update a progress bar programmatically you can use
		hProgressbar.SetCurrent(value)
		hProgressbar.GetCurrent()
		hProgressbar.Min = 5
		hProgressbar.Max = 10
	*/
	// Add the progress bars as a child of their container.
	progressBarsContainer.AddChild(hProgressbar)
	progressBarsContainer.AddChild(vProgressbar)
	// Add the progress bars container as a child of the root container.
	rootContainer.AddChild(progressBarsContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - ProgressBar")

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
