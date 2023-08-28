package main

import (
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

	// construct a new container for top of window
	topContainer := widget.NewContainer(
		// the container will use a plain red color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0xff, 0x00, 0x00, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			// Set how far apart to space the children
			widget.RowLayoutOpts.Spacing(15),
			// Padding between elements
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
		)),
	)

	// construct a new container for middle left of window
	middleLeftContainer := widget.NewContainer(
		// the container will use a plain red/green color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0xff, 0xff, 0x00, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			// Which direction to layout children
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			// Set how far apart to space the children
			widget.RowLayoutOpts.Spacing(15),
			// Padding between elements
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
	)

	// construct a new container for middle middle of window
	middleMiddleContainer := widget.NewContainer(
		// the container will use a plain red/green color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0xff, 0xff, 0xff, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			// Which direction to layout children
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			// Set how far apart to space the children
			widget.RowLayoutOpts.Spacing(15),
			// Padding between elements
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(5)),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
		),
	)

	// construct a new container for middle right of window
	middleRightContainer := widget.NewContainer(
		// the container will use a plain green/blue color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x00, 0xff, 0xff, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			//Which direction to layout children
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			// Spacing between elements
			widget.RowLayoutOpts.Spacing(15),
			// Padding between elements
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
		)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Stretch: true,
			}),
			widget.WidgetOpts.MinSize(100, 100),
		),
	)

	// construct a new container for middle of window
	middleContainer := widget.NewContainer(
		// the container will use a plain green color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x00, 0xff, 0x00, 0xff})),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			// Which direction to layout children
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
		)),
	)
	middleContainer.AddChild(middleLeftContainer)
	middleContainer.AddChild(middleMiddleContainer)
	middleContainer.AddChild(middleRightContainer)

	// construct top button
	topButton := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Hide top container (blocking)", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   5,
			Right:  5,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			topContainer.GetWidget().Visibility = widget.Visibility_Hide_Blocking
		}),
	)
	topContainer.AddChild(topButton)

	// construct middle buttons
	middleShowButton := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Show blue container (non blocking)", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   5,
			Right:  5,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			middleRightContainer.GetWidget().Visibility = widget.Visibility_Show
		}),
	)
	middleHideButton := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Hide blue container (non blocking)", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   5,
			Right:  5,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			middleRightContainer.GetWidget().Visibility = widget.Visibility_Hide
		}),
	)
	middleLeftContainer.AddChild(middleShowButton)
	middleLeftContainer.AddChild(middleHideButton)

	// construct top button
	buttonShowTop := widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Show top container (blocking)", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   5,
			Right:  5,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			topContainer.GetWidget().Visibility = widget.Visibility_Show
		}),
	)
	middleLeftContainer.AddChild(buttonShowTop)

	// build up rootContainer
	rootContainer := widget.NewContainer(
		// the container will use a plain white color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0xff, 0xff, 0xff, 0xff})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			// Define number of columns in the grid
			widget.GridLayoutOpts.Columns(1),
			// Define how to stretch the rows and columns. Note it is required to
			// specify the Stretch for each row and column.
			widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{false, true}),
			// Padding between elements
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    10,
				Bottom: 10,
				Left:   10,
				Right:  10,
			}),
		)),
	)

	rootContainer.AddChild(topContainer)
	rootContainer.AddChild(middleContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("Ebiten UI - Visibility")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

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
