package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game object used by ebiten
type game struct {
	ui       *ebitenui.UI
	checkBox *widget.Checkbox
}

func main() {
	game := game{}
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

	// load button text font
	// face, _ := loadFont(20)

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	uncheckedImage := ebiten.NewImage(20, 20)
	uncheckedImage.Fill(color.White)

	checkedImage := ebiten.NewImage(20, 20)
	checkedImage.Fill(color.NRGBA{255, 255, 0, 255})

	game.checkBox = widget.NewCheckbox(
		widget.CheckboxOpts.ButtonOpts(
			widget.ButtonOpts.WidgetOpts(
				// Set the location of the checkbox
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
				// Set the minimum size of the checkbox
				widget.WidgetOpts.MinSize(30, 30),
			),
			// Set the background images - idle, hover, pressed
			widget.ButtonOpts.Image(buttonImage),

			// This disables space and enter triggering the checkbox
			// widget.ButtonOpts.DisableDefaultKeys(),
		),
		// Set the check object images
		widget.CheckboxOpts.Image(&widget.CheckboxGraphicImage{
			// When the checkbox is unchecked
			Unchecked: &widget.ButtonImageImage{
				Idle: uncheckedImage,
			},
			// When the checkbox is checked
			Checked: &widget.ButtonImageImage{
				Idle: checkedImage,
			},
		}),
		// Set the state change handler
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			if args.State == widget.WidgetChecked {
				fmt.Println("Checkbox is Checked")
			} else {
				fmt.Println("Checkbox is Unchecked")
			}
		}),
	)

	rootContainer.AddChild(game.checkBox)

	// construct the UI
	game.ui = &ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Checkbox")

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

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.checkBox.Click()
	}
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
