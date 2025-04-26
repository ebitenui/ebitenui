package main

import (
	"bytes"
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
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

//go:embed assets
var embeddedAssets embed.FS

// Game object used by ebiten.
type game struct {
	ui               *ebitenui.UI
	labeledCheckBox1 *widget.Checkbox
}

func main() {
	g := game{}
	// load images for button states: idle, hover, and pressed
	checkboxImage, _ := loadCheckboxImage()

	// load button text font
	face, _ := loadFont(20)

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(35),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			),
		),
	)

	g.labeledCheckBox1 = widget.NewCheckbox(
		// Set the labeled checkbox's position
		widget.CheckboxOpts.WidgetOpts(
			// Set the location of the checkbox
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
			widget.WidgetOpts.MinSize(30, 30),
		),

		// Set the images
		widget.CheckboxOpts.Image(checkboxImage),

		// Set the label
		widget.CheckboxOpts.Text("Labeled Checkbox1", &face, &widget.LabelColor{
			Idle:     color.White,
			Disabled: color.White,
		}),

		// Set the spacing between the label and the checkbox
		widget.CheckboxOpts.Spacing(15),

		// Set the label to be before the checkbox.
		// widget.CheckboxOpts.LabelFirst(),

		// Set the state change handler
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			if args.State == widget.WidgetChecked {
				fmt.Println("Checkbox1 is Checked")
			} else {
				fmt.Println("Checkbox1 is Unchecked")
			}
		}),
	)
	rootContainer.AddChild(g.labeledCheckBox1)

	labeledCheckBox2 := widget.NewCheckbox(
		// Set the labeled checkbox's position
		widget.CheckboxOpts.WidgetOpts(
			// Set the location of the checkbox
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
			// Set the minimum size of the checkbox
			widget.WidgetOpts.MinSize(30, 30),
		),

		// Set the images
		widget.CheckboxOpts.Image(checkboxImage),

		// Set the label
		widget.CheckboxOpts.Text("Labeled Checkbox2", &face, &widget.LabelColor{
			Idle:     color.White,
			Disabled: color.White,
		}),

		// Set the spacing between the label and the checkbox
		widget.CheckboxOpts.Spacing(15),

		// Set the label to be before the checkbox.
		widget.CheckboxOpts.LabelFirst(),

		// Set the state change handler
		widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
			if args.State == widget.WidgetChecked {
				fmt.Println("Checkbox2 is Checked")
			} else {
				fmt.Println("Checkbox2 is Unchecked")
			}
		}),
	)
	// Set this checkbox as Checked by default
	labeledCheckBox2.SetState(widget.WidgetChecked)

	rootContainer.AddChild(labeledCheckBox2)
	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Labeled  Checkbox")

	g.ui = &ui

	// run Ebiten main loop
	err := ebiten.RunGame(&g)
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
		g.labeledCheckBox1.Click()
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

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
