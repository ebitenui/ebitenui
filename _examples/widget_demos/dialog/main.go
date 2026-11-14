package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

// Game object used by ebiten
type game struct {
	ui *ebitenui.UI
}

func main() {
	face, err := loadFont(18)
	if err != nil {
		log.Fatal(err)
	}
	buttonImage := &widget.ButtonImage{
		Idle:    eimage.NewNineSliceColor(color.NRGBA{R: 40, G: 52, B: 80, A: 255}),
		Hover:   eimage.NewNineSliceColor(color.NRGBA{R: 58, G: 76, B: 116, A: 255}),
		Pressed: eimage.NewNineSliceColor(color.NRGBA{R: 30, G: 40, B: 62, A: 255}),
	}
	buttonText := &widget.ButtonTextColor{Idle: color.White}

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0xff})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)
	ui := &ebitenui.UI{Container: rootContainer}

	// a button that opens a confirmation dialog
	rootContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
		})),
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Open dialog", &face, buttonText),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(10)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			dialog := widget.NewDialog(
				// Set the title, the message and the choices, one button each
				widget.DialogOpts.Title("Dust storm"),
				widget.DialogOpts.Message("A storm is coming over the landing site. Recall the rover?"),
				widget.DialogOpts.Choices("Recall", "Keep going"),
				// Wrap the title and the message at a width
				widget.DialogOpts.MaxWidth(280),
				// Set the look of the dialog
				widget.DialogOpts.TitleText(&face, color.White),
				widget.DialogOpts.MessageText(&face, color.NRGBA{R: 190, G: 200, B: 215, A: 255}),
				widget.DialogOpts.ButtonParams(&widget.ButtonParams{
					Image:       buttonImage,
					TextFace:    &face,
					TextColor:   buttonText,
					TextPadding: widget.NewInsetsSimple(8),
				}),
				widget.DialogOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{R: 16, G: 20, B: 32, A: 255})),
				// Ease the dialog in from below, fading in, and back out when it closes
				widget.DialogOpts.WindowOpts(widget.WindowOpts.Transition(widget.WindowTransition{
					Duration: 15,
					Offset:   image.Pt(0, 40),
					Fade:     true,
				})),
				// Set the callback to call when a choice is made; the dialog then closes
				widget.DialogOpts.ChoiceHandler(func(args *widget.DialogChoiceEventArgs) {
					fmt.Println("chose", args.Choice, args.Label)
				}),
			)
			// Place the dialog in the middle of the screen and show it
			dialog.CenterIn(image.Rect(0, 0, 640, 480))
			ui.AddWindow(dialog.Window())
		}),
	))

	// Ebiten setup
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Ebiten UI - Dialog")

	game := game{
		ui: ui,
	}

	// run Ebiten main loop
	err = ebiten.RunGame(&game)
	if err != nil {
		log.Println(err)
	}
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}
	return &text.GoTextFace{Source: s, Size: size}, nil
}

// Layout implements Game.
func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return 640, 480
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
