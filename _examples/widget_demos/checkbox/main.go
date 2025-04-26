package main

import (
	"embed"
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed assets
var embeddedAssets embed.FS

// Game object used by ebiten.
type game struct {
	ui       *ebitenui.UI
	checkBox *widget.Checkbox
}

func main() {
	game := game{}
	// load images for button states: idle, hover, and pressed
	checkboxImage, _ := loadCheckboxImage()

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	game.checkBox = widget.NewCheckbox(
		widget.CheckboxOpts.WidgetOpts(
			// Set the location of the checkbox
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		// Set the check object images
		widget.CheckboxOpts.Image(checkboxImage),
		// Set the state change handler
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			switch args.State {
			case widget.WidgetChecked:
				fmt.Println("Checkbox is Checked")
			case widget.WidgetGreyed:
				fmt.Println("Checkbox is Greyed")
			case widget.WidgetUnchecked:
				fmt.Println("Checkbox is Unchecked")
			}
		}),
		widget.CheckboxOpts.TriState(),
		widget.CheckboxOpts.InitialState(widget.WidgetChecked),
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

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.checkBox.GetWidget().Disabled = !g.checkBox.GetWidget().Disabled
	}

	return nil
}

// Draw implements Ebiten's Draw method.
func (g *game) Draw(screen *ebiten.Image) {
	// draw the UI onto the screen
	g.ui.Draw(screen)
}

func loadCheckboxImage() (*widget.CheckboxImage, error) {
	f1, err := embeddedAssets.Open("assets/checkbox-idle.png")
	if err != nil {
		return nil, err
	}
	defer f1.Close()
	idle, _, _ := ebitenutil.NewImageFromReader(f1)

	f2, err := embeddedAssets.Open("assets/checkbox-checked.png")
	if err != nil {
		return nil, err
	}
	defer f2.Close()
	checked, _, _ := ebitenutil.NewImageFromReader(f2)

	f3, err := embeddedAssets.Open("assets/checkbox-greyed.png")
	if err != nil {
		return nil, err
	}
	defer f3.Close()
	greyed, _, _ := ebitenutil.NewImageFromReader(f3)

	f4, err := embeddedAssets.Open("assets/checkbox-hover.png")
	if err != nil {
		return nil, err
	}
	defer f4.Close()
	idle_hovered, _, _ := ebitenutil.NewImageFromReader(f4)

	f5, err := embeddedAssets.Open("assets/checkbox-checked-hover.png")
	if err != nil {
		return nil, err
	}
	defer f5.Close()
	checked_hovered, _, _ := ebitenutil.NewImageFromReader(f5)

	f6, err := embeddedAssets.Open("assets/checkbox-greyed-hover.png")
	if err != nil {
		return nil, err
	}
	defer f6.Close()
	greyed_hovered, _, _ := ebitenutil.NewImageFromReader(f6)

	f7, err := embeddedAssets.Open("assets/checkbox-disabled.png")
	if err != nil {
		return nil, err
	}
	defer f7.Close()
	idle_disabled, _, _ := ebitenutil.NewImageFromReader(f7)

	f8, err := embeddedAssets.Open("assets/checkbox-checked-disabled.png")
	if err != nil {
		return nil, err
	}
	defer f8.Close()
	checked_disabled, _, _ := ebitenutil.NewImageFromReader(f8)

	f9, err := embeddedAssets.Open("assets/checkbox-greyed-disabled.png")
	if err != nil {
		return nil, err
	}
	defer f9.Close()
	greyed_disabled, _, _ := ebitenutil.NewImageFromReader(f9)

	return &widget.CheckboxImage{
		Unchecked:         image.NewFixedNineSlice(idle),
		Checked:           image.NewFixedNineSlice(checked),
		Greyed:            image.NewFixedNineSlice(greyed),
		UncheckedHovered:  image.NewFixedNineSlice(idle_hovered),
		CheckedHovered:    image.NewFixedNineSlice(checked_hovered),
		GreyedHovered:     image.NewFixedNineSlice(greyed_hovered),
		UncheckedDisabled: image.NewFixedNineSlice(idle_disabled),
		CheckedDisabled:   image.NewFixedNineSlice(checked_disabled),
		GreyedDisabled:    image.NewFixedNineSlice(greyed_disabled),
	}, nil
}
