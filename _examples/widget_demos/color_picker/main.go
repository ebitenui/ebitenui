package main

import (
	"fmt"
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

		// the container will use an anchor layout to layout its children
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// a swatch showing the color picked
	swatch := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x40, 0xa0, 0xff, 0xff})),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
			widget.WidgetOpts.MinSize(200, 40),
		),
	)

	// construct a color picker
	picker := widget.NewColorPicker(
		// Set the wheel's diameter
		widget.ColorPickerOpts.WheelSize(200),
		// Set the color the picker starts with
		widget.ColorPickerOpts.InitialColor(color.NRGBA{0x40, 0xa0, 0xff, 0xff}),
		// Add a row of preset colors under the slider
		widget.ColorPickerOpts.Swatches(24,
			color.NRGBA{0x40, 0xa0, 0xff, 0xff}, color.NRGBA{0xff, 0x9f, 0x43, 0xff}, color.NRGBA{0x5f, 0xd3, 0x8d, 0xff},
			color.NRGBA{0xff, 0x6a, 0xd5, 0xff}, color.NRGBA{0xff, 0xd7, 0x5e, 0xff}, color.White),
		// Set the brightness slider's handle image and size
		widget.ColorPickerOpts.SliderHandle(&widget.ButtonImage{
			Idle:    image.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
			Hover:   image.NewNineSliceColor(color.NRGBA{230, 230, 230, 255}),
			Pressed: image.NewNineSliceColor(color.NRGBA{255, 255, 255, 255}),
		}, 12),
		widget.ColorPickerOpts.WidgetOpts(
			// Set the picker to layout in the center on the screen
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		// Set the callback to call when the color changes
		widget.ColorPickerOpts.ChangedHandler(func(args *widget.ColorPickerChangedEventArgs) {
			swatch.SetBackgroundImage(image.NewNineSliceColor(args.Color))
			fmt.Println(args.Color, "dragging", args.Dragging)
		}),
	)
	rootContainer.AddChild(picker, swatch)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Color Picker")

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
