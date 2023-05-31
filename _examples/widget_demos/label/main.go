package main

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

func main() {

	// load button text font
	face, _ := loadFont(20)

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
			widget.RowLayoutOpts.Spacing(20),
		)),
	)

	/**
		There are two ways to create a label within ebitenui.
		1) widget.NewText - This is a simple method to create a text label.
		   You will want to use this one unless you have a need to provide a 
			 separate Disabled color for the text label you're displaying.

		2) widget.NewLabel - This is a more complex method to create a text label that can be disabled.
			You will want to use this one when you want to be able to provide a separate 
			Disabled color for the text label you're displaying

	*/
	
	label1 := widget.NewText(
		widget.TextOpts.Text("Label 1 (NewText)", face, color.White),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
	)
	//Set the widget as Disabled. This does not affect NewText
	label1.GetWidget().Disabled = true
	// Add the first Text as a child of the container
	rootContainer.AddChild(label1)

	// Create a new label
	label2 := widget.NewLabel(
		widget.LabelOpts.Text("Label 2 (NewLabel - Enabled)", face, &widget.LabelColor{
			Idle:     color.White,
			Disabled: color.NRGBA{100, 100, 100, 255},
		}),
		widget.LabelOpts.TextOpts(
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
			widget.TextOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		),
	)
	// Add the first label as a child of the container
	rootContainer.AddChild(label2)

	// Create a new label
	label3 := widget.NewLabel(
		widget.LabelOpts.Text("Label 3 (NewLabel - Disabled)", face, &widget.LabelColor{
			Idle:     color.White,
			Disabled: color.NRGBA{100, 100, 100, 255},
		}),
		widget.LabelOpts.TextOpts(
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
			widget.TextOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),
		),
	)
	// Add the second label as a child of the container
	rootContainer.AddChild(label3)
	// Set this label as disabled and tells it to use the 
	// Disabled color.
	label3.GetWidget().Disabled = true

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Label")

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

func loadFont(size float64) (font.Face, error) {
	ttfFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(ttfFont, &truetype.Options{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}
