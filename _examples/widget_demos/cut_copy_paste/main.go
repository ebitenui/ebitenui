package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/utilities/mobile"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.design/x/clipboard"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten.
type game struct {
	ui *ebitenui.UI
	// This parameter is so you can keep track of the textInput widget to update and retrieve
	// its values in other parts of your game.
	standardTextInput *widget.TextInput
}

func main() {
	// Ebiten setup
	ebiten.SetWindowSize(400, 400)
	ebiten.SetWindowTitle("Ebiten UI - Cut Copy Paste")

	game := game{}

	// load the font
	face, _ := loadFont(20)

	// construct a new container that serves as the root of the UI hierarchy.
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),

		// the container will use a row layout to layout the textinput widgets
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)))),
	)

	// construct a standard textinput widget
	game.standardTextInput = widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			// Set the layout information to center the textbox in the parent.
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),

		// Set the keyboard type when opened on mobile devices.
		widget.TextInputOpts.MobileInputMode(mobile.TEXT),

		// Set the Idle and Disabled background image for the text input.
		// If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),

		// Set the font face and size for the widget.
		widget.TextInputOpts.Face(face),

		// Set the colors for the text and caret.
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{254, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),

		// Set how much padding there is between the edge of the input and the text.
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),

		// Set the font and width of the caret.
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 2),
		),

		// This text is displayed if the input is empty.
		widget.TextInputOpts.Placeholder("Standard Textbox"),

		// This is called when the user hits the "Enter" key.
		// There are other options that can configure this behavior.
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),

		// This is called whenver there is a change to the text.
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Changed: ", args.InputText)
		}),
	)

	rootContainer.AddChild(game.standardTextInput)

	// construct the UI.
	ui := ebitenui.UI{
		Container: rootContainer,
	}
	game.ui = &ui

	// Init returns an error if the package is not ready for use.
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	// run Ebiten main loop.
	err = ebiten.RunGame(&game)
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
	// Select all
	if input.KeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.standardTextInput.SelectAll()
	}

	g.HandleCCP()

	g.ui.Update()
	return nil
}

func (g *game) HandleCCP() {
	// Copy
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyC) {
		text := g.standardTextInput.SelectedText()
		if len(text) > 0 {
			clipboard.Write(clipboard.FmtText, []byte(text))
		}
	}

	// Cut
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyX) {
		text := g.standardTextInput.SelectedText()
		if len(text) > 0 {
			clipboard.Write(clipboard.FmtText, []byte(text))
			g.standardTextInput.DeleteSelectedText()
		}
	}

	// Paste
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
		clipVal := string(clipboard.Read(clipboard.FmtText))
		if len(clipVal) > 0 {
			g.standardTextInput.Insert(clipVal)
		}
	}
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
