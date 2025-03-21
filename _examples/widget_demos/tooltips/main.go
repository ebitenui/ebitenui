package main

import (
	"bytes"
	"fmt"
	img "image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

func main() {
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	// load button text font
	face, _ := loadFont(20)

	tooltipContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
		widget.ContainerOpts.AutoDisableChildren(),
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 230, A: 255})),
	)

	var button *widget.Button

	for i := 1; i < 6; i++ {
		iter := i
		Option1 := widget.NewText(
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
			widget.TextOpts.Text(fmt.Sprint("Label: ", iter), face, color.White),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)),
		)
		tooltipContainer.AddChild(Option1)
	}

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(30),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
		)),
	)

	// construct a button
	button = widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.ToolTip(widget.NewToolTip(
				widget.ToolTipOpts.Content(tooltipContainer),
				//widget.WidgetToolTipOpts.Delay(1*time.Second),
				widget.ToolTipOpts.Offset(img.Point{-5, 5}),
				widget.ToolTipOpts.Position(widget.TOOLTIP_POS_WIDGET),
				// When the Position is set to TOOLTIP_POS_WIDGET, you can configure where it opens with the optional parameters below
				// They will default to what you see below if you do not provide them
				widget.ToolTipOpts.AnchorOriginHorizontal(widget.TOOLTIP_ANCHOR_END),
				widget.ToolTipOpts.AnchorOriginVertical(widget.TOOLTIP_ANCHOR_END),
				widget.ToolTipOpts.ContentOriginHorizontal(widget.TOOLTIP_ANCHOR_END),
				widget.ToolTipOpts.ContentOriginVertical(widget.TOOLTIP_ANCHOR_START),
			)),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Hover for tooltip #1", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button clicked")
		}),
	)
	// add the button as a child of the container
	rootContainer.AddChild(button)

	// Use the NewTextToolTip convenience method to create the tooltip
	btn2ToolTip := widget.NewTextToolTip("Label: 1\nLabel: 2\nLabel: 3\nLabel: 4\nLabel: 5",
		face, color.White,
		image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 230, A: 255}))

	// The NewTextToolTip defaults to follow the cursor
	// But every parameter is available to update after it has been created
	btn2ToolTip.Position = widget.TOOLTIP_POS_CURSOR_STICKY

	// construct a button
	button2 := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.ToolTip(btn2ToolTip),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Hover for tooltip #2", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button2 clicked")
		}),
	)
	// add the button2 as a child of the container
	rootContainer.AddChild(button2)

	// Use the NewTextToolTip convenience method to create the tooltip
	btn3ToolTip := widget.NewTextToolTip("This Tooltip uses\n'widget.TOOLTIP_POS_ABSOLUTE'\nto always appear at X: 200 / Y: 100!",
		face, color.White,
		image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 230, A: 255}))

	// The NewTextToolTip defaults to follow the cursor
	// But every parameter is available to update after it has been created
	btn3ToolTip.Position = widget.TOOLTIP_POS_ABSOLUTE
	btn3ToolTip.Offset.X = 200
	btn3ToolTip.Offset.Y = 100
	btn3ToolTip.ContentOriginHorizontal = widget.TOOLTIP_ANCHOR_MIDDLE
	btn3ToolTip.ContentOriginVertical = widget.TOOLTIP_ANCHOR_MIDDLE
	btn3ToolTip.Delay = 0

	// construct a button
	button3 := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.ToolTip(btn3ToolTip),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Hover for tooltip #3", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button3 clicked")
		}),
	)
	// add the button2 as a child of the container
	rootContainer.AddChild(button3)

	// Use the NewTextToolTip convenience method to create the tooltip
	btn4ToolTip := widget.NewTextToolTip("This Tooltip uses\n'widget.TOOLTIP_POS_SCREEN'\nto always appear center screen!",
		face, color.White,
		image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 230, A: 255}))

	// The NewTextToolTip defaults to follow the cursor
	// But every parameter is available to update after it has been created
	btn4ToolTip.Position = widget.TOOLTIP_POS_SCREEN
	btn4ToolTip.Offset.X = 0
	btn4ToolTip.Offset.Y = 0
	btn4ToolTip.ContentOriginHorizontal = widget.TOOLTIP_ANCHOR_MIDDLE
	btn4ToolTip.ContentOriginVertical = widget.TOOLTIP_ANCHOR_MIDDLE
	btn4ToolTip.AnchorOriginHorizontal = widget.TOOLTIP_ANCHOR_MIDDLE
	btn4ToolTip.AnchorOriginVertical = widget.TOOLTIP_ANCHOR_MIDDLE
	btn4ToolTip.Delay = 0

	// construct a button
	button4 := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.ToolTip(btn4ToolTip),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Hover for tooltip #4", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button4 clicked")
		}),
	)
	// add the button2 as a child of the container
	rootContainer.AddChild(button4)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Tooltips")

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

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("loadFont: %w", err)
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
