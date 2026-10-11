package main

import (
	"embed"
	"image/color"
	"log"
	"strings"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/event"
	uihtml "github.com/ebitenui/ebitenui/html"
	"github.com/ebitenui/ebitenui/themes"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

// Embed the separate HTML file so the example works from any working directory
// and can also run in a browser. For an external file, use uihtml.LoadFile.
//
//go:embed layout.html
var layouts embed.FS

const initialMessage = "Enter your name, then click Say hello."

func main() {
	g, err := newGame()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(760, 420)
	ebiten.SetWindowTitle("Ebiten UI HTML Layout")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func newGame() (*game, error) {
	file, err := layouts.Open("layout.html")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Load creates the entire widget tree; the theme supplies its visual resources.
	doc, err := uihtml.Load(file, uihtml.Options{Theme: themes.GetBasicDarkTheme()})
	if err != nil {
		return nil, err
	}

	// IDs return the actual concrete objects, so normal ebitenui APIs still work.
	name := doc.GetByID("name").(*widget.TextInput)
	enthusiastic := doc.GetByID("enthusiastic").(*widget.Checkbox)
	message := doc.GetByID("message").(*widget.Text)
	greet := doc.GetByID("greet").(*widget.Button)
	reset := doc.GetByID("reset").(*widget.Button)

	greet.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		who := strings.TrimSpace(name.GetText())
		if who == "" {
			who = "world"
		}
		punctuation := "."
		if enthusiastic.State() == widget.WidgetChecked {
			punctuation = "!"
		}
		message.Label = "Hello, " + who + punctuation
		doc.Root.RequestRelayout()
	}))
	reset.ClickedEvent.AddHandler(event.WrapHandler(func(args *widget.ButtonClickedEventArgs) {
		name.SetText("")
		enthusiastic.SetState(widget.WidgetChecked)
		message.Label = initialMessage
		doc.Root.RequestRelayout()
	}))

	return &game{ui: &ebitenui.UI{Container: doc.Root}}, nil
}

type game struct {
	ui *ebitenui.UI
}

func (g *game) Update() error {
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{R: 24, G: 26, B: 32, A: 255})
	g.ui.Draw(screen)
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}
