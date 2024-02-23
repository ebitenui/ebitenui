package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

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

		// the container will use an grid layout to layout its ScrollableContainer and Slider
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Spacing(2, 0),
			widget.GridLayoutOpts.Stretch([]bool{true, false}, []bool{true}),
		)),
	)

	//Create the container with the content that should be scrolled
	content := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewRowLayout(
		widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		widget.RowLayoutOpts.Spacing(20),
	)))

	//Add 20 buttons to the scrollable content container
	for x := 0; x < 20; x++ {
		//Capture x for use in callback
		x := x
		// construct a button
		button := widget.NewButton(
			// set general widget options
			widget.ButtonOpts.WidgetOpts(
				// instruct the container's anchor layout to center the button both horizontally and vertically
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
			),

			// specify the images to use
			widget.ButtonOpts.Image(buttonImage),

			// specify the button's text, the font face, and the color
			widget.ButtonOpts.Text(fmt.Sprintf("Hello, World! - %d", x), face, &widget.ButtonTextColor{
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
				println(fmt.Sprintf("Button %d Clicked!", x))
			}),
		)

		// add the button as a child of the container
		content.AddChild(button)
	}

	//Create the new ScrollContainer object
	scrollContainer := widget.NewScrollContainer(
		//Set the content that will be scrolled
		widget.ScrollContainerOpts.Content(content),
		//Tell the container to stretch the content width to match available space
		widget.ScrollContainerOpts.StretchContentWidth(),
		//Set the background images for the scrollable container
		widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
			Idle: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
			Mask: image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff}),
		}),
	)
	//Add the scrollable container to the left side of the window
	rootContainer.AddChild(scrollContainer)

	//Create a function to return the page size used by the slider
	pageSizeFunc := func() int {
		return int(math.Round(float64(scrollContainer.ViewRect().Dy()) / float64(content.GetWidget().Rect.Dy()) * 1000))
	}
	//Create a vertical Slider bar to control the ScrollableContainer
	vSlider := widget.NewSlider(
		widget.SliderOpts.Direction(widget.DirectionVertical),
		widget.SliderOpts.MinMax(0, 1000),
		widget.SliderOpts.PageSizeFunc(pageSizeFunc),
		//On change update scroll location based on the Slider's value
		widget.SliderOpts.ChangedHandler(func(args *widget.SliderChangedEventArgs) {
			scrollContainer.ScrollTop = float64(args.Slider.Current) / 1000
		}),
		widget.SliderOpts.Images(
			// Set the track images
			&widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			},
			// Set the handle images
			&widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		),
	)
	//Set the slider's position if the scrollContainer is scrolled by other means than the slider
	scrollContainer.GetWidget().ScrolledEvent.AddHandler(func(args interface{}) {
		a := args.(*widget.WidgetScrolledEventArgs)
		p := pageSizeFunc() / 3
		if p < 1 {
			p = 1
		}
		vSlider.Current -= int(math.Round(a.Y * float64(p)))
	})

	//Add the slider to the second slot in the root container
	rootContainer.AddChild(vSlider)
	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Scroll Container")

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
