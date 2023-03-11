package main

import (
	"fmt"
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
	// load images for button states: idle, hover, and pressed
	buttonImage, _ := loadButtonImage()

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

	uncheckedImage := ebiten.NewImage(20, 20)
	uncheckedImage.Fill(color.White)

	checkedImage := ebiten.NewImage(20, 20)
	checkedImage.Fill(color.NRGBA{255, 255, 0, 255})

	labeledCheckBox1 := widget.NewLabeledCheckbox(
		//Set the labeled checkbox's position
		widget.LabeledCheckboxOpts.WidgetOpts(
			//Set the location of the checkbox
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
		),
		//Set the checkbox Opts
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(
				widget.ButtonOpts.WidgetOpts(
					//Set the minimum size of the checkbox
					widget.WidgetOpts.MinSize(30, 30),
				),
				//Set the background images - idle, hover, pressed
				widget.ButtonOpts.Image(buttonImage),
			),
			//Set the check object images
			widget.CheckboxOpts.Image(&widget.CheckboxGraphicImage{
				//When the checkbox is unchecked
				Unchecked: &widget.ButtonImageImage{
					Idle: uncheckedImage,
				},
				//When the checkbox is checked
				Checked: &widget.ButtonImageImage{
					Idle: checkedImage,
				},
			}),
			//Set the state change handler
			widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				if args.State == widget.WidgetChecked {
					fmt.Println("Checkbox1 is Checked")
				} else {
					fmt.Println("Checkbox1 is Unchecked")
				}
			}),
		),
		//Set the label
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Labeled Checkbox1", face, &widget.LabelColor{
			Idle:     color.White,
			Disabled: color.White,
		})),
		//Set the spacing between the label and the checkbox
		widget.LabeledCheckboxOpts.Spacing(15),
		//Set the label to be before the checkbox.
		//widget.LabeledCheckboxOpts.LabelFirst(),
	)
	rootContainer.AddChild(labeledCheckBox1)

	labeledCheckBox2 := widget.NewLabeledCheckbox(
		//Set the labeled checkbox's position
		widget.LabeledCheckboxOpts.WidgetOpts(
			//Set the location of the checkbox
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
		),
		//Set the checkbox Opts
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(
				widget.ButtonOpts.WidgetOpts(
					//Set the minimum size of the checkbox
					widget.WidgetOpts.MinSize(30, 30),
				),
				//Set the background images - idle, hover, pressed
				widget.ButtonOpts.Image(buttonImage),
			),
			//Set the check object images
			widget.CheckboxOpts.Image(&widget.CheckboxGraphicImage{
				//When the checkbox is unchecked
				Unchecked: &widget.ButtonImageImage{
					Idle: uncheckedImage,
				},
				//When the checkbox is checked
				Checked: &widget.ButtonImageImage{
					Idle: checkedImage,
				},
			}),
			//Set the state change handler
			widget.CheckboxOpts.StateChangedHandler(func(args *widget.CheckboxChangedEventArgs) {
				if args.State == widget.WidgetChecked {
					fmt.Println("Checkbox2 is Checked")
				} else {
					fmt.Println("Checkbox2 is Unchecked")
				}
			}),
		),
		//Set the label
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Labeled Checkbox2", face, &widget.LabelColor{
			Idle:     color.White,
			Disabled: color.White,
		})),
		//Set the spacing between the label and the checkbox
		widget.LabeledCheckboxOpts.Spacing(15),
		//Set the label to be before the checkbox.
		widget.LabeledCheckboxOpts.LabelFirst(),
	)
	//Set this checkbox as Checked by default
	labeledCheckBox2.SetState(widget.WidgetChecked)

	rootContainer.AddChild(labeledCheckBox2)
	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Labeled  Checkbox")

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
