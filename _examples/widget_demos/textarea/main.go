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
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten,
type game struct {
	ui *ebitenui.UI
}

func main() {
	face, _ := loadFont(14)
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)),
	)

	// construct a textarea
	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				// Set the layout data for the textarea
				//including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
					// MaxWidth:  300,
					MaxHeight: 100,
					Stretch:   true,
				}),
				// Set the minimum size for the widget
				widget.WidgetOpts.MinSize(300, 100),
			),
		),
		// Set gap between scrollbar and text
		widget.TextAreaOpts.ControlWidgetSpacing(2),
		// Tell the textarea to display bbcodes
		widget.TextAreaOpts.ProcessBBCode(true),
		// Tell the textarea to remove any unknown BBCodes
		widget.TextAreaOpts.StripBBCode(true),
		// Set the font color
		widget.TextAreaOpts.FontColor(color.Black),
		// Set the font face (size) to use
		widget.TextAreaOpts.FontFace(&face),
		widget.TextAreaOpts.TextPadding(widget.Insets{
			Right: 16,
		}),
		// Set the initial text for the textarea
		// It will automatically line wrap and process newlines characters
		// If ProcessBBCode is true it will parse out bbcode
		widget.TextAreaOpts.Text("[link=a]Hello[/link] [color=#FFF000]World[/color] Hello [b]World[/b]\n Hello [b]World[/b]\n Hello [b]World[/b]\n Hello [b]World[/b]\n Hello \n[b]World[/b]\n[link=b]Hello[/link] \n[link=c]Hello[/link] "),
		// Tell the TextArea to show the vertical scrollbar
		widget.TextAreaOpts.ShowVerticalScrollbar(),
		// Set padding between edge of the widget and where the text is drawn
		// widget.TextAreaOpts.TextPadding(widget.NewInsetsSimple(20)),
		// This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerImage(&widget.ScrollContainerImage{
			Idle: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			Mask: image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
		}),
		// This sets the images to use for the sliders
		widget.TextAreaOpts.SliderParams(&widget.SliderParams{
			// Set the track images
			TrackImage: &widget.SliderTrackImage{
				Idle:  image.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
				Hover: image.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
			},
			// Set the handle images
			HandleImage: &widget.ButtonImage{
				Idle:    image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Hover:   image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				Pressed: image.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
			},
		}),
		widget.TextAreaOpts.LinkClickedEvent(func(args *widget.LinkEventArgs) {
			fmt.Println("Link Clicked Id: ", args.Id, " value: ", args.Value, " args: ", args.Args,
				" offsetX/offsetY ", args.OffsetX, "/", args.OffsetY)
		}),
		widget.TextAreaOpts.LinkCursorEnteredEvent(func(args *widget.LinkEventArgs) {
			fmt.Println("Link Entered Id: ", args.Id, " value: ", args.Value, " args: ", args.Args,
				" offsetX/offsetY ", args.OffsetX, "/", args.OffsetY)
		}),
		widget.TextAreaOpts.LinkCursorExitedEvent(func(args *widget.LinkEventArgs) {
			fmt.Println("Link Exited Id: ", args.Id, " value: ", args.Value, " args: ", args.Args,
				" offsetX/offsetY ", args.OffsetX, "/", args.OffsetY)
		}),
	)

	// Add text to the end of the textarea
	// textarea.AppendText("\nLast Row")
	// Add text to the beginning of the textarea
	// textarea.PrependText("First Row\n")
	// Replace the current text with the new value
	// textarea.SetText("New Value!")
	// Retrieve the current value of the text area text
	fmt.Println(textarea.GetText())
	// add the textarea as a child of the container
	rootContainer.AddChild(textarea)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - TextArea")
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
