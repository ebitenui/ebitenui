package main

import (
	"embed"
	"image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	e_image "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

//go:embed assets
var embeddedAssets embed.FS

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

func main() {
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	// load button text font
	face, _ := loadFont(20)

	ui := ebitenui.UI{}
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	winContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(e_image.NewNineSliceColor(color.NRGBA{155, 155, 0, 255})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()), // the container will use a plain color as its background
	)
	winContainer.AddChild(widget.NewText(widget.TextOpts.Text("Click outside to close.\nResizable.\nUses System cursor for E/W.\nUses Custom cursor for N/S", face, color.White)))

	win := widget.NewWindow(widget.WindowOpts.CloseMode(widget.CLICK_OUT), widget.WindowOpts.Contents(winContainer), widget.WindowOpts.Resizeable())

	// construct a button
	button := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
			widget.WidgetOpts.CursorHovered("buttonHover"),
			widget.WidgetOpts.CursorPressed("buttonPressed"),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
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
			win.SetLocation(image.Rect(50, 50, 350, 150))
			ui.AddWindow(win)
		}),
	)

	// add the button as a child of the container
	rootContainer.AddChild(button)

	// construct the UI
	ui.Container = rootContainer

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Cursor Shape")

	game := game{
		ui: &ui,
	}

	//Set the main cursor used within the application
	input.SetCursorImage(input.CURSOR_DEFAULT, loadNormalCursorImage())

	//Set the custom hover image
	input.SetCursorImage("buttonHover", loadHoverCursorImage())

	input.SetCursorImage("buttonPressed", loadPressedCursorImage())

	//Set the NS resize cursor with an offset so that it shows up a little above the cursor point
	input.SetCursorImageWithOffset(input.CURSOR_NSRESIZE, loadNSCursorImage(), image.Point{0, -6})

	//Disable cursor management by ebitenui
	//input.CursorManagementEnabled = false

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
	idle := e_image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := e_image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := e_image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

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

func loadNormalCursorImage() *ebiten.Image {
	f, err := embeddedAssets.Open("assets/cursor.png")
	if err != nil {
		return nil
	}
	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return ebiten.NewImageFromImage(i.SubImage(image.Rect(0, 0, 16, 16)))
}

func loadNSCursorImage() *ebiten.Image {
	f, err := embeddedAssets.Open("assets/cursor.png")
	if err != nil {
		return nil
	}
	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return ebiten.NewImageFromImage(i.SubImage(image.Rect(0, 16, 16, 32)))
}

func loadHoverCursorImage() *ebiten.Image {
	f, err := embeddedAssets.Open("assets/cursor.png")
	if err != nil {
		return nil
	}
	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return ebiten.NewImageFromImage(i.SubImage(image.Rect(16, 0, 32, 16)))
}

func loadPressedCursorImage() *ebiten.Image {
	f, err := embeddedAssets.Open("assets/cursor.png")
	if err != nil {
		return nil
	}
	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return ebiten.NewImageFromImage(i.SubImage(image.Rect(32, 0, 48, 16)))
}
