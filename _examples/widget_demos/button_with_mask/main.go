package main

import (
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten
type game struct {
	ui  *ebitenui.UI
	btn *widget.Button
}

func main() {
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	// load button text font
	face, _ := loadFont(20)

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	// construct a button
	button := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify to ignore transparent pixels for mouse events
		widget.ButtonOpts.IgnoreTransparentPixels(true),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Hello!", face, &widget.ButtonTextColor{
			Idle:  color.NRGBA{255, 255, 255, 255},
			Hover: color.NRGBA{0, 0, 0, 255},
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button clicked")
		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor entered button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
		}),

		// add a handler that reacts to moving the cursor on the button
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor moved on button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY, "diffX =", args.DiffX, "diffY =", args.DiffY)
		}),

		// add a handler that reacts to exiting the button with the cursor
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			println("cursor exited button: entered =", args.Entered, "offsetX =", args.OffsetX, "offsetY =", args.OffsetY)
		}),

		// Indicate that this button should not be submitted when enter or space are pressed
		// widget.ButtonOpts.DisableDefaultKeys(),
	)

	// add the button as a child of the container
	rootContainer.AddChild(button)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Buttons with mask")

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

	//Test that you can call Click on the focused widget.
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		if btn, ok := g.ui.GetFocusedWidget().(*widget.Button); ok {
			btn.Click()
		}
	}

	return nil
}

// Draw implements Ebiten's Draw method.
func (g *game) Draw(screen *ebiten.Image) {
	// draw the UI onto the screen
	g.ui.Draw(screen)
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idleImage := ebiten.NewImage(100, 100)
	vector.DrawFilledCircle(idleImage, float32(idleImage.Bounds().Dx())/2, float32(idleImage.Bounds().Dy())/2, 45, color.Black, true)
	idle := image.NewNineSlice(idleImage, [3]int{0, idleImage.Bounds().Dx(), 0}, [3]int{0, idleImage.Bounds().Dy(), 0})

	hoverImage := ebiten.NewImage(100, 100)
	vector.DrawFilledCircle(hoverImage, float32(hoverImage.Bounds().Dx())/2, float32(hoverImage.Bounds().Dy())/2, 45, color.White, true)
	hover := image.NewNineSlice(hoverImage, [3]int{0, hoverImage.Bounds().Dx(), 0}, [3]int{0, hoverImage.Bounds().Dy(), 0})

	pressedImage := ebiten.NewImage(100, 100)
	vector.DrawFilledCircle(pressedImage, float32(pressedImage.Bounds().Dx())/2, float32(pressedImage.Bounds().Dy())/2, 45, color.NRGBA{R: 255, G: 0, B: 0, A: 255}, true)
	pressed := image.NewNineSlice(pressedImage, [3]int{0, pressedImage.Bounds().Dx(), 0}, [3]int{0, pressedImage.Bounds().Dy(), 0})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
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
