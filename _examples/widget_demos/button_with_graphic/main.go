package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten.
type game struct {
	ui  *ebitenui.UI
	btn *widget.Button
}

func main() {
	// load images for button states: idle, hover, and pressed.
	buttonImage, _ := loadButtonImage()
	buttonIcon := loadButtonIcon()
	buttonDisabledIcon := loadDisabledButtonIcon()

	// load button text font.
	face, _ := loadFont(20)

	// construct a new container that serves as the root of the UI hierarchy.
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background.
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget.
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(10),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)))),
	)

	var button *widget.Button
	// construct a button.
	button = widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically.
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		// specify the images to use.
		widget.ButtonOpts.Image(buttonImage),

		// specify the button graphics
		widget.ButtonOpts.Graphic(&widget.GraphicImage{
			Idle:     buttonIcon,
			Disabled: buttonDisabledIcon,
		}),

		// specify the button's text, the font face, and the color.
		widget.ButtonOpts.Text("Button With Graphic", &face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextProcessBBCode(false),
		// specify that the button's text needs some padding for correct display.
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button.
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button clicked")
		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor entered button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
		}),

		// add a handler that reacts to entering the button with the cursor.
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor entered button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
		}),

		// add a handler that reacts to moving the cursor on the button.
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor moved on button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY, "diffX =", args.DiffX, "diffY =", args.DiffY)
		}),

		// add a handler that reacts to exiting the button with the cursor.
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor exited button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
		}),

		// Indicate that this button should not be submitted when enter or space are pressed
		// widget.ButtonOpts.DisableDefaultKeys(),
	)

	// add the button as a child of the container.
	rootContainer.AddChild(button)

	button2 := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically.
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(70, 70),
		),
		// specify the images to use.
		widget.ButtonOpts.Image(buttonImage),

		// specify the button graphics
		widget.ButtonOpts.Graphic(&widget.GraphicImage{
			Idle:     buttonIcon,
			Disabled: buttonDisabledIcon,
		}),
		widget.ButtonOpts.TextProcessBBCode(false),
		// specify that the button's text needs some padding for correct display.
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(4)),

		// add a handler that reacts to clicking the button.
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button2 clicked")
		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor entered button2: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
		}),

		// add a handler that reacts to entering the button with the cursor.
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor entered button2: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
		}),

		// add a handler that reacts to moving the cursor on the button.
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor moved on button2: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY, "diffX =", args.DiffX, "diffY =", args.DiffY)
		}),

		// add a handler that reacts to exiting the button with the cursor.
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor exited button2: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
		}),

		// Indicate that this button should not be submitted when enter or space are pressed
		// widget.ButtonOpts.DisableDefaultKeys(),
	)

	// add the button as a child of the container.
	rootContainer.AddChild(button2)
	// construct the UI.
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Buttons")

	game := game{
		ui:  &ui,
		btn: button,
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
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {
		g.btn.Click()
	}

	// Test that you can call Click on the focused widget.
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if btn, ok := g.ui.GetFocusedWidget().(*widget.Button); ok {
			btn.Click()
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.btn.Press()
	} else if inpututil.IsKeyJustReleased(ebiten.KeyG) {
		g.btn.Release()
	}

	return nil
}

// Draw implements Ebiten's Draw method.
func (g *game) Draw(screen *ebiten.Image) {
	// draw the UI onto the screen
	g.ui.Draw(screen)
}

func loadButtonImage() (*widget.ButtonImage, error) {

	idle := image.NewBorderedNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255}, color.NRGBA{90, 90, 90, 255}, 3)

	hover := image.NewBorderedNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255}, color.NRGBA{70, 70, 70, 255}, 3)

	pressed := image.NewAdvancedNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255}, image.NewBorder(3, 2, 2, 2, color.NRGBA{70, 70, 70, 255}))

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
		return nil, fmt.Errorf("Error loading font: %w", err)
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func loadButtonIcon() *ebiten.Image {
	// we'll use a circle as an icon image
	// in reality it could be an arbitrary *ebiten.Image
	icon := ebiten.NewImage(32, 32)
	ebitenutil.DrawCircle(icon, 16, 16, 16, color.RGBA{R: 0x71, G: 0x56, B: 0xbd, A: 255})
	return icon
}

func loadDisabledButtonIcon() *ebiten.Image {
	// we'll use a circle as an icon image
	// in reality it could be an arbitrary *ebiten.Image
	icon := ebiten.NewImage(32, 32)
	ebitenutil.DrawCircle(icon, 16, 16, 16, color.RGBA{R: 250, G: 0x56, B: 0xbd, A: 255})
	return icon
}
