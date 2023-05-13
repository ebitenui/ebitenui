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
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

//go:embed assets
var embeddedAssets embed.FS

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

type cursor_updater struct {
	currentPosition image.Point
	systemPosition  image.Point
	cursorImages    map[string]*ebiten.Image
}

func CreateUpdater() *cursor_updater {
	cu := cursor_updater{}
	X, Y := ebiten.CursorPosition()
	cu.systemPosition = image.Point{X, Y}
	cu.currentPosition = image.Point{X, Y}

	cu.cursorImages = make(map[string]*ebiten.Image)
	cu.cursorImages[input.CURSOR_DEFAULT] = loadNormalCursorImage()
	cu.cursorImages["buttonHover"] = loadHoverCursorImage()
	cu.cursorImages["buttonPressed"] = loadPressedCursorImage()
	return &cu
}

// Called every Update call from Ebiten
// Note that before this is called the current cursor shape is reset to DEFAULT every cycle
func (cu *cursor_updater) Update() {
	X, Y := ebiten.CursorPosition()
	diffX := cu.systemPosition.X - X
	diffY := cu.systemPosition.Y - Y
	cu.currentPosition.X -= diffX
	cu.currentPosition.Y -= diffY

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		cu.currentPosition.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		cu.currentPosition.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		cu.currentPosition.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		cu.currentPosition.Y += 2
	}

	cu.systemPosition = image.Point{X, Y}

}
func (cu *cursor_updater) Draw(screen *ebiten.Image) {
}
func (cu *cursor_updater) AfterDraw(screen *ebiten.Image) {
}

// MouseButtonPressed returns whether mouse button b is currently pressed.
func (cu *cursor_updater) MouseButtonPressed(b ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(b) || ebiten.IsKeyPressed(ebiten.KeySpace)
}

// MouseButtonJustPressed returns whether mouse button b has just been pressed.
// It only returns true during the first frame that the button is pressed.
func (cu *cursor_updater) MouseButtonJustPressed(b ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(b) || inpututil.IsKeyJustPressed(ebiten.KeySpace)
}

// CursorPosition returns the current cursor position.
// If you define a CursorPosition that doesn't align with a system cursor you will need to
// set the CursorDrawMode to Custom. This is because ebiten doesn't have a way to set the
// cursor location manually
func (cu *cursor_updater) CursorPosition() (int, int) {
	return cu.currentPosition.X, cu.currentPosition.Y
}

// Returns the image to use as the cursor
// EbitenUI by default will look for the following cursors:
//
//	"EWResize"
//	"NSResize"
//	"Default"
func (cu *cursor_updater) GetCursorImage(name string) *ebiten.Image {
	return cu.cursorImages[name]
}

// Returns how far from the CursorPosition to offset the cursor image.
// This is best used with cursors such as resizing.
func (cu *cursor_updater) GetCursorOffset(name string) image.Point {
	return image.Point{}
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
		}),
	)

	// add the button as a child of the container
	rootContainer.AddChild(button)

	// construct the UI
	ui.Container = rootContainer

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Cursor Updater")

	//Set the new updater
	input.SetCursorUpdater(CreateUpdater())

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
